// chat.go 实现 POST /api/chat：一个转发到 DeepSeek API 的聊天代理，给首页的
// 悬浮聊天助手用（回答"怎么捐款""里程碑是什么"之类的使用问题）。
// chat.go implements POST /api/chat: a proxy to the DeepSeek API for the
// floating chat widget on the homepage (answers "how do I donate", "what's
// a milestone" and similar usage questions).
//
// 这是 server.go 包注释里说的"三个例外"之外的第四个例外，但性质不一样——
// 前三个是往数据库里写字段，这个完全不碰数据库，只是转发一次 HTTP 请求。
// 之所以单独开一个后端接口而不是让前端直接调 DeepSeek API：API key 绝不能
// 出现在前端代码里（会被浏览器直接看到）。
// This is a fourth exception beyond the three in server.go's package doc, but
// a different kind — the first three write to the database; this one never
// touches it, it just forwards one HTTP request. It needs its own backend
// endpoint (instead of the frontend calling DeepSeek directly) because the
// API key can never live in frontend code — the browser would expose it.
//
// 加了两层保护，因为这个接口和前面几个不一样：调用是真的要花钱的（LLM API
// 按 token 计费），滥用会直接产生账单，不像"改错一个活动标题"那样零成本。
// Two safeguards here that the other endpoints don't need, because this call
// actually costs money (LLM APIs bill per token) — abuse has a real bill
// attached, unlike "someone sets the wrong campaign title".
//
// DeepSeek 的接口是 OpenAI 兼容格式（system prompt 放在 messages 数组第一条，
// 鉴权用 "Authorization: Bearer <key>"），跟 Anthropic 的 Messages API 不是
// 一套协议——这个文件只认 DeepSeek 这一种格式。
// DeepSeek's API is OpenAI-compatible (the system prompt is the first entry
// in the messages array; auth is "Authorization: Bearer <key>") — a
// different wire format from Anthropic's Messages API. This file only
// speaks the DeepSeek shape.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	maxChatMessageLen = 2000 // 单条消息最长字符数 / max chars per message
	maxChatHistory    = 12   // 只保留最近 N 条历史消息，控制 token 花费 / cap history sent, to bound token cost
	chatMaxTokens     = 400  // 回复最长 token 数 / max tokens in the reply

	chatRateLimit  = 20               // 每个 IP 每个窗口最多几次 / max requests per IP per window
	chatRateWindow = 10 * time.Minute // 限流窗口 / rate-limit window
)

// systemPrompt 让模型只谈 TrustChain 是什么、怎么用，别的随便应付一句就转回来。
// Keeps the model on-topic: what TrustChain is and how to use it.
const systemPrompt = `You are the help assistant embedded on TrustChain, a blockchain charity donation
platform (a university capstone project, running on the Sepolia test network — not real money).

How it works:
- Anyone can connect a MetaMask wallet to browse campaigns and donate ETH.
- A platform admin (the CharityRegistry contract owner) verifies charity addresses on-chain.
- A verified charity creates a campaign with one or more milestones, each with its own ETH amount.
- Donors send ETH to a campaign; funds sit locked in the TrustChainEscrow contract, not with the charity.
- A milestone is only released to the charity after they submit a receipt/proof on-chain — funds move
  one milestone at a time (Locked -> Released -> Proven), never as one lump sum up front.
- Every donation, release and receipt is a real on-chain event, visible on the Transparency Dashboard
  and on Sepolia Etherscan — that's the platform's core pitch: donors can verify where funds went
  instead of trusting the charity's word.

Pages on this site. When your answer involves one of them, point the user there with a Markdown link
written exactly as [label](/path) — use ONLY these paths, never invent others, never use external URLs:
- [Campaigns](/) — home page: how-it-works guide and the list of all campaigns; browse and donate here
- [Create a campaign](/create) — for verified charities: set the goal and milestones (with optional plan
  descriptions) for a new campaign
- [Transparency Dashboard](/transparency) — platform-wide totals and the live feed of on-chain events
  with Etherscan links
- [Charity verification](/admin/charities) — admin-only page for verifying charity addresses
- A single campaign's page is /campaigns/<id> (e.g. [Campaign 7](/campaigns/7)); it shows donations,
  milestones, receipts and that campaign's on-chain activity timeline. Only link to one if the user
  named a campaign number; otherwise send them to [Campaigns](/) to pick one.
Whenever your answer touches anything a page above covers (browsing or donating to campaigns, creating
one, checking where funds went, verifying a charity), you MUST include the matching link in that
answer — don't just name the page, link it. Link each page at most once per reply, inline in the
sentence.

Answer questions about how to use the site (connecting a wallet, donating, creating a campaign,
verifying a charity, what a milestone or "Confirming..." badge means) and about what makes this
platform different from a normal donation site. Keep answers short — a few sentences, no long essays.
If asked something unrelated to TrustChain or blockchain donations, say briefly that you can only help
with questions about this platform.`

type chatMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

// chatRequest is the body of POST /api/chat.
type chatRequest struct {
	Message string        `json:"message"`
	History []chatMessage `json:"history"` // prior turns, oldest first; optional
}

type chatResponse struct {
	Reply string `json:"reply"`
}

// --- 极简的按 IP 限流 / a minimal per-IP rate limiter ---
//
// 内存里存，重启就清空；单进程够用，不需要 Redis 之类的外部依赖——这个项目
// 的规模完全用不上更重的方案。
// In-memory, resets on restart; fine for a single process — this project's
// scale doesn't call for anything heavier (e.g. Redis).
type chatRateLimiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
}

func newChatRateLimiter() *chatRateLimiter {
	return &chatRateLimiter{hits: make(map[string][]time.Time)}
}

func (l *chatRateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-chatRateWindow)
	kept := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= chatRateLimit {
		l.hits[key] = kept
		return false
	}
	l.hits[key] = append(kept, now)
	return true
}

// handleChat: POST /api/chat — proxies one turn to the DeepSeek API.
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if s.deepSeekAPIKey == "" {
		writeError(w, http.StatusServiceUnavailable, "聊天功能未配置 / chat isn't configured on this deployment")
		return
	}
	if !s.chatLimiter.allow(clientIP(r)) {
		writeError(w, http.StatusTooManyRequests, "问得太快了，等一下再试 / too many requests, slow down")
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体不是合法 JSON / malformed JSON body")
		return
	}
	message := strings.TrimSpace(req.Message)
	if message == "" {
		writeError(w, http.StatusBadRequest, "message 不能为空 / message must not be empty")
		return
	}
	if len(message) > maxChatMessageLen {
		writeError(w, http.StatusBadRequest, "message 太长 / message too long")
		return
	}

	history := req.History
	if len(history) > maxChatHistory {
		history = history[len(history)-maxChatHistory:]
	}
	// DeepSeek 是 OpenAI 兼容格式：system prompt 是 messages 数组的第一条，
	// 不是像 Anthropic 那样单独一个字段。
	// DeepSeek is OpenAI-compatible: the system prompt is the first entry in
	// the messages array, not a separate top-level field like Anthropic's.
	messages := make([]chatMessage, 0, len(history)+2)
	messages = append(messages, chatMessage{Role: "system", Content: systemPrompt})
	for _, m := range history {
		if m.Role != "user" && m.Role != "assistant" {
			continue // 忽略非法角色，别把整条历史丢给上游报错 / skip bad roles rather than failing the whole call
		}
		messages = append(messages, chatMessage{Role: m.Role, Content: m.Content})
	}
	messages = append(messages, chatMessage{Role: "user", Content: message})

	reply, err := s.callDeepSeek(r.Context(), messages)
	if err != nil {
		writeError(w, http.StatusBadGateway, "聊天服务暂时不可用 / chat service is temporarily unavailable")
		return
	}
	writeJSON(w, http.StatusOK, chatResponse{Reply: reply})
}

type deepSeekRequest struct {
	Model     string        `json:"model"`
	MaxTokens int           `json:"max_tokens"`
	Messages  []chatMessage `json:"messages"`
}

type deepSeekResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// callDeepSeek sends one request to the DeepSeek chat-completions API and
// returns the assistant's text reply.
func (s *Server) callDeepSeek(ctx context.Context, messages []chatMessage) (string, error) {
	body, err := json.Marshal(deepSeekRequest{
		Model:     s.deepSeekModel,
		MaxTokens: chatMaxTokens,
		Messages:  messages,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.deepSeekAPIKey)

	httpClient := http.Client{Timeout: 20 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed deepSeekResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("deepseek: bad response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil {
			return "", fmt.Errorf("deepseek: %s", parsed.Error.Message)
		}
		return "", fmt.Errorf("deepseek: status %d", resp.StatusCode)
	}
	if len(parsed.Choices) == 0 || parsed.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("deepseek: no content in response")
	}
	return parsed.Choices[0].Message.Content, nil
}

// clientIP 取限流用的客户端标识：优先 X-Forwarded-For 的第一段，否则用
// RemoteAddr。不追求在恶意伪造头部时依然精确——这只是控制误用的软限流，
// 不是安全边界。
// clientIP picks an identifier for rate-limiting: X-Forwarded-For's first
// entry if present, otherwise RemoteAddr. Not hardened against a forged
// header — this is a soft abuse limiter, not a security boundary.
func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.IndexByte(fwd, ','); i != -1 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}
	return r.RemoteAddr
}

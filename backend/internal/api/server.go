// Package api 是读侧的 HTTP 层：把 store 的查询方法暴露成 REST 接口给前端。
// Package api is the read-side HTTP layer: it exposes store queries as REST
// endpoints for the frontend.
//
// 基本只读，三个例外 / mostly read-only, three exceptions:
//
//	所有链上状态的写（认证机构、建活动、捐款、放款）都由前端通过 MetaMask
//	直接发交易上链，再由索引器同步回库，这条边界没变。三个例外分两类：
//	  - POST /api/users/connect：记录"谁连过/什么角色/何时"，链上事件本来
//	    就不携带这个，索引器管不到。见 store/users.go。
//	  - POST /api/campaigns/{id}/metadata：活动的标题/详情/主题，同样是链上
//	    事件从不携带的数据（合约的 metadataHash 只是哈希，不是可解析指针）。
//	    见 store/campaign_metadata.go。同一类还有
//	    POST /api/campaigns/{id}/milestones/metadata（每个里程碑的计划说明），
//	    见 store/milestone_metadata.go。
//	  - POST /api/campaigns：**这个不太一样**——它写的是本该由索引器写的
//	    同一张 campaigns 表，只是抢在索引器追上之前，把前端已经从自己发出
//	    的交易里拿到的数据先乐观地存一份，让活动创建后立刻能在列表里看到，
//	    不用等 CONFIRMATIONS 个块+下一次轮询（1~2 分钟）。索引器真正处理到
//	    这个事件后会覆盖这行数据并标记为 confirmed，见 store.InsertCampaign
//	    对"confirmed"的说明。
//	All on-chain-state writes still happen on-chain via MetaMask, synced back
//	by the indexer — that boundary is unchanged. The three exceptions split
//	into two kinds:
//	  - POST /api/users/connect: records who connected, as what role, and
//	    when — data on-chain events never carried. See store/users.go.
//	  - POST /api/campaigns/{id}/metadata: a campaign's title/description/
//	    theme, also never carried on-chain (metadataHash is just a hash, not
//	    a resolvable pointer). See store/campaign_metadata.go. The same kind:
//	    POST /api/campaigns/{id}/milestones/metadata (each milestone's plan
//	    description), see store/milestone_metadata.go.
//	  - POST /api/campaigns: **different from the other two** — it writes
//	    the very same campaigns table the indexer owns, just optimistically,
//	    ahead of the indexer catching up, using data the frontend already
//	    has from the transaction it just sent — so a new campaign shows up
//	    in the list immediately instead of after CONFIRMATIONS blocks + the
//	    next poll (1-2 minutes). Once the indexer processes the real event
//	    it overwrites this row and marks it confirmed — see
//	    store.InsertCampaign's comment on "confirmed".
//
// 还有第四个例外，性质又不一样：
//   - POST /api/chat：不碰数据库，只是把请求转发给 DeepSeek API（key 只能
//     放后端，不能进前端代码），给首页的聊天助手用。见 chat.go。
//
// One more exception, of a different kind again:
//   - POST /api/chat: touches no database at all — it just proxies a
//     request to the DeepSeek API (the key must live server-side, never in
//     frontend code) for the homepage chat assistant. See chat.go.
//
// 第五个例外，写的是本地磁盘而不是数据库：
//   - POST /api/campaigns/{id}/milestones/{idx}/receipt：机构上传里程碑的
//     支出凭证文件（链上 receiptHash 只是个哈希，从没存过真正的文件）。
//     见 receipt_upload.go。
//
// A fifth exception, writing local disk instead of the database:
//   - POST /api/campaigns/{id}/milestones/{idx}/receipt: the charity uploads
//     the actual spending-receipt document (the on-chain receiptHash was
//     always just a hash, never the file). See receipt_upload.go.
package api

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/zoujunhui1/trust-chain/backend/internal/store"
)

// Server 持有存储层，并把各处理器挂到一个 http.Handler 上。
// Server holds the store and assembles the handlers into one http.Handler.
type Server struct {
	store *store.Store

	// 聊天助手用的配置，见 chat.go。deepSeekAPIKey 为空时 /api/chat 直接 503。
	// Chat assistant config, see chat.go. Blank deepSeekAPIKey => /api/chat returns 503.
	deepSeekAPIKey string
	deepSeekModel  string
	chatLimiter    *chatRateLimiter

	// 里程碑凭证文件的存放目录，见 receipt_upload.go。
	// Directory milestone receipt files are stored under. See receipt_upload.go.
	uploadsDir string
}

// New 构造 Server。deepSeekAPIKey 留空则聊天助手功能关闭（其它接口不受影响）。
// New builds a Server. A blank deepSeekAPIKey disables the chat assistant
// without affecting any other endpoint.
func New(st *store.Store, deepSeekAPIKey, deepSeekModel, uploadsDir string) *Server {
	return &Server{
		store:          st,
		deepSeekAPIKey: deepSeekAPIKey,
		deepSeekModel:  deepSeekModel,
		uploadsDir:     uploadsDir,
		chatLimiter:    newChatRateLimiter(),
	}
}

// Handler 注册所有路由并返回带中间件的 http.Handler。
// Handler registers all routes and returns the http.Handler with middleware.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// Go 1.22+ 的模式：方法 + 路径 + {占位符}，占位符用 r.PathValue 取。
	// Go 1.22+ patterns: method + path + {placeholder}, read via r.PathValue.
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/campaigns", s.handleListCampaigns)
	mux.HandleFunc("GET /api/campaigns/{id}", s.handleGetCampaign)
	mux.HandleFunc("GET /api/campaigns/{id}/donations", s.handleListDonations)
	mux.HandleFunc("GET /api/campaigns/{id}/activity", s.handleListCampaignActivity)
	mux.HandleFunc("GET /api/charities", s.handleListCharities)
	mux.HandleFunc("GET /api/charities/{address}", s.handleGetCharity)
	mux.HandleFunc("GET /api/activity", s.handleListActivity)
	mux.HandleFunc("POST /api/users/connect", s.handleConnectUser)
	mux.HandleFunc("POST /api/campaigns", s.handleCreateCampaign)
	mux.HandleFunc("POST /api/campaigns/{id}/metadata", s.handleSetCampaignMetadata)
	mux.HandleFunc("POST /api/campaigns/{id}/milestones/metadata", s.handleSetMilestoneDescriptions)
	mux.HandleFunc("POST /api/chat", s.handleChat)
	mux.HandleFunc("POST /api/campaigns/{id}/milestones/{idx}/receipt", s.handleUploadMilestoneReceipt)
	mux.HandleFunc("GET /api/campaigns/{id}/milestones/{idx}/receipt/file", s.handleGetMilestoneReceiptFile)

	// 用 CORS 中间件包一层，允许浏览器前端跨域访问。
	// Wrap with CORS so the browser frontend can call across origins.
	return withCORS(mux)
}

// ============================ 处理器 / handlers ============================

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleListCampaigns: GET /api/campaigns?limit=&offset=
func (s *Server) handleListCampaigns(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePaging(r)
	list, err := s.store.ListCampaigns(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// handleGetCampaign: GET /api/campaigns/{id} —— 返回活动详情 + 里程碑列表。
// Returns campaign detail plus its milestones.
func (s *Server) handleGetCampaign(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	camp, found, err := s.store.GetCampaign(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "活动不存在 / campaign not found")
		return
	}
	milestones, err := s.store.ListMilestones(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	// 把活动和它的里程碑合成一个对象返回。/ Combine campaign + milestones into one object.
	writeJSON(w, http.StatusOK, map[string]any{
		"campaign":   camp,
		"milestones": milestones,
	})
}

// handleListDonations: GET /api/campaigns/{id}/donations?limit=&offset=
func (s *Server) handleListDonations(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	limit, offset := parsePaging(r)
	list, err := s.store.ListDonations(r.Context(), id, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// handleListCampaignActivity: GET /api/campaigns/{id}/activity?limit=&offset=
// —— 单个活动的链上活动流水，按发生顺序正序，给活动详情页做真实交易时间线用。
// The on-chain activity feed for one campaign, oldest first — powers the
// real-transaction timeline on the campaign detail page.
func (s *Server) handleListCampaignActivity(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	limit, offset := parsePaging(r)
	list, err := s.store.ListActivityByCampaign(r.Context(), id, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// handleGetCharity: GET /api/charities/{address} —— 查机构认证状态。
// Looks up a charity's verification status.
func (s *Server) handleGetCharity(w http.ResponseWriter, r *http.Request) {
	raw := r.PathValue("address")
	// 校验并规范成小写（库里存的是小写），非法地址直接 400。
	// Validate and normalize to lowercase (DB stores lowercase); reject bad input with 400.
	if !common.IsHexAddress(raw) {
		writeError(w, http.StatusBadRequest, "地址格式错误 / invalid address")
		return
	}
	addr := strings.ToLower(common.HexToAddress(raw).Hex())

	c, found, err := s.store.GetCharity(r.Context(), addr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	if !found {
		// 查无记录：不算错误，返回「未认证」的默认视图，前端好处理。
		// Not found isn't an error: return an unverified default so the frontend has a shape.
		writeJSON(w, http.StatusOK, store.CharityView{Address: addr, Verified: false})
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// handleListCharities: GET /api/charities?limit=&offset= —— 所有见过的机构（含已撤销）。
// Lists every charity ever seen (including revoked ones).
func (s *Server) handleListCharities(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePaging(r)
	list, err := s.store.ListCharities(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// handleListActivity: GET /api/activity?limit=&offset= —— 全局链上活动流水，最新在前。
// Global on-chain activity feed, newest first.
func (s *Server) handleListActivity(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePaging(r)
	list, err := s.store.ListActivity(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// connectUserRequest is the body of POST /api/users/connect.
type connectUserRequest struct {
	Address string `json:"address"`
	Role    string `json:"role"` // admin / charity / donor — see store/users.go
}

var validUserRoles = map[string]bool{"admin": true, "charity": true, "donor": true}

// handleConnectUser: POST /api/users/connect — records that a wallet
// connected, with the role the frontend computed for it (lib/role.ts). See
// the package doc comment above for why this is the one write on this API.
func (s *Server) handleConnectUser(w http.ResponseWriter, r *http.Request) {
	var req connectUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体不是合法 JSON / malformed JSON body")
		return
	}
	if !common.IsHexAddress(req.Address) {
		writeError(w, http.StatusBadRequest, "地址格式错误 / invalid address")
		return
	}
	if !validUserRoles[req.Role] {
		writeError(w, http.StatusBadRequest, "role 必须是 admin/charity/donor 之一 / role must be admin, charity or donor")
		return
	}
	addr := strings.ToLower(common.HexToAddress(req.Address).Hex())

	if err := s.store.UpsertUser(r.Context(), addr, req.Role); err != nil {
		writeError(w, http.StatusInternalServerError, "写入失败 / write failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// setCampaignMetadataRequest is the body of POST /api/campaigns/{id}/metadata.
type setCampaignMetadataRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"` // optional
	Theme       *string `json:"theme"`       // optional — a key from frontend/src/lib/theme.ts
}

const (
	maxCampaignTitleLen       = 200
	maxCampaignDescriptionLen = 5000
	maxCampaignThemeLen       = 32
)

// handleSetCampaignMetadata: POST /api/campaigns/{id}/metadata — stores a
// campaign's title (required), description and theme (both optional; see
// store/campaign_metadata.go for why none of this is on-chain).
// No ownership check: like /api/users/connect, this API has no auth/session
// layer, so anyone who knows a campaign id can set its metadata. Acceptable
// for this project's scope — nothing of real value moves through it, and
// every action with real stakes (verify/create/donate/release) is still
// gated on-chain — but it's a real gap if this API is ever exposed beyond a
// demo.
func (s *Server) handleSetCampaignMetadata(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req setCampaignMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体不是合法 JSON / malformed JSON body")
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "title 不能为空 / title must not be empty")
		return
	}
	if len(title) > maxCampaignTitleLen {
		writeError(w, http.StatusBadRequest, "title 太长 / title too long")
		return
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		if trimmed == "" {
			req.Description = nil // treat blank as "not set"
		} else if len(trimmed) > maxCampaignDescriptionLen {
			writeError(w, http.StatusBadRequest, "description 太长 / description too long")
			return
		} else {
			req.Description = &trimmed
		}
	}
	if req.Theme != nil {
		trimmed := strings.TrimSpace(*req.Theme)
		if trimmed == "" {
			req.Theme = nil
		} else if len(trimmed) > maxCampaignThemeLen {
			writeError(w, http.StatusBadRequest, "theme 太长 / theme too long")
			return
		} else {
			req.Theme = &trimmed
		}
	}

	if err := s.store.SetCampaignMetadata(r.Context(), id, title, req.Description, req.Theme); err != nil {
		writeError(w, http.StatusInternalServerError, "写入失败 / write failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// setMilestoneDescriptionsRequest is the body of POST /api/campaigns/{id}/milestones/metadata.
type setMilestoneDescriptionsRequest struct {
	Descriptions []string `json:"descriptions"` // index-aligned with milestone idx; blank entries are fine
}

const maxMilestoneDescriptionLen = 2000

// handleSetMilestoneDescriptions: POST /api/campaigns/{id}/milestones/metadata
// — stores each milestone's optional plan description (see
// store/milestone_metadata.go). Same trust level and no-ownership-check as
// handleSetCampaignMetadata above — same reasoning applies.
func (s *Server) handleSetMilestoneDescriptions(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req setMilestoneDescriptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体不是合法 JSON / malformed JSON body")
		return
	}
	for i, d := range req.Descriptions {
		if len(d) > maxMilestoneDescriptionLen {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("里程碑 %d 说明太长 / milestone %d description too long", i, i))
			return
		}
	}
	if err := s.store.SetMilestoneDescriptions(r.Context(), id, req.Descriptions); err != nil {
		writeError(w, http.StatusInternalServerError, "写入失败 / write failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// createCampaignRequest is the body of POST /api/campaigns.
type createCampaignRequest struct {
	ID               *uint64  `json:"id"`
	Charity          string   `json:"charity"`
	MilestoneAmounts []string `json:"milestoneAmounts"` // wei, decimal strings, one per milestone
	CreatedBlock     uint64   `json:"createdBlock"`
	CreatedTx        string   `json:"createdTx"`
}

var txHashPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)

// zeroMetadataHash matches what the frontend passes on-chain for
// metadataHash (it doesn't build real off-chain metadata bundles — title/
// description/theme all live in campaign_metadata instead).
var zeroMetadataHash = "0x" + strings.Repeat("0", 64)

// handleCreateCampaign: POST /api/campaigns — optimistically records a
// campaign the moment its create-campaign transaction confirms, ahead of the
// indexer. See the package doc comment above for the reasoning, and
// store.InsertCampaign for how this reconciles with the indexer's own
// (authoritative) write once it catches up.
//
// Not a way to fabricate a campaign that never happened on-chain in any way
// that matters: nothing here lets money move (donate/release are separate,
// on-chain-gated actions), and if no matching CampaignCreated event ever
// arrives, the row just stays confirmed=false forever — see the frontend's
// handling of that flag.
func (s *Server) handleCreateCampaign(w http.ResponseWriter, r *http.Request) {
	var req createCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体不是合法 JSON / malformed JSON body")
		return
	}
	if req.ID == nil {
		writeError(w, http.StatusBadRequest, "id 不能为空 / id is required")
		return
	}
	if !common.IsHexAddress(req.Charity) {
		writeError(w, http.StatusBadRequest, "charity 地址格式错误 / invalid charity address")
		return
	}
	if len(req.MilestoneAmounts) == 0 {
		writeError(w, http.StatusBadRequest, "至少要有一个里程碑 / at least one milestone is required")
		return
	}
	if !txHashPattern.MatchString(req.CreatedTx) {
		writeError(w, http.StatusBadRequest, "createdTx 格式错误 / invalid createdTx")
		return
	}
	if req.CreatedBlock == 0 {
		writeError(w, http.StatusBadRequest, "createdBlock 不能为空 / createdBlock is required")
		return
	}

	// goal = 各里程碑金额之和，跟合约 createCampaign 的算法一致，不接受前端
	// 单独传 goal（没必要多一个可能对不上的输入）。
	// goal = sum of the milestone amounts, mirroring the contract's own
	// createCampaign — not accepted as a separate field from the client,
	// there's no reason to allow one that could disagree with the sum.
	goal := new(big.Int)
	ms := make([]store.Milestone, len(req.MilestoneAmounts))
	for i, raw := range req.MilestoneAmounts {
		amount, ok := new(big.Int).SetString(raw, 10)
		if !ok || amount.Sign() <= 0 {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("里程碑 %d 金额非法 / milestone %d has an invalid amount", i, i))
			return
		}
		goal.Add(goal, amount)
		ms[i] = store.Milestone{CampaignID: *req.ID, Idx: uint32(i), Amount: amount}
	}

	c := store.Campaign{
		ID:             *req.ID,
		Charity:        strings.ToLower(common.HexToAddress(req.Charity).Hex()),
		Goal:           goal,
		MilestoneCount: uint32(len(ms)),
		MetadataHash:   zeroMetadataHash,
		CreatedBlock:   req.CreatedBlock,
		CreatedTx:      strings.ToLower(req.CreatedTx),
	}
	if err := s.store.InsertCampaign(r.Context(), c, ms, false); err != nil {
		writeError(w, http.StatusInternalServerError, "写入失败 / write failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ============================ 小工具 / helpers ============================

// parseID 从路径取 {id} 并转成 uint64，非法则写 400 并返回 ok=false。
// parseID reads {id}, parses to uint64, writing 400 and ok=false on bad input.
func parseID(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id 非法 / invalid id")
		return 0, false
	}
	return id, true
}

// parsePaging 解析 ?limit=&offset=，给出安全默认值与上限，避免一次拉太多。
// parsePaging reads ?limit=&offset= with safe defaults and a cap.
func parsePaging(r *http.Request) (limit, offset int) {
	limit, offset = 20, 0 // 默认 / defaults
	if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if limit > 100 {
		limit = 100 // 上限，防止一次查太多 / cap to protect the DB
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v > 0 {
		offset = v
	}
	return limit, offset
}

// writeJSON 统一以 JSON 返回，设置好 Content-Type 和状态码。
// writeJSON writes a JSON response with the right Content-Type and status.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// 头已发出，无法再改状态码，只能记录（此处简单忽略）。
		// Headers already sent; can't change status now.
		_ = err
	}
}

// writeError 返回统一格式的错误 JSON：{"error": "..."}。
// writeError writes a uniform error JSON body.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// withCORS 允许任意来源的跨域请求（无凭证；几个 POST 接口各自在 handler
// 里做校验/限流，没有共享的敏感操作）。
// withCORS allows cross-origin requests from any origin (no credentials; each
// POST endpoint validates/limits itself in its handler — nothing shared to abuse).
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// 浏览器的预检请求（OPTIONS）直接回 204。/ Answer CORS preflight with 204.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

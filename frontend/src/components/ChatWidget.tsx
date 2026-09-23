import { Fragment, useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { sendChatMessage } from '../lib/api'
import type { ChatMessage } from '../lib/api'

const GREETING =
  "Hi! I'm the TrustChain assistant — ask me how donating, milestones, or verified charities work."

// A few one-tap starters so a first-time visitor isn't staring at an empty
// input box wondering what's worth asking.
const SUGGESTIONS = ['How do I donate?', 'What is a milestone?', "What makes TrustChain different?"]

// Floating help widget, mounted once in AppShell so it's on every page. Talks
// to POST /api/chat (backend/internal/api/chat.go), which proxies to
// DeepSeek — this component never sees an API key, it just calls our own
// backend.
export default function ChatWidget() {
  const [open, setOpen] = useState(false)
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState('')
  const [pending, setPending] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const scrollRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: 'smooth' })
  }, [messages, pending])

  useEffect(() => {
    if (open) {
      // Autofocus after the panel's open transition starts, not before.
      const t = setTimeout(() => inputRef.current?.focus(), 150)
      return () => clearTimeout(t)
    }
  }, [open])

  // A link that lands on the page you're already on changes nothing visible,
  // so close the panel and scroll up to make the click read as "it worked".
  function closeAndTop() {
    setOpen(false)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  async function submit(text: string) {
    const trimmed = text.trim()
    if (!trimmed || pending) return
    const history = messages
    setMessages([...history, { role: 'user', content: trimmed }])
    setInput('')
    setError(null)
    setPending(true)
    try {
      const reply = await sendChatMessage(trimmed, history)
      setMessages((prev) => [...prev, { role: 'assistant', content: reply }])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'The assistant is unavailable right now.')
    } finally {
      setPending(false)
    }
  }

  return (
    // The wrapper is always as big as the (possibly hidden) panel, so it must
    // not catch clicks itself — only the open panel and the button do.
    <div className="pointer-events-none fixed bottom-6 right-6 z-20 flex flex-col items-end gap-3">
      <div
        className={`flex h-[28rem] w-80 origin-bottom-right flex-col overflow-hidden rounded-2xl border border-border bg-white shadow-2xl transition-all duration-200 sm:w-96 ${
          open ? 'pointer-events-auto scale-100 opacity-100' : 'pointer-events-none scale-90 opacity-0'
        }`}
      >
        <div
          className="flex items-center gap-3 px-4 py-3.5 text-white"
          style={{ background: 'linear-gradient(90deg, #1e3a5f, #7c3aed)' }}
        >
          <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-white/15 text-lg">
            🤝
          </span>
          <div className="min-w-0 flex-1">
            <p className="text-sm font-semibold">TrustChain Assistant</p>
            <p className="flex items-center gap-1.5 text-xs text-white/75">
              <span className="h-1.5 w-1.5 rounded-full bg-emerald-400" aria-hidden="true" />
              Usually replies in a few seconds
            </p>
          </div>
          <button
            type="button"
            onClick={() => setOpen(false)}
            className="shrink-0 rounded-full p-1 text-white/80 transition-colors hover:bg-white/10 hover:text-white"
            aria-label="Close chat"
          >
            ✕
          </button>
        </div>

        <div ref={scrollRef} className="flex-1 space-y-3 overflow-y-auto px-4 py-4">
          <ChatBubble role="assistant" content={GREETING} onNavigate={closeAndTop} />

          {messages.length === 0 && (
            <div className="flex flex-col items-start gap-2 pl-1">
              {SUGGESTIONS.map((s) => (
                <button
                  key={s}
                  type="button"
                  onClick={() => submit(s)}
                  className="rounded-full border border-border bg-white px-3 py-1.5 text-xs font-medium text-accent transition-colors hover:bg-accent-tint"
                >
                  {s}
                </button>
              ))}
            </div>
          )}

          {messages.map((m, i) => (
            <ChatBubble key={i} role={m.role} content={m.content} onNavigate={closeAndTop} />
          ))}
          {pending && <TypingBubble />}
          {error && (
            <p className="rounded-lg bg-red-50 px-3 py-2 text-xs text-red-600">{error}</p>
          )}
        </div>

        <form
          onSubmit={(e) => {
            e.preventDefault()
            submit(input)
          }}
          className="flex items-center gap-2 border-t border-border p-3"
        >
          <input
            ref={inputRef}
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ask a question…"
            disabled={pending}
            maxLength={2000}
            className="flex-1 rounded-full border border-border px-4 py-2 text-sm text-ink outline-none transition-colors focus:border-accent disabled:opacity-50"
          />
          <button
            type="submit"
            disabled={pending || !input.trim()}
            aria-label="Send message"
            className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-white shadow-sm transition-all hover:opacity-90 active:scale-95 disabled:opacity-40"
            style={{ background: 'linear-gradient(135deg, #1e3a5f, #7c3aed)' }}
          >
            <svg viewBox="0 0 20 20" fill="none" className="h-4 w-4 -translate-x-px">
              <path
                d="M3 10l14-6.5-4 6.5 4 6.5L3 10z"
                fill="currentColor"
              />
            </svg>
          </button>
        </form>
      </div>

      <button
        type="button"
        onClick={() => setOpen((prev) => !prev)}
        aria-label={open ? 'Close chat assistant' : 'Open chat assistant'}
        className="pointer-events-auto flex h-14 w-14 items-center justify-center rounded-full text-2xl text-white shadow-lg transition-transform hover:scale-105 active:scale-95"
        style={{ background: 'linear-gradient(135deg, #1e3a5f, #7c3aed)' }}
      >
        {open ? '✕' : '💬'}
      </button>
    </div>
  )
}

function ChatBubble({
  role,
  content,
  onNavigate,
}: {
  role: 'user' | 'assistant'
  content: string
  onNavigate: () => void
}) {
  const isUser = role === 'user'
  return (
    <div className={`flex animate-[fade-in_0.15s_ease-out] ${isUser ? 'justify-end' : 'justify-start'}`}>
      <p
        className={`max-w-[85%] whitespace-pre-line px-3.5 py-2.5 text-sm leading-relaxed ${
          isUser
            ? 'rounded-2xl rounded-br-md bg-accent text-white'
            : 'rounded-2xl rounded-bl-md bg-accent-tint text-ink'
        }`}
      >
        {isUser ? content : renderWithLinks(content, onNavigate)}
      </p>
    </div>
  )
}

// The assistant is told (backend/internal/api/chat.go) to link site pages as
// [label](/path). Only same-site paths become links — anything else stays as
// plain text, so a model slip can't turn into an outbound link.
const LINK_RE = /\[([^\]]+)\]\((\/(?!\/)[^)\s]*)\)/g

function renderWithLinks(text: string, onNavigate: () => void) {
  const out: React.ReactNode[] = []
  let last = 0
  for (const m of text.matchAll(LINK_RE)) {
    const idx = m.index ?? 0
    if (idx > last) out.push(<Fragment key={`t${idx}`}>{text.slice(last, idx)}</Fragment>)
    out.push(
      <Link key={`l${idx}`} to={m[2]} onClick={onNavigate} className="font-medium text-accent underline underline-offset-2">
        {m[1]}
      </Link>,
    )
    last = idx + m[0].length
  }
  if (last < text.length) out.push(<Fragment key="tend">{text.slice(last)}</Fragment>)
  return out
}

function TypingBubble() {
  return (
    <div className="flex justify-start">
      <div className="flex items-center gap-1 rounded-2xl rounded-bl-md bg-accent-tint px-4 py-3">
        {[0, 1, 2].map((i) => (
          <span
            key={i}
            className="h-1.5 w-1.5 animate-bounce rounded-full bg-accent/60"
            style={{ animationDelay: `${i * 0.15}s` }}
          />
        ))}
      </div>
    </div>
  )
}

// Typed client for the backend REST API (backend/internal/api) — mostly
// read-only; connectUser(), setCampaignMetadata() and createCampaignRecord()
// below are the writes (see that package's doc comment for why).
// All amounts come back as decimal-wei strings (uint256 can overflow JS numbers).

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8090'

export interface Campaign {
  id: number
  charity: string
  goal: string // wei
  raised: string // wei
  released: string // wei
  milestoneCount: number
  metadataHash: string
  completed: boolean
  createdBlock: number
  createdTx: string
  createdAt: string
  // false right after createCampaignRecord() optimistically inserts this row
  // — the indexer hasn't confirmed it against a real on-chain event yet
  // (usually takes ~1-2 minutes). See backend/internal/store's InsertCampaign
  // comment on "confirmed".
  confirmed: boolean
  // Off-chain — the contract's metadataHash isn't a resolvable pointer, so
  // these come from campaign_metadata via the API. null until set.
  title: string | null
  description: string | null
  theme: string | null // a CampaignTheme.key from lib/theme.ts, or null
}

// Milestone state as stored on-chain: 0 = Locked, 1 = Released, 2 = Proven.
export type MilestoneState = 0 | 1 | 2

export interface Milestone {
  idx: number
  amount: string // wei
  state: MilestoneState
  receiptHash: string | null // on-chain fingerprint only — see the fields below for the actual file
  // The uploaded expense document's metadata (milestone_receipts table),
  // null until the charity uploads one via uploadMilestoneReceipt(). The
  // file itself is at milestoneReceiptFileUrl(campaignId, idx).
  receiptFileName: string | null
  receiptContentType: string | null
  receiptFileSize: number | null
  receiptSha256: string | null // hex, no 0x — compare against receiptHash (0x + this) to verify
  receiptNote: string | null
  receiptUploadedAt: string | null
  // Optional "what this milestone will accomplish", entered at creation —
  // distinct from receiptNote ("what was actually spent"), set at receipt
  // time. See lib/api's setMilestoneDescriptions().
  description: string | null
}

export interface Donation {
  donor: string
  amount: string // wei
  blockNumber: number
  txHash: string
  logIndex: number
  createdAt: string
}

export interface Charity {
  address: string
  verified: boolean
  updatedBlock: number
}

export interface CampaignDetail {
  campaign: Campaign
  milestones: Milestone[]
}

export type ActivityEventType =
  | 'CampaignCreated'
  | 'DonationReceived'
  | 'MilestoneReleased'
  | 'ReceiptSubmitted'
  | 'CampaignCompleted'

export interface ActivityEvent {
  id: number
  campaignId: number
  eventType: ActivityEventType
  amount: string | null // wei
  milestoneIdx: number | null
  blockNumber: number
  txHash: string
  logIndex: number
  createdAt: string
}

export class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

async function getJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`)
  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw new ApiError(res.status, body?.error ?? `${res.status} ${res.statusText}`)
  }
  return res.json() as Promise<T>
}

async function postJSON(path: string, body: unknown): Promise<void> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    const errBody = await res.json().catch(() => null)
    throw new ApiError(res.status, errBody?.error ?? `${res.status} ${res.statusText}`)
  }
}

export function listCampaigns(): Promise<Campaign[]> {
  return getJSON('/api/campaigns')
}

export function getCampaign(id: number | string): Promise<CampaignDetail> {
  return getJSON(`/api/campaigns/${id}`)
}

export function listDonations(id: number | string): Promise<Donation[]> {
  return getJSON(`/api/campaigns/${id}/donations`)
}

export function getCharity(address: string): Promise<Charity> {
  return getJSON(`/api/charities/${address}`)
}

export function listCharities(): Promise<Charity[]> {
  return getJSON('/api/charities')
}

export function listActivity(): Promise<ActivityEvent[]> {
  return getJSON('/api/activity')
}

// One campaign's on-chain activity, oldest first — a real transaction
// timeline for the campaign detail page (see backend's ListActivityByCampaign).
export function listCampaignActivity(id: number | string): Promise<ActivityEvent[]> {
  return getJSON(`/api/campaigns/${id}/activity`)
}

// Records that a wallet connected, with the role the frontend computed for
// it (lib/role.ts). Not a source of truth for permissions — just a "who
// connected, as what, when" log. Callers should fire-and-forget this; a
// failure here shouldn't block the wallet UI.
export function connectUser(address: string, role: 'admin' | 'charity' | 'donor'): Promise<void> {
  return postJSON('/api/users/connect', { address, role })
}

// Stores a campaign's title/description/theme (see the Campaign fields'
// comments for why these aren't on-chain). Call right after the
// create-campaign tx confirms — no need to wait for the indexer to pick up
// the campaign row first.
export function setCampaignMetadata(
  id: number,
  metadata: { title: string; description?: string; theme?: string },
): Promise<void> {
  return postJSON(`/api/campaigns/${id}/metadata`, metadata)
}

// Stores each milestone's optional plan description, index-aligned with
// milestone idx (descriptions[i] describes milestone i; blank entries are
// fine and just don't get stored). Call alongside setCampaignMetadata, right
// after the create-campaign tx confirms.
export function setMilestoneDescriptions(id: number, descriptions: string[]): Promise<void> {
  return postJSON(`/api/campaigns/${id}/milestones/metadata`, { descriptions })
}

// Optimistically records a campaign the moment its create-campaign tx
// confirms, so it shows up in listCampaigns()/getCampaign() immediately
// instead of after the indexer catches up (~1-2 minutes: CONFIRMATIONS
// blocks + the next poll). Comes back as Campaign.confirmed === false until
// the indexer reconciles it against the real on-chain event.
export function createCampaignRecord(params: {
  id: number
  charity: string
  milestoneAmountsWei: string[]
  createdBlock: number
  createdTx: string
}): Promise<void> {
  return postJSON('/api/campaigns', {
    id: params.id,
    charity: params.charity,
    milestoneAmounts: params.milestoneAmountsWei,
    createdBlock: params.createdBlock,
    createdTx: params.createdTx,
  })
}

export interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

// Sends one turn to the backend's DeepSeek-API proxy (backend/internal/api/chat.go
// — the API key lives server-side, never in this file). `history` is the prior
// turns, oldest first; the caller owns the running conversation and just appends
// the new reply once it comes back.
export async function sendChatMessage(message: string, history: ChatMessage[]): Promise<string> {
  const res = await fetch(`${BASE_URL}/api/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message, history }),
  })
  if (!res.ok) {
    const errBody = await res.json().catch(() => null)
    throw new ApiError(res.status, errBody?.error ?? `${res.status} ${res.statusText}`)
  }
  const body = (await res.json()) as { reply: string }
  return body.reply
}

// URL for a milestone's uploaded receipt file (backend/internal/api/receipt_upload.go).
// Safe to use directly as an <img src> or <a href> — 404s if nothing was uploaded yet.
export function milestoneReceiptFileUrl(campaignId: number | string, idx: number): string {
  return `${BASE_URL}/api/campaigns/${campaignId}/milestones/${idx}/receipt/file`
}

// Uploads the actual expense document for a milestone, after its
// submitReceipt transaction has already confirmed on-chain (see lib/hash.ts
// and lib/escrow.ts's submitReceipt — the hash committed there and the hash
// this computes server-side from the same bytes should match). Not JSON, so
// this doesn't go through postJSON.
export async function uploadMilestoneReceipt(
  campaignId: number,
  idx: number,
  file: File,
  note?: string,
): Promise<{ sha256: string; fileName: string; contentType: string; fileSize: number }> {
  const form = new FormData()
  form.append('file', file)
  if (note) form.append('note', note)
  const res = await fetch(`${BASE_URL}/api/campaigns/${campaignId}/milestones/${idx}/receipt`, {
    method: 'POST',
    body: form,
  })
  if (!res.ok) {
    const errBody = await res.json().catch(() => null)
    throw new ApiError(res.status, errBody?.error ?? `${res.status} ${res.statusText}`)
  }
  return res.json()
}

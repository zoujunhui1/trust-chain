// Typed client for the backend REST API (backend/internal/api) — mostly
// read-only; connectUser() below is the one write (see that package's doc
// comment for why).
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
}

// Milestone state as stored on-chain: 0 = Locked, 1 = Released, 2 = Proven.
export type MilestoneState = 0 | 1 | 2

export interface Milestone {
  idx: number
  amount: string // wei
  state: MilestoneState
  receiptHash: string | null
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

// Records that a wallet connected, with the role the frontend computed for
// it (lib/role.ts). Not a source of truth for permissions — just a "who
// connected, as what, when" log. Callers should fire-and-forget this; a
// failure here shouldn't block the wallet UI.
export async function connectUser(address: string, role: 'admin' | 'charity' | 'donor'): Promise<void> {
  const res = await fetch(`${BASE_URL}/api/users/connect`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ address, role }),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw new ApiError(res.status, body?.error ?? `${res.status} ${res.statusText}`)
  }
}

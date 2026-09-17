import { formatEther } from 'ethers'
import type { ActivityEvent, ActivityEventType, Campaign, Milestone } from './api'

// wei decimal string -> "1.500" (3dp, matches the Figma mockup's amount style).
export function weiToEth(wei: string, decimals = 3): string {
  return Number(formatEther(wei)).toFixed(decimals)
}

// Campaign.title is off-chain and only set once someone calls
// setCampaignTitle() — falls back to the id for anything created before
// this existed, or if the title save failed.
export function campaignTitle(campaign: Pick<Campaign, 'id' | 'title'>): string {
  return campaign.title ?? `Campaign #${campaign.id}`
}

// "0x092dD42fFc79E13217985C1bB4A9155Fa085cD9A" -> "0x092d…cD9A"
export function shortAddress(address: string): string {
  if (address.length <= 12) return address
  return `${address.slice(0, 6)}…${address.slice(-4)}`
}

// wei goal/raised -> 0-100, clamped (goal can be 0 before it's indexed).
export function progressPercent(raised: string, goal: string): number {
  const goalEth = Number(formatEther(goal))
  if (goalEth <= 0) return 0
  const raisedEth = Number(formatEther(raised))
  return Math.min(100, Math.round((raisedEth / goalEth) * 100))
}

// ISO timestamp -> "2h ago" / "3d ago" (falls back to a date once it's old).
export function relativeTime(iso: string): string {
  const diffMs = Date.now() - new Date(iso).getTime()
  const minutes = Math.round(diffMs / 60_000)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.round(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.round(hours / 24)
  if (days < 30) return `${days}d ago`
  return new Date(iso).toLocaleDateString()
}

// Human-readable status line for a milestone row, mirroring the on-chain
// unlock rules in TrustChainEscrow (releaseMilestone/submitReceipt):
// milestone i only unlocks once milestone i-1 is Proven, and milestone 0
// only once the campaign goal is fully raised.
export function milestoneStatusText(milestone: Milestone, idx: number, campaign: Campaign): string {
  if (milestone.state === 2) {
    return milestone.receiptHash ? `Receipt ${shortAddress(milestone.receiptHash)}` : 'Proven'
  }
  if (milestone.state === 1) return 'Awaiting completion receipt'
  if (idx === 0) {
    return BigInt(campaign.raised) >= BigInt(campaign.goal)
      ? 'Awaiting milestone release'
      : 'Unlocks once the funding goal is reached'
  }
  return `Unlocks after Milestone ${idx} is proven`
}

// txHash -> a real, live Sepolia Etherscan URL — the whole point of
// TrustChain's transparency pitch is that this isn't a screenshot, it's a
// link anyone can open and check for themselves.
export function etherscanTxUrl(txHash: string): string {
  return `https://sepolia.etherscan.io/tx/${txHash}`
}

// Shared with TransparencyDashboard's global feed and CampaignDetail's
// per-campaign timeline, so the two read the same event the same way.
const ACTIVITY_LABEL: Record<ActivityEventType, (a: ActivityEvent) => string> = {
  CampaignCreated: () => 'Campaign created',
  DonationReceived: () => 'Donation received',
  MilestoneReleased: (a) => `Milestone ${(a.milestoneIdx ?? 0) + 1} released`,
  ReceiptSubmitted: (a) => `Receipt submitted (Milestone ${(a.milestoneIdx ?? 0) + 1})`,
  CampaignCompleted: () => 'Campaign completed',
}

export function activityEventLabel(a: ActivityEvent): string {
  return ACTIVITY_LABEL[a.eventType](a)
}

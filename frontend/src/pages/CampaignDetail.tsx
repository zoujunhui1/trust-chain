import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ApiError, getCampaign, getCharity, listCampaignActivity, listDonations, milestoneReceiptFileUrl } from '../lib/api'
import type { ActivityEvent, Campaign, Donation, Milestone, MilestoneState } from '../lib/api'
import {
  activityEventLabel,
  campaignTitle,
  etherscanTxUrl,
  milestoneStatusText,
  progressPercent,
  relativeTime,
  shortAddress,
  weiToEth,
} from '../lib/format'
import { campaignTheme } from '../lib/theme'
import { useWallet } from '../lib/wallet'
import MilestoneChip from '../components/MilestoneChip'
import MilestoneReceiptPanel from '../components/MilestoneReceiptPanel'
import DonationPanel from '../components/DonationPanel'
import FundsFlow from '../components/FundsFlow'

// Same palette as MilestoneChip, just as a filled dot for the timeline below.
const MILESTONE_DOT_CLASSES: Record<MilestoneState, string> = {
  0: 'bg-locked-tint text-locked',
  1: 'bg-released-tint text-released',
  2: 'bg-proven-tint text-proven',
}

export default function CampaignDetail() {
  const { id } = useParams()
  const wallet = useWallet()

  const [campaign, setCampaign] = useState<Campaign | null>(null)
  const [milestones, setMilestones] = useState<Milestone[]>([])
  const [donations, setDonations] = useState<Donation[]>([])
  const [activity, setActivity] = useState<ActivityEvent[]>([])
  const [verified, setVerified] = useState(false)
  const [notFound, setNotFound] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const theme = campaign ? campaignTheme(campaign) : null

  useEffect(() => {
    if (!id) return
    let cancelled = false
    setCampaign(null)
    setNotFound(false)
    setError(null)

    getCampaign(id)
      .then(async ({ campaign, milestones }) => {
        if (cancelled) return
        setCampaign(campaign)
        setMilestones(milestones)

        const [charityRes, donationsRes, activityRes] = await Promise.allSettled([
          getCharity(campaign.charity),
          listDonations(id),
          listCampaignActivity(id),
        ])
        if (cancelled) return
        if (charityRes.status === 'fulfilled') setVerified(charityRes.value.verified)
        if (donationsRes.status === 'fulfilled') setDonations(donationsRes.value)
        if (activityRes.status === 'fulfilled') setActivity(activityRes.value)
      })
      .catch((err: Error) => {
        if (cancelled) return
        if (err instanceof ApiError && err.status === 404) {
          setNotFound(true)
        } else {
          setError(err.message)
        }
      })

    return () => {
      cancelled = true
    }
  }, [id])

  return (
    <div className="mx-auto max-w-7xl px-6 sm:px-10 lg:px-16">
      <div className="pt-10">
        <Link to="/" className="text-sm text-muted transition-colors hover:text-ink">
          ← Back to Campaigns
        </Link>
      </div>

      {error && <p className="py-12 text-sm text-red-600">Couldn't load campaign: {error}</p>}
      {notFound && <p className="py-12 text-muted">Campaign not found.</p>}
      {!error && !notFound && !campaign && <p className="py-12 text-muted">Loading campaign…</p>}

      {campaign && theme && (
        <>
          <div className="pt-6">
            <div
              className="flex h-36 items-center justify-between rounded-xl px-8 shadow-sm"
              style={{ background: theme.gradient }}
            >
              <span className="text-5xl" aria-hidden="true">
                {theme.emoji}
              </span>
              <span className="rounded-full bg-white/80 px-3 py-1.5 text-sm font-medium" style={{ color: theme.accent }}>
                {theme.label}
              </span>
            </div>
          </div>

          <div className="pb-10 pt-6">
            <div className="flex items-center gap-1.5 text-sm text-muted">
              {verified && (
                <svg viewBox="0 0 14 14" fill="none" className="h-3.5 w-3.5 shrink-0 text-accent">
                  <circle cx="7" cy="7" r="7" fill="currentColor" />
                  <path d="M4 7.2l2 2 4-4.4" stroke="white" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              )}
              <span>
                {shortAddress(campaign.charity)}
                {verified ? ' · Verified' : ''}
              </span>
            </div>

            <div className="mt-3 flex items-center gap-3">
              <h1 className="text-3xl font-bold tracking-tight text-ink">{campaignTitle(campaign)}</h1>
              {campaign.completed && (
                <span className="rounded-full bg-proven-tint px-2.5 py-1 text-xs font-medium text-proven">
                  Completed
                </span>
              )}
              {!campaign.confirmed && (
                <span className="rounded-full bg-released-tint px-2.5 py-1 text-xs font-medium text-released">
                  Confirming on-chain…
                </span>
              )}
            </div>

            <div className="mt-6 h-2.5 w-full max-w-[700px] overflow-hidden rounded-full bg-border">
              <div
                className="h-full rounded-full"
                style={{ width: `${progressPercent(campaign.raised, campaign.goal)}%`, backgroundColor: theme.accent }}
              />
            </div>
            <p className="mt-3 text-sm text-muted">
              {weiToEth(campaign.raised)} ETH raised of {weiToEth(campaign.goal)} ETH goal ·{' '}
              {progressPercent(campaign.raised, campaign.goal)}% funded
            </p>
          </div>

          <FundsFlow campaign={campaign} milestones={milestones} />

          <div className="flex flex-col gap-10 pb-16 lg:flex-row">
            <div className="flex-1">
              <h2 className="text-lg font-semibold text-ink">About this campaign</h2>
              <p className="mt-2 max-w-[760px] whitespace-pre-line text-sm text-muted">
                {campaign.description ??
                  "This campaign hasn't published a description yet — only the on-chain fields below are currently tracked."}
              </p>

              <div className="mt-6 rounded-xl border border-border bg-white p-5 shadow-sm">
                <ol className="relative">
                  <span className="absolute left-4 top-4 bottom-4 w-px bg-border" aria-hidden="true" />
                  {milestones.map((m, i) => (
                    <li key={m.idx} className={`relative flex gap-4 ${i > 0 ? 'mt-6' : ''}`}>
                      <span
                        className={`relative z-10 flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-semibold ${MILESTONE_DOT_CLASSES[m.state]}`}
                      >
                        {m.state === 2 ? '✓' : m.idx + 1}
                      </span>
                      <div className="flex-1 pt-1">
                        <div className="flex items-center justify-between gap-4">
                          <p className="font-medium text-ink">Milestone {m.idx + 1}</p>
                          <MilestoneChip state={m.state} />
                        </div>
                        <p className="mt-1 text-sm text-muted">
                          {weiToEth(m.amount)} ETH · {milestoneStatusText(m, m.idx, campaign)}
                        </p>
                        {m.description && <p className="mt-1 text-sm text-ink">{m.description}</p>}

                        {m.receiptFileName && (
                          <ReceiptFileCard campaignId={campaign.id} milestone={m} accent={theme.accent} />
                        )}

                        <MilestoneReceiptPanel
                          campaignId={campaign.id}
                          milestone={m}
                          prevState={i > 0 ? milestones[i - 1].state : undefined}
                          goalReached={BigInt(campaign.raised) >= BigInt(campaign.goal)}
                          isCharity={wallet.address?.toLowerCase() === campaign.charity.toLowerCase()}
                          accent={theme.accent}
                        />
                      </div>
                    </li>
                  ))}
                </ol>
              </div>
            </div>

            <DonationPanel campaignId={campaign.id} donations={donations} accent={theme.accent} />
          </div>

          <div className="pb-16">
            <h2 className="text-lg font-semibold text-ink">On-Chain Activity</h2>
            <p className="mt-1 max-w-[760px] text-sm text-muted">
              This isn't a summary we wrote — it's this campaign's real transaction history. Every row
              links to the actual transaction on Sepolia Etherscan.
            </p>

            <div className="mt-5 overflow-hidden rounded-xl border border-border bg-white shadow-sm">
              {activity.length === 0 ? (
                <p className="px-5 py-6 text-sm text-muted">
                  No on-chain activity indexed yet — it can take a minute or two to appear after a
                  transaction confirms.
                </p>
              ) : (
                <ul className="divide-y divide-border">
                  {activity.map((a) => (
                    <li key={a.id} className="flex flex-wrap items-center justify-between gap-2 px-5 py-3.5 text-sm">
                      <div className="flex items-center gap-3">
                        <span
                          className="h-2 w-2 shrink-0 rounded-full"
                          style={{ backgroundColor: theme.accent }}
                          aria-hidden="true"
                        />
                        <span className="text-ink">{activityEventLabel(a)}</span>
                        {a.amount && <span className="text-muted">· {weiToEth(a.amount)} ETH</span>}
                      </div>
                      <div className="flex items-center gap-4 text-xs text-muted">
                        <span>{relativeTime(a.createdAt)}</span>
                        <a
                          href={etherscanTxUrl(a.txHash)}
                          target="_blank"
                          rel="noreferrer"
                          className="font-medium underline"
                          style={{ color: theme.accent }}
                        >
                          View on Etherscan ↗
                        </a>
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  )
}

// The persisted, already-uploaded receipt (from the initial fetch) — shown
// to every visitor, not just the charity. MilestoneReceiptPanel handles the
// charity's own upload flow and its own "just submitted" confirmation.
function ReceiptFileCard({
  campaignId,
  milestone,
  accent,
}: {
  campaignId: number
  milestone: Milestone
  accent: string
}) {
  const verified =
    milestone.receiptHash != null &&
    milestone.receiptSha256 != null &&
    milestone.receiptHash.toLowerCase() === `0x${milestone.receiptSha256}`.toLowerCase()

  return (
    <div className="mt-2 rounded-lg border border-border bg-page/60 p-3 text-sm">
      <a
        href={milestoneReceiptFileUrl(campaignId, milestone.idx)}
        target="_blank"
        rel="noreferrer"
        className="font-medium underline"
        style={{ color: accent }}
      >
        📎 {milestone.receiptFileName}
      </a>
      {milestone.receiptNote && <p className="mt-1 text-muted">{milestone.receiptNote}</p>}
      <p className="mt-1 text-xs">
        {verified ? (
          <span className="text-proven">✓ Verified — file hash matches the on-chain receipt</span>
        ) : (
          <span className="text-muted">Waiting on indexer confirmation to verify against the on-chain hash…</span>
        )}
      </p>
    </div>
  )
}

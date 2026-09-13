import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ApiError, getCampaign, getCharity, listDonations } from '../lib/api'
import type { Campaign, Donation, Milestone, MilestoneState } from '../lib/api'
import { campaignTitle, milestoneStatusText, progressPercent, shortAddress, weiToEth } from '../lib/format'
import { campaignTheme } from '../lib/theme'
import MilestoneChip from '../components/MilestoneChip'
import DonationPanel from '../components/DonationPanel'

// Same palette as MilestoneChip, just as a filled dot for the timeline below.
const MILESTONE_DOT_CLASSES: Record<MilestoneState, string> = {
  0: 'bg-locked-tint text-locked',
  1: 'bg-released-tint text-released',
  2: 'bg-proven-tint text-proven',
}

export default function CampaignDetail() {
  const { id } = useParams()

  const [campaign, setCampaign] = useState<Campaign | null>(null)
  const [milestones, setMilestones] = useState<Milestone[]>([])
  const [donations, setDonations] = useState<Donation[]>([])
  const [verified, setVerified] = useState(false)
  const [notFound, setNotFound] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const theme = campaign ? campaignTheme(campaign.id) : null

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

        const [charityRes, donationsRes] = await Promise.allSettled([
          getCharity(campaign.charity),
          listDonations(id),
        ])
        if (cancelled) return
        if (charityRes.status === 'fulfilled') setVerified(charityRes.value.verified)
        if (donationsRes.status === 'fulfilled') setDonations(donationsRes.value)
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

          <div className="flex flex-col gap-10 pb-16 lg:flex-row">
            <div className="flex-1">
              <h2 className="text-lg font-semibold text-ink">About this campaign</h2>
              <p className="mt-2 max-w-[760px] text-sm text-muted">
                This campaign hasn't published a description yet — only the on-chain fields below
                are currently tracked.
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
                      </div>
                    </li>
                  ))}
                </ol>
              </div>
            </div>

            <DonationPanel campaignId={campaign.id} donations={donations} accent={theme.accent} />
          </div>
        </>
      )}
    </div>
  )
}

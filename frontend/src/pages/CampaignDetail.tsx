import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ApiError, getCampaign, getCharity, listDonations } from '../lib/api'
import type { Campaign, Donation, Milestone } from '../lib/api'
import { milestoneStatusText, progressPercent, shortAddress, weiToEth } from '../lib/format'
import MilestoneChip from '../components/MilestoneChip'
import DonationPanel from '../components/DonationPanel'

export default function CampaignDetail() {
  const { id } = useParams()

  const [campaign, setCampaign] = useState<Campaign | null>(null)
  const [milestones, setMilestones] = useState<Milestone[]>([])
  const [donations, setDonations] = useState<Donation[]>([])
  const [verified, setVerified] = useState(false)
  const [notFound, setNotFound] = useState(false)
  const [error, setError] = useState<string | null>(null)

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
    <div>
      <div className="px-16 pt-10">
        <Link to="/" className="text-sm text-muted hover:text-ink">
          ← Back to Campaigns
        </Link>
      </div>

      {error && <p className="px-16 py-12 text-sm text-red-600">Couldn't load campaign: {error}</p>}
      {notFound && <p className="px-16 py-12 text-muted">Campaign not found.</p>}
      {!error && !notFound && !campaign && <p className="px-16 py-12 text-muted">Loading campaign…</p>}

      {campaign && (
        <>
          <div className="px-16 pb-10 pt-6">
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
              <h1 className="text-3xl font-bold text-ink">Campaign #{campaign.id}</h1>
              {campaign.completed && (
                <span className="rounded-full bg-proven-tint px-2.5 py-1 text-xs font-medium text-proven">
                  Completed
                </span>
              )}
            </div>

            <div className="mt-6 h-2.5 w-full max-w-[700px] overflow-hidden rounded-full bg-border">
              <div
                className="h-full rounded-full bg-accent"
                style={{ width: `${progressPercent(campaign.raised, campaign.goal)}%` }}
              />
            </div>
            <p className="mt-3 text-sm text-muted">
              {weiToEth(campaign.raised)} ETH raised of {weiToEth(campaign.goal)} ETH goal ·{' '}
              {progressPercent(campaign.raised, campaign.goal)}% funded
            </p>
          </div>

          <div className="flex flex-col gap-10 px-16 pb-16 lg:flex-row">
            <div className="flex-1">
              <h2 className="text-lg font-semibold text-ink">About this campaign</h2>
              <p className="mt-2 max-w-[760px] text-sm text-muted">
                This campaign hasn't published off-chain details (title, description) yet — only the
                on-chain fields below are currently tracked.
              </p>

              <div className="mt-6 divide-y divide-border rounded-xl border border-border bg-white">
                {milestones.map((m) => (
                  <div key={m.idx} className="flex items-center justify-between gap-4 px-5 py-4">
                    <div>
                      <p className="font-medium text-ink">Milestone {m.idx + 1}</p>
                      <p className="mt-1 text-sm text-muted">
                        {weiToEth(m.amount)} ETH · {milestoneStatusText(m, m.idx, campaign)}
                      </p>
                    </div>
                    <MilestoneChip state={m.state} />
                  </div>
                ))}
              </div>
            </div>

            <DonationPanel campaignId={campaign.id} donations={donations} />
          </div>
        </>
      )}
    </div>
  )
}

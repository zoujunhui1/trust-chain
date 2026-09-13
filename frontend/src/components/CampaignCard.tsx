import { Link } from 'react-router-dom'
import type { Campaign, Milestone } from '../lib/api'
import { campaignTitle, progressPercent, shortAddress, weiToEth } from '../lib/format'
import { campaignTheme } from '../lib/theme'
import MilestoneChip from './MilestoneChip'

interface CampaignCardProps {
  campaign: Campaign
  verified?: boolean
  milestones?: Milestone[]
}

export default function CampaignCard({ campaign, verified, milestones }: CampaignCardProps) {
  const percent = progressPercent(campaign.raised, campaign.goal)
  const theme = campaignTheme(campaign.id)

  return (
    <div className="flex flex-col overflow-hidden rounded-xl border border-border bg-white shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:shadow-lg">
      <div className="flex items-center justify-between px-6 py-8" style={{ background: theme.gradient }}>
        <span className="text-4xl" aria-hidden="true">
          {theme.emoji}
        </span>
        <span className="rounded-full bg-white/80 px-2.5 py-1 text-xs font-medium" style={{ color: theme.accent }}>
          {theme.label}
        </span>
      </div>

      <div className="flex flex-1 flex-col p-6">
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

        <h3 className="mt-2 text-lg font-semibold text-ink">{campaignTitle(campaign)}</h3>

        <div className="mt-6 h-2 w-full overflow-hidden rounded-full bg-border">
          <div className="h-full rounded-full" style={{ width: `${percent}%`, backgroundColor: theme.accent }} />
        </div>
        <p className="mt-3 text-sm text-muted">
          {weiToEth(campaign.raised)} / {weiToEth(campaign.goal)} ETH raised
        </p>

        {milestones && milestones.length > 0 && (
          <div className="mt-5 flex flex-wrap gap-2">
            {milestones.map((m) => (
              <MilestoneChip key={m.idx} state={m.state} />
            ))}
          </div>
        )}

        <Link
          to={`/campaigns/${campaign.id}`}
          className="mt-6 rounded-lg px-4 py-2.5 text-center text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98]"
          style={{ backgroundColor: theme.accent }}
        >
          View Campaign
        </Link>
      </div>
    </div>
  )
}

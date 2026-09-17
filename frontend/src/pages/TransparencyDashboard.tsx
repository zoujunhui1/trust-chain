import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { listActivity, listCampaigns, listCharities } from '../lib/api'
import type { ActivityEvent, Campaign, Charity } from '../lib/api'
import { activityEventLabel, campaignTitle, etherscanTxUrl, shortAddress, weiToEth } from '../lib/format'

export default function TransparencyDashboard() {
  const [campaigns, setCampaigns] = useState<Campaign[] | null>(null)
  const [charities, setCharities] = useState<Charity[] | null>(null)
  const [activity, setActivity] = useState<ActivityEvent[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    Promise.all([listCampaigns(), listCharities(), listActivity()])
      .then(([c, ch, a]) => {
        if (cancelled) return
        setCampaigns(c)
        setCharities(ch)
        setActivity(a)
      })
      .catch((err: Error) => {
        if (!cancelled) setError(err.message)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const titleById = new Map(campaigns?.map((c) => [c.id, campaignTitle(c)]))
  const totalRaisedWei = campaigns?.reduce((sum, c) => sum + BigInt(c.raised), 0n) ?? 0n
  const activeCount = campaigns?.filter((c) => !c.completed).length ?? 0
  const completedCount = campaigns?.filter((c) => c.completed).length ?? 0
  const verifiedCount = charities?.filter((c) => c.verified).length ?? 0

  return (
    <div className="mx-auto max-w-7xl px-6 py-12 sm:px-10 lg:px-16">
      <h1 className="text-3xl font-bold tracking-tight text-ink">Transparency Dashboard</h1>
      <p className="mt-2 max-w-2xl text-muted">
        Every donation, milestone release and receipt is recorded on-chain and verifiable by anyone.
      </p>

      {error && <p className="mt-8 text-sm text-red-600">Couldn't load dashboard data: {error}</p>}

      {!error && (!campaigns || !charities || !activity) && (
        <p className="mt-8 text-muted">Loading…</p>
      )}

      {campaigns && charities && activity && (
        <>
          <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-3">
            <StatCard
              label="TOTAL RAISED"
              value={`${weiToEth(totalRaisedWei.toString())} ETH`}
              detail={`across ${campaigns.length} campaign${campaigns.length === 1 ? '' : 's'}`}
            />
            <StatCard
              label="ACTIVE CAMPAIGNS"
              value={String(activeCount)}
              detail={`${completedCount} completed`}
            />
            <StatCard
              label="VERIFIED CHARITIES"
              value={String(verifiedCount)}
              detail="registry-approved"
            />
          </div>

          <div className="mt-10 overflow-x-auto rounded-xl border border-border bg-white shadow-sm">
            <table className="w-full min-w-[800px] text-left text-sm">
              <thead>
                <tr className="border-b border-border text-xs font-medium text-muted">
                  <th className="px-5 py-3">CAMPAIGN</th>
                  <th className="px-5 py-3">EVENT</th>
                  <th className="px-5 py-3">AMOUNT</th>
                  <th className="px-5 py-3">TX HASH</th>
                  <th className="px-5 py-3">STATUS</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {activity.length === 0 && (
                  <tr>
                    <td colSpan={5} className="px-5 py-6 text-center text-muted">
                      No activity yet.
                    </td>
                  </tr>
                )}
                {activity.map((a) => (
                  <tr key={a.id} className="transition-colors hover:bg-page/60">
                    <td className="px-5 py-3.5">
                      <Link to={`/campaigns/${a.campaignId}`} className="text-ink hover:text-accent">
                        {titleById.get(a.campaignId) ?? `Campaign #${a.campaignId}`}
                      </Link>
                    </td>
                    <td className="px-5 py-3.5 text-ink">{activityEventLabel(a)}</td>
                    <td className="px-5 py-3.5 text-ink">{a.amount ? `${weiToEth(a.amount)} ETH` : '—'}</td>
                    <td className="px-5 py-3.5">
                      <a
                        href={etherscanTxUrl(a.txHash)}
                        target="_blank"
                        rel="noreferrer"
                        className="text-accent underline"
                      >
                        {shortAddress(a.txHash)}
                      </a>
                    </td>
                    <td className="px-5 py-3.5">
                      {a.eventType === 'ReceiptSubmitted' ? (
                        <span className="rounded-full bg-proven-tint px-2.5 py-1 text-xs font-medium text-proven">
                          Proven
                        </span>
                      ) : (
                        <span className="rounded-full bg-locked-tint px-2.5 py-1 text-xs font-medium text-locked">
                          Confirmed
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  )
}

function StatCard({ label, value, detail }: { label: string; value: string; detail: string }) {
  return (
    <div className="rounded-xl border border-border bg-white p-6 shadow-sm transition-shadow hover:shadow-md">
      <p className="text-xs font-medium tracking-wide text-muted">{label}</p>
      <p className="mt-2 text-3xl font-bold text-ink">{value}</p>
      <p className="mt-2 text-sm text-muted">{detail}</p>
    </div>
  )
}

import type { Campaign, Milestone } from '../lib/api'
import { weiToEth } from '../lib/format'

// "Where the money is" — splits what a campaign has raised into the three
// places it can be, all derived from data we already have:
//   Proven   = milestones in state 2: released AND backed by an on-chain receipt
//   Released = milestones in state 1: sent to the charity, receipt still owed
//   Escrow   = everything else raised: still locked in the contract, not with the charity
// Escrow is raised minus the released milestones, so an over-funded or
// not-yet-released campaign still adds up to exactly `raised`.
export default function FundsFlow({ campaign, milestones }: { campaign: Campaign; milestones: Milestone[] }) {
  const raised = BigInt(campaign.raised)
  if (raised === 0n) return null

  const sumByState = (state: number) =>
    milestones.filter((m) => m.state === state).reduce((acc, m) => acc + BigInt(m.amount), 0n)
  const proven = sumByState(2)
  const released = sumByState(1)
  let escrow = raised - proven - released
  if (escrow < 0n) escrow = 0n

  const segments = [
    { key: 'proven', label: 'Spent & proven', hint: 'released, receipt on-chain', value: proven, bar: 'bg-proven', dot: 'bg-proven' },
    { key: 'released', label: 'With the charity', hint: 'released, receipt still owed', value: released, bar: 'bg-released', dot: 'bg-released' },
    { key: 'escrow', label: 'Locked in escrow', hint: 'not with the charity yet', value: escrow, bar: 'bg-locked', dot: 'bg-locked' },
  ]
  const pct = (v: bigint) => Number((v * 10000n) / raised) / 100

  return (
    <div className="mb-10 rounded-xl border border-border bg-white p-5 shadow-sm">
      <h2 className="text-lg font-semibold text-ink">Where the money is</h2>
      <p className="mt-1 text-sm text-muted">
        Every ETH donated to this campaign is in exactly one of these three places.
      </p>

      <div className="mt-4 flex h-3 w-full overflow-hidden rounded-full bg-border" role="img" aria-label="Funds split by stage">
        {segments.map((s) => (
          <div key={s.key} className={s.bar} style={{ width: `${pct(s.value)}%` }} />
        ))}
      </div>

      <ul className="mt-4 grid gap-3 sm:grid-cols-3">
        {segments.map((s) => (
          <li key={s.key} className="flex items-start gap-2.5">
            <span className={`mt-1.5 h-2.5 w-2.5 shrink-0 rounded-full ${s.dot}`} aria-hidden="true" />
            <div>
              <p className="text-sm font-medium text-ink">
                {weiToEth(s.value.toString(), 4)} ETH <span className="font-normal text-muted">· {Math.round(pct(s.value))}%</span>
              </p>
              <p className="text-sm text-ink">{s.label}</p>
              <p className="text-xs text-muted">{s.hint}</p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}

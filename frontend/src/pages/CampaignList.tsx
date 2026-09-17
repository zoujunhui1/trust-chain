import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import CampaignCard from '../components/CampaignCard'
import { getCampaign, getCharity, listCampaigns } from '../lib/api'
import type { Campaign, Milestone } from '../lib/api'

export default function CampaignList() {
  const [campaigns, setCampaigns] = useState<Campaign[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [verifiedByAddress, setVerifiedByAddress] = useState<Record<string, boolean>>({})
  const [milestonesById, setMilestonesById] = useState<Record<number, Milestone[]>>({})

  useEffect(() => {
    let cancelled = false

    listCampaigns()
      .then(async (list) => {
        if (cancelled) return
        setCampaigns(list)

        // Milestone states and charity verification aren't in the list payload,
        // so fetch them per campaign (dedup charities) once the list is in.
        const uniqueCharities = [...new Set(list.map((c) => c.charity))]
        const [charityResults, detailResults] = await Promise.all([
          Promise.allSettled(uniqueCharities.map((addr) => getCharity(addr))),
          Promise.allSettled(list.map((c) => getCampaign(c.id))),
        ])
        if (cancelled) return

        const verified: Record<string, boolean> = {}
        charityResults.forEach((res, i) => {
          if (res.status === 'fulfilled') verified[uniqueCharities[i]] = res.value.verified
        })
        setVerifiedByAddress(verified)

        const milestones: Record<number, Milestone[]> = {}
        detailResults.forEach((res, i) => {
          if (res.status === 'fulfilled') milestones[list[i].id] = res.value.milestones
        })
        setMilestonesById(milestones)
      })
      .catch((err: Error) => {
        if (!cancelled) setError(err.message)
      })

    return () => {
      cancelled = true
    }
  }, [])

  return (
    <div className="mx-auto max-w-7xl px-6 py-12 sm:px-10 lg:px-16">
      <h1 className="text-3xl font-bold tracking-tight text-ink">Campaigns</h1>
      <p className="mt-2 max-w-2xl text-muted">
        Support verified charities — funds are only released once each milestone is proven on-chain.
      </p>

      <section className="mt-10 grid grid-cols-1 gap-6 rounded-xl border border-border bg-white p-6 shadow-sm sm:p-8 lg:grid-cols-[1.15fr_1fr]">
        <div>
          <h2 className="text-xs font-semibold uppercase tracking-wide text-muted">How it works</h2>
          <ol className="mt-4 space-y-4">
            <Step
              n={1}
              title="Connect your wallet"
              text="Click Connect Wallet and approve in MetaMask — no account or signup needed."
            />
            <Step
              n={2}
              title="Donate, or create a campaign"
              text="Donors send ETH to any campaign. A registry-verified charity can create one and split its goal into milestones."
            />
            <Step
              n={3}
              title="Funds release milestone by milestone"
              text="A charity only receives each milestone's ETH after submitting proof on-chain — never the full amount up front."
            />
          </ol>
        </div>

        <div
          className="rounded-lg p-5"
          style={{ background: 'linear-gradient(135deg, rgba(30,58,95,0.06), rgba(124,58,237,0.06))' }}
        >
          <h2 className="text-xs font-semibold uppercase tracking-wide text-muted">Why TrustChain</h2>
          <p className="mt-3 text-sm text-ink">
            Your donation never sits in a charity's bank account. It's locked in a smart contract and
            released one milestone at a time, only after the charity proves it on-chain — every step
            is a public transaction you can check yourself on the{' '}
            <Link to="/transparency" className="font-medium text-accent underline">
              Transparency Dashboard
            </Link>{' '}
            or Sepolia Etherscan.
          </p>
        </div>
      </section>

      {error && (
        <p className="mt-8 text-sm text-red-600">Couldn't load campaigns: {error}</p>
      )}

      {!error && campaigns === null && (
        <p className="mt-8 text-muted">Loading campaigns…</p>
      )}

      {campaigns !== null && campaigns.length === 0 && (
        <p className="mt-8 text-muted">No campaigns yet.</p>
      )}

      {campaigns !== null && campaigns.length > 0 && (
        <div className="mt-8 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
          {campaigns.map((c) => (
            <CampaignCard
              key={c.id}
              campaign={c}
              verified={verifiedByAddress[c.charity]}
              milestones={milestonesById[c.id]}
            />
          ))}
        </div>
      )}
    </div>
  )
}

function Step({ n, title, text }: { n: number; title: string; text: string }) {
  return (
    <li className="flex gap-3">
      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent-tint text-xs font-semibold text-accent">
        {n}
      </span>
      <div>
        <p className="text-sm font-medium text-ink">{title}</p>
        <p className="mt-0.5 text-sm text-muted">{text}</p>
      </div>
    </li>
  )
}

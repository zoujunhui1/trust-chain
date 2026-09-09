import { useEffect, useState } from 'react'
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

import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getCharity } from '../lib/api'
import { createCampaign, parseCreatedCampaignId } from '../lib/escrow'
import { shortAddress } from '../lib/format'
import { useWallet } from '../lib/wallet'

type TxStatus = 'idle' | 'pending' | 'success' | 'error'

export default function CreateCampaign() {
  const wallet = useWallet()
  const [verified, setVerified] = useState<boolean | null>(null)
  const [milestones, setMilestones] = useState<string[]>([''])
  const [status, setStatus] = useState<TxStatus>('idle')
  const [txHash, setTxHash] = useState<string | null>(null)
  const [createdId, setCreatedId] = useState<number | null>(null)
  const [txError, setTxError] = useState<string | null>(null)

  useEffect(() => {
    if (!wallet.address) {
      setVerified(null)
      return
    }
    let cancelled = false
    setVerified(null)
    getCharity(wallet.address)
      .then((c) => {
        if (!cancelled) setVerified(c.verified)
      })
      .catch(() => {
        if (!cancelled) setVerified(false)
      })
    return () => {
      cancelled = true
    }
  }, [wallet.address])

  const amounts = milestones.map((m) => m.trim())
  const validAmounts = amounts.length > 0 && amounts.every((a) => Number(a) > 0)
  const canSubmit = wallet.address && verified && !wallet.isWrongNetwork && validAmounts && status !== 'pending'

  function updateMilestone(i: number, value: string) {
    setMilestones((prev) => prev.map((m, idx) => (idx === i ? value : m)))
  }

  function addMilestone() {
    setMilestones((prev) => [...prev, ''])
  }

  function removeMilestone(i: number) {
    setMilestones((prev) => prev.filter((_, idx) => idx !== i))
  }

  async function handleSubmit() {
    setStatus('pending')
    setTxError(null)
    setTxHash(null)
    setCreatedId(null)
    try {
      const signer = await wallet.getSigner()
      const tx = await createCampaign(signer, amounts)
      setTxHash(tx.hash)
      const receipt = await tx.wait()
      setCreatedId(parseCreatedCampaignId(receipt))
      setStatus('success')
      setMilestones([''])
    } catch (err) {
      setStatus('error')
      setTxError(err instanceof Error ? err.message : 'Transaction failed.')
    }
  }

  return (
    <div className="px-16 py-12">
      <h1 className="text-3xl font-bold text-ink">Create a Campaign</h1>
      <p className="mt-2 max-w-2xl text-muted">
        Only registry-verified charities can create campaigns. Set the payout amount for each
        milestone — the goal is their sum, and funds release one milestone at a time as you submit
        proof.
      </p>

      <div className="mt-8 max-w-xl rounded-xl border border-border bg-white p-6">
        {!wallet.address && (
          <>
            <p className="text-sm text-muted">Connect your charity's wallet to get started.</p>
            <button
              type="button"
              onClick={wallet.connect}
              disabled={wallet.connecting}
              className="mt-4 rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
            >
              {wallet.connecting ? 'Connecting…' : 'Connect Wallet'}
            </button>
            {wallet.error && <p className="mt-2 text-xs text-red-600">{wallet.error}</p>}
          </>
        )}

        {wallet.address && wallet.isWrongNetwork && (
          <>
            <p className="text-sm text-muted">Wrong network — this app runs on Sepolia.</p>
            <button
              type="button"
              onClick={wallet.switchNetwork}
              className="mt-4 rounded-lg bg-released px-4 py-2.5 text-sm font-medium text-white hover:opacity-90"
            >
              Switch to Sepolia
            </button>
          </>
        )}

        {wallet.address && !wallet.isWrongNetwork && verified === null && (
          <p className="text-sm text-muted">Checking verification status…</p>
        )}

        {wallet.address && !wallet.isWrongNetwork && verified === false && (
          <div>
            <p className="text-sm text-ink">
              <span className="font-medium">{shortAddress(wallet.address)}</span> isn't a
              registry-verified charity yet.
            </p>
            <p className="mt-2 text-sm text-muted">
              The platform admin needs to call <code className="text-xs">verifyCharity</code> for this
              address on <code className="text-xs">CharityRegistry</code> before it can create
              campaigns.
            </p>
          </div>
        )}

        {wallet.address && !wallet.isWrongNetwork && verified && (
          <div>
            <p className="text-sm font-medium text-ink">Milestones</p>
            <div className="mt-3 space-y-3">
              {milestones.map((m, i) => (
                <div key={i} className="flex items-center gap-2">
                  <span className="w-6 text-sm text-muted">{i + 1}.</span>
                  <div className="flex flex-1 items-center rounded-lg border border-border px-3 py-2">
                    <input
                      type="number"
                      min="0"
                      step="0.001"
                      placeholder="0.00"
                      value={m}
                      onChange={(e) => updateMilestone(i, e.target.value)}
                      disabled={status === 'pending'}
                      className="w-full text-ink outline-none"
                    />
                    <span className="text-sm text-muted">ETH</span>
                  </div>
                  <button
                    type="button"
                    onClick={() => removeMilestone(i)}
                    disabled={milestones.length === 1 || status === 'pending'}
                    className="px-2 text-sm text-muted hover:text-red-600 disabled:opacity-30"
                    aria-label="Remove milestone"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>

            <button
              type="button"
              onClick={addMilestone}
              disabled={status === 'pending'}
              className="mt-3 text-sm font-medium text-accent hover:opacity-80"
            >
              + Add Milestone
            </button>

            <button
              type="button"
              onClick={handleSubmit}
              disabled={!canSubmit}
              className="mt-6 w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {status === 'pending' ? 'Confirming…' : 'Create Campaign'}
            </button>

            {status === 'error' && txError && <p className="mt-3 text-xs text-red-600">{txError}</p>}

            {status === 'pending' && txHash && (
              <p className="mt-3 text-xs text-muted">
                Waiting for confirmation —{' '}
                <a
                  href={`https://sepolia.etherscan.io/tx/${txHash}`}
                  target="_blank"
                  rel="noreferrer"
                  className="text-accent underline"
                >
                  view on Etherscan
                </a>
              </p>
            )}

            {status === 'success' && txHash && (
              <p className="mt-3 text-xs text-proven">
                Campaign created —{' '}
                <a
                  href={`https://sepolia.etherscan.io/tx/${txHash}`}
                  target="_blank"
                  rel="noreferrer"
                  className="underline"
                >
                  view on Etherscan
                </a>
                .{' '}
                {createdId !== null ? (
                  <>
                    It'll show up on{' '}
                    <Link to={`/campaigns/${createdId}`} className="underline">
                      its campaign page
                    </Link>{' '}
                    once the indexer picks it up.
                  </>
                ) : (
                  "It'll show up on the campaigns list once the indexer picks it up."
                )}
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

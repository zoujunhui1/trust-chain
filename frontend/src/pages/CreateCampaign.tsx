import { useEffect, useState } from 'react'
import { parseEther } from 'ethers'
import { useNavigate } from 'react-router-dom'
import { createCampaignRecord, getCharity, setCampaignMetadata, setMilestoneDescriptions } from '../lib/api'
import { createCampaign, parseCreatedCampaignId } from '../lib/escrow'
import { shortAddress } from '../lib/format'
import { listThemes } from '../lib/theme'
import { useWallet } from '../lib/wallet'

type TxStatus = 'idle' | 'pending' | 'success' | 'error'

export default function CreateCampaign() {
  const wallet = useWallet()
  const navigate = useNavigate()
  const [verified, setVerified] = useState<boolean | null>(null)
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [theme, setTheme] = useState('') // a CampaignTheme.key, or '' for "no preference"
  const [milestones, setMilestones] = useState<string[]>([''])
  // Index-aligned with milestones — descriptions[i] describes milestone i.
  // Optional, so an empty string is a valid "nothing to say" value.
  const [milestoneDescriptions, setMilestoneDescriptionsState] = useState<string[]>([''])
  const [status, setStatus] = useState<TxStatus>('idle')
  const [txHash, setTxHash] = useState<string | null>(null)
  const [createdId, setCreatedId] = useState<number | null>(null)
  const [txError, setTxError] = useState<string | null>(null)
  // The on-chain tx can still succeed even if these off-chain saves fail —
  // tracked separately so a save hiccup doesn't look like the campaign
  // itself failed to create.
  const [recordSaveError, setRecordSaveError] = useState<string | null>(null)
  const [metadataSaveError, setMetadataSaveError] = useState<string | null>(null)
  const [milestoneDescSaveError, setMilestoneDescSaveError] = useState<string | null>(null)

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
  const canSubmit =
    wallet.address && verified && !wallet.isWrongNetwork && title.trim() && validAmounts && status !== 'pending'

  function updateMilestone(i: number, value: string) {
    setMilestones((prev) => prev.map((m, idx) => (idx === i ? value : m)))
  }

  function updateMilestoneDescription(i: number, value: string) {
    setMilestoneDescriptionsState((prev) => prev.map((d, idx) => (idx === i ? value : d)))
  }

  function addMilestone() {
    setMilestones((prev) => [...prev, ''])
    setMilestoneDescriptionsState((prev) => [...prev, ''])
  }

  function removeMilestone(i: number) {
    setMilestones((prev) => prev.filter((_, idx) => idx !== i))
    setMilestoneDescriptionsState((prev) => prev.filter((_, idx) => idx !== i))
  }

  async function handleSubmit() {
    setStatus('pending')
    setTxError(null)
    setTxHash(null)
    setCreatedId(null)
    setRecordSaveError(null)
    setMetadataSaveError(null)
    setMilestoneDescSaveError(null)
    try {
      const charity = wallet.address!
      const signer = await wallet.getSigner()
      const tx = await createCampaign(signer, amounts)
      setTxHash(tx.hash)
      const receipt = await tx.wait()
      const id = parseCreatedCampaignId(receipt)
      setCreatedId(id)
      setStatus('success')
      const savedTitle = title.trim()
      const savedDescription = description.trim()
      const savedTheme = theme
      const savedMilestoneDescriptions = milestoneDescriptions
      setMilestones([''])
      setMilestoneDescriptionsState([''])
      setTitle('')
      setDescription('')
      setTheme('')

      if (id !== null && receipt) {
        // Both calls are independent inserts (no FK between campaigns and
        // campaign_metadata — see that table's schema comment), so they can
        // go out in parallel and fail independently.
        createCampaignRecord({
          id,
          charity,
          milestoneAmountsWei: amounts.map((a) => parseEther(a).toString()),
          createdBlock: receipt.blockNumber,
          createdTx: receipt.hash,
        }).catch((err) => {
          setRecordSaveError(err instanceof Error ? err.message : 'Failed to save the campaign record.')
        })
        setCampaignMetadata(id, {
          title: savedTitle,
          description: savedDescription || undefined,
          theme: savedTheme || undefined,
        }).catch((err) => {
          setMetadataSaveError(err instanceof Error ? err.message : 'Failed to save the campaign details.')
        })
        setMilestoneDescriptions(id, savedMilestoneDescriptions).catch((err) => {
          setMilestoneDescSaveError(err instanceof Error ? err.message : 'Failed to save the milestone descriptions.')
        })
        // It's already in the campaigns list (see createCampaignRecord above)
        // — go straight there instead of making the charity click through.
        navigate('/')
      }
    } catch (err) {
      setStatus('error')
      setTxError(err instanceof Error ? err.message : 'Transaction failed.')
    }
  }

  return (
    <div className="mx-auto max-w-7xl px-6 py-12 sm:px-10 lg:px-16">
      <h1 className="text-3xl font-bold tracking-tight text-ink">Create a Campaign</h1>
      <p className="mt-2 max-w-2xl text-muted">
        Only registry-verified charities can create campaigns. Set the payout amount for each
        milestone — the goal is their sum, and funds release one milestone at a time as you submit
        proof.
      </p>

      <div className="mt-8 max-w-xl rounded-xl border border-border bg-white p-6 shadow-sm">
        {!wallet.address && (
          <>
            <p className="text-sm text-muted">Connect your charity's wallet to get started.</p>
            <button
              type="button"
              onClick={wallet.connect}
              disabled={wallet.connecting}
              className="mt-4 rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:opacity-50"
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
              className="mt-4 rounded-lg bg-released px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98]"
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
            <p className="text-sm font-medium text-ink">Title</p>
            <input
              type="text"
              placeholder="e.g. Clean Water for Riverside Village"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={status === 'pending'}
              maxLength={200}
              className="mt-3 w-full rounded-lg border border-border px-3 py-2 text-sm text-ink outline-none transition-colors focus:border-accent"
            />

            <p className="mt-6 text-sm font-medium text-ink">
              Details <span className="font-normal text-muted">(optional)</span>
            </p>
            <textarea
              placeholder="What is this campaign for, and how will the funds be used?"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={status === 'pending'}
              rows={3}
              maxLength={5000}
              className="mt-3 w-full resize-none rounded-lg border border-border px-3 py-2 text-sm text-ink outline-none transition-colors focus:border-accent"
            />

            <p className="mt-6 text-sm font-medium text-ink">
              Theme <span className="font-normal text-muted">(optional — picks one for you otherwise)</span>
            </p>
            <div className="mt-3 flex flex-wrap gap-2">
              {listThemes().map((t) => (
                <button
                  key={t.key}
                  type="button"
                  onClick={() => setTheme((prev) => (prev === t.key ? '' : t.key))}
                  disabled={status === 'pending'}
                  className="flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-xs font-medium transition-colors disabled:opacity-50"
                  style={
                    theme === t.key
                      ? { borderColor: t.accent, color: t.accent, backgroundColor: `${t.accent}14` }
                      : { borderColor: 'var(--color-border)', color: 'var(--color-muted)' }
                  }
                >
                  <span aria-hidden="true">{t.emoji}</span>
                  {t.label}
                </button>
              ))}
            </div>

            <p className="mt-6 text-sm font-medium text-ink">Milestones</p>
            <div className="mt-3 space-y-3">
              {milestones.map((m, i) => (
                <div key={i} className="rounded-lg border border-border p-2.5">
                  <div className="flex items-center gap-2">
                    <span className="w-6 text-sm text-muted">{i + 1}.</span>
                    <div className="flex flex-1 items-center rounded-lg border border-border px-3 py-2 transition-colors focus-within:border-accent">
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
                  <input
                    type="text"
                    placeholder="What will this milestone accomplish? (optional)"
                    value={milestoneDescriptions[i] ?? ''}
                    onChange={(e) => updateMilestoneDescription(i, e.target.value)}
                    disabled={status === 'pending'}
                    maxLength={2000}
                    className="mt-2 ml-8 w-[calc(100%-2rem)] rounded-lg border border-border px-3 py-1.5 text-xs text-ink outline-none transition-colors focus:border-accent"
                  />
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
              className="mt-6 w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 disabled:shadow-none disabled:active:scale-100"
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

            {/* createdId !== null redirects to the campaigns list immediately (see
                handleSubmit) — this only has time to show for the rare case where
                the id couldn't be parsed from the receipt. */}
            {status === 'success' && txHash && createdId === null && (
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
                . It'll show up on the campaigns list once the indexer picks it up.
              </p>
            )}

            {recordSaveError && (
              <p className="mt-2 text-xs text-red-600">
                Campaign created, but it may take longer than usual to appear: {recordSaveError}
              </p>
            )}
            {metadataSaveError && (
              <p className="mt-2 text-xs text-red-600">
                Campaign created, but saving its details failed: {metadataSaveError}
              </p>
            )}
            {milestoneDescSaveError && (
              <p className="mt-2 text-xs text-red-600">
                Campaign created, but saving milestone descriptions failed: {milestoneDescSaveError}
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

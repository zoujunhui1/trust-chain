import { useState } from 'react'
import type { Milestone, MilestoneState } from '../lib/api'
import { milestoneReceiptFileUrl, uploadMilestoneReceipt } from '../lib/api'
import { sha256File } from '../lib/hash'
import { releaseMilestone, submitReceipt } from '../lib/escrow'
import { useWallet } from '../lib/wallet'

interface MilestoneReceiptPanelProps {
  campaignId: number
  milestone: Milestone
  prevState: MilestoneState | undefined // undefined for milestone 0
  goalReached: boolean
  isCharity: boolean // connected wallet === this campaign's charity
  accent: string
  // Tell the page a tx of ours just confirmed on-chain, so it can show the new
  // state now instead of waiting for the indexer (see CampaignDetail).
  onStateChange: (idx: number, state: MilestoneState, receiptHash?: string) => void
}

type ActionStatus = 'idle' | 'pending' | 'success' | 'error'

// ethers' errors carry a long dump (calldata, tx object…); show only the
// short reason, e.g. "user rejected action" or the contract's revert text.
function txErrorMessage(err: unknown, fallback: string): string {
  if (err && typeof err === 'object') {
    const e = err as { reason?: string; shortMessage?: string }
    if (e.reason) return e.reason
    if (e.shortMessage) return e.shortMessage
  }
  return err instanceof Error ? err.message.slice(0, 200) : fallback
}

// The charity-only controls for one milestone: "Release" while Locked, then
// a file upload for "Submit Receipt" once Released. Renders nothing for a
// visitor who isn't this campaign's charity, or once the milestone is Proven
// with nothing left to do (CampaignDetail shows the uploaded file itself).
export default function MilestoneReceiptPanel({
  campaignId,
  milestone,
  prevState,
  goalReached,
  isCharity,
  accent,
  onStateChange,
}: MilestoneReceiptPanelProps) {
  const wallet = useWallet()
  const [releaseStatus, setReleaseStatus] = useState<ActionStatus>('idle')
  const [releaseError, setReleaseError] = useState<string | null>(null)

  const [file, setFile] = useState<File | null>(null)
  const [note, setNote] = useState('')
  const [receiptStatus, setReceiptStatus] = useState<ActionStatus>('idle')
  const [receiptError, setReceiptError] = useState<string | null>(null)
  // Set as soon as the upload succeeds, so the "uploaded" view appears
  // immediately — no need to wait for the indexer to confirm the on-chain
  // side before showing what was just uploaded.
  const [justUploaded, setJustUploaded] = useState<{ fileName: string; sha256: string } | null>(null)

  if (!isCharity) return null

  const canRelease = milestone.state === 0 && goalReached && (prevState === undefined || prevState === 2)
  const canSubmitReceipt = milestone.state === 1

  async function handleRelease() {
    setReleaseStatus('pending')
    setReleaseError(null)
    try {
      const signer = await wallet.getSigner()
      const tx = await releaseMilestone(signer, campaignId, milestone.idx)
      await tx.wait()
      setReleaseStatus('success')
      onStateChange(milestone.idx, 1)
    } catch (err) {
      setReleaseStatus('error')
      setReleaseError(txErrorMessage(err, 'Transaction failed.'))
    }
  }

  async function handleSubmitReceipt() {
    if (!file) return
    setReceiptStatus('pending')
    setReceiptError(null)
    try {
      const hash = await sha256File(file)
      const signer = await wallet.getSigner()
      const tx = await submitReceipt(signer, campaignId, milestone.idx, hash)
      await tx.wait()
      onStateChange(milestone.idx, 2, hash)
      const uploaded = await uploadMilestoneReceipt(campaignId, milestone.idx, file, note.trim() || undefined)
      setJustUploaded({ fileName: uploaded.fileName, sha256: uploaded.sha256 })
      setReceiptStatus('success')
      setFile(null)
      setNote('')
    } catch (err) {
      setReceiptStatus('error')
      setReceiptError(txErrorMessage(err, 'Submitting the receipt failed.'))
    }
  }

  if (justUploaded) {
    return (
      <div className="mt-3 rounded-lg border border-border bg-page/60 p-3 text-sm">
        <p className="font-medium text-proven">Receipt submitted and uploaded ✓</p>
        <p className="mt-1 text-muted">
          {justUploaded.fileName} —{' '}
          <a
            href={milestoneReceiptFileUrl(campaignId, milestone.idx)}
            target="_blank"
            rel="noreferrer"
            className="underline"
            style={{ color: accent }}
          >
            view file
          </a>
        </p>
        <p className="mt-1 text-xs text-muted">
          It can take a minute or two for the on-chain confirmation to catch up and mark this milestone
          Proven.
        </p>
      </div>
    )
  }

  if (!canRelease && !canSubmitReceipt) return null

  return (
    <div className="mt-3 rounded-lg border border-border bg-page/60 p-3">
      {canRelease && (
        <>
          <button
            type="button"
            onClick={handleRelease}
            disabled={releaseStatus === 'pending'}
            className="rounded-lg px-3 py-1.5 text-xs font-medium text-white shadow-sm transition-all hover:opacity-90 disabled:opacity-50"
            style={{ backgroundColor: accent }}
          >
            {releaseStatus === 'pending' ? 'Releasing…' : `Release Milestone ${milestone.idx + 1}`}
          </button>
          {releaseStatus === 'error' && releaseError && (
            <p className="mt-2 text-xs text-red-600">{releaseError}</p>
          )}
        </>
      )}

      {canSubmitReceipt && (
        <div>
          <p className="text-xs font-medium text-ink">Submit the spending receipt</p>
          <p className="mt-1 text-xs text-muted">
            Upload a photo or PDF of the receipt — its hash goes on-chain, the file itself gets stored so
            donors can open it.
          </p>
          <input
            type="file"
            accept="image/png,image/jpeg,image/webp,application/pdf"
            onChange={(e) => setFile(e.target.files?.[0] ?? null)}
            disabled={receiptStatus === 'pending'}
            className="mt-2 block w-full text-xs text-muted file:mr-3 file:rounded-lg file:border-0 file:bg-accent-tint file:px-3 file:py-1.5 file:text-xs file:font-medium file:text-accent"
          />
          <textarea
            placeholder="What was this spent on? (optional)"
            value={note}
            onChange={(e) => setNote(e.target.value)}
            disabled={receiptStatus === 'pending'}
            rows={2}
            maxLength={2000}
            className="mt-2 w-full resize-none rounded-lg border border-border px-3 py-2 text-xs text-ink outline-none transition-colors focus:border-accent"
          />
          <button
            type="button"
            onClick={handleSubmitReceipt}
            disabled={!file || receiptStatus === 'pending'}
            className="mt-2 rounded-lg px-3 py-1.5 text-xs font-medium text-white shadow-sm transition-all hover:opacity-90 disabled:opacity-50"
            style={{ backgroundColor: accent }}
          >
            {receiptStatus === 'pending' ? 'Submitting…' : 'Submit Receipt'}
          </button>
          {receiptStatus === 'error' && receiptError && (
            <p className="mt-2 text-xs text-red-600">{receiptError}</p>
          )}
        </div>
      )}
    </div>
  )
}

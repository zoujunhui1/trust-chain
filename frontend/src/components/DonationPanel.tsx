import { useState } from 'react'
import type { Donation } from '../lib/api'
import { relativeTime, shortAddress, weiToEth } from '../lib/format'
import { donate } from '../lib/escrow'
import { useWallet } from '../lib/wallet'

interface DonationPanelProps {
  campaignId: number
  donations: Donation[]
  // Ties the panel's accent to the campaign's theme (lib/theme.ts) so the
  // donate CTA visually matches the card/hero the donor arrived from.
  accent: string
}

type TxStatus = 'idle' | 'pending' | 'success' | 'error'

const QUICK_AMOUNTS = ['0.01', '0.05', '0.1', '0.25']

// Cheap deterministic color for a donor avatar, so "Recent Donations" isn't a
// wall of identical gray initials — same idea as lib/theme.ts's id-cycling.
const AVATAR_PALETTE = ['#2f6b3f', '#9a5b13', '#8a4b8f', '#1d6f8c', '#b0331a', '#1e3a5f']
function avatarColor(address: string): string {
  const sum = [...address.slice(2, 10)].reduce((acc, ch) => acc + ch.charCodeAt(0), 0)
  return AVATAR_PALETTE[sum % AVATAR_PALETTE.length]
}

export default function DonationPanel({ campaignId, donations, accent }: DonationPanelProps) {
  const wallet = useWallet()
  const [amount, setAmount] = useState('')
  const [status, setStatus] = useState<TxStatus>('idle')
  const [txHash, setTxHash] = useState<string | null>(null)
  const [txError, setTxError] = useState<string | null>(null)

  const canSubmit = wallet.address && !wallet.isWrongNetwork && Number(amount) > 0 && status !== 'pending'

  async function handleSubmit() {
    setStatus('pending')
    setTxError(null)
    setTxHash(null)
    try {
      const signer = await wallet.getSigner()
      const tx = await donate(signer, campaignId, amount)
      setTxHash(tx.hash)
      await tx.wait()
      setStatus('success')
      setAmount('')
    } catch (err) {
      setStatus('error')
      setTxError(err instanceof Error ? err.message : 'Transaction failed.')
    }
  }

  return (
    <div className="w-full shrink-0 rounded-xl border border-border bg-white p-6 shadow-sm lg:w-[380px]">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-ink">Make a Donation</h2>
        {donations.length > 0 && (
          <span className="text-xs font-medium text-muted">
            {donations.length} supporter{donations.length === 1 ? '' : 's'}
          </span>
        )}
      </div>

      <div
        className="mt-5 flex items-center rounded-lg border border-border px-4 py-3 transition-colors focus-within:border-2 focus-within:px-[15px] focus-within:py-[11px]"
        style={{ borderColor: amount ? accent : undefined }}
      >
        <input
          type="number"
          min="0"
          step="0.001"
          placeholder="0.00"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          disabled={status === 'pending'}
          className="w-full text-lg font-medium text-ink outline-none disabled:bg-transparent disabled:text-muted"
        />
        <span className="text-sm text-muted">ETH</span>
      </div>

      <div className="mt-2.5 flex gap-2">
        {QUICK_AMOUNTS.map((preset) => (
          <button
            key={preset}
            type="button"
            onClick={() => setAmount(preset)}
            disabled={status === 'pending'}
            className="flex-1 rounded-lg border px-2 py-1.5 text-xs font-medium transition-colors disabled:opacity-50"
            style={
              amount === preset
                ? { borderColor: accent, color: accent, backgroundColor: `${accent}14` }
                : { borderColor: 'var(--color-border)', color: 'var(--color-muted)' }
            }
          >
            {preset}
          </button>
        ))}
      </div>

      {!wallet.address ? (
        <button
          type="button"
          onClick={wallet.connect}
          disabled={wallet.connecting}
          className="mt-4 w-full rounded-lg px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:opacity-50"
          style={{ backgroundColor: accent }}
        >
          {wallet.connecting ? 'Connecting…' : 'Connect Wallet to Donate'}
        </button>
      ) : wallet.isWrongNetwork ? (
        <button
          type="button"
          onClick={wallet.switchNetwork}
          className="mt-4 w-full rounded-lg bg-released px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98]"
        >
          Switch to Sepolia
        </button>
      ) : (
        <button
          type="button"
          onClick={handleSubmit}
          disabled={!canSubmit}
          className="mt-4 w-full rounded-lg px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 disabled:shadow-none disabled:active:scale-100"
          style={{ backgroundColor: accent }}
        >
          {status === 'pending' ? 'Confirming…' : 'Donate Now'}
        </button>
      )}

      <p className="mt-3 text-xs text-muted">
        Connects via MetaMask — funds go directly to the escrow contract.
      </p>

      {wallet.error && <p className="mt-2 text-xs text-red-600">{wallet.error}</p>}
      {status === 'error' && txError && <p className="mt-2 text-xs text-red-600">{txError}</p>}
      {status === 'pending' && txHash && (
        <p className="mt-2 text-xs text-muted">
          Waiting for confirmation —{' '}
          <a
            href={`https://sepolia.etherscan.io/tx/${txHash}`}
            target="_blank"
            rel="noreferrer"
            className="underline"
            style={{ color: accent }}
          >
            view on Etherscan
          </a>
        </p>
      )}
      {status === 'success' && txHash && (
        <p className="mt-2 text-xs text-proven">
          Donation confirmed —{' '}
          <a
            href={`https://sepolia.etherscan.io/tx/${txHash}`}
            target="_blank"
            rel="noreferrer"
            className="underline"
          >
            view on Etherscan
          </a>
          . It'll show up below once the indexer picks it up.
        </p>
      )}

      <div className="mt-6 border-t border-border pt-5">
        <h3 className="text-sm font-semibold text-ink">Recent Donations</h3>
        {donations.length === 0 ? (
          <p className="mt-3 text-sm text-muted">No donations yet — be the first to support this campaign.</p>
        ) : (
          <ul className="mt-3 space-y-3">
            {donations.map((d) => (
              <li key={`${d.txHash}-${d.logIndex}`} className="flex items-center gap-3 text-sm">
                <span
                  className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-[10px] font-semibold uppercase text-white"
                  style={{ backgroundColor: avatarColor(d.donor) }}
                  aria-hidden="true"
                >
                  {d.donor.slice(2, 4)}
                </span>
                <span className="flex-1 text-ink">{shortAddress(d.donor)}</span>
                <span className="text-muted">
                  {weiToEth(d.amount)} ETH · {relativeTime(d.createdAt)}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

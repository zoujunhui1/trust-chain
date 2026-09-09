import { useState } from 'react'
import type { Donation } from '../lib/api'
import { relativeTime, shortAddress, weiToEth } from '../lib/format'
import { donate } from '../lib/escrow'
import { useWallet } from '../lib/wallet'

interface DonationPanelProps {
  campaignId: number
  donations: Donation[]
}

type TxStatus = 'idle' | 'pending' | 'success' | 'error'

export default function DonationPanel({ campaignId, donations }: DonationPanelProps) {
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
      <h2 className="text-lg font-semibold text-ink">Make a Donation</h2>

      <div className="mt-5 flex items-center rounded-lg border border-border px-4 py-3 transition-colors focus-within:border-accent">
        <input
          type="number"
          min="0"
          step="0.001"
          placeholder="0.00"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          disabled={status === 'pending'}
          className="w-full text-ink outline-none disabled:bg-transparent disabled:text-muted"
        />
        <span className="text-sm text-muted">ETH</span>
      </div>

      {!wallet.address ? (
        <button
          type="button"
          onClick={wallet.connect}
          disabled={wallet.connecting}
          className="mt-4 w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:opacity-50"
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
          className="mt-4 w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 disabled:shadow-none disabled:active:scale-100"
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
            className="text-accent underline"
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
          <p className="mt-3 text-sm text-muted">No donations yet.</p>
        ) : (
          <ul className="mt-3 space-y-3">
            {donations.map((d) => (
              <li key={`${d.txHash}-${d.logIndex}`} className="flex items-center justify-between text-sm">
                <span className="text-ink">{shortAddress(d.donor)}</span>
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

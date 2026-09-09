import { useCallback, useEffect, useState } from 'react'
import { listCharities } from '../lib/api'
import type { Charity } from '../lib/api'
import { getOwner, revokeCharity, verifyCharity } from '../lib/registry'
import { shortAddress } from '../lib/format'
import { useWallet } from '../lib/wallet'

type TxAction = 'verify' | 'revoke'
type TxStatus = 'idle' | 'pending' | 'success' | 'error'

interface TxState {
  action: TxAction
  address: string
  status: TxStatus
  hash: string | null
  error: string | null
}

// Platform-admin page for CharityRegistry.verifyCharity/revokeCharity — the
// two owner-only actions that previously had no UI at all (cast send only).
export default function AdminCharities() {
  const wallet = useWallet()
  const [owner, setOwner] = useState<string | null>(null)
  const [checkingOwner, setCheckingOwner] = useState(false)
  const [charities, setCharities] = useState<Charity[] | null>(null)
  const [charitiesError, setCharitiesError] = useState<string | null>(null)
  const [newAddress, setNewAddress] = useState('')
  const [tx, setTx] = useState<TxState | null>(null)

  const loadCharities = useCallback(() => {
    setCharitiesError(null)
    listCharities()
      .then(setCharities)
      .catch((err) => setCharitiesError(err instanceof Error ? err.message : 'Failed to load charities.'))
  }, [])

  useEffect(() => {
    loadCharities()
  }, [loadCharities])

  useEffect(() => {
    if (!wallet.address || wallet.isWrongNetwork) {
      setOwner(null)
      return
    }
    let cancelled = false
    setCheckingOwner(true)
    wallet
      .getSigner()
      .then((signer) => getOwner(signer))
      .then((addr) => {
        if (!cancelled) setOwner(addr)
      })
      .catch(() => {
        if (!cancelled) setOwner(null)
      })
      .finally(() => {
        if (!cancelled) setCheckingOwner(false)
      })
    return () => {
      cancelled = true
    }
  }, [wallet.address, wallet.isWrongNetwork])

  const isOwner = !!wallet.address && !!owner && wallet.address.toLowerCase() === owner.toLowerCase()
  const isValidAddress = /^0x[a-fA-F0-9]{40}$/.test(newAddress.trim())

  async function runAction(action: TxAction, address: string) {
    setTx({ action, address, status: 'pending', hash: null, error: null })
    try {
      const signer = await wallet.getSigner()
      const tx = action === 'verify' ? await verifyCharity(signer, address) : await revokeCharity(signer, address)
      setTx({ action, address, status: 'pending', hash: tx.hash, error: null })
      await tx.wait()
      setTx({ action, address, status: 'success', hash: tx.hash, error: null })
      if (action === 'verify') setNewAddress('')
      // The registry updates immediately, but /api/charities only reflects it
      // once the indexer has processed the event — refresh anyway so a
      // refresh a few seconds later (or a manual reload) picks it up.
      loadCharities()
    } catch (err) {
      setTx({ action, address, status: 'error', hash: null, error: err instanceof Error ? err.message : 'Transaction failed.' })
    }
  }

  const pendingFor = (address: string) => tx?.address.toLowerCase() === address.toLowerCase() && tx.status === 'pending'

  return (
    <div className="mx-auto max-w-7xl px-6 py-12 sm:px-10 lg:px-16">
      <h1 className="text-3xl font-bold tracking-tight text-ink">Charity Registry Admin</h1>
      <p className="mt-2 max-w-2xl text-muted">
        Verify or revoke charity addresses on <code className="text-xs">CharityRegistry</code>. Only the
        contract owner can call these — everyone else can view this page but every action will fail on-chain.
      </p>

      <div className="mt-8 max-w-2xl rounded-xl border border-border bg-white p-6 shadow-sm">
        {!wallet.address && (
          <>
            <p className="text-sm text-muted">Connect the platform admin wallet to manage charities.</p>
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

        {wallet.address && !wallet.isWrongNetwork && checkingOwner && (
          <p className="text-sm text-muted">Checking admin status…</p>
        )}

        {wallet.address && !wallet.isWrongNetwork && !checkingOwner && !isOwner && (
          <p className="text-sm text-ink">
            <span className="font-medium">{shortAddress(wallet.address)}</span> isn't the registry owner
            (<span className="font-medium">{owner ? shortAddress(owner) : 'unknown'}</span>). Connect the
            admin wallet to verify or revoke charities.
          </p>
        )}

        {wallet.address && !wallet.isWrongNetwork && !checkingOwner && isOwner && (
          <div>
            <p className="text-sm font-medium text-ink">Verify a new charity</p>
            <div className="mt-3 flex items-center gap-2">
              <input
                type="text"
                placeholder="0x..."
                value={newAddress}
                onChange={(e) => setNewAddress(e.target.value)}
                disabled={pendingFor(newAddress)}
                className="flex-1 rounded-lg border border-border px-3 py-2 text-sm text-ink outline-none transition-colors focus:border-accent"
              />
              <button
                type="button"
                onClick={() => runAction('verify', newAddress.trim())}
                disabled={!isValidAddress || pendingFor(newAddress)}
                className="rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 disabled:shadow-none disabled:active:scale-100"
              >
                {pendingFor(newAddress) ? 'Confirming…' : 'Verify'}
              </button>
            </div>

            {tx && tx.address.toLowerCase() === newAddress.trim().toLowerCase() && (
              <TxFeedback tx={tx} />
            )}

            <p className="mt-8 text-sm font-medium text-ink">Known charities</p>
            {charitiesError && <p className="mt-2 text-xs text-red-600">{charitiesError}</p>}
            {!charitiesError && charities === null && <p className="mt-2 text-sm text-muted">Loading…</p>}
            {!charitiesError && charities !== null && charities.length === 0 && (
              <p className="mt-2 text-sm text-muted">No charities indexed yet.</p>
            )}
            {!charitiesError && charities !== null && charities.length > 0 && (
              <div className="mt-3 divide-y divide-border rounded-lg border border-border">
                {charities.map((c) => (
                  <div key={c.address} className="flex items-center justify-between px-4 py-3 transition-colors hover:bg-page/60">
                    <div>
                      <p className="text-sm font-medium text-ink">{shortAddress(c.address)}</p>
                      <p className="text-xs text-muted">{c.verified ? 'Verified' : 'Not verified'}</p>
                    </div>
                    <div className="flex flex-col items-end gap-1">
                      <button
                        type="button"
                        onClick={() => runAction(c.verified ? 'revoke' : 'verify', c.address)}
                        disabled={pendingFor(c.address)}
                        className={
                          c.verified
                            ? 'rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 transition-colors hover:bg-red-50 disabled:opacity-50'
                            : 'rounded-lg bg-accent px-3 py-1.5 text-xs font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:opacity-50 disabled:shadow-none'
                        }
                      >
                        {pendingFor(c.address) ? 'Confirming…' : c.verified ? 'Revoke' : 'Verify'}
                      </button>
                      {tx && tx.address.toLowerCase() === c.address.toLowerCase() && <TxFeedback tx={tx} compact />}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

function TxFeedback({ tx, compact = false }: { tx: TxState; compact?: boolean }) {
  const cls = compact ? 'text-right text-[11px]' : 'mt-3 text-xs'
  if (tx.status === 'pending' && tx.hash) {
    return (
      <p className={`${cls} text-muted`}>
        Waiting —{' '}
        <a href={`https://sepolia.etherscan.io/tx/${tx.hash}`} target="_blank" rel="noreferrer" className="text-accent underline">
          view
        </a>
      </p>
    )
  }
  if (tx.status === 'success') {
    return (
      <p className={`${cls} text-proven`}>
        {tx.action === 'verify' ? 'Verified' : 'Revoked'}
        {tx.hash && (
          <>
            {' — '}
            <a href={`https://sepolia.etherscan.io/tx/${tx.hash}`} target="_blank" rel="noreferrer" className="underline">
              view
            </a>
          </>
        )}
        . Indexer will pick it up shortly.
      </p>
    )
  }
  if (tx.status === 'error' && tx.error) {
    return <p className={`${cls} text-red-600`}>{tx.error}</p>
  }
  return null
}

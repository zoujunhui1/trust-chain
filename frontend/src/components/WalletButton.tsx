import { useWallet } from '../lib/wallet'
import { useRole, ROLE_LABEL } from '../lib/role'
import { shortAddress } from '../lib/format'

const ROLE_BADGE_CLASS: Record<string, string> = {
  admin: 'bg-accent-tint text-accent',
  charity: 'bg-proven-tint text-proven',
  donor: 'bg-locked-tint text-locked',
}

export default function WalletButton() {
  const { hasWallet, address, connecting, isWrongNetwork, connect, switchNetwork } = useWallet()
  const { role } = useRole()

  if (!hasWallet) {
    return (
      <a
        href="https://metamask.io/download/"
        target="_blank"
        rel="noreferrer"
        className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-muted transition-colors hover:text-ink"
      >
        Install MetaMask
      </a>
    )
  }

  if (address && isWrongNetwork) {
    return (
      <button
        type="button"
        onClick={switchNetwork}
        className="rounded-lg bg-released px-4 py-2 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98]"
      >
        Switch to Sepolia
      </button>
    )
  }

  if (address) {
    return (
      <span className="flex items-center gap-2 rounded-lg border border-border bg-white px-4 py-2 text-sm font-medium text-ink shadow-sm">
        {role && (
          <span className={`rounded-full px-2.5 py-1 text-xs font-medium ${ROLE_BADGE_CLASS[role]}`}>
            {ROLE_LABEL[role]}
          </span>
        )}
        {shortAddress(address)}
      </span>
    )
  }

  return (
    <button
      type="button"
      onClick={connect}
      disabled={connecting}
      className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white shadow-sm transition-all hover:opacity-90 hover:shadow-md active:scale-[0.98] disabled:opacity-50 disabled:shadow-none disabled:active:scale-100"
    >
      {connecting ? 'Connecting…' : 'Connect Wallet'}
    </button>
  )
}

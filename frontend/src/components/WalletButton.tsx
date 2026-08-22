import { useWallet } from '../lib/wallet'
import { shortAddress } from '../lib/format'

export default function WalletButton() {
  const { hasWallet, address, connecting, isWrongNetwork, connect, switchNetwork } = useWallet()

  if (!hasWallet) {
    return (
      <a
        href="https://metamask.io/download/"
        target="_blank"
        rel="noreferrer"
        className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-muted hover:text-ink"
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
        className="rounded-lg bg-released px-4 py-2 text-sm font-medium text-white hover:opacity-90"
      >
        Switch to Sepolia
      </button>
    )
  }

  if (address) {
    return (
      <span className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-ink">
        {shortAddress(address)}
      </span>
    )
  }

  return (
    <button
      type="button"
      onClick={connect}
      disabled={connecting}
      className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
    >
      {connecting ? 'Connecting…' : 'Connect Wallet'}
    </button>
  )
}

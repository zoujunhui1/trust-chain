import { useEffect, useState } from 'react'
import { connectUser, getCharity } from './api'
import { getOwner } from './registry'
import { useWallet } from './wallet'

export type Role = 'admin' | 'charity' | 'donor'

interface RoleState {
  // null = wallet not connected (or on the wrong network) — role is unknown.
  role: Role | null
  loading: boolean
}

// useRole() is called from more than one component (WalletButton, App's nav),
// each running its own effect. Report a given address+role to the backend at
// most once regardless of how many callers computed it, instead of posting
// once per caller on every connect.
let lastReported: string | null = null
function reportConnectOnce(address: string, role: Role) {
  const key = `${address.toLowerCase()}:${role}`
  if (key === lastReported) return
  lastReported = key
  // Fire-and-forget: this is a "who connected" log, not something that
  // should ever block or break the wallet UI if the API is unreachable.
  connectUser(address, role).catch(() => {
    lastReported = null // let a later render retry
  })
}

// There's no off-chain "users" table backing permissions — role for gating
// the UI is derived from the same on-chain/indexed signals the Admin and
// Create Campaign pages already read: the registry owner is the admin, a
// verified charity address is a charity, everyone else with a connected
// wallet is a donor. Once computed, it's also reported to POST
// /api/users/connect purely as a connect-time record (see backend/internal
// /store/users.go) — that table is never consulted for the check above.
export function useRole(): RoleState {
  const wallet = useWallet()
  const [role, setRole] = useState<Role | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!wallet.address || wallet.isWrongNetwork) {
      setRole(null)
      return
    }
    const address = wallet.address
    let cancelled = false
    setLoading(true)
    ;(async () => {
      try {
        const signer = await wallet.getSigner()
        const owner = await getOwner(signer)
        if (cancelled) return
        if (owner.toLowerCase() === address.toLowerCase()) {
          setRole('admin')
          reportConnectOnce(address, 'admin')
          return
        }
        const charity = await getCharity(address)
        if (cancelled) return
        const resolved: Role = charity.verified ? 'charity' : 'donor'
        setRole(resolved)
        reportConnectOnce(address, resolved)
      } catch {
        // Registry/API read failed (e.g. address never indexed as a charity)
        // — fall back to the least-privileged role rather than block the UI.
        if (!cancelled) setRole('donor')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [wallet.address, wallet.isWrongNetwork])

  return { role, loading }
}

export const ROLE_LABEL: Record<Role, string> = {
  admin: 'Admin',
  charity: 'Charity',
  donor: 'Donor',
}

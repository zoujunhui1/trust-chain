import { useEffect, useState } from 'react'
import { getCharity } from './api'
import { getOwner } from './registry'
import { useWallet } from './wallet'

export type Role = 'admin' | 'charity' | 'donor'

interface RoleState {
  // null = wallet not connected (or on the wrong network) — role is unknown.
  role: Role | null
  loading: boolean
}

// There's no off-chain "users" table — role is derived from the same
// on-chain/indexed signals the Admin and Create Campaign pages already read:
// the registry owner is the admin, a verified charity address is a charity,
// everyone else with a connected wallet is a donor.
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
          return
        }
        const charity = await getCharity(address)
        if (cancelled) return
        setRole(charity.verified ? 'charity' : 'donor')
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

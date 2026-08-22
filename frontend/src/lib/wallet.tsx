import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import { BrowserProvider } from 'ethers'
import type { Signer } from 'ethers'
import { CHAIN_ID } from './escrow'

interface WalletState {
  address: string | null
  chainId: bigint | null
  connecting: boolean
  error: string | null
}

interface WalletContextValue extends WalletState {
  hasWallet: boolean
  isWrongNetwork: boolean
  connect: () => Promise<void>
  switchNetwork: () => Promise<void>
  getSigner: () => Promise<Signer>
}

const WalletContext = createContext<WalletContextValue | null>(null)

export function WalletProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<WalletState>({
    address: null,
    chainId: null,
    connecting: false,
    error: null,
  })

  const hasWallet = typeof window !== 'undefined' && !!window.ethereum
  const provider = useMemo(() => (hasWallet ? new BrowserProvider(window.ethereum!) : null), [hasWallet])

  // Pick up an already-authorized account without prompting MetaMask.
  const refresh = useCallback(async () => {
    if (!provider) return
    const [accounts, network] = await Promise.all([provider.listAccounts(), provider.getNetwork()])
    setState((s) => ({ ...s, address: accounts[0]?.address ?? null, chainId: network.chainId }))
  }, [provider])

  useEffect(() => {
    refresh()
    if (!window.ethereum) return
    const onChange = () => refresh()
    window.ethereum.on('accountsChanged', onChange)
    window.ethereum.on('chainChanged', onChange)
    return () => {
      window.ethereum?.removeListener('accountsChanged', onChange)
      window.ethereum?.removeListener('chainChanged', onChange)
    }
  }, [refresh])

  const connect = useCallback(async () => {
    if (!provider) {
      setState((s) => ({ ...s, error: 'No wallet found — install MetaMask.' }))
      return
    }
    setState((s) => ({ ...s, connecting: true, error: null }))
    try {
      await provider.send('eth_requestAccounts', [])
      await refresh()
    } catch (err) {
      setState((s) => ({ ...s, error: err instanceof Error ? err.message : 'Failed to connect wallet.' }))
    } finally {
      setState((s) => ({ ...s, connecting: false }))
    }
  }, [provider, refresh])

  const switchNetwork = useCallback(async () => {
    if (!window.ethereum) return
    try {
      await window.ethereum.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId: `0x${CHAIN_ID.toString(16)}` }],
      })
      await refresh()
    } catch (err) {
      setState((s) => ({ ...s, error: err instanceof Error ? err.message : 'Failed to switch network.' }))
    }
  }, [refresh])

  const getSigner = useCallback(async () => {
    if (!provider) throw new Error('No wallet found — install MetaMask.')
    return provider.getSigner()
  }, [provider])

  const value: WalletContextValue = {
    ...state,
    hasWallet,
    isWrongNetwork: state.address !== null && state.chainId !== null && state.chainId !== CHAIN_ID,
    connect,
    switchNetwork,
    getSigner,
  }

  return <WalletContext.Provider value={value}>{children}</WalletContext.Provider>
}

export function useWallet(): WalletContextValue {
  const ctx = useContext(WalletContext)
  if (!ctx) throw new Error('useWallet must be used within a WalletProvider')
  return ctx
}

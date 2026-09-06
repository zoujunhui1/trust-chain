/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_ESCROW_ADDRESS: string
  readonly VITE_CHARITY_REGISTRY_ADDRESS: string
  readonly VITE_CHAIN_ID: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

// Minimal MetaMask injected-provider shape (EIP-1193 request + event emitter).
interface EthereumProvider {
  request(args: { method: string; params?: unknown[] | object }): Promise<unknown>
  on(event: string, handler: (...args: unknown[]) => void): void
  removeListener(event: string, handler: (...args: unknown[]) => void): void
}

interface Window {
  ethereum?: EthereumProvider
}

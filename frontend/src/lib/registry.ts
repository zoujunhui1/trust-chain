import { Contract } from 'ethers'
import type { Signer, TransactionResponse } from 'ethers'

export const CHARITY_REGISTRY_ADDRESS = import.meta.env.VITE_CHARITY_REGISTRY_ADDRESS

// Only the fragments the admin page needs: who the owner is, and the two
// owner-only actions. isVerified is covered by the backend's /api/charities
// instead (it's just the indexed mirror of the same state).
const REGISTRY_ABI = [
  'function owner() view returns (address)',
  'function verifyCharity(address charity) external',
  'function revokeCharity(address charity) external',
]

export async function getOwner(signer: Signer): Promise<string> {
  const contract = new Contract(CHARITY_REGISTRY_ADDRESS, REGISTRY_ABI, signer)
  return (await contract.owner()) as string
}

export function verifyCharity(signer: Signer, address: string): Promise<TransactionResponse> {
  const contract = new Contract(CHARITY_REGISTRY_ADDRESS, REGISTRY_ABI, signer)
  return contract.verifyCharity(address) as Promise<TransactionResponse>
}

export function revokeCharity(signer: Signer, address: string): Promise<TransactionResponse> {
  const contract = new Contract(CHARITY_REGISTRY_ADDRESS, REGISTRY_ABI, signer)
  return contract.revokeCharity(address) as Promise<TransactionResponse>
}

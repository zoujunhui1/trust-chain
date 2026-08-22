import { Contract, parseEther } from 'ethers'
import type { Signer, TransactionResponse } from 'ethers'

export const ESCROW_ADDRESS = import.meta.env.VITE_ESCROW_ADDRESS
export const CHAIN_ID = BigInt(import.meta.env.VITE_CHAIN_ID)

// Only the fragment the frontend needs: TrustChainEscrow.donate(campaignId) payable.
const ESCROW_ABI = ['function donate(uint256 campaignId) external payable']

export function donate(signer: Signer, campaignId: number | string, amountEth: string): Promise<TransactionResponse> {
  const contract = new Contract(ESCROW_ADDRESS, ESCROW_ABI, signer)
  return contract.donate(campaignId, { value: parseEther(amountEth) }) as Promise<TransactionResponse>
}

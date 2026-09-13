import { Contract, Interface, parseEther, ZeroHash } from 'ethers'
import type { Signer, TransactionReceipt, TransactionResponse } from 'ethers'

export const ESCROW_ADDRESS = import.meta.env.VITE_ESCROW_ADDRESS
export const CHAIN_ID = BigInt(import.meta.env.VITE_CHAIN_ID)

// Only the fragments the frontend needs: donate, createCampaign, and the
// event createCampaign emits (so we can read back the new campaign id).
const ESCROW_ABI = [
  'function donate(uint256 campaignId) external payable',
  'function createCampaign(uint256[] milestoneAmounts, bytes32 metadataHash) external returns (uint256 campaignId)',
  'event CampaignCreated(uint256 indexed campaignId, address indexed charity, uint256 goal, uint256 milestoneCount)',
]

const escrowInterface = new Interface(ESCROW_ABI)

export function donate(signer: Signer, campaignId: number | string, amountEth: string): Promise<TransactionResponse> {
  const contract = new Contract(ESCROW_ADDRESS, ESCROW_ABI, signer)
  return contract.donate(campaignId, { value: parseEther(amountEth) }) as Promise<TransactionResponse>
}

// metadataHash is left as the zero hash — title/description/theme live in
// the backend's campaign_metadata table (set via setCampaignMetadata) rather
// than in anything this hash could point to.
export function createCampaign(signer: Signer, milestoneAmountsEth: string[]): Promise<TransactionResponse> {
  const contract = new Contract(ESCROW_ADDRESS, ESCROW_ABI, signer)
  const amountsWei = milestoneAmountsEth.map((a) => parseEther(a))
  return contract.createCampaign(amountsWei, ZeroHash) as Promise<TransactionResponse>
}

// Pulls the new campaign id out of a createCampaign receipt's logs, if present.
export function parseCreatedCampaignId(receipt: TransactionReceipt | null): number | null {
  for (const log of receipt?.logs ?? []) {
    try {
      const parsed = escrowInterface.parseLog(log)
      if (parsed?.name === 'CampaignCreated') return Number(parsed.args.campaignId)
    } catch {
      // not one of our events, skip
    }
  }
  return null
}

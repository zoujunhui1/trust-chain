// SHA-256 of a File's bytes, as a 0x-prefixed 32-byte hex string — exactly
// the shape submitReceipt's bytes32 receiptHash parameter expects. Computed
// in the browser (Web Crypto) so the hash committed on-chain is provably the
// hash of the file the charity is about to upload, not something the
// backend could substitute.
export async function sha256File(file: File): Promise<string> {
  const buffer = await file.arrayBuffer()
  const digest = await crypto.subtle.digest('SHA-256', buffer)
  const hex = [...new Uint8Array(digest)].map((b) => b.toString(16).padStart(2, '0')).join('')
  return `0x${hex}`
}

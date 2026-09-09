export interface CampaignTheme {
  label: string
  emoji: string
  accent: string
  gradient: string
}

// Campaign title/description are still placeholders (no off-chain metadata —
// see the project's decision to skip that infrastructure), so there's no
// real category to theme off of. Cycling a fixed palette by campaign id still
// fixes the actual complaint (every card looking identical) without
// reopening the metadata question — swap this for a real category field if
// one ever gets added.
const THEMES: CampaignTheme[] = [
  { label: 'Environment', emoji: '🌳', accent: '#2f6b3f', gradient: 'linear-gradient(135deg, #eaf5ec 0%, #cfe8d6 100%)' },
  { label: 'Community', emoji: '🤝', accent: '#9a5b13', gradient: 'linear-gradient(135deg, #faf1e2 0%, #f2dcb8 100%)' },
  { label: 'Animal Welfare', emoji: '🐾', accent: '#8a4b8f', gradient: 'linear-gradient(135deg, #f4ecf6 0%, #e3cbe9 100%)' },
  { label: 'Accessibility', emoji: '♿', accent: '#1d6f8c', gradient: 'linear-gradient(135deg, #e6f4f8 0%, #bfe3ee 100%)' },
  { label: 'Education', emoji: '📚', accent: '#b0331a', gradient: 'linear-gradient(135deg, #fbe9e6 0%, #f3c4ba 100%)' },
]

export function campaignTheme(id: number): CampaignTheme {
  return THEMES[id % THEMES.length]
}

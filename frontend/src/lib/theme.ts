export interface CampaignTheme {
  key: string
  label: string
  emoji: string
  accent: string
  gradient: string
}

const THEMES: CampaignTheme[] = [
  {
    key: 'environment',
    label: 'Environment',
    emoji: '🌳',
    accent: '#2f6b3f',
    gradient: 'linear-gradient(135deg, #eaf5ec 0%, #cfe8d6 100%)',
  },
  {
    key: 'community',
    label: 'Community',
    emoji: '🤝',
    accent: '#9a5b13',
    gradient: 'linear-gradient(135deg, #faf1e2 0%, #f2dcb8 100%)',
  },
  {
    key: 'animal-welfare',
    label: 'Animal Welfare',
    emoji: '🐾',
    accent: '#8a4b8f',
    gradient: 'linear-gradient(135deg, #f4ecf6 0%, #e3cbe9 100%)',
  },
  {
    key: 'accessibility',
    label: 'Accessibility',
    emoji: '♿',
    accent: '#1d6f8c',
    gradient: 'linear-gradient(135deg, #e6f4f8 0%, #bfe3ee 100%)',
  },
  {
    key: 'education',
    label: 'Education',
    emoji: '📚',
    accent: '#b0331a',
    gradient: 'linear-gradient(135deg, #fbe9e6 0%, #f3c4ba 100%)',
  },
]

// Charities can pick a theme at creation time (stored in campaign_metadata,
// keyed by CampaignTheme.key). If they don't, or for campaigns created
// before this existed, fall back to a deterministic pick by id — this still
// fixes the original complaint (every card looking identical) without
// forcing a choice.
export function campaignTheme(campaign: { id: number; theme?: string | null }): CampaignTheme {
  const chosen = campaign.theme ? THEMES.find((t) => t.key === campaign.theme) : undefined
  return chosen ?? THEMES[campaign.id % THEMES.length]
}

// For a theme picker UI (e.g. Create Campaign's form).
export function listThemes(): CampaignTheme[] {
  return THEMES
}

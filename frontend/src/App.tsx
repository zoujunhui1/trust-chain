import { Routes, Route, Link } from 'react-router-dom'
import CampaignList from './pages/CampaignList'
import CampaignDetail from './pages/CampaignDetail'
import TransparencyDashboard from './pages/TransparencyDashboard'
import CreateCampaign from './pages/CreateCampaign'
import AdminCharities from './pages/AdminCharities'
import WalletButton from './components/WalletButton'
import { WalletProvider } from './lib/wallet'
import { useRole } from './lib/role'

function App() {
  return (
    <WalletProvider>
      <AppShell />
    </WalletProvider>
  )
}

// Split out so it renders inside WalletProvider — useRole needs wallet context.
function AppShell() {
  const { role } = useRole()

  return (
    <div className="min-h-svh">
      <div className="hero-glow" aria-hidden="true" />

      <nav className="sticky top-0 z-10 border-b border-border bg-white/85 backdrop-blur-sm">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-5 sm:px-10 lg:px-16">
          <span className="text-xl font-bold tracking-tight text-accent">TrustChain</span>
          <div className="flex items-center gap-8 text-sm font-medium text-muted">
            <Link to="/" className="transition-colors hover:text-ink">
              Campaigns
            </Link>
            <Link to="/transparency" className="transition-colors hover:text-ink">
              Transparency
            </Link>
            {/* Nav is role-gated: only a verified charity can create campaigns,
                only the registry owner can admin charities. Everyone else never
                sees these links (they'd fail on-chain anyway). */}
            {role === 'charity' && (
              <Link to="/create" className="transition-colors hover:text-ink">
                Create Campaign
              </Link>
            )}
            {role === 'admin' && (
              <Link to="/admin/charities" className="transition-colors hover:text-ink">
                Admin
              </Link>
            )}
            <WalletButton />
          </div>
        </div>
      </nav>

      <Routes>
        <Route path="/" element={<CampaignList />} />
        <Route path="/campaigns/:id" element={<CampaignDetail />} />
        <Route path="/transparency" element={<TransparencyDashboard />} />
        <Route path="/create" element={<CreateCampaign />} />
        <Route path="/admin/charities" element={<AdminCharities />} />
      </Routes>
    </div>
  )
}

export default App

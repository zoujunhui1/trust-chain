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

      <nav
        className="sticky top-0 z-10 backdrop-blur-sm"
        style={{
          background:
            'linear-gradient(rgba(255,255,255,0.88), rgba(255,255,255,0.88)), ' +
            'linear-gradient(90deg, rgba(59,130,246,0.14), rgba(168,85,247,0.14), rgba(236,72,153,0.14))',
        }}
      >
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-5 sm:px-10 lg:px-16">
          <span
            className="text-xl font-bold tracking-tight"
            style={{
              backgroundImage: 'linear-gradient(90deg, #1e3a5f, #7c3aed)',
              WebkitBackgroundClip: 'text',
              backgroundClip: 'text',
              color: 'transparent',
            }}
          >
            TrustChain
          </span>
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
        <div
          className="h-[3px] w-full"
          style={{ background: 'linear-gradient(90deg, #3b82f6, #a855f7, #ec4899, #22c55e)' }}
          aria-hidden="true"
        />
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

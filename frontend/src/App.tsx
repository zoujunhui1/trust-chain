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
      <nav className="flex items-center justify-between border-b border-border bg-white px-16 py-5">
        <span className="text-xl font-bold text-accent">TrustChain</span>
        <div className="flex items-center gap-8 text-sm font-medium text-muted">
          <Link to="/">Campaigns</Link>
          <Link to="/transparency">Transparency</Link>
          {/* Nav is role-gated: only a verified charity can create campaigns,
              only the registry owner can admin charities. Everyone else never
              sees these links (they'd fail on-chain anyway). */}
          {role === 'charity' && <Link to="/create">Create Campaign</Link>}
          {role === 'admin' && <Link to="/admin/charities">Admin</Link>}
          <WalletButton />
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

import { Routes, Route, Link } from 'react-router-dom'
import CampaignList from './pages/CampaignList'
import CampaignDetail from './pages/CampaignDetail'
import TransparencyDashboard from './pages/TransparencyDashboard'
import WalletButton from './components/WalletButton'
import { WalletProvider } from './lib/wallet'

function App() {
  return (
    <WalletProvider>
      <div className="min-h-svh">
        <nav className="flex items-center justify-between border-b border-border bg-white px-16 py-5">
          <span className="text-xl font-bold text-accent">TrustChain</span>
          <div className="flex items-center gap-8 text-sm font-medium text-muted">
            <Link to="/">Campaigns</Link>
            <Link to="/transparency">Transparency</Link>
            <WalletButton />
          </div>
        </nav>

        <Routes>
          <Route path="/" element={<CampaignList />} />
          <Route path="/campaigns/:id" element={<CampaignDetail />} />
          <Route path="/transparency" element={<TransparencyDashboard />} />
        </Routes>
      </div>
    </WalletProvider>
  )
}

export default App

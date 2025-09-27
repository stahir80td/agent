export default function Navbar() {
  return (
    <nav className="navbar">
      <div className="nav-content">
        <div className="nav-logo">
          <h2>📈 TradePro</h2>
        </div>
        <div className="nav-links">
          <a href="#features">Features</a>
          <a href="#pricing">Pricing</a>
          <a href="#about">About</a>
          <button className="nav-cta">Sign Up</button>
        </div>
      </div>
    </nav>
  )
}
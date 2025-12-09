// Layout component with header and navigation
// This provides the main structure for all pages

import { Outlet, Link, useLocation } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faNetworkWired } from '@fortawesome/free-solid-svg-icons';
import { ThemeToggle } from './ThemeToggle';
import Logo from './Logo';
import Footer from './Footer';

function Layout() {
  const location = useLocation();
  
  const isActive = (path: string) => {
    if (path === '/' || path === '/subnets') {
      return location.pathname === '/' || location.pathname === '/subnets';
    }
    return location.pathname.startsWith(path);
  };

  return (
    <div className="layout">
      <header className="app-header">
        <div className="header-content">
          <Link to="/" className="logo-link">
            <Logo variant="compact" size="medium" showText={true} />
          </Link>
          <nav className="main-nav">
            <Link 
              to="/" 
              className={`nav-link ${isActive('/') ? 'active' : ''}`}
            >
              <FontAwesomeIcon icon={faNetworkWired} />
              <span>Subnets</span>
            </Link>
            <Link 
              to="/subnets/create" 
              className={`nav-link nav-link-create ${isActive('/subnets/create') ? 'active' : ''}`}
            >
              <FontAwesomeIcon icon={faPlus} />
              <span>Create</span>
            </Link>
            <div className="nav-divider"></div>
            <ThemeToggle />
          </nav>
        </div>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
      <Footer />
    </div>
  );
}

export default Layout;

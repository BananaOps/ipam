// Layout component with header and navigation
// This provides the main structure for all pages

import { Outlet, Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faList } from '@fortawesome/free-solid-svg-icons';
import { ThemeToggle } from './ThemeToggle';
import Logo from './Logo';

function Layout() {
  return (
    <div className="layout">
      <header className="app-header">
        <div className="header-content">
          <Link to="/" className="logo-link">
            <Logo variant="compact" size="medium" showText={true} />
          </Link>
          <nav className="main-nav">
            <Link to="/" className="nav-link">
              <FontAwesomeIcon icon={faList} />
              <span>Subnets</span>
            </Link>
            <Link to="/subnets/create" className="nav-link">
              <FontAwesomeIcon icon={faPlus} />
              <span>Create</span>
            </Link>
            <ThemeToggle />
          </nav>
        </div>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;

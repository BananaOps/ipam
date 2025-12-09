// Logo component for IPAM by BananaOps
// Modern logo with network icon and gradient background

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faNetworkWired } from '@fortawesome/free-solid-svg-icons';
import './Logo.css';

interface LogoProps {
  variant?: 'full' | 'compact';
  size?: 'small' | 'medium' | 'large';
  showText?: boolean;
  className?: string;
}

function Logo({ 
  variant = 'compact', 
  size = 'medium', 
  showText = true,
  className = '' 
}: LogoProps) {
  return (
    <div className={`logo-container ${size} ${className}`}>
      <div className="logo-icon-box">
        <FontAwesomeIcon icon={faNetworkWired} className="logo-icon" />
      </div>
      {showText && (
        <div className="logo-text">
          <span className="logo-title">
            IPAM
          </span>
          <span className="logo-subtitle">
            by Banana<span className="ops-highlight">Ops</span>
          </span>
        </div>
      )}
    </div>
  );
}

export default Logo;

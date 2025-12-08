// Logo component for IPAM by BananaOps
// Displays the Cyber Minimal style logo with proper theming

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
  const logoSrc = variant === 'full' ? '/logo.svg' : '/logo-horizontal.svg';
  
  return (
    <div className={`logo-container ${size} ${className}`}>
      <img 
        src={logoSrc} 
        alt="IPAM by BananaOps Logo" 
        className="logo-image"
      />
      {showText && (
        <div className="logo-text">
          <span className="logo-title">IPAM</span>
          <span className="logo-subtitle">by BananaOps</span>
        </div>
      )}
    </div>
  );
}

export default Logo;

import { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faXmark, faHeart } from '@fortawesome/free-solid-svg-icons';
import { faGithub } from '@fortawesome/free-brands-svg-icons';
import './OpenSourceBanner.css';

export default function OpenSourceBanner() {
  const [isVisible, setIsVisible] = useState(() => {
    // Vérifier si l'utilisateur a déjà fermé la bannière
    return localStorage.getItem('openSourceBannerDismissed') !== 'true';
  });

  const handleDismiss = () => {
    localStorage.setItem('openSourceBannerDismissed', 'true');
    setIsVisible(false);
  };

  if (!isVisible) return null;

  return (
    <div className="opensource-banner">
      <div className="banner-content">
        <div className="banner-message">
          <FontAwesomeIcon icon={faGithub} className="banner-icon" />
          <p className="banner-text">
            <span>This is an open source project!</span>
            <span className="banner-text-secondary">
              Star us on GitHub
              <FontAwesomeIcon icon={faHeart} className="heart-icon" />
            </span>
          </p>
        </div>
        
        <div className="banner-actions">
          <a
            href="https://github.com/BananaOps/ipam"
            target="_blank"
            rel="noopener noreferrer"
            className="banner-button"
          >
            <FontAwesomeIcon icon={faGithub} />
            <span>View on GitHub</span>
          </a>
          
          <button
            onClick={handleDismiss}
            className="banner-dismiss"
            aria-label="Dismiss banner"
          >
            <FontAwesomeIcon icon={faXmark} />
          </button>
        </div>
      </div>
    </div>
  );
}

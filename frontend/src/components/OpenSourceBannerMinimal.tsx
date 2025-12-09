import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faStar } from '@fortawesome/free-solid-svg-icons';
import { faGithub } from '@fortawesome/free-brands-svg-icons';
import './OpenSourceBannerMinimal.css';

/**
 * Version minimaliste de la bannière open source
 * Affichée en permanence dans le footer ou en haut
 */
export default function OpenSourceBannerMinimal() {
  return (
    <div className="opensource-banner-minimal">
      <div className="banner-minimal-content">
        <span className="banner-minimal-label">
          <FontAwesomeIcon icon={faGithub} />
          Open Source Project
        </span>
        <span className="banner-minimal-separator">•</span>
        <a
          href="https://github.com/BananaOps/ipam"
          target="_blank"
          rel="noopener noreferrer"
          className="banner-minimal-link"
        >
          <FontAwesomeIcon icon={faStar} />
          Star on GitHub
        </a>
      </div>
    </div>
  );
}

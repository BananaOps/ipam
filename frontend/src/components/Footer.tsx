import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faStar, faBug, faComments, faCode } from '@fortawesome/free-solid-svg-icons';
import { faGithub } from '@fortawesome/free-brands-svg-icons';
import './Footer.css';

export default function Footer() {
  return (
    <footer className="app-footer">
      <div className="footer-content">
        <div className="footer-main">
          {/* Message principal avec banane animée */}
          <p className="footer-message">
            Made with{' '}
            <span className="banana-emoji">🍌</span>{' '}
            by the{' '}
            <a
              href="https://github.com/BananaOps"
              target="_blank"
              rel="noopener noreferrer"
              className="banana-link"
            >
              BananaOps
            </a>{' '}
            community
          </p>

          {/* Liens avec icônes */}
          <div className="footer-links">
            <a
              href="https://github.com/BananaOps/ipam"
              target="_blank"
              rel="noopener noreferrer"
              className="footer-link footer-link-github"
            >
              <FontAwesomeIcon icon={faGithub} />
              <span>View on GitHub</span>
            </a>

            <span className="footer-separator">•</span>

            <a
              href="https://github.com/BananaOps/ipam/stargazers"
              target="_blank"
              rel="noopener noreferrer"
              className="footer-link footer-link-star"
            >
              <FontAwesomeIcon icon={faStar} />
              <span>Star us</span>
            </a>

            <span className="footer-separator">•</span>

            <a
              href="https://github.com/BananaOps/ipam/issues"
              target="_blank"
              rel="noopener noreferrer"
              className="footer-link footer-link-bug"
            >
              <FontAwesomeIcon icon={faBug} />
              <span>Report a Bug</span>
            </a>

            <span className="footer-separator">•</span>

            <a
              href="https://github.com/BananaOps/ipam/discussions"
              target="_blank"
              rel="noopener noreferrer"
              className="footer-link footer-link-discuss"
            >
              <FontAwesomeIcon icon={faComments} />
              <span>Discussions</span>
            </a>
          </div>

          {/* Licence */}
          <div className="footer-license">
            <FontAwesomeIcon icon={faCode} />
            <a
              href="https://github.com/BananaOps/ipam/blob/main/LICENSE"
              target="_blank"
              rel="noopener noreferrer"
            >
              Apache 2.0 License
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}

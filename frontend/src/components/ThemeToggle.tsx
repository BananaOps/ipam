import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSun, faMoon, faCircleHalfStroke } from '@fortawesome/free-solid-svg-icons';
import { useTheme } from '../contexts/ThemeContext';
import type { Theme } from '../config/theme';

/**
 * ThemeToggle component for switching between dark, light, and auto themes
 */
export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  const handleToggle = () => {
    const themes: Theme[] = ['light', 'dark', 'auto'];
    const currentIndex = themes.indexOf(theme);
    const nextIndex = (currentIndex + 1) % themes.length;
    setTheme(themes[nextIndex]);
  };

  const getIcon = () => {
    switch (theme) {
      case 'light':
        return faSun;
      case 'dark':
        return faMoon;
      case 'auto':
        return faCircleHalfStroke;
    }
  };

  const getLabel = () => {
    switch (theme) {
      case 'light':
        return 'Light';
      case 'dark':
        return 'Dark';
      case 'auto':
        return 'Auto';
    }
  };

  return (
    <button
      className="theme-toggle"
      onClick={handleToggle}
      aria-label={`Current theme: ${getLabel()}. Click to change theme.`}
      title={`Switch theme (current: ${getLabel()})`}
    >
      <FontAwesomeIcon icon={getIcon()} className="theme-toggle-icon" />
      <span>{getLabel()}</span>
    </button>
  );
}

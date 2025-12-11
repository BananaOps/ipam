import './ErrorMessage.css';
import type { APIError } from '../types';

interface ErrorMessageProps {
  error: APIError | Error | string;
  onRetry?: () => void;
  onDismiss?: () => void;
  showDetails?: boolean;
}

/**
 * ErrorMessage component displays user-friendly error messages
 * Supports APIError objects, Error objects, and plain strings
 * Optionally shows retry and dismiss actions
 */
function ErrorMessage({ 
  error, 
  onRetry, 
  onDismiss, 
  showDetails = false 
}: ErrorMessageProps) {
  // Extract error information based on error type
  const getErrorInfo = () => {
    if (typeof error === 'string') {
      return {
        title: 'Erreur',
        message: error,
        suggestion: undefined,
        code: undefined,
        original: undefined,
      };
    }

    if ('code' in error && 'message' in error) {
      // APIError with enhanced details
      const details = error.details as any;
      return {
        title: details?.title || 'Erreur',
        message: error.message,
        suggestion: details?.suggestion,
        code: error.code,
        original: details?.original,
      };
    }

    // Standard Error object
    return {
      title: 'Erreur inattendue',
      message: error.message || 'Une erreur inattendue s\'est produite',
      suggestion: 'Veuillez réessayer ou contacter le support.',
      code: undefined,
      original: undefined,
    };
  };

  const { title, message, suggestion, code, original } = getErrorInfo();

  return (
    <div className="error-message-container">
      <div className="error-message-header">
        <span className="error-message-icon">⚠️</span>
        <div className="error-message-content">
          <h4 className="error-message-title">{title}</h4>
          <p className="error-message-text">{message}</p>
          {suggestion && (
            <p className="error-message-suggestion">💡 {suggestion}</p>
          )}
        </div>
      </div>

      {showDetails && (code || original) && (
        <div className="error-message-details">
          <strong>Détails techniques :</strong>
          <ul>
            {code && (
              <li><strong>Code :</strong> {code}</li>
            )}
            {original && (
              <li><strong>Message original :</strong> {original}</li>
            )}
          </ul>
        </div>
      )}

      {(onRetry || onDismiss) && (
        <div className="error-message-actions">
          {onRetry && (
            <button
              className="btn-primary error-message-retry"
              onClick={onRetry}
            >
              Réessayer
            </button>
          )}
          {onDismiss && (
            <button
              className="btn-secondary error-message-dismiss"
              onClick={onDismiss}
            >
              Fermer
            </button>
          )}
        </div>
      )}
    </div>
  );
}

export default ErrorMessage;

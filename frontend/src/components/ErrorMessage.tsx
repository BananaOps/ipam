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
        message: error,
        code: undefined,
        details: undefined,
      };
    }

    if ('code' in error && 'message' in error) {
      // APIError
      return {
        message: error.message,
        code: error.code,
        details: error.details,
      };
    }

    // Standard Error object
    return {
      message: error.message || 'An unexpected error occurred',
      code: undefined,
      details: undefined,
    };
  };

  const { message, code, details } = getErrorInfo();

  return (
    <div className="error-message-container">
      <div className="error-message-header">
        <span className="error-message-icon">⚠️</span>
        <div className="error-message-content">
          <p className="error-message-text">{message}</p>
          {code && (
            <p className="error-message-code">Error Code: {code}</p>
          )}
        </div>
      </div>

      {showDetails && details && Object.keys(details).length > 0 && (
        <div className="error-message-details">
          <strong>Additional Details:</strong>
          <ul>
            {Object.entries(details).map(([key, value]) => (
              <li key={key}>
                <strong>{key}:</strong> {value}
              </li>
            ))}
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
              Retry
            </button>
          )}
          {onDismiss && (
            <button
              className="btn-secondary error-message-dismiss"
              onClick={onDismiss}
            >
              Dismiss
            </button>
          )}
        </div>
      )}
    </div>
  );
}

export default ErrorMessage;

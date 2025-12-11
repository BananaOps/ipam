// Error message translations for better user experience

interface ErrorTranslation {
  title: string;
  message: string;
  suggestion?: string;
}

/**
 * Translate technical error messages to user-friendly messages
 */
export function translateError(error: any): ErrorTranslation {
  // If error has a response with data, try to extract meaningful info
  if (error.response?.data) {
    const errorData = error.response.data;
    
    // Handle structured API errors
    if (errorData.error) {
      return translateAPIError(errorData.error, error.response.status);
    }
    
    // Handle direct error messages
    if (typeof errorData === 'string') {
      return translateErrorMessage(errorData, error.response.status);
    }
  }
  
  // Handle network and other errors
  if (error.code) {
    return translateNetworkError(error);
  }
  
  // Fallback for unknown errors
  return {
    title: 'Erreur inattendue',
    message: error.message || 'Une erreur inattendue s\'est produite.',
    suggestion: 'Veuillez réessayer ou contacter le support si le problème persiste.'
  };
}

/**
 * Translate structured API errors
 */
function translateAPIError(apiError: any, status: number): ErrorTranslation {
  const code = apiError.code || '';
  const message = apiError.message || '';
  
  // CIDR/Subnet specific errors
  if (message.includes('already exists') || message.includes('duplicate') || code === 'CONFLICT') {
    return {
      title: 'Subnet déjà utilisé',
      message: 'Ce CIDR est déjà utilisé par un autre subnet.',
      suggestion: 'Choisissez une plage d\'adresses IP différente.'
    };
  }
  
  if (message.includes('invalid CIDR') || message.includes('INVALID_CIDR')) {
    return {
      title: 'Format CIDR invalide',
      message: 'Le format du CIDR n\'est pas valide.',
      suggestion: 'Utilisez le format correct : 192.168.1.0/24'
    };
  }
  
  if (message.includes('overlaps') || message.includes('overlap')) {
    return {
      title: 'Conflit de plages IP',
      message: 'Ce subnet chevauche avec un subnet existant.',
      suggestion: 'Choisissez une plage d\'adresses qui ne chevauche pas avec les subnets existants.'
    };
  }
  
  // Validation errors
  if (message.includes('required') || message.includes('missing')) {
    return {
      title: 'Champs obligatoires manquants',
      message: 'Certains champs obligatoires ne sont pas renseignés.',
      suggestion: 'Vérifiez que tous les champs marqués comme obligatoires sont remplis.'
    };
  }
  
  if (message.includes('too long') || message.includes('length')) {
    return {
      title: 'Données trop longues',
      message: 'Certaines données dépassent la longueur maximale autorisée.',
      suggestion: 'Raccourcissez le nom ou la description.'
    };
  }
  
  // Cloud provider errors
  if (message.includes('cloud provider') || message.includes('region')) {
    return {
      title: 'Erreur de configuration cloud',
      message: 'Les informations du fournisseur cloud sont incorrectes.',
      suggestion: 'Vérifiez le fournisseur, la région et l\'ID de compte.'
    };
  }
  
  // Database errors
  if (message.includes('database') || message.includes('DB_ERROR')) {
    return {
      title: 'Erreur de base de données',
      message: 'Une erreur s\'est produite lors de l\'accès aux données.',
      suggestion: 'Veuillez réessayer dans quelques instants.'
    };
  }
  
  // Permission errors
  if (status === 403 || message.includes('permission') || message.includes('forbidden')) {
    return {
      title: 'Accès refusé',
      message: 'Vous n\'avez pas les permissions nécessaires pour cette action.',
      suggestion: 'Contactez votre administrateur pour obtenir les droits appropriés.'
    };
  }
  
  // Not found errors
  if (status === 404 || message.includes('not found')) {
    return {
      title: 'Ressource introuvable',
      message: 'Le subnet demandé n\'existe pas ou a été supprimé.',
      suggestion: 'Vérifiez que le subnet existe toujours ou rafraîchissez la page.'
    };
  }
  
  // Server errors
  if (status >= 500) {
    return {
      title: 'Erreur serveur',
      message: 'Le serveur rencontre des difficultés temporaires.',
      suggestion: 'Veuillez réessayer dans quelques minutes.'
    };
  }
  
  // Default API error
  return {
    title: 'Erreur de l\'API',
    message: message || 'Une erreur s\'est produite lors de la communication avec le serveur.',
    suggestion: 'Veuillez réessayer ou vérifier votre connexion.'
  };
}

/**
 * Translate simple error messages
 */
function translateErrorMessage(message: string, status: number): ErrorTranslation {
  return translateAPIError({ message }, status);
}

/**
 * Translate network errors
 */
function translateNetworkError(error: any): ErrorTranslation {
  switch (error.code) {
    case 'ECONNABORTED':
    case 'NETWORK_TIMEOUT':
      return {
        title: 'Délai d\'attente dépassé',
        message: 'La requête a pris trop de temps à répondre.',
        suggestion: 'Vérifiez votre connexion internet et réessayez.'
      };
      
    case 'ECONNREFUSED':
    case 'NETWORK_ERROR':
      return {
        title: 'Connexion impossible',
        message: 'Impossible de se connecter au serveur.',
        suggestion: 'Vérifiez que le serveur est démarré et accessible.'
      };
      
    case 'ENOTFOUND':
      return {
        title: 'Serveur introuvable',
        message: 'L\'adresse du serveur est incorrecte ou inaccessible.',
        suggestion: 'Vérifiez la configuration de l\'URL du serveur.'
      };
      
    default:
      return {
        title: 'Erreur de réseau',
        message: 'Une erreur de réseau s\'est produite.',
        suggestion: 'Vérifiez votre connexion internet et réessayez.'
      };
  }
}

/**
 * Get a short error message for toasts
 */
export function getShortErrorMessage(error: any): string {
  const translation = translateError(error);
  return translation.title;
}

/**
 * Get a detailed error message for error components
 */
export function getDetailedErrorMessage(error: any): { title: string; message: string; suggestion?: string } {
  return translateError(error);
}

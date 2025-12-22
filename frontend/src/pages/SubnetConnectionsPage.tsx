// SubnetConnectionsPage - Page for managing subnet connections
import { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faLink } from '@fortawesome/free-solid-svg-icons';
import { apiClient } from '../services/api';
import { 
  Subnet, 
  SubnetConnection, 
  CreateConnectionRequest, 
  UpdateConnectionRequest,
  APIError 
} from '../types';
import SubnetConnectionManager from '../components/SubnetConnectionManager';
import ErrorMessage from '../components/ErrorMessage';
import { useToast } from '../contexts/ToastContext';
import './SubnetConnectionsPage.css';

function SubnetConnectionsPage() {
  const [subnets, setSubnets] = useState<Subnet[]>([]);
  const [connections, setConnections] = useState<SubnetConnection[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<APIError | Error | null>(null);
  const { showError: showToastError, showSuccess } = useToast();

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [subnetsResponse, connectionsResponse] = await Promise.all([
        apiClient.listSubnets(),
        apiClient.listConnections()
      ]);
      
      setSubnets(subnetsResponse.subnets);
      setConnections(connectionsResponse.connections);
    } catch (err: any) {
      setError(err);
      showToastError(err.message || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateConnection = async (data: CreateConnectionRequest) => {
    const newConnection = await apiClient.createConnection(data);
    setConnections(prev => [...prev, newConnection]);
  };

  const handleUpdateConnection = async (id: string, data: UpdateConnectionRequest) => {
    const updatedConnection = await apiClient.updateConnection(id, data);
    setConnections(prev => prev.map(c => c.id === id ? updatedConnection : c));
  };

  const handleDeleteConnection = async (id: string) => {
    await apiClient.deleteConnection(id);
    setConnections(prev => prev.filter(c => c.id !== id));
  };

  if (loading) {
    return (
      <div className="subnet-connections-loading">
        <div className="loading-spinner"></div>
        <p>Loading connections...</p>
      </div>
    );
  }

  if (error) {
    return (
      <ErrorMessage
        error={error}
        onRetry={loadData}
        onDismiss={() => setError(null)}
        showDetails={true}
      />
    );
  }

  return (
    <div className="subnet-connections-page">
      <div className="page-header">
        <div className="header-title">
          <FontAwesomeIcon icon={faLink} />
          <h1>Connexions Réseau</h1>
        </div>
        <p className="header-subtitle">
          Gérez les connexions entre vos sous-réseaux (VPN, Peering, NAT Gateway, etc.)
        </p>
      </div>

      <SubnetConnectionManager
        subnets={subnets}
        connections={connections}
        onCreateConnection={handleCreateConnection}
        onUpdateConnection={handleUpdateConnection}
        onDeleteConnection={handleDeleteConnection}
      />
    </div>
  );
}

export default SubnetConnectionsPage;

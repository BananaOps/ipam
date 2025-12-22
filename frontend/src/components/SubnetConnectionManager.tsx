// SubnetConnectionManager - Component for managing connections between subnets
import { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faPlus, 
  faEdit, 
  faTrash, 
  faLink,
  faTimes,
  faCheck,
  faExclamationTriangle
} from '@fortawesome/free-solid-svg-icons';
import { 
  Subnet, 
  SubnetConnection, 
  ConnectionType, 
  ConnectionStatus,
  CreateConnectionRequest,
  UpdateConnectionRequest
} from '../types';
import { useToast } from '../contexts/ToastContext';
import './SubnetConnectionManager.css';

interface SubnetConnectionManagerProps {
  subnets: Subnet[];
  connections: SubnetConnection[];
  onCreateConnection: (connection: CreateConnectionRequest) => Promise<void>;
  onUpdateConnection: (id: string, connection: UpdateConnectionRequest) => Promise<void>;
  onDeleteConnection: (id: string) => Promise<void>;
  selectedSubnet?: Subnet;
}

const CONNECTION_TYPE_LABELS: Record<ConnectionType, string> = {
  [ConnectionType.VPN_SITE_TO_SITE]: 'VPN Site-à-Site',
  [ConnectionType.OPENVPN_CLIENT]: 'Client OpenVPN',
  [ConnectionType.NAT_GATEWAY]: 'NAT Gateway',
  [ConnectionType.INTERNET_GATEWAY]: 'Internet Gateway',
  [ConnectionType.PEERING]: 'Peering',
  [ConnectionType.TRANSIT_GATEWAY]: 'Transit Gateway',
  [ConnectionType.DIRECT_CONNECT]: 'Direct Connect',
  [ConnectionType.EXPRESSROUTE]: 'ExpressRoute',
  [ConnectionType.CLOUD_INTERCONNECT]: 'Cloud Interconnect',
  [ConnectionType.LOAD_BALANCER]: 'Load Balancer',
  [ConnectionType.FIREWALL]: 'Firewall',
  [ConnectionType.CUSTOM]: 'Personnalisé'
};

const CONNECTION_STATUS_LABELS: Record<ConnectionStatus, string> = {
  [ConnectionStatus.ACTIVE]: 'Actif',
  [ConnectionStatus.INACTIVE]: 'Inactif',
  [ConnectionStatus.PENDING]: 'En attente',
  [ConnectionStatus.ERROR]: 'Erreur'
};

function SubnetConnectionManager({
  subnets,
  connections,
  onCreateConnection,
  onUpdateConnection,
  onDeleteConnection,
  selectedSubnet
}: SubnetConnectionManagerProps) {
  const [isCreating, setIsCreating] = useState(false);
  const [editingConnection, setEditingConnection] = useState<SubnetConnection | null>(null);
  const [formData, setFormData] = useState<CreateConnectionRequest>({
    sourceSubnetId: selectedSubnet?.id || '',
    targetSubnetId: '',
    connectionType: ConnectionType.VPN_SITE_TO_SITE,
    name: '',
    description: '',
    bandwidth: '',
    latency: undefined,
    cost: undefined,
  });
  const { showError, showSuccess } = useToast();

  useEffect(() => {
    if (selectedSubnet) {
      setFormData(prev => ({ ...prev, sourceSubnetId: selectedSubnet.id }));
    }
  }, [selectedSubnet]);

  const handleCreateConnection = async () => {
    try {
      if (!formData.sourceSubnetId || !formData.targetSubnetId || !formData.name) {
        showError('Veuillez remplir tous les champs obligatoires');
        return;
      }

      if (formData.sourceSubnetId === formData.targetSubnetId) {
        showError('Un sous-réseau ne peut pas se connecter à lui-même');
        return;
      }

      await onCreateConnection(formData);
      showSuccess('Connexion créée avec succès');
      setIsCreating(false);
      resetForm();
    } catch (error: any) {
      showError(error.message || 'Erreur lors de la création de la connexion');
    }
  };

  const handleUpdateConnection = async () => {
    if (!editingConnection) return;

    try {
      const updateData: UpdateConnectionRequest = {
        name: formData.name,
        description: formData.description,
        connectionType: formData.connectionType,
        bandwidth: formData.bandwidth,
        latency: formData.latency,
        cost: formData.cost,
      };

      await onUpdateConnection(editingConnection.id, updateData);
      showSuccess('Connexion mise à jour avec succès');
      setEditingConnection(null);
      resetForm();
    } catch (error: any) {
      showError(error.message || 'Erreur lors de la mise à jour de la connexion');
    }
  };

  const handleDeleteConnection = async (connection: SubnetConnection) => {
    if (!window.confirm(`Êtes-vous sûr de vouloir supprimer la connexion "${connection.name}" ?`)) {
      return;
    }

    try {
      await onDeleteConnection(connection.id);
      showSuccess('Connexion supprimée avec succès');
    } catch (error: any) {
      showError(error.message || 'Erreur lors de la suppression de la connexion');
    }
  };

  const resetForm = () => {
    setFormData({
      sourceSubnetId: selectedSubnet?.id || '',
      targetSubnetId: '',
      connectionType: ConnectionType.VPN_SITE_TO_SITE,
      name: '',
      description: '',
      bandwidth: '',
      latency: undefined,
      cost: undefined,
    });
  };

  const startEditing = (connection: SubnetConnection) => {
    setEditingConnection(connection);
    setFormData({
      sourceSubnetId: connection.sourceSubnetId,
      targetSubnetId: connection.targetSubnetId,
      connectionType: connection.connectionType,
      name: connection.name,
      description: connection.description || '',
      bandwidth: connection.bandwidth || '',
      latency: connection.latency,
      cost: connection.cost,
    });
  };

  const cancelEditing = () => {
    setEditingConnection(null);
    setIsCreating(false);
    resetForm();
  };

  const getSubnetName = (subnetId: string): string => {
    if (subnetId === 'internet') {
      return '🌐 Internet';
    }
    const subnet = subnets.find(s => s.id === subnetId);
    return subnet ? `${subnet.name} (${subnet.cidr})` : 'Destination inconnue';
  };

  const getStatusIcon = (status: ConnectionStatus) => {
    switch (status) {
      case ConnectionStatus.ACTIVE:
        return <FontAwesomeIcon icon={faCheck} className="status-icon active" />;
      case ConnectionStatus.INACTIVE:
        return <FontAwesomeIcon icon={faTimes} className="status-icon inactive" />;
      case ConnectionStatus.PENDING:
        return <FontAwesomeIcon icon={faExclamationTriangle} className="status-icon pending" />;
      case ConnectionStatus.ERROR:
        return <FontAwesomeIcon icon={faExclamationTriangle} className="status-icon error" />;
      default:
        return null;
    }
  };

  // Filter connections for selected subnet
  const relevantConnections = selectedSubnet 
    ? connections.filter(c => c.sourceSubnetId === selectedSubnet.id || c.targetSubnetId === selectedSubnet.id)
    : connections;

  return (
    <div className="subnet-connection-manager">
      <div className="connection-header">
        <h3>
          <FontAwesomeIcon icon={faLink} />
          Connexions Réseau
          {selectedSubnet && (
            <span className="selected-subnet">
              - {selectedSubnet.name}
            </span>
          )}
        </h3>
        <button
          onClick={() => setIsCreating(true)}
          className="create-connection-btn"
          disabled={isCreating || !!editingConnection}
        >
          <FontAwesomeIcon icon={faPlus} />
          Nouvelle Connexion
        </button>
      </div>

      {/* Create/Edit Form */}
      {(isCreating || editingConnection) && (
        <div className="connection-form">
          <div className="form-header">
            <h4>{editingConnection ? 'Modifier la Connexion' : 'Nouvelle Connexion'}</h4>
            <button onClick={cancelEditing} className="cancel-btn">
              <FontAwesomeIcon icon={faTimes} />
            </button>
          </div>

          <div className="form-grid">
            <div className="form-group">
              <label>Sous-réseau source *</label>
              <select
                value={formData.sourceSubnetId}
                onChange={(e) => setFormData(prev => ({ ...prev, sourceSubnetId: e.target.value }))}
                disabled={!!selectedSubnet}
              >
                <option value="">Sélectionner un sous-réseau</option>
                {subnets.map(subnet => (
                  <option key={subnet.id} value={subnet.id}>
                    {subnet.name} ({subnet.cidr})
                  </option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label>Sous-réseau cible *</label>
              <select
                value={formData.targetSubnetId}
                onChange={(e) => setFormData(prev => ({ ...prev, targetSubnetId: e.target.value }))}
              >
                <option value="">Sélectionner une destination</option>
                <option value="internet">🌐 Internet</option>
                {subnets
                  .filter(subnet => subnet.id !== formData.sourceSubnetId)
                  .map(subnet => (
                    <option key={subnet.id} value={subnet.id}>
                      {subnet.name} ({subnet.cidr})
                    </option>
                  ))}
              </select>
            </div>

            <div className="form-group">
              <label>Type de connexion *</label>
              <select
                value={formData.connectionType}
                onChange={(e) => setFormData(prev => ({ ...prev, connectionType: e.target.value as ConnectionType }))}
              >
                {Object.entries(CONNECTION_TYPE_LABELS).map(([value, label]) => (
                  <option key={value} value={value}>{label}</option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label>Nom de la connexion *</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                placeholder="ex: VPN Paris-Londres"
              />
            </div>

            <div className="form-group full-width">
              <label>Description</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                placeholder="Description de la connexion..."
                rows={2}
              />
            </div>

            <div className="form-group">
              <label>Bande passante</label>
              <input
                type="text"
                value={formData.bandwidth}
                onChange={(e) => setFormData(prev => ({ ...prev, bandwidth: e.target.value }))}
                placeholder="ex: 1Gbps, 100Mbps"
              />
            </div>

            <div className="form-group">
              <label>Latence (ms)</label>
              <input
                type="number"
                value={formData.latency || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, latency: e.target.value ? parseInt(e.target.value) : undefined }))}
                placeholder="ex: 50"
              />
            </div>

            <div className="form-group">
              <label>Coût mensuel (€)</label>
              <input
                type="number"
                step="0.01"
                value={formData.cost || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, cost: e.target.value ? parseFloat(e.target.value) : undefined }))}
                placeholder="ex: 99.99"
              />
            </div>
          </div>

          <div className="form-actions">
            <button onClick={cancelEditing} className="cancel-button">
              Annuler
            </button>
            <button
              onClick={editingConnection ? handleUpdateConnection : handleCreateConnection}
              className="save-button"
            >
              {editingConnection ? 'Mettre à jour' : 'Créer'}
            </button>
          </div>
        </div>
      )}

      {/* Connections List */}
      <div className="connections-list">
        {relevantConnections.length === 0 ? (
          <div className="no-connections">
            <FontAwesomeIcon icon={faLink} className="empty-icon" />
            <p>Aucune connexion trouvée</p>
            {selectedSubnet && (
              <p className="empty-subtitle">
                Créez une connexion pour lier ce sous-réseau à d'autres
              </p>
            )}
          </div>
        ) : (
          <div className="connections-grid">
            {relevantConnections.map(connection => (
              <div key={connection.id} className="connection-card">
                <div className="connection-header-card">
                  <div className="connection-title">
                    <h4>{connection.name}</h4>
                    {getStatusIcon(connection.status)}
                  </div>
                  <div className="connection-actions">
                    <button
                      onClick={() => startEditing(connection)}
                      className="edit-btn"
                      title="Modifier"
                    >
                      <FontAwesomeIcon icon={faEdit} />
                    </button>
                    <button
                      onClick={() => handleDeleteConnection(connection)}
                      className="delete-btn"
                      title="Supprimer"
                    >
                      <FontAwesomeIcon icon={faTrash} />
                    </button>
                  </div>
                </div>

                <div className="connection-details">
                  <div className="connection-route">
                    <span className="source">{getSubnetName(connection.sourceSubnetId)}</span>
                    <FontAwesomeIcon icon={faLink} className="link-icon" />
                    <span className="target">{getSubnetName(connection.targetSubnetId)}</span>
                  </div>

                  <div className="connection-type">
                    {CONNECTION_TYPE_LABELS[connection.connectionType]}
                  </div>

                  {connection.description && (
                    <div className="connection-description">
                      {connection.description}
                    </div>
                  )}

                  <div className="connection-metadata">
                    {connection.bandwidth && (
                      <span className="metadata-item">
                        <strong>Bande passante:</strong> {connection.bandwidth}
                      </span>
                    )}
                    {connection.latency && (
                      <span className="metadata-item">
                        <strong>Latence:</strong> {connection.latency}ms
                      </span>
                    )}
                    {connection.cost && (
                      <span className="metadata-item">
                        <strong>Coût:</strong> {connection.cost}€/mois
                      </span>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default SubnetConnectionManager;

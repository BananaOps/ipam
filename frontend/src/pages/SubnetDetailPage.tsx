import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import SubnetDetail from '../components/SubnetDetail';
import ConfirmDialog from '../components/ConfirmDialog';
import Toast from '../components/Toast';
import { apiClient } from '../services/api';
import type { Subnet, APIError } from '../types';

function SubnetDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [subnet, setSubnet] = useState<Subnet | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<APIError | null>(null);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' | 'info' | 'warning' } | null>(null);

  useEffect(() => {
    const fetchSubnet = async () => {
      if (!id) {
        setError({
          code: 'INVALID_ID',
          message: 'No subnet ID provided',
          timestamp: Date.now(),
        });
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError(null);
        const data = await apiClient.getSubnet(id);
        setSubnet(data);
      } catch (err) {
        setError(err as APIError);
      } finally {
        setLoading(false);
      }
    };

    fetchSubnet();
  }, [id]);

  const handleDeleteClick = () => {
    setShowDeleteDialog(true);
  };

  const handleDeleteConfirm = async () => {
    if (!id) return;

    try {
      setIsDeleting(true);
      await apiClient.deleteSubnet(id);
      
      // Show success message
      setToast({
        message: `Subnet "${subnet?.name || id}" has been successfully deleted.`,
        type: 'success',
      });

      // Navigate back to list after a short delay
      setTimeout(() => {
        navigate('/subnets');
      }, 1500);
    } catch (err) {
      const apiError = err as APIError;
      setToast({
        message: `Failed to delete subnet: ${apiError.message}`,
        type: 'error',
      });
      setIsDeleting(false);
      setShowDeleteDialog(false);
    }
  };

  const handleDeleteCancel = () => {
    setShowDeleteDialog(false);
  };

  const handleToastClose = () => {
    setToast(null);
  };

  if (loading) {
    return (
      <div className="subnet-detail-page">
        <div className="loading-state">
          <p>Loading subnet details...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="subnet-detail-page">
        <div className="error-state">
          <h2>Error Loading Subnet</h2>
          <p className="error-message">{error.message}</p>
          <p className="error-code">Error Code: {error.code}</p>
          <div className="error-actions">
            <button onClick={() => navigate('/subnets')} className="btn-secondary">
              Back to List
            </button>
            <button onClick={() => window.location.reload()} className="btn-primary">
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!subnet) {
    return (
      <div className="subnet-detail-page">
        <div className="empty-state">
          <h2>Subnet Not Found</h2>
          <p>The requested subnet could not be found.</p>
          <button onClick={() => navigate('/subnets')} className="btn-primary">
            Back to List
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="subnet-detail-page">
      <div className="page-header">
        <div className="header-actions">
          <Link to="/subnets" className="btn-back">
            ← Back to List
          </Link>
          <div className="action-buttons">
            <Link to={`/subnets/${id}/edit`} className="btn-secondary">
              Edit Subnet
            </Link>
            <button
              onClick={handleDeleteClick}
              className="btn-danger"
              disabled={isDeleting}
            >
              {isDeleting ? 'Deleting...' : 'Delete Subnet'}
            </button>
          </div>
        </div>
      </div>
      
      <SubnetDetail subnet={subnet} />

      <ConfirmDialog
        isOpen={showDeleteDialog}
        title="Delete Subnet"
        message={
          <>
            <p>Are you sure you want to delete the subnet <strong>{subnet?.name}</strong> ({subnet?.cidr})?</p>
            <p style={{ marginTop: '1rem', color: 'var(--color-error)' }}>
              This action cannot be undone.
            </p>
          </>
        }
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleDeleteConfirm}
        onCancel={handleDeleteCancel}
        variant="danger"
      />

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          isVisible={true}
          onClose={handleToastClose}
        />
      )}
    </div>
  );
}

export default SubnetDetailPage;

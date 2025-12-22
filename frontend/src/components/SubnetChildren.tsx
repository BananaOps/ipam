import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faNetworkWired, faChevronDown, faChevronRight } from '@fortawesome/free-solid-svg-icons';
import { Subnet, CloudResourceType } from '../types';
import CloudProviderIcon from './CloudProviderIcon';
import { useToast } from '../contexts/ToastContext';
import './SubnetChildren.css';

interface SubnetChildrenProps {
  parentSubnet: Subnet;
  expanded?: boolean;
}

const SubnetChildren: React.FC<SubnetChildrenProps> = ({ parentSubnet, expanded = false }) => {
  const [children, setChildren] = useState<Subnet[]>([]);
  const [loading, setLoading] = useState(false);
  const [isExpanded, setIsExpanded] = useState(expanded);
  const { showError } = useToast();

  useEffect(() => {
    if (isExpanded && children.length === 0) {
      loadChildren();
    }
  }, [isExpanded, parentSubnet.id]);

  const loadChildren = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/api/v1/subnets/${parentSubnet.id}/children`);
      if (response.ok) {
        const data = await response.json();
        setChildren(data.children || []);
      } else {
        showError('Failed to load child subnets');
      }
    } catch (error) {
      showError('Failed to load child subnets');
      console.error('Error loading children:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };

  // Only show if parent is a VPC or has potential children
  const isVPC = parentSubnet.cloudInfo?.resourceType === CloudResourceType.VPC;
  if (!isVPC && !parentSubnet.children?.length) {
    return null;
  }

  return (
    <div className="subnet-children">
      <button 
        className="children-toggle"
        onClick={toggleExpanded}
        disabled={loading}
      >
        <FontAwesomeIcon 
          icon={isExpanded ? faChevronDown : faChevronRight} 
          className={loading ? 'loading' : ''}
        />
        <span>
          {isVPC ? 'VPC Subnets' : 'Child Subnets'}
          {children.length > 0 && ` (${children.length})`}
        </span>
      </button>

      {isExpanded && (
        <div className="children-container">
          {loading ? (
            <div className="children-loading">
              <FontAwesomeIcon icon={faNetworkWired} spin />
              <span>Loading child subnets...</span>
            </div>
          ) : children.length > 0 ? (
            <div className="children-list">
              {children.map((child) => (
                <div key={child.id} className="child-subnet">
                  <div className="child-subnet-header">
                    <Link to={`/subnets/${child.id}`} className="child-subnet-link">
                      <span className="child-cidr">{child.cidr}</span>
                      <span className="child-name">{child.name}</span>
                    </Link>
                    
                    {child.cloudInfo && (
                      <div className="child-cloud-info">
                        <CloudProviderIcon
                          provider={child.cloudInfo.provider}
                          size="sm"
                          title={child.cloudInfo.provider.toUpperCase()}
                        />
                        {child.cloudInfo.resourceType && (
                          <span className={`resource-type-badge small ${child.cloudInfo.resourceType}`}>
                            {child.cloudInfo.resourceType.toUpperCase()}
                          </span>
                        )}
                      </div>
                    )}
                  </div>

                  <div className="child-subnet-details">
                    <span className="child-location">{child.location}</span>
                    {child.utilization && (
                      <div className="child-utilization">
                        <div className="utilization-bar small">
                          <div
                            className={`utilization-fill ${
                              child.utilization.utilizationPercent >= 80 ? 'high' : ''
                            }`}
                            style={{ width: `${child.utilization.utilizationPercent}%` }}
                          ></div>
                        </div>
                        <span className="utilization-percent">
                          {child.utilization.utilizationPercent.toFixed(1)}%
                        </span>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="no-children">
              <FontAwesomeIcon icon={faNetworkWired} />
              <span>No child subnets found</span>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default SubnetChildren;

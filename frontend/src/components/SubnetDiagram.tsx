// SubnetDiagram - Interactive network diagram component with zoom and pan
import { useState, useEffect, useRef, useCallback } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faInfoCircle,
  faSearchPlus,
  faSearchMinus,
  faHome,
  faExpand
} from '@fortawesome/free-solid-svg-icons';
import { Subnet, CloudProviderType } from '../types';
import CloudProviderIcon from './CloudProviderIcon';
import './SubnetDiagram.css';

interface SubnetDiagramProps {
  subnets: Subnet[];
  viewMode: 'hierarchy' | 'network' | 'cloud';
  isFullscreen: boolean;
}

interface DiagramNode {
  id: string;
  subnet: Subnet;
  x: number;
  y: number;
  width: number;
  height: number;
  level: number;
  children: DiagramNode[];
  parent?: DiagramNode;
}

interface DiagramConnection {
  from: DiagramNode;
  to: DiagramNode;
  type: 'parent-child' | 'network' | 'cloud';
}

interface Transform {
  x: number;
  y: number;
  scale: number;
}

/**
 * Parse CIDR to get network information
 */
function parseCIDR(cidr: string) {
  const [ip, prefixLength] = cidr.split('/');
  const prefix = parseInt(prefixLength);
  const parts = ip.split('.').map(Number);
  const ipNumber = (parts[0] << 24) + (parts[1] << 16) + (parts[2] << 8) + parts[3];
  
  return {
    ip,
    prefix,
    ipNumber,
    networkSize: Math.pow(2, 32 - prefix),
  };
}

/**
 * Check if one subnet contains another
 */
function isSubnetContained(parent: string, child: string): boolean {
  const parentInfo = parseCIDR(parent);
  const childInfo = parseCIDR(child);
  
  // Parent must have smaller prefix (larger network)
  if (parentInfo.prefix >= childInfo.prefix) {
    return false;
  }
  
  // Check if child network is within parent network
  const parentMask = ~((1 << (32 - parentInfo.prefix)) - 1);
  return (parentInfo.ipNumber & parentMask) === (childInfo.ipNumber & parentMask);
}

/**
 * Build hierarchy tree from subnets
 */
function buildHierarchy(subnets: Subnet[]): DiagramNode[] {
  const nodes: DiagramNode[] = subnets.map(subnet => ({
    id: subnet.id,
    subnet,
    x: 0,
    y: 0,
    width: 200,
    height: 80,
    level: 0,
    children: [],
  }));

  // Sort by prefix length (smaller prefix = larger network = higher in hierarchy)
  nodes.sort((a, b) => {
    const aPrefix = parseInt(a.subnet.cidr.split('/')[1]);
    const bPrefix = parseInt(b.subnet.cidr.split('/')[1]);
    return aPrefix - bPrefix;
  });

  // Build parent-child relationships
  const rootNodes: DiagramNode[] = [];
  
  for (const node of nodes) {
    let parent: DiagramNode | undefined;
    
    // Find the most specific parent (smallest network that contains this one)
    for (const potentialParent of nodes) {
      if (potentialParent.id !== node.id && 
          isSubnetContained(potentialParent.subnet.cidr, node.subnet.cidr)) {
        if (!parent || 
            parseInt(potentialParent.subnet.cidr.split('/')[1]) > parseInt(parent.subnet.cidr.split('/')[1])) {
          parent = potentialParent;
        }
      }
    }
    
    if (parent) {
      parent.children.push(node);
      node.parent = parent;
    } else {
      rootNodes.push(node);
    }
  }

  return rootNodes;
}

/**
 * Calculate positions for hierarchy layout
 */
function calculateHierarchyLayout(nodes: DiagramNode[], startX = 50, startY = 50): void {
  let currentY = startY;
  
  function layoutNode(node: DiagramNode, level: number, x: number): number {
    node.level = level;
    node.x = x;
    node.y = currentY;
    
    currentY += node.height + 30;
    
    if (node.children.length > 0) {
      const childX = x + 250;
      for (const child of node.children) {
        layoutNode(child, level + 1, childX);
      }
    }
    
    return currentY;
  }
  
  for (const rootNode of nodes) {
    layoutNode(rootNode, 0, startX);
    currentY += 50; // Extra space between root nodes
  }
}

/**
 * Calculate positions for network layout (by IP ranges)
 */
function calculateNetworkLayout(subnets: Subnet[]): DiagramNode[] {
  const nodes: DiagramNode[] = subnets.map(subnet => ({
    id: subnet.id,
    subnet,
    x: 0,
    y: 0,
    width: 200,
    height: 80,
    level: 0,
    children: [],
  }));

  // Sort by IP address
  nodes.sort((a, b) => {
    const aInfo = parseCIDR(a.subnet.cidr);
    const bInfo = parseCIDR(b.subnet.cidr);
    return aInfo.ipNumber - bInfo.ipNumber;
  });

  // Arrange in a grid based on IP ranges
  const cols = Math.ceil(Math.sqrt(nodes.length));
  nodes.forEach((node, index) => {
    const row = Math.floor(index / cols);
    const col = index % cols;
    node.x = 50 + col * 250;
    node.y = 50 + row * 120;
  });

  return nodes;
}

/**
 * Calculate positions for cloud layout (grouped by provider)
 */
function calculateCloudLayout(subnets: Subnet[]): DiagramNode[] {
  const nodes: DiagramNode[] = subnets.map(subnet => ({
    id: subnet.id,
    subnet,
    x: 0,
    y: 0,
    width: 200,
    height: 80,
    level: 0,
    children: [],
  }));

  // Group by cloud provider
  const groups: { [key: string]: DiagramNode[] } = {};
  
  nodes.forEach(node => {
    const provider = node.subnet.cloudInfo?.provider || 'on-premise';
    if (!groups[provider]) {
      groups[provider] = [];
    }
    groups[provider].push(node);
  });

  // Layout each group
  let currentX = 50;
  Object.entries(groups).forEach(([provider, groupNodes]) => {
    let currentY = 100;
    
    groupNodes.forEach(node => {
      node.x = currentX;
      node.y = currentY;
      currentY += node.height + 20;
    });
    
    currentX += 300;
  });

  return nodes;
}

function SubnetDiagram({ subnets, viewMode, isFullscreen }: SubnetDiagramProps) {
  const [nodes, setNodes] = useState<DiagramNode[]>([]);
  const [connections, setConnections] = useState<DiagramConnection[]>([]);
  const [selectedNode, setSelectedNode] = useState<DiagramNode | null>(null);
  const [hoveredNode, setHoveredNode] = useState<DiagramNode | null>(null);
  const [transform, setTransform] = useState<Transform>({ x: 0, y: 0, scale: 1 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const [tooltipPosition, setTooltipPosition] = useState({ x: 0, y: 0 });
  
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Zoom and pan functionality
  const handleWheel = useCallback((e: React.WheelEvent) => {
    e.preventDefault();
    
    if (!svgRef.current || !containerRef.current) return;
    
    const rect = containerRef.current.getBoundingClientRect();
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    
    // Mouse position relative to container
    const mouseX = e.clientX - rect.left;
    const mouseY = e.clientY - rect.top;
    
    const scaleFactor = e.deltaY > 0 ? 0.9 : 1.1;
    const newScale = Math.max(0.1, Math.min(5, transform.scale * scaleFactor));
    
    // Calculate new position to zoom towards mouse
    const scaleChange = newScale / transform.scale;
    const newX = mouseX - (mouseX - transform.x) * scaleChange;
    const newY = mouseY - (mouseY - transform.y) * scaleChange;
    
    setTransform({
      x: newX,
      y: newY,
      scale: newScale
    });
  }, [transform]);

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    if (e.target === svgRef.current) {
      setIsDragging(true);
      setDragStart({ x: e.clientX - transform.x, y: e.clientY - transform.y });
    }
  }, [transform]);

  const handleMouseMove = useCallback((e: React.MouseEvent) => {
    if (isDragging) {
      setTransform(prev => ({
        ...prev,
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y
      }));
    }
    
    // Update tooltip position
    if (hoveredNode && containerRef.current) {
      const rect = containerRef.current.getBoundingClientRect();
      setTooltipPosition({
        x: e.clientX - rect.left,
        y: e.clientY - rect.top
      });
    }
  }, [isDragging, dragStart, hoveredNode]);

  const handleMouseUp = useCallback(() => {
    setIsDragging(false);
  }, []);

  // Zoom controls
  const zoomIn = () => {
    const newScale = Math.min(5, transform.scale * 1.2);
    setTransform(prev => ({ ...prev, scale: newScale }));
  };

  const zoomOut = () => {
    const newScale = Math.max(0.1, transform.scale * 0.8);
    setTransform(prev => ({ ...prev, scale: newScale }));
  };

  const resetZoom = () => {
    setTransform({ x: 0, y: 0, scale: 1 });
  };

  const fitToScreen = () => {
    if (!containerRef.current || nodes.length === 0) return;
    
    const rect = containerRef.current.getBoundingClientRect();
    const padding = 50;
    
    const minX = Math.min(...nodes.map(n => n.x)) - padding;
    const maxX = Math.max(...nodes.map(n => n.x + n.width)) + padding;
    const minY = Math.min(...nodes.map(n => n.y)) - padding;
    const maxY = Math.max(...nodes.map(n => n.y + n.height)) + padding;
    
    const contentWidth = maxX - minX;
    const contentHeight = maxY - minY;
    
    const scaleX = rect.width / contentWidth;
    const scaleY = rect.height / contentHeight;
    const scale = Math.min(scaleX, scaleY, 1);
    
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    const contentCenterX = (minX + maxX) / 2;
    const contentCenterY = (minY + maxY) / 2;
    
    setTransform({
      x: centerX - contentCenterX * scale,
      y: centerY - contentCenterY * scale,
      scale
    });
  };

  useEffect(() => {
    let calculatedNodes: DiagramNode[];
    
    switch (viewMode) {
      case 'hierarchy':
        calculatedNodes = buildHierarchy(subnets);
        calculateHierarchyLayout(calculatedNodes);
        // Flatten for rendering
        const flatNodes: DiagramNode[] = [];
        function collectNodes(nodeList: DiagramNode[]) {
          for (const node of nodeList) {
            flatNodes.push(node);
            collectNodes(node.children);
          }
        }
        collectNodes(calculatedNodes);
        calculatedNodes = flatNodes;
        break;
      case 'network':
        calculatedNodes = calculateNetworkLayout(subnets);
        break;
      case 'cloud':
        calculatedNodes = calculateCloudLayout(subnets);
        break;
      default:
        calculatedNodes = [];
    }
    
    setNodes(calculatedNodes);
    
    // Calculate connections based on view mode
    const newConnections: DiagramConnection[] = [];
    if (viewMode === 'hierarchy') {
      calculatedNodes.forEach(node => {
        if (node.parent) {
          newConnections.push({
            from: node.parent,
            to: node,
            type: 'parent-child'
          });
        }
      });
    }
    
    setConnections(newConnections);
  }, [subnets, viewMode]);

  const handleNodeClick = (node: DiagramNode) => {
    setSelectedNode(selectedNode?.id === node.id ? null : node);
  };

  const getNodeColor = (subnet: Subnet): string => {
    if (subnet.cloudInfo?.provider) {
      switch (subnet.cloudInfo.provider) {
        case CloudProviderType.AWS: return '#FF9900';
        case CloudProviderType.AZURE: return '#0078D4';
        case CloudProviderType.GCP: return '#4285F4';
        case CloudProviderType.SCALEWAY: return '#4F0599';
        case CloudProviderType.OVH: return '#123F6D';
        default: return '#6B7280';
      }
    }
    return '#6B7280';
  };

  const getUtilizationColor = (percent: number): string => {
    if (percent >= 80) return '#EF4444';
    if (percent >= 60) return '#F59E0B';
    if (percent >= 40) return '#10B981';
    return '#6B7280';
  };

  // Calculate SVG dimensions
  const maxX = Math.max(...nodes.map(n => n.x + n.width), 500);
  const maxY = Math.max(...nodes.map(n => n.y + n.height), 400);

  // Reset zoom when view mode changes
  useEffect(() => {
    resetZoom();
  }, [viewMode]);

  // Fit to screen when nodes change
  useEffect(() => {
    if (nodes.length > 0) {
      setTimeout(fitToScreen, 100);
    }
  }, [nodes]);

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!containerRef.current?.contains(document.activeElement)) return;
      
      switch (e.key) {
        case '+':
        case '=':
          e.preventDefault();
          zoomIn();
          break;
        case '-':
          e.preventDefault();
          zoomOut();
          break;
        case '0':
          e.preventDefault();
          resetZoom();
          break;
        case 'f':
        case 'F':
          e.preventDefault();
          fitToScreen();
          break;
        case 'Escape':
          setSelectedNode(null);
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, []);

  return (
    <div 
      className="subnet-diagram"
      ref={containerRef}
      tabIndex={0}
      onWheel={handleWheel}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
    >
      <svg
        ref={svgRef}
        width="100%"
        height="100%"
        viewBox={`0 0 ${maxX + 100} ${maxY + 100}`}
        className={`diagram-svg ${isDragging ? 'dragging' : ''}`}
        style={{
          transform: `translate(${transform.x}px, ${transform.y}px) scale(${transform.scale})`,
          transformOrigin: '0 0'
        }}
      >
        {/* Connections */}
        {connections.map((connection, index) => (
          <line
            key={index}
            x1={connection.from.x + connection.from.width / 2}
            y1={connection.from.y + connection.from.height}
            x2={connection.to.x + connection.to.width / 2}
            y2={connection.to.y}
            stroke="#6B7280"
            strokeWidth="2"
            strokeDasharray={connection.type === 'parent-child' ? '0' : '5,5'}
            className="connection-line"
          />
        ))}

        {/* Nodes */}
        {nodes.map((node) => (
          <g key={node.id} className="diagram-node">
            {/* Node background */}
            <rect
              x={node.x}
              y={node.y}
              width={node.width}
              height={node.height}
              fill={selectedNode?.id === node.id ? '#3B82F6' : '#FFFFFF'}
              stroke={getNodeColor(node.subnet)}
              strokeWidth="2"
              rx="8"
              className="node-background"
              onClick={() => handleNodeClick(node)}
              onMouseEnter={() => setHoveredNode(node)}
              onMouseLeave={() => setHoveredNode(null)}
            />

            {/* Utilization bar */}
            <rect
              x={node.x + 5}
              y={node.y + 5}
              width={(node.width - 10) * (node.subnet.utilization.utilizationPercent / 100)}
              height="4"
              fill={getUtilizationColor(node.subnet.utilization.utilizationPercent)}
              rx="2"
            />

            {/* Node content */}
            <text
              x={node.x + 10}
              y={node.y + 25}
              fontSize="14"
              fontWeight="bold"
              fill={selectedNode?.id === node.id ? '#FFFFFF' : '#1F2937'}
            >
              {node.subnet.cidr}
            </text>
            
            <text
              x={node.x + 10}
              y={node.y + 40}
              fontSize="12"
              fill={selectedNode?.id === node.id ? '#E5E7EB' : '#6B7280'}
            >
              {node.subnet.name}
            </text>

            <text
              x={node.x + 10}
              y={node.y + 55}
              fontSize="10"
              fill={selectedNode?.id === node.id ? '#E5E7EB' : '#9CA3AF'}
            >
              {node.subnet.location}
            </text>

            {/* Cloud provider icon */}
            {node.subnet.cloudInfo?.provider && (
              <foreignObject
                x={node.x + node.width - 30}
                y={node.y + 10}
                width="20"
                height="20"
              >
                <CloudProviderIcon
                  provider={node.subnet.cloudInfo.provider}
                  size="sm"
                />
              </foreignObject>
            )}

            {/* Utilization percentage */}
            <text
              x={node.x + node.width - 10}
              y={node.y + node.height - 10}
              fontSize="10"
              textAnchor="end"
              fill={selectedNode?.id === node.id ? '#E5E7EB' : '#6B7280'}
            >
              {node.subnet.utilization.utilizationPercent.toFixed(1)}%
            </text>
          </g>
        ))}
      </svg>

      {/* Zoom and pan controls */}
      <div className="diagram-controls">
        <button
          onClick={zoomIn}
          className="diagram-control-btn"
          title="Zoom In"
        >
          <FontAwesomeIcon icon={faSearchPlus} />
        </button>
        <button
          onClick={zoomOut}
          className="diagram-control-btn"
          title="Zoom Out"
        >
          <FontAwesomeIcon icon={faSearchMinus} />
        </button>
        <button
          onClick={resetZoom}
          className="diagram-control-btn"
          title="Reset Zoom"
        >
          <FontAwesomeIcon icon={faHome} />
        </button>
        <button
          onClick={fitToScreen}
          className="diagram-control-btn"
          title="Fit to Screen"
        >
          <FontAwesomeIcon icon={faExpand} />
        </button>
      </div>

      {/* Zoom indicator */}
      <div className="zoom-indicator">
        {Math.round(transform.scale * 100)}%
      </div>

      {/* Keyboard shortcuts legend */}
      <div className="keyboard-shortcuts">
        <div className="shortcuts-title">Raccourcis clavier:</div>
        <div className="shortcut-item">+ / - : Zoom</div>
        <div className="shortcut-item">0 : Reset zoom</div>
        <div className="shortcut-item">F : Ajuster à l'écran</div>
        <div className="shortcut-item">Esc : Fermer détails</div>
      </div>

      {/* Node details panel */}
      {selectedNode && (
        <div className="node-details-panel">
          <div className="panel-header">
            <FontAwesomeIcon icon={faInfoCircle} />
            <h3>Subnet Details</h3>
            <button
              onClick={() => setSelectedNode(null)}
              className="close-button"
            >
              ×
            </button>
          </div>
          
          <div className="panel-content">
            <div className="detail-row">
              <label>CIDR:</label>
              <span>{selectedNode.subnet.cidr}</span>
            </div>
            
            <div className="detail-row">
              <label>Name:</label>
              <span>{selectedNode.subnet.name}</span>
            </div>
            
            <div className="detail-row">
              <label>Location:</label>
              <span>{selectedNode.subnet.location}</span>
            </div>
            
            <div className="detail-row">
              <label>Type:</label>
              <span>{selectedNode.subnet.locationType}</span>
            </div>
            
            {selectedNode.subnet.cloudInfo && (
              <>
                <div className="detail-row">
                  <label>Provider:</label>
                  <span className="provider-info">
                    <CloudProviderIcon
                      provider={selectedNode.subnet.cloudInfo.provider}
                      size="sm"
                    />
                    {selectedNode.subnet.cloudInfo.provider.toUpperCase()}
                  </span>
                </div>
                
                <div className="detail-row">
                  <label>Region:</label>
                  <span>{selectedNode.subnet.cloudInfo.region}</span>
                </div>
                
                <div className="detail-row">
                  <label>Account:</label>
                  <span>{selectedNode.subnet.cloudInfo.accountId}</span>
                </div>
              </>
            )}
            
            <div className="detail-row">
              <label>Utilization:</label>
              <span>
                {selectedNode.subnet.utilization.allocatedIps} / {selectedNode.subnet.utilization.totalIps} IPs
                ({selectedNode.subnet.utilization.utilizationPercent.toFixed(1)}%)
              </span>
            </div>
            
            {selectedNode.subnet.details && (
              <>
                <div className="detail-row">
                  <label>Network:</label>
                  <span>{selectedNode.subnet.details.network}</span>
                </div>
                
                <div className="detail-row">
                  <label>Broadcast:</label>
                  <span>{selectedNode.subnet.details.broadcast}</span>
                </div>
                
                <div className="detail-row">
                  <label>Host Range:</label>
                  <span>{selectedNode.subnet.details.hostMin} - {selectedNode.subnet.details.hostMax}</span>
                </div>
              </>
            )}
          </div>
        </div>
      )}

      {/* Tooltip for hovered node */}
      {hoveredNode && !selectedNode && (
        <div 
          className="node-tooltip"
          style={{
            left: tooltipPosition.x,
            top: tooltipPosition.y - 10,
          }}
        >
          <strong>{hoveredNode.subnet.name}</strong><br />
          {hoveredNode.subnet.cidr}<br />
          {hoveredNode.subnet.utilization.utilizationPercent.toFixed(1)}% used
        </div>
      )}
    </div>
  );
}

export default SubnetDiagram;

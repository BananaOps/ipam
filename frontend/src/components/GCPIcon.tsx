import gcpLogo from '../assets/gcp.svg';

interface GCPIconProps {
  size?: 'sm' | 'lg' | '1x' | '2x' | '3x';
  className?: string;
}

function GCPIcon({ size = '1x', className = '' }: GCPIconProps) {
  const sizeMap = {
    'sm': '14px',
    'lg': '20px',
    '1x': '16px',
    '2x': '32px',
    '3x': '48px',
  };

  const iconSize = sizeMap[size] || sizeMap['1x'];

  return (
    <img 
      src={gcpLogo}
      alt="Google Cloud Platform"
      width={iconSize}
      height={iconSize}
      className={className}
      style={{ display: 'inline-block', verticalAlign: 'middle', objectFit: 'contain' }}
    />
  );
}

export default GCPIcon;

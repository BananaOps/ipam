import azureLogo from '../assets/azure.png';

interface AzureIconProps {
  size?: 'sm' | 'lg' | '1x' | '2x' | '3x';
  className?: string;
}

function AzureIcon({ size = '1x', className = '' }: AzureIconProps) {
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
      src={azureLogo}
      alt="Azure"
      width={iconSize}
      height={iconSize}
      className={className}
      style={{ display: 'inline-block', verticalAlign: 'middle', objectFit: 'contain' }}
    />
  );
}

export default AzureIcon;

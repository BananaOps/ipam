# IPAM by BananaOps - Logo Assets

## Overview
This directory contains the logo assets for IPAM by BananaOps, designed in the Cyber Minimal style.

## Design Concept
The logo represents a network topology with geometric precision:
- **Hexagon frame**: Represents the network perimeter/boundary
- **Network nodes**: Five circular nodes representing IP addresses in a subnet
- **Central hub**: The gateway/router connecting all nodes
- **Connection lines**: Network paths between nodes
- **CIDR indicator**: Shows the subnet notation (/24)

## Color Palette (Cyber Minimal)
- **Bleu nuit (Dark Primary)**: #0A1A2F - Used for filled elements
- **Bleu cyan (Cyan Accent)**: #0EA5E9 - Primary brand color, used for strokes and highlights
- **Gris clair (Light Gray)**: #F3F4F6 - Light theme background
- **Blanc (White)**: #FFFFFF - Light theme text and elements

## Available Formats

### SVG (Recommended for Web)
- **logo.svg** (200x200px) - Full logo for splash screens, about pages, etc.
- **logo-horizontal.svg** (40x40px) - Compact version for headers and navigation

### PNG (For External Use)
To generate PNG versions from SVG:

```bash
# Using ImageMagick (if available)
convert -background none logo.svg -resize 512x512 logo-512.png
convert -background none logo.svg -resize 256x256 logo-256.png
convert -background none logo.svg -resize 128x128 logo-128.png
convert -background none logo.svg -resize 64x64 logo-64.png
convert -background none logo.svg -resize 32x32 logo-32.png

# Using Inkscape (if available)
inkscape logo.svg --export-type=png --export-width=512 --export-filename=logo-512.png
```

### Recommended Sizes
- **512x512**: App icons, social media
- **256x256**: Large displays
- **128x128**: Standard icons
- **64x64**: Small icons
- **32x32**: Favicons

## Usage in Application

### React Component
```tsx
import Logo from './components/Logo';

// Compact version with text (header)
<Logo variant="compact" size="medium" showText={true} />

// Full version without text
<Logo variant="full" size="large" showText={false} />

// Small version
<Logo variant="compact" size="small" showText={false} />
```

### Direct SVG Usage
```html
<img src="/logo.svg" alt="IPAM by BananaOps" />
<img src="/logo-horizontal.svg" alt="IPAM by BananaOps" />
```

## Theme Compatibility
The logo uses `currentColor` in the compact version, making it automatically adapt to the current theme (dark/light). The full version uses explicit colors from the Cyber Minimal palette.

## License
© BananaOps - All rights reserved

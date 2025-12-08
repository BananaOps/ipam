# IPAM by BananaOps - Logo Design Documentation

## Design Overview

The IPAM by BananaOps logo follows the **Cyber Minimal** design style, combining geometric precision with a modern, tech-focused aesthetic. The logo represents network topology and IP address management through clean, minimalist shapes.

## Design Concept

### Visual Elements

1. **Hexagonal Frame**
   - Represents the network perimeter/boundary
   - Geometric shape symbolizing structure and organization
   - Uses stroke-only rendering for a lightweight, modern look

2. **Network Nodes (5 circles)**
   - Represent individual IP addresses within a subnet
   - Strategically positioned to show network distribution
   - Filled with cyan accent color for visibility

3. **Central Hub (circle with stroke)**
   - Represents the gateway/router
   - Central position shows its role as the connection point
   - Larger than nodes to indicate importance

4. **Connection Lines**
   - Show network paths between nodes and hub
   - Semi-transparent to create depth
   - Illustrate the interconnected nature of IP networks

5. **CIDR Notation (/24)**
   - Included in full logo version
   - Directly references IP subnet notation
   - Uses monospace font for technical authenticity

## Color Palette (Cyber Minimal)

| Color Name | Hex Code | Usage |
|------------|----------|-------|
| Bleu nuit (Dark Primary) | `#0A1A2F` | Filled elements, dark theme background |
| Bleu cyan (Cyan Accent) | `#0EA5E9` | Primary brand color, strokes, highlights |
| Gris clair (Light Gray) | `#F3F4F6` | Light theme background |
| Blanc (White) | `#FFFFFF` | Light theme text and elements |

## Logo Variants

### 1. Full Logo (200x200px)
- **File**: `frontend/public/logo.svg`
- **Usage**: Splash screens, about pages, marketing materials
- **Features**: Complete design with CIDR notation
- **Colors**: Explicit Cyber Minimal palette

### 2. Compact Logo (40x40px)
- **File**: `frontend/public/logo-horizontal.svg`
- **Usage**: Headers, navigation, small UI elements
- **Features**: Simplified design without text
- **Colors**: Uses `currentColor` for theme adaptation

### 3. Favicon (32x32px)
- **File**: `frontend/public/favicon.svg`
- **Usage**: Browser tabs, bookmarks
- **Features**: Optimized for small sizes

## Implementation

### React Component

The logo is implemented as a flexible React component:

```tsx
<Logo 
  variant="compact"  // or "full"
  size="medium"      // "small", "medium", or "large"
  showText={true}    // show/hide text
  className=""       // custom classes
/>
```

### Usage Examples

**Header/Navigation:**
```tsx
<Logo variant="compact" size="medium" showText={true} />
```

**Loading Screen:**
```tsx
<Logo variant="full" size="large" showText={false} />
```

**Footer:**
```tsx
<Logo variant="compact" size="small" showText={true} />
```

## Theme Compatibility

The logo is designed to work seamlessly with both dark and light themes:

- **Dark Theme**: Cyan accent (#0EA5E9) stands out against dark backgrounds
- **Light Theme**: Dark elements (#0A1A2F) provide contrast on light backgrounds
- **Auto Theme**: Compact variant uses `currentColor` for automatic adaptation

## Design Principles

1. **Minimalism**: Clean lines, no unnecessary decoration
2. **Geometric Precision**: Perfect shapes, consistent spacing
3. **Technical Authenticity**: CIDR notation, network topology
4. **Scalability**: SVG format ensures crisp rendering at any size
5. **Accessibility**: High contrast ratios, clear visual hierarchy

## File Formats

### SVG (Primary)
- Vector format for perfect scaling
- Small file size
- CSS-animatable
- Theme-adaptable with `currentColor`

### PNG (Optional)
Generate from SVG for external use:
```bash
# Using ImageMagick
convert -background none logo.svg -resize 512x512 logo-512.png

# Using Inkscape
inkscape logo.svg --export-type=png --export-width=512
```

## Brand Guidelines

### Do's ✓
- Use the logo on dark or light backgrounds with sufficient contrast
- Maintain aspect ratio when resizing
- Use provided color palette
- Keep clear space around logo (minimum 10px)

### Don'ts ✗
- Don't distort or skew the logo
- Don't change the color palette
- Don't add effects (shadows, gradients, etc.)
- Don't place on busy backgrounds without a solid backdrop

## Technical Specifications

- **Format**: SVG (Scalable Vector Graphics)
- **Viewbox**: 200x200 (full), 40x40 (compact), 32x32 (favicon)
- **Stroke Width**: 2-3px depending on size
- **Font**: Monospace for CIDR notation
- **Opacity**: 0.6 for connection lines

## Integration Checklist

- [x] SVG logo files created (full, compact, favicon)
- [x] Logo component implemented
- [x] CSS styles with theme support
- [x] Integrated into Layout header
- [x] Favicon updated in index.html
- [x] Component tests written and passing
- [x] Documentation created

## Future Enhancements

- Animated version for loading states
- PNG exports in multiple sizes
- Social media variants (square, wide)
- Dark/light specific variants
- Monochrome version for print

---

**Design Status**: ✅ Complete
**Last Updated**: December 2024
**Designer**: IPAM by BananaOps Team

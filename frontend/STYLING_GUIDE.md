# IPAM by BananaOps - Styling Guide

## Cyber Minimal Design System

This document describes the styling approach and design system used in IPAM by BananaOps.

## Color Palette

### Primary Colors
- **Bleu nuit (Dark Primary)**: `#0A1A2F` - Main dark background
- **Bleu cyan (Cyan Accent)**: `#0EA5E9` - Primary accent color for interactive elements
- **Gris clair (Light Gray)**: `#F3F4F6` - Light background and secondary elements
- **Blanc pur (White)**: `#FFFFFF` - Pure white for text and highlights

### Secondary Colors
- **Dark Gray**: `#1F2937` - Secondary dark background
- **Medium Gray**: `#6B7280` - Secondary text and borders
- **Light Border**: `#E5E7EB` - Light theme borders
- **Dark Border**: `#374151` - Dark theme borders

### Status Colors
- **Success**: `#10B981` - Success states and positive indicators
- **Warning**: `#F59E0B` - Warning states and high utilization
- **Error**: `#EF4444` - Error states and destructive actions
- **Info**: `#3B82F6` - Informational elements

## Theme System

### Light Theme (Default)
- Background: White and Light Gray
- Text: Dark Primary
- Borders: Light Border
- Accent: Cyan Accent

### Dark Theme
- Background: Dark Primary and Dark Gray
- Text: White
- Borders: Dark Border
- Accent: Cyan Accent

### Auto Theme
Automatically synchronizes with system preferences using `prefers-color-scheme` media query.

## Typography

### Font Family
- Primary: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif
- Monospace: 'Courier New', monospace (for code and technical data)

### Font Sizes
- Extra Small: 0.75rem (12px)
- Small: 0.875rem (14px)
- Base: 1rem (16px)
- Large: 1.125rem (18px)
- XL: 1.25rem (20px)
- 2XL: 1.5rem (24px)
- 3XL: 1.875rem (30px)

### Responsive Typography
- Desktop (>1024px): 16px base
- Tablet (769-1024px): 15px base
- Mobile (481-768px): 14px base
- Small Mobile (<480px): 13px base

## Component Styling

### Cards
Cards use gradient backgrounds and subtle shadows for depth:
```css
background: linear-gradient(135deg, var(--bg-secondary) 0%, var(--bg-primary) 100%);
box-shadow: 0 2px 8px var(--shadow-color);
border: 1px solid var(--border-color);
```

### Buttons
Three main button styles:
- **Primary**: Cyan accent background, white text
- **Secondary**: Transparent background, bordered
- **Danger**: Red background for destructive actions

All buttons have hover effects with transform and shadow changes.

### Forms
Form inputs have:
- Focus states with cyan accent border and shadow
- Error states with red border
- Disabled states with reduced opacity
- Smooth transitions on all state changes

### Tables
Tables feature:
- Hover effects on rows with left border accent
- Sticky headers for long lists
- Responsive horizontal scrolling on mobile
- Zebra striping for better readability

## Animations

### Standard Animations
- **fadeIn**: Smooth opacity transition (0.3s)
- **slideInUp**: Slide from bottom with fade (0.3s)
- **slideInDown**: Slide from top with fade (0.3s)
- **spin**: Continuous rotation for loading states
- **pulse**: Breathing effect for attention

### Hover Effects
- **hover-lift**: Translates element up with shadow
- **hover-scale**: Scales element to 1.05
- All interactive elements have 0.2s transitions

### Reduced Motion
Respects `prefers-reduced-motion` media query for accessibility.

## Responsive Design

### Breakpoints
- **Mobile**: < 480px
- **Tablet**: 481px - 768px
- **Desktop**: 769px - 1024px
- **Large Desktop**: > 1024px

### Mobile-First Approach
Base styles are mobile-optimized, with progressive enhancement for larger screens.

### Key Responsive Features
1. **Navigation**: Collapses to icon-only on mobile
2. **Tables**: Horizontal scroll on mobile with minimum width
3. **Forms**: Full-width inputs on mobile
4. **Grids**: Auto-fit columns that stack on mobile
5. **Typography**: Scales down on smaller screens

## Accessibility

### Focus States
All interactive elements have visible focus indicators using `outline` with cyan accent color.

### Color Contrast
All text meets WCAG AA standards for contrast ratios:
- Normal text: 4.5:1 minimum
- Large text: 3:1 minimum

### Screen Reader Support
- Semantic HTML elements
- ARIA labels where needed
- `.sr-only` class for screen reader-only content

### Keyboard Navigation
All interactive elements are keyboard accessible with proper tab order.

## Utility Classes

### Layout
- Flexbox: `.flex`, `.flex-col`, `.items-center`, `.justify-between`
- Grid: `.grid`, `.grid-cols-2`, `.grid-auto-fit`
- Spacing: `.p-{1-4}`, `.m-{1-4}`, `.gap-{1-3}`

### Display
- `.block`, `.inline-block`, `.hidden`
- `.w-full`, `.h-full`, `.min-h-screen`

### Text
- Sizes: `.text-xs` to `.text-3xl`
- Weight: `.font-normal`, `.font-medium`, `.font-semibold`, `.font-bold`
- Alignment: `.text-left`, `.text-center`, `.text-right`

### Effects
- Shadows: `.shadow-sm` to `.shadow-xl`
- Rounded: `.rounded-sm` to `.rounded-full`
- Opacity: `.opacity-{0,25,50,75,100}`

## Best Practices

### 1. Use CSS Variables
Always use CSS variables for colors and theme-dependent values:
```css
color: var(--text-primary);
background-color: var(--bg-primary);
border-color: var(--border-color);
```

### 2. Consistent Spacing
Use the spacing scale (0.5rem increments) for consistent layout:
- Small: 0.5rem
- Medium: 1rem
- Large: 1.5rem
- XL: 2rem

### 3. Smooth Transitions
Add transitions to interactive elements:
```css
transition: all 0.2s ease;
```

### 4. Hover States
Provide visual feedback on hover:
```css
.element:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px var(--shadow-color);
}
```

### 5. Mobile-First
Write base styles for mobile, then add media queries for larger screens:
```css
.element {
  /* Mobile styles */
}

@media (min-width: 769px) {
  .element {
    /* Desktop styles */
  }
}
```

### 6. Semantic Class Names
Use descriptive, purpose-based class names:
- Good: `.subnet-list-filters`, `.utilization-progress-bar`
- Bad: `.blue-box`, `.big-text`

### 7. Component Isolation
Keep component styles scoped to their CSS files to avoid conflicts.

### 8. Performance
- Use `transform` and `opacity` for animations (GPU-accelerated)
- Avoid animating `width`, `height`, or `top/left` properties
- Use `will-change` sparingly for complex animations

## Print Styles

Print styles hide navigation and interactive elements, focusing on content:
```css
@media print {
  .app-header,
  .main-nav,
  .theme-toggle,
  .btn-primary,
  .btn-secondary {
    display: none;
  }
}
```

## Browser Support

- Chrome/Edge: Latest 2 versions
- Firefox: Latest 2 versions
- Safari: Latest 2 versions
- Mobile browsers: iOS Safari 12+, Chrome Android 90+

## Future Enhancements

1. **Dark Mode Improvements**: Add more granular control over dark mode colors
2. **Animation Library**: Create reusable animation components
3. **Component Variants**: Add more button and card variants
4. **Theming System**: Allow custom color schemes
5. **CSS-in-JS**: Consider migrating to styled-components or emotion for better TypeScript integration

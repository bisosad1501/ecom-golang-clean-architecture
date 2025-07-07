// ===== CLASS NAME UTILITIES =====

import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

/**
 * Combines class names using clsx and tailwind-merge
 * This is the most commonly used utility for combining Tailwind classes
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Creates conditional class names based on variants
 * Useful for component variants
 */
export function createVariants<T extends Record<string, Record<string, string>>>(
  variants: T
) {
  return function getVariantClasses<K extends keyof T>(
    variant: K,
    value: keyof T[K]
  ): string {
    return variants[variant]?.[value] || ''
  }
}

/**
 * Merges base classes with conditional classes
 */
export function mergeClasses(
  baseClasses: string,
  conditionalClasses: Record<string, boolean | undefined>,
  className?: string
): string {
  const conditionalClassArray = Object.entries(conditionalClasses)
    .filter(([, condition]) => condition)
    .map(([cls]) => cls)

  return cn(baseClasses, ...conditionalClassArray, className)
}

/**
 * Creates responsive class names
 */
export function responsive(classes: {
  base?: string
  sm?: string
  md?: string
  lg?: string
  xl?: string
  '2xl'?: string
}): string {
  const { base, sm, md, lg, xl, '2xl': xl2 } = classes
  
  return cn(
    base,
    sm && `sm:${sm}`,
    md && `md:${md}`,
    lg && `lg:${lg}`,
    xl && `xl:${xl}`,
    xl2 && `2xl:${xl2}`
  )
}

/**
 * Creates state-based class names (hover, focus, active, etc.)
 */
export function stateClasses(classes: {
  base?: string
  hover?: string
  focus?: string
  active?: string
  disabled?: string
  loading?: string
}): string {
  const { base, hover, focus, active, disabled, loading } = classes
  
  return cn(
    base,
    hover && `hover:${hover}`,
    focus && `focus:${focus}`,
    active && `active:${active}`,
    disabled && `disabled:${disabled}`,
    loading && `loading:${loading}`
  )
}

/**
 * Creates theme-based class names (light/dark mode)
 */
export function themeClasses(classes: {
  base?: string
  light?: string
  dark?: string
}): string {
  const { base, light, dark } = classes
  
  return cn(
    base,
    light && `light:${light}`,
    dark && `dark:${dark}`
  )
}

/**
 * Utility for creating component size variants
 */
export function sizeVariants<T extends string>(
  sizes: Record<T, string>,
  defaultSize: T
) {
  return function getSizeClasses(size: T = defaultSize): string {
    return sizes[size] || sizes[defaultSize]
  }
}

/**
 * Utility for creating component color variants
 */
export function colorVariants<T extends string>(
  colors: Record<T, string>,
  defaultColor: T
) {
  return function getColorClasses(color: T = defaultColor): string {
    return colors[color] || colors[defaultColor]
  }
}

/**
 * Creates animation class names with optional delays and durations
 */
export function animationClasses(options: {
  animation?: string
  delay?: string
  duration?: string
  easing?: string
  fillMode?: string
}): string {
  const { animation, delay, duration, easing, fillMode } = options
  
  return cn(
    animation && `animate-${animation}`,
    delay && `animation-delay-${delay}`,
    duration && `animation-duration-${duration}`,
    easing && `animation-${easing}`,
    fillMode && `animation-fill-${fillMode}`
  )
}

/**
 * Creates grid layout class names
 */
export function gridClasses(options: {
  cols?: number | string
  rows?: number | string
  gap?: string
  colSpan?: number | string
  rowSpan?: number | string
  colStart?: number | string
  rowStart?: number | string
}): string {
  const { cols, rows, gap, colSpan, rowSpan, colStart, rowStart } = options
  
  return cn(
    cols && `grid-cols-${cols}`,
    rows && `grid-rows-${rows}`,
    gap && `gap-${gap}`,
    colSpan && `col-span-${colSpan}`,
    rowSpan && `row-span-${rowSpan}`,
    colStart && `col-start-${colStart}`,
    rowStart && `row-start-${rowStart}`
  )
}

/**
 * Creates flexbox layout class names
 */
export function flexClasses(options: {
  direction?: 'row' | 'col' | 'row-reverse' | 'col-reverse'
  wrap?: 'wrap' | 'nowrap' | 'wrap-reverse'
  justify?: 'start' | 'end' | 'center' | 'between' | 'around' | 'evenly'
  align?: 'start' | 'end' | 'center' | 'baseline' | 'stretch'
  gap?: string
  grow?: boolean | number
  shrink?: boolean | number
  basis?: string
}): string {
  const { direction, wrap, justify, align, gap, grow, shrink, basis } = options
  
  return cn(
    'flex',
    direction && `flex-${direction}`,
    wrap && `flex-${wrap}`,
    justify && `justify-${justify}`,
    align && `items-${align}`,
    gap && `gap-${gap}`,
    grow === true && 'flex-grow',
    grow === false && 'flex-grow-0',
    typeof grow === 'number' && `flex-grow-${grow}`,
    shrink === true && 'flex-shrink',
    shrink === false && 'flex-shrink-0',
    typeof shrink === 'number' && `flex-shrink-${shrink}`,
    basis && `flex-basis-${basis}`
  )
}

/**
 * Creates spacing class names (padding, margin)
 */
export function spacingClasses(options: {
  p?: string | number
  px?: string | number
  py?: string | number
  pt?: string | number
  pr?: string | number
  pb?: string | number
  pl?: string | number
  m?: string | number
  mx?: string | number
  my?: string | number
  mt?: string | number
  mr?: string | number
  mb?: string | number
  ml?: string | number
}): string {
  const { p, px, py, pt, pr, pb, pl, m, mx, my, mt, mr, mb, ml } = options
  
  return cn(
    p && `p-${p}`,
    px && `px-${px}`,
    py && `py-${py}`,
    pt && `pt-${pt}`,
    pr && `pr-${pr}`,
    pb && `pb-${pb}`,
    pl && `pl-${pl}`,
    m && `m-${m}`,
    mx && `mx-${mx}`,
    my && `my-${my}`,
    mt && `mt-${mt}`,
    mr && `mr-${mr}`,
    mb && `mb-${mb}`,
    ml && `ml-${ml}`
  )
}

/**
 * Creates border class names
 */
export function borderClasses(options: {
  width?: string | number
  color?: string
  style?: 'solid' | 'dashed' | 'dotted' | 'double' | 'none'
  radius?: string | number
  side?: 't' | 'r' | 'b' | 'l' | 'x' | 'y'
}): string {
  const { width, color, style, radius, side } = options
  
  const borderPrefix = side ? `border-${side}` : 'border'
  
  return cn(
    width && `${borderPrefix}-${width}`,
    color && `border-${color}`,
    style && `border-${style}`,
    radius && `rounded-${radius}`
  )
}

/**
 * Creates text styling class names
 */
export function textClasses(options: {
  size?: string
  weight?: string
  color?: string
  align?: 'left' | 'center' | 'right' | 'justify'
  decoration?: 'underline' | 'overline' | 'line-through' | 'no-underline'
  transform?: 'uppercase' | 'lowercase' | 'capitalize' | 'normal-case'
  leading?: string
  tracking?: string
}): string {
  const { size, weight, color, align, decoration, transform, leading, tracking } = options
  
  return cn(
    size && `text-${size}`,
    weight && `font-${weight}`,
    color && `text-${color}`,
    align && `text-${align}`,
    decoration && decoration,
    transform && transform,
    leading && `leading-${leading}`,
    tracking && `tracking-${tracking}`
  )
}

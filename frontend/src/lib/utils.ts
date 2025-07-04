import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

/**
 * Combines CSS class names using clsx and tailwind-merge.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Formats bytes into human-readable units.
 */
export function formatBytes(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

/**
 * Formats a number as a percentage with specified decimal places.
 */
export function formatPercentage(value: number, decimals = 1): string {
  return `${value.toFixed(decimals)}%`
}

/**
 * Formats temperature in Celsius with degree symbol.
 */
export function formatTemperature(celsius: number): string {
  return `${celsius.toFixed(1)}°C`
}

/**
 * Formats timestamp string into localized date-time string.
 */
export function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp)
  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

/**
 * Returns CSS color class based on GPU utilization percentage.
 */
export function getUtilizationColor(utilization: number): string {
  if (utilization < 30) return 'text-green-500'
  if (utilization < 70) return 'text-yellow-500'
  return 'text-red-500'
}

/**
 * Returns background color class for utilization progress bars.
 */
export function getUtilizationBarColor(utilization: number): string {
  if (utilization < 30) return 'bg-green-500'
  if (utilization < 70) return 'bg-yellow-500'
  return 'bg-red-500'
}

/**
 * Returns CSS color class based on GPU temperature.
 */
export function getTemperatureColor(temperature: number): string {
  if (temperature < 60) return 'text-green-500'
  if (temperature < 80) return 'text-yellow-500'
  return 'text-red-500'
} 
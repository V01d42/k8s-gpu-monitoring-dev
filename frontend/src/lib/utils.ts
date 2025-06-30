import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

// CSS クラス名を結合するユーティリティ
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// 数値を適切な単位で表示するユーティリティ
export function formatBytes(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

// パーセンテージを適切にフォーマット
export function formatPercentage(value: number, decimals = 1): string {
  return `${value.toFixed(decimals)}%`
}

// 温度を適切にフォーマット
export function formatTemperature(celsius: number): string {
  return `${celsius.toFixed(1)}°C`
}

// 電力を適切にフォーマット
export function formatPower(watts: number): string {
  return `${watts.toFixed(1)}W`
}

// 時間を適切にフォーマット
export function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp)
  return date.toLocaleString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

// 相対時間を表示
export function formatRelativeTime(timestamp: string): string {
  const now = new Date()
  const past = new Date(timestamp)
  const diffMs = now.getTime() - past.getTime()
  
  const seconds = Math.floor(diffMs / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 0) return `${days}日前`
  if (hours > 0) return `${hours}時間前`
  if (minutes > 0) return `${minutes}分前`
  return `${seconds}秒前`
}

// 利用率に基づいて色を取得
export function getUtilizationColor(utilization: number): string {
  if (utilization < 30) return 'text-green-500'
  if (utilization < 70) return 'text-yellow-500'
  return 'text-red-500'
}

// 利用率に基づいてバーの色を取得
export function getUtilizationBarColor(utilization: number): string {
  if (utilization < 30) return 'bg-green-500'
  if (utilization < 70) return 'bg-yellow-500'
  return 'bg-red-500'
}

// 温度に基づいて色を取得
export function getTemperatureColor(temperature: number): string {
  if (temperature < 60) return 'text-green-500'
  if (temperature < 80) return 'text-yellow-500'
  return 'text-red-500'
}

// デバウンス関数
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout
  return (...args: Parameters<T>) => {
    clearTimeout(timeout)
    timeout = setTimeout(() => func(...args), wait)
  }
}

// 値の範囲チェック
export function isInRange(value: number, range: [number, number]): boolean {
  return value >= range[0] && value <= range[1]
} 
export interface SafeStorageError {
  code: string
  safeMessage: string
  retryable: boolean
}

export interface AppInfo {
  name: string
  tagline: string
  version: string
  commit: string
  buildTime: string
  runtime: string
  status: string
  error?: SafeStorageError
}

export type LoadAppInfo = () => Promise<AppInfo>

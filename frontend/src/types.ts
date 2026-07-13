export interface AppInfo {
  name: string
  tagline: string
  version: string
  commit: string
  buildTime: string
  runtime: string
  status: string
}

export type LoadAppInfo = () => Promise<AppInfo>

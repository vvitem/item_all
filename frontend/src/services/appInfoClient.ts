import { GetAppInfo } from '../../wailsjs/go/main/App'
import type { AppInfo } from '../types'

export async function getAppInfo(): Promise<AppInfo> {
  return (await GetAppInfo()) as AppInfo
}

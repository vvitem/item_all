import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import App from './App'
import type { AppInfo } from './types'

const readyInfo: AppInfo = {
  name: 'ItemAll',
  tagline: 'Local-first SafeOps',
  version: '0.0.0-dev',
  commit: 'abc123def456',
  buildTime: '2026-07-13T08:00:00Z',
  runtime: 'go1.26.5 windows/amd64',
  status: 'ready',
}

describe('App', () => {
  it('shows loading while app information is pending', () => {
    const pending = new Promise<AppInfo>(() => undefined)

    render(<App loadAppInfo={() => pending} />)

    expect(screen.getByRole('status')).toHaveTextContent('正在读取应用信息')
  })

  it('shows trusted app information when loading succeeds', async () => {
    render(<App loadAppInfo={async () => readyInfo} />)

    expect(await screen.findByRole('heading', { name: 'ItemAll' })).toBeInTheDocument()
    expect(screen.getByText('Local-first SafeOps')).toBeInTheDocument()
    expect(screen.getByText('0.0.0-dev')).toBeInTheDocument()
    expect(screen.getByText('abc123def456')).toBeInTheDocument()
    expect(screen.getByText('2026-07-13T08:00:00Z')).toBeInTheDocument()
    expect(screen.getByText('go1.26.5 windows/amd64')).toBeInTheDocument()
    expect(screen.getByText('运行正常')).toBeInTheDocument()
  })

  it('shows a safe error without exposing the thrown message', async () => {
    render(<App loadAppInfo={async () => Promise.reject(new Error('C:\\Users\\secret\\token.txt'))} />)

    expect(await screen.findByRole('alert')).toHaveTextContent('无法读取应用信息')
    expect(screen.queryByText(/secret|token\.txt/i)).not.toBeInTheDocument()
  })

  it('retries after a failed request', async () => {
    const loadAppInfo = vi
      .fn<() => Promise<AppInfo>>()
      .mockRejectedValueOnce(new Error('first failure'))
      .mockResolvedValueOnce(readyInfo)

    render(<App loadAppInfo={loadAppInfo} />)

    const retryButton = await screen.findByRole('button', { name: '重新读取' })
    fireEvent.click(retryButton)

    await waitFor(() => expect(loadAppInfo).toHaveBeenCalledTimes(2))
    expect(await screen.findByRole('heading', { name: 'ItemAll' })).toBeInTheDocument()
  })
})

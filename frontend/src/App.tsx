import { useEffect, useState } from 'react'
import { getAppInfo } from './services/appInfoClient'
import type { AppInfo, LoadAppInfo } from './types'

type ViewState =
  | { kind: 'loading' }
  | { kind: 'ready'; info: AppInfo }
  | { kind: 'error' }

interface AppProps {
  loadAppInfo?: LoadAppInfo
}

export default function App({ loadAppInfo = getAppInfo }: AppProps) {
  const [attempt, setAttempt] = useState(0)
  const [state, setState] = useState<ViewState>({ kind: 'loading' })

  useEffect(() => {
    let active = true
    setState({ kind: 'loading' })

    loadAppInfo()
      .then((info) => {
        if (active) {
          setState({ kind: 'ready', info })
        }
      })
      .catch(() => {
        if (active) {
          setState({ kind: 'error' })
        }
      })

    return () => {
      active = false
    }
  }, [attempt, loadAppInfo])

  if (state.kind === 'loading') {
    return (
      <main className="shell shell--centered">
        <p role="status" className="status-card">
          正在读取应用信息
        </p>
      </main>
    )
  }

  if (state.kind === 'error') {
    return (
      <main className="shell shell--centered">
        <section role="alert" className="status-card status-card--error">
          <h1>无法读取应用信息</h1>
          <p>应用未能完成本地初始化，请重新读取。</p>
          <button type="button" onClick={() => setAttempt((value) => value + 1)}>
            重新读取
          </button>
        </section>
      </main>
    )
  }

  const { info } = state

  return (
    <main className="shell">
      <header className="hero">
        <div>
          <p className="eyebrow">Local desktop operations workspace</p>
          <h1>{info.name}</h1>
          <p className="tagline">{info.tagline}</p>
        </div>
        <span className="health-badge">
          {info.status === 'ready' ? '运行正常' : '状态异常'}
        </span>
      </header>

      <section className="metadata" aria-label="应用构建信息">
        <article>
          <span>Version</span>
          <strong>{info.version}</strong>
        </article>
        <article>
          <span>Commit</span>
          <strong>{info.commit}</strong>
        </article>
        <article>
          <span>Build Time</span>
          <strong>{info.buildTime}</strong>
        </article>
        <article>
          <span>Runtime</span>
          <strong>{info.runtime}</strong>
        </article>
      </section>
    </main>
  )
}

import { useState, useEffect, useCallback } from 'react'

export type TimeRange = '1h' | '6h' | '24h'

const RANGE_SECONDS: Record<TimeRange, number> = {
  '1h': 3600,
  '6h': 21600,
  '24h': 86400,
}

const REFRESH_INTERVAL = 30_000

interface UseMetricsResult {
  data: Record<string, any>[]
  loading: boolean
  range: TimeRange
  setRange: (r: TimeRange) => void
}

export function useMetrics(
  env: string | null,
  endpoint: string,
  params?: Record<string, string>
): UseMetricsResult {
  const [data, setData] = useState<Record<string, any>[]>([])
  const [loading, setLoading] = useState(true)
  const [range, setRange] = useState<TimeRange>('1h')

  const fetchData = useCallback(async () => {
    if (!env) return
    const now = Math.floor(Date.now() / 1000)
    const from = now - RANGE_SECONDS[range]

    const search = new URLSearchParams({ from: from.toString(), to: now.toString() })
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        if (v) search.set(k, v)
      }
    }

    try {
      const res = await fetch(`/api/environments/${env}/${endpoint}?${search}`)
      if (res.ok) {
        const json = await res.json()
        setData(json.points || [])
      }
    } catch {
      // ignore fetch errors
    }
    setLoading(false)
  }, [env, endpoint, range, params])

  useEffect(() => {
    setLoading(true)
    fetchData()
  }, [fetchData])

  useEffect(() => {
    if (!env) return
    const id = setInterval(fetchData, REFRESH_INTERVAL)
    return () => clearInterval(id)
  }, [env, fetchData])

  return { data, loading, range, setRange }
}

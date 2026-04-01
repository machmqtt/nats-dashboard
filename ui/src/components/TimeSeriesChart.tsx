import { memo } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from 'recharts'

export interface LineDef {
  key: string
  color: string
  label: string
}

interface Props {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: Record<string, any>[]
  lines: LineDef[]
  yFormatter?: (v: number) => string
  height?: number
}

function formatTime(ts: number): string {
  const d = new Date(ts * 1000)
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

export const TimeSeriesChart = memo(function TimeSeriesChart({ data, lines, yFormatter, height = 200 }: Props) {
  if (!data || data.length === 0) {
    return (
      <div className="flex items-center justify-center text-gray-400 text-sm" style={{ height }}>
        No data yet
      </div>
    )
  }

  const fmt = yFormatter || ((v: number) => v.toFixed(1))

  return (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={data} margin={{ top: 5, right: 10, bottom: 5, left: 0 }}>
        <CartesianGrid strokeDasharray="3 3" strokeOpacity={0.15} />
        <XAxis
          dataKey="ts"
          tickFormatter={formatTime}
          tick={{ fontSize: 11 }}
          stroke="currentColor"
          strokeOpacity={0.3}
        />
        <YAxis
          tickFormatter={fmt}
          tick={{ fontSize: 11 }}
          stroke="currentColor"
          strokeOpacity={0.3}
          width={50}
        />
        <Tooltip
          labelFormatter={(ts) => new Date(Number(ts) * 1000).toLocaleString()}
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          formatter={(value: any, name: any) => {
            const line = lines.find((l) => l.key === name)
            return [fmt(Number(value)), line?.label || String(name)]
          }}
          contentStyle={{
            backgroundColor: 'var(--tooltip-bg, #1f2937)',
            border: 'none',
            borderRadius: '6px',
            fontSize: '12px',
            color: 'var(--tooltip-text, #e5e7eb)',
          }}
        />
        {lines.map((l) => (
          <Line
            key={l.key}
            type="monotone"
            dataKey={l.key}
            stroke={l.color}
            strokeWidth={1.5}
            dot={false}
            name={l.key}
          />
        ))}
      </LineChart>
    </ResponsiveContainer>
  )
})

import { memo } from 'react'
import { useStore } from '../store/store'
import { X } from 'lucide-react'

const colors = {
  info: 'bg-blue-500',
  error: 'bg-red-500',
  success: 'bg-green-500',
}

export const Toasts = memo(function Toasts() {
  const { toasts, removeToast } = useStore()

  if (toasts.length === 0) return null

  return (
    <div className="fixed bottom-4 right-4 z-[100] flex flex-col gap-2">
      {toasts.map((t) => (
        <div
          key={t.id}
          className={`${colors[t.type]} text-white px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 min-w-[280px] animate-[slideIn_0.2s_ease-out]`}
        >
          <span className="flex-1 text-sm">{t.message}</span>
          <button onClick={() => removeToast(t.id)} className="text-white/70 hover:text-white">
            <X className="w-4 h-4" />
          </button>
        </div>
      ))}
    </div>
  )
})

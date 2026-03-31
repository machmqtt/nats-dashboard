import { useState, useEffect, useCallback } from 'react'

export interface User {
  id: number
  username: string
  role: 'admin' | 'viewer'
}

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const checkSession = useCallback(async () => {
    try {
      const res = await fetch('/api/me')
      if (res.ok) {
        setUser(await res.json())
      } else {
        setUser(null)
      }
    } catch {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { checkSession() }, [checkSession])

  const login = async (username: string, password: string) => {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    })
    if (!res.ok) throw new Error('Invalid credentials')
    const u = await res.json()
    setUser(u)
    return u
  }

  const logout = async () => {
    await fetch('/api/logout', { method: 'POST' })
    setUser(null)
  }

  return { user, loading, login, logout }
}

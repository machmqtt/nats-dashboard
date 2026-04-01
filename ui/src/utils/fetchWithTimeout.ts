const DEFAULT_TIMEOUT = 15_000

export function fetchWithTimeout(
  url: string,
  opts?: RequestInit & { timeout?: number },
): Promise<Response> {
  const { timeout = DEFAULT_TIMEOUT, ...fetchOpts } = opts || {}
  const controller = new AbortController()
  const id = setTimeout(() => controller.abort(), timeout)
  return fetch(url, { ...fetchOpts, signal: controller.signal }).finally(() =>
    clearTimeout(id),
  )
}

export type CalculateOutcome = { result: number } | { error: string }

export const CONNECTION_ERROR = 'Could not reach the server'

// Rejects only when the request was aborted through signal; every other
// failure is reported as an outcome so callers can display it.
export async function calculate(
  expression: string,
  signal: AbortSignal,
): Promise<CalculateOutcome> {
  try {
    const response = await fetch('/api/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ expression }),
      signal,
    })
    const body: unknown = await response.json()
    if (typeof body === 'object' && body !== null) {
      if (response.ok && 'result' in body && typeof body.result === 'number') {
        return { result: body.result }
      }
      if ('error' in body && typeof body.error === 'string') {
        return { error: body.error }
      }
    }
  } catch (err) {
    if (signal.aborted) throw err
  }
  return { error: CONNECTION_ERROR }
}

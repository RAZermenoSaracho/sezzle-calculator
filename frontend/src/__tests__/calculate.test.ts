import { afterEach, expect, it, vi } from 'vitest'
import { calculate, CONNECTION_ERROR } from '../calculate'

function stubFetch(impl: () => Promise<Response>) {
  const fetchMock = vi.fn(impl)
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status })

afterEach(() => vi.unstubAllGlobals())

it('posts the expression and returns the result', async () => {
  const fetchMock = stubFetch(async () => json({ result: 7 }))

  const outcome = await calculate('1 + 2 * 3', new AbortController().signal)

  expect(outcome).toEqual({ result: 7 })
  const [url, init] = fetchMock.mock.calls[0] as unknown as [string, RequestInit]
  expect(url).toBe('/api/calculate')
  expect(init.method).toBe('POST')
  expect(JSON.parse(init.body as string)).toEqual({ expression: '1 + 2 * 3' })
})

it('returns the API error message', async () => {
  stubFetch(async () => json({ error: 'division by zero' }, 400))

  expect(await calculate('1/0', new AbortController().signal)).toEqual({
    error: 'division by zero',
  })
})

it('reports a connection error for non-JSON responses', async () => {
  stubFetch(async () => new Response('<html>Bad Gateway</html>', { status: 502 }))

  expect(await calculate('1', new AbortController().signal)).toEqual({
    error: CONNECTION_ERROR,
  })
})

it('reports a connection error when the network fails', async () => {
  stubFetch(async () => {
    throw new TypeError('Failed to fetch')
  })

  expect(await calculate('1', new AbortController().signal)).toEqual({
    error: CONNECTION_ERROR,
  })
})

it('rejects when the request was aborted', async () => {
  const controller = new AbortController()
  stubFetch(async () => {
    controller.abort()
    throw new DOMException('Aborted', 'AbortError')
  })

  await expect(calculate('1', controller.signal)).rejects.toThrow()
})

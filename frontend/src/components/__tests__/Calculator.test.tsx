import { act, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { CONNECTION_ERROR } from '../../calculate'
import Calculator from '../Calculator'

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status })

// Fetch mock that deliberately ignores abort signals, so late replies must be
// discarded by the component itself.
function stubFetch(impl: (expression: string) => Promise<Response>) {
  const fetchMock = vi.fn((_url: string, init: RequestInit) =>
    impl(JSON.parse(init.body as string).expression),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

// Types one character at a time, like a user. fireEvent is used instead of
// user-event because user-event's async helpers hang under fake timers.
const user = {
  type(input: HTMLElement, text: string) {
    for (let i = 1; i <= text.length; i++) {
      fireEvent.change(input, { target: { value: text.slice(0, i) } })
    }
  },
  clear(input: HTMLElement) {
    fireEvent.change(input, { target: { value: '' } })
  },
}

function setup() {
  render(<Calculator />)
  return { user, input: screen.getByLabelText('Expression') }
}

const settle = (ms = 300) => act(() => vi.advanceTimersByTimeAsync(ms))

beforeEach(() => vi.useFakeTimers())
afterEach(() => {
  vi.useRealTimers()
  vi.unstubAllGlobals()
})

it('renders a focused input and no button', () => {
  stubFetch(async () => json({ result: 0 }))
  const { input } = setup()

  expect(input).toHaveFocus()
  expect(screen.queryByRole('button')).not.toBeInTheDocument()
})

it('sends one request with the final text after typing stops', async () => {
  const fetchMock = stubFetch(async () => json({ result: 7 }))
  const { user, input } = setup()

  user.type(input, '1+2*3')
  expect(fetchMock).not.toHaveBeenCalled()
  await settle()

  expect(fetchMock).toHaveBeenCalledTimes(1)
  expect(JSON.parse(fetchMock.mock.calls[0][1].body as string)).toEqual({
    expression: '1+2*3',
  })
  expect(screen.getByRole('status')).toHaveTextContent('7')
})

it('hides floating point noise', async () => {
  stubFetch(async () => json({ result: 0.30000000000000004 }))
  const { user, input } = setup()

  user.type(input, '0.1+0.2')
  await settle()

  expect(screen.getByRole('status')).toHaveTextContent(/^0\.3$/)
})

it('shows the API error message', async () => {
  stubFetch(async () => json({ error: 'division by zero' }, 400))
  const { user, input } = setup()

  user.type(input, '1/0')
  await settle()

  expect(screen.getByRole('status')).toHaveTextContent('division by zero')
})

it('shows a connection error when the request fails', async () => {
  stubFetch(async () => {
    throw new TypeError('Failed to fetch')
  })
  const { user, input } = setup()

  user.type(input, '1+1')
  await settle()

  expect(screen.getByRole('status')).toHaveTextContent(CONNECTION_ERROR)
})

it('clears the output and sends no request for empty input', async () => {
  const fetchMock = stubFetch(async () => json({ result: 2 }))
  const { user, input } = setup()

  user.type(input, '1+1')
  await settle()
  expect(screen.getByRole('status')).toHaveTextContent('2')

  user.clear(input)
  await settle()

  expect(screen.getByRole('status')).toBeEmptyDOMElement()
  expect(fetchMock).toHaveBeenCalledTimes(1)
})

it('never lets an older response replace a newer one', async () => {
  const replies: Record<string, (r: Response) => void> = {}
  stubFetch(
    (expression) =>
      new Promise((resolve) => {
        replies[expression] = resolve
      }),
  )
  const { user, input } = setup()

  user.type(input, '1')
  await settle()
  user.type(input, '1+1')
  await settle()

  await act(async () => replies['1+1'](json({ result: 2 })))
  await act(async () => replies['1'](json({ result: 1 })))

  expect(screen.getByRole('status')).toHaveTextContent(/^2$/)
})

it('is a multi-line field that ignores Enter and flattens pasted newlines', () => {
  stubFetch(async () => json({ result: 3 }))
  const { input } = setup()

  expect(input.tagName).toBe('TEXTAREA')
  expect(fireEvent.keyDown(input, { key: 'Enter' })).toBe(false)

  fireEvent.change(input, { target: { value: '1 +\n2' } })
  expect(input).toHaveValue('1 + 2')
})

it('sends long expressions unchanged', async () => {
  const long = Array(40).fill('12345').join('+')
  const fetchMock = stubFetch(async () => json({ result: 1 }))
  const { input } = setup()

  fireEvent.change(input, { target: { value: long } })
  await settle()

  expect(JSON.parse(fetchMock.mock.calls[0][1].body as string)).toEqual({
    expression: long,
  })
})

it('shows errors in a color distinct from results', async () => {
  stubFetch(async (expression) =>
    expression === '1/0' ? json({ error: 'division by zero' }, 400) : json({ result: 1 }),
  )
  const { user, input } = setup()

  user.type(input, '1')
  await settle()
  const resultClass = screen.getByRole('status').className

  user.type(input, '1/0')
  await settle()
  const errorClass = screen.getByRole('status').className

  expect(resultClass).toContain('text-green-400')
  expect(errorClass).toContain('text-red-400')
  expect(errorClass).not.toContain('green')
})

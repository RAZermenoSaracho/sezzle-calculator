import { useEffect, useLayoutEffect, useRef, useState } from 'react'
import { calculate } from '../calculate'

const DEBOUNCE_MS = 300

type Output = { text: string; isError: boolean }

// Hides binary floating-point noise (0.1 + 0.2) without changing semantics.
function formatResult(value: number): string {
  return Number(value.toPrecision(12)).toString()
}

function Calculator() {
  const [expression, setExpression] = useState('')
  const [output, setOutput] = useState<Output | null>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  // Grow the textarea to fit its wrapped content.
  useLayoutEffect(() => {
    const el = inputRef.current
    if (!el) return
    el.style.height = 'auto'
    el.style.height = `${el.scrollHeight}px`
  }, [expression])

  useEffect(() => {
    if (expression.trim() === '') return

    const controller = new AbortController()
    const timer = setTimeout(async () => {
      try {
        const outcome = await calculate(expression, controller.signal)
        // The effect cleanup aborts superseded requests; ignore their late replies.
        if (controller.signal.aborted) return
        setOutput(
          'result' in outcome
            ? { text: formatResult(outcome.result), isError: false }
            : { text: outcome.error, isError: true },
        )
      } catch {
        // Aborted: a newer expression owns the output.
      }
    }, DEBOUNCE_MS)

    return () => {
      clearTimeout(timer)
      controller.abort()
    }
  }, [expression])

  const shown = expression.trim() === '' ? null : output

  return (
    <main className="flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-xl">
        <textarea
          ref={inputRef}
          rows={1}
          className="peer block w-full resize-none overflow-hidden bg-transparent pb-3 text-2xl leading-tight text-green-400 caret-green-400 outline-none placeholder:text-green-900 sm:text-4xl"
          value={expression}
          // The expression is a single logical line; newlines only wrap visually.
          onChange={(e) => setExpression(e.target.value.replace(/\r?\n/g, ' '))}
          onKeyDown={(e) => {
            if (e.key === 'Enter') e.preventDefault()
          }}
          placeholder="1 + 2 * 3"
          aria-label="Expression"
          autoFocus
          autoComplete="off"
          autoCapitalize="off"
          spellCheck={false}
        />
        <div className="border-t border-green-900 peer-focus:border-green-500" />
        <output
          aria-live="polite"
          className={`mt-3 block min-h-10 text-2xl sm:text-4xl ${
            shown?.isError ? 'text-red-400' : 'text-green-400'
          }`}
        >
          {shown?.text}
        </output>
      </div>
    </main>
  )
}

export default Calculator

import { render, screen, within } from '@testing-library/react'
import { expect, it } from 'vitest'
import SyntaxHelp from '../SyntaxHelp'

function renderHelp() {
  render(<SyntaxHelp />)
  return within(screen.getByRole('region', { name: 'Supported syntax' }))
}

it('lists every supported operation', () => {
  const help = renderHelp()

  for (const symbol of ['+ -', '* /', '^', 'sqrt(x)', 'x%', '( )', '+x -x', '1.5 .5']) {
    expect(help.getByText(symbol)).toBeInTheDocument()
  }
})

it('explains chained signs and that signs alone are not an expression', () => {
  const help = renderHelp()

  expect(help.getByRole('heading', { name: /allowed\?/ })).toHaveTextContent(
    'Why are ++++ and ---- allowed?',
  )
  expect(help.getByText(/A sign in front of a value/)).toHaveTextContent(
    /\+\+\+\+5 is valid, but \+\+\+\+ on its own is not/,
  )
})

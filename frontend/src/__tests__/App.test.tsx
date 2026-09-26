import { render, screen } from '@testing-library/react'
import { expect, it } from 'vitest'
import App from '../App'

it('shows the calculator with the syntax help below it', () => {
  render(<App />)

  const input = screen.getByLabelText('Expression')
  const help = screen.getByRole('region', { name: 'Supported syntax' })

  expect(input).toBeInTheDocument()
  expect(
    input.compareDocumentPosition(help) & Node.DOCUMENT_POSITION_FOLLOWING,
  ).toBeTruthy()
})

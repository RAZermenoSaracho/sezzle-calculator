const syntax = [
  ['+ -', 'addition and subtraction'],
  ['* /', 'multiplication and division'],
  ['^', 'exponentiation'],
  ['sqrt(x)', 'square root'],
  ['x%', 'percentage (divides by 100)'],
  ['( )', 'grouping'],
  ['+x -x', 'unary plus and minus'],
  ['1.5 .5', 'decimal numbers'],
]

function Code({ children }: { children: string }) {
  return <code className="font-mono text-green-500">{children}</code>
}

function SyntaxHelp() {
  return (
    <section
      aria-label="Supported syntax"
      className="w-full max-w-xl pb-10 text-sm text-neutral-400"
    >
      <h2 className="mb-2 text-xs font-medium tracking-wider text-neutral-500 uppercase">
        Supported syntax
      </h2>
      <dl className="grid grid-cols-[7rem_1fr] gap-x-4 gap-y-1">
        {syntax.map(([symbol, meaning]) => (
          <div key={symbol} className="contents">
            <dt>
              <Code>{symbol}</Code>
            </dt>
            <dd>{meaning}</dd>
          </div>
        ))}
      </dl>

      <h2 className="mt-6 mb-2 text-xs font-medium tracking-wider text-neutral-500 uppercase">
        Why are <span className="font-mono normal-case">++++</span> and{' '}
        <span className="font-mono normal-case">----</span> allowed?
      </h2>
      <p className="leading-relaxed">
        A sign in front of a value is an operator, and operators can be stacked.{' '}
        <Code>+</Code> leaves a value unchanged and <Code>-</Code> flips its sign,
        so <Code>----5</Code> is 5, <Code>---5</Code> is -5, and{' '}
        <Code>5 + ++++5</Code> is 10. The signs must be followed by a value:{' '}
        <Code>++++5</Code> is valid, but <Code>++++</Code> on its own is not, and
        you will see “invalid expression” until you add a number.
      </p>
    </section>
  )
}

export default SyntaxHelp

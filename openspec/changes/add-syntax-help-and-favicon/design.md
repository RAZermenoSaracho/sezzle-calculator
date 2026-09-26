# Design: Syntax Help and Favicon

## Accuracy source

Content mirrors `backend/internal/calculator/parser.go`. Verified against the running API:

| Input | Result |
|---|---|
| `++++5`, `----5` | `5` |
| `---5` | `-5` |
| `5 + ++++5` | `10` |
| `54+++++6666` | `6720` (first `+` is the binary operator, four are signs) |
| `++++`, `----`, `+` | `invalid expression` |

The grammar is `unary := ('+' | '-') unary | power`, so signs chain, but the chain must end in a value. A sign-only string is incomplete, not valid. While typing, the live UI may briefly show `invalid expression` for such input; that is the backend rejecting an incomplete expression.

## Layout

`Calculator` becomes only the calculator block (`w-full max-w-xl`). `App` owns the page layout: a column where the calculator sits in a vertically centered flexible region and `SyntaxHelp` follows below it in a muted, small style. On narrow screens the page scrolls; there is no horizontal overflow.

## SyntaxHelp component

Two areas in `src/components/SyntaxHelp.tsx`, Tailwind only, muted neutral text with dim green code samples (no new colors, no animation, no dependency):

1. Supported syntax as a definition list: `+ -`, `* /`, `^`, `sqrt(x)`, `x%`, `( )`, `+x -x`, decimals.
2. "Why are `++++` and `----` allowed?" with this text:

> A sign in front of a value is an operator, and operators can be stacked. `+` leaves a value unchanged and `-` flips its sign, so `----5` is 5 and `---5` is -5, and `5 + ++++5` is 10. The signs must be followed by a value: `++++5` is valid, but `++++` on its own is not, and you will see "invalid expression" until you add a number.

## Favicon

`public/favicon.svg`, already linked from `index.html` (`<link rel="icon" type="image/svg+xml" href="/favicon.svg">`). A dark rounded square with a green plus above an equals sign, drawn with basic shapes only, echoing the two calculator keys and the green-on-black theme.

## Testing

Component test asserts the operations and the key explanation claims (valid chained example, bare signs invalid). App test asserts calculator and help render together. No prose duplication beyond that.

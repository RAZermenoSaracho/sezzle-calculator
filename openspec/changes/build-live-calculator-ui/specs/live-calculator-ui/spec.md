## ADDED Requirements

### Requirement: Minimal dark calculator layout
The UI SHALL present a centered, dark-only layout consisting of an expression input without a four-sided border, a thin horizontal divider, and a result/error area below the divider, with no Calculate button, navigation, cards, gradients, animations, or theme switcher.

#### Scenario: Initial render
- **WHEN** the page loads
- **THEN** the input is focused, the divider is visible beneath it, the output area is empty, and no button is present

#### Scenario: Responsive
- **WHEN** the viewport is 360 px wide or 1280 px wide
- **THEN** the calculator remains centered, readable, and free of horizontal scrolling

### Requirement: Live debounced calculation
The UI SHALL send the current expression to `POST /api/calculate` automatically 300 ms after the user stops typing, without requiring Enter or any click.

#### Scenario: Debounce
- **WHEN** the user types `1+2*3` as several rapid keystrokes
- **THEN** exactly one request is sent, after the last keystroke's debounce delay, containing `1+2*3`

#### Scenario: Result display
- **WHEN** the API returns `{"result": 7}`
- **THEN** `7` is shown below the divider

#### Scenario: Float noise hidden
- **WHEN** the API returns `0.30000000000000004`
- **THEN** `0.3` is shown

### Requirement: Error display
The UI SHALL show the `error` message returned by the API below the divider in a visually distinct error style, and SHALL show `Could not reach the server` when the request fails or the response is unusable.

#### Scenario: API error
- **WHEN** the API returns 400 with `{"error":"division by zero"}`
- **THEN** `division by zero` is shown as an error

#### Scenario: Network failure
- **WHEN** the request fails with a network error
- **THEN** `Could not reach the server` is shown

### Requirement: Empty input
The UI SHALL clear the output and send no request when the input is empty or whitespace-only.

#### Scenario: Clearing input
- **WHEN** the user deletes all text
- **THEN** the output is cleared and no request is sent

### Requirement: Stale response protection
The UI SHALL NOT let a response for an earlier expression replace the output for a later expression, regardless of the order in which responses arrive.

#### Scenario: Out-of-order responses
- **WHEN** the response for an older expression resolves after the response for a newer expression
- **THEN** the output shows only the newer expression's result

### Requirement: Backend is the source of truth
The frontend SHALL NOT parse or evaluate expressions and SHALL communicate only through the backend API using a relative `/api` URL.

#### Scenario: No local evaluation
- **WHEN** the backend is unreachable
- **THEN** the UI shows the connectivity error rather than a locally computed value

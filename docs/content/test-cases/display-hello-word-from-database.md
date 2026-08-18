# Test Cases — Display Hello Word from database

Module: `content`
Function: Display Hello Word from database
Source: `docs/content/SRS.md` §4.1, acceptance criteria AC-1 … AC-8
Risk level: **low-medium** — this screen is the entire product, and the value must
prove it originates from the database, so AC-2 is verified by source inspection in
addition to runtime checks.

---

## Acceptance criteria coverage

### AC-1 — page renders exactly `Hello Word`, centred, dark on white, no animation

**Scenario**: Site loads and displays the database value
**Given**: the database is seeded with a row whose text value is `Hello Word`
**When**: the Guest opens the site and the fetch completes
**Then**: the page shows exactly the string `Hello Word`; the text is centred on a
`#FFFFFF` background with dark text (`#111827` per the design system) and no
animation on the rendered state

Traceability: CONTENT-001 behaviour 3–4, AC-1

### AC-2 — value originates from the database, not hardcoded in the frontend

**Scenario**: Frontend source does not hardcode the value
**Given**: the frontend source is available for inspection
**When**: the frontend source code is inspected
**Then**: the string `Hello Word` appears nowhere in the frontend source
(`code/frontend/`); the frontend fetches the value from the backend API and
renders whatever the API returns

Traceability: CONTENT-001 behaviour 4, AC-2

### AC-3 — loading state shown while fetch is in progress

**Scenario**: Loading state on first load
**Given**: the fetch is in progress (e.g. the API response is artificially delayed
in the test environment)
**When**: the Guest loads the site
**Then**: a loading state is shown (not a blank screen and not the loaded text) and
persists until the fetch completes

Traceability: CONTENT-002 behaviour 1, 3; AC-3

### AC-4 — loading state replaced by loaded text when fetch succeeds

**Scenario**: Loading state resolves to loaded text
**Given**: the fetch succeeds and returns the value `Hello Word`
**When**: the fetch completes
**Then**: the loading state is replaced by the text `Hello Word`; no loading
indicator remains on screen

Traceability: CONTENT-002 behaviour 2, AC-4

### AC-5 — error state shown when fetch fails

**Scenario**: Failed fetch shows an error state, not a blank screen
**Given**: the backend is unavailable or the API returns a non-success status
**When**: the Guest loads or reloads the site
**Then**: an error state is shown with a way to retry; the screen is not blank and
no partial or cached text is rendered

Traceability: CONTENT-003 behaviour 1–2, AC-5; failure behaviour "Upstream failure"

### AC-6 — retry that succeeds replaces the error state with the loaded text

**Scenario**: Retry after a failure succeeds
**Given**: the fetch failed and the error state with a retry control is showing
**When**: the Guest triggers retry and the retry fetch succeeds
**Then**: the error state is replaced by the loaded text `Hello Word`

Traceability: CONTENT-003 behaviour 3, AC-6

### AC-7 — no row in the database results in an error state, not a blank screen

**Scenario**: Database holds no row for the text
**Given**: the database has no row for the text (e.g. the table is empty)
**When**: the Guest loads the site
**Then**: an error state with a retry control is shown, not a blank screen and not
the loaded text

Traceability: AC-7, failure behaviour "Not found"

### AC-8 — empty stored text renders as an error state, not as loaded success

**Scenario**: Stored text is an empty string
**Given**: the database row for the text stores an empty string
**When**: the Guest loads the site
**Then**: an error state with a retry control is shown; the empty string is not
rendered as if it were the loaded text

Traceability: AC-8, failure behaviour "Empty value"

---

## Boundary and permission coverage

### Slow response — loading persists until the fetch completes

**Scenario**: Slow but eventually successful fetch
**Given**: the API responds slowly but eventually succeeds
**When**: the Guest loads the site
**Then**: the loading state persists for the whole fetch; once the response
arrives, the loaded text `Hello Word` is shown

Traceability: failure behaviour "Slow response"

### Permission — retry is always allowed

**Scenario**: Guest retries without any permission check
**Given**: the error state with a retry control is showing
**When**: the Guest triggers retry
**Then**: the retry is performed with no sign-in and no permission gate — there is
no permission check in this module, and any Guest may retry

Traceability: failure behaviour "Permission"; §2 actor "Guest"

---

## Manual checks

### M-1 — visual presentation (manual)

**Scenario**: Visual layout matches the design
**Given**: the fetch succeeded and the page shows `Hello Word`
**When**: the page is inspected visually at widths from 320px up
**Then**: the text is centred horizontally and vertically on a white background
with dark text, there is no animation, and no horizontal page scroll occurs

Reason manual: appearance, contrast and scroll behaviour are visual properties no
automated assertion in this repo checks.

### M-2 — retry control focus visibility (manual)

**Scenario**: Retry control is focusable and visibly focused
**Given**: the error state with a retry control is showing
**When**: the Guest tabs to the retry control
**Then**: the retry control receives keyboard focus and shows a visible focus
indicator

Reason manual: focus visibility is a rendered, browser-level behaviour not covered
by an automated check here.

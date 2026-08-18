# Story — Display Hello Word from database

- Module: `content`
- Plan item: Display Hello Word from database
- Requirements: CONTENT-001, CONTENT-002, CONTENT-003
- Test cases: `docs/content/test-cases/display-hello-word.md`

## 1. User story

As a Guest, I want to load the site and see the text `Hello Word`, so that I
receive the product's only message, which originates from the PostgreSQL
database rather than the frontend.

## 2. In scope

- A Go API endpoint that reads the content text from PostgreSQL and returns it.
- A frontend that fetches that value on load and renders it centred, dark text
  on white, no animation.
- Three UI states: loading (before the fetch completes), loaded (the fetched
  text), error (fetch failed, with a retry control).
- Backend self-migration and a seed row whose text is exactly `Hello Word`
  (seeding is a migration per `architecture/overview.md` §8).
- The frontend must not hardcode the string `Hello Word`; it arrives from the API.

## 3. Out of scope

- Editing or writing the text — it is fixed to `Hello Word`.
- Authentication, accounts, or any permission check — the only actor is a Guest.
- Any second screen or navigation — the product is a single screen.
- The design's `Loading` / `Loaded` / `Empty` / `Error` demo pills and caption —
  preview-only controls, not shipped (`design/design-system.md` Known deviations
  and project memory `design.preview_note`).
- A separate empty-value UI state — an empty or absent stored value renders as
  the error state per SRS §4.

## 4. UI scope

Single content screen, matching the approved design. It touches the `Status
panel` component and the `Button` component (primary, for retry).

| State | What the Guest sees | Tokens |
|---|---|---|
| Loading | "Loading" plus three pulsing dots | `--color-text-muted`, `--duration-pulse` |
| Loaded | `h1` showing the fetched text `Hello Word` | `--color-text`, `--text-display` |
| Error | "Could not load the text from the database." plus a retry button | `--color-danger`, `--color-primary` |

The panel carries `role="status"` and `aria-live="polite"`; the error message
carries `role="alert"`. State switches toggle a `.hidden` class on the state
nodes. The retry button is `type="button"` with a visible focus ring.

## 5. Acceptance criteria

1. Given the database holds the value `Hello Word`, when the Guest loads the
   site, then the page shows exactly `Hello Word`, centred, dark text on white,
   no animation.
2. Given the page is loaded, when the frontend source is inspected, then the
   string `Hello Word` is not hardcoded in frontend source; it arrives from the API.
3. Given the fetch is in progress, when the Guest loads the site, then a loading
   state is shown until the fetch completes.
4. Given the fetch succeeds, when the fetch completes, then the loading state is
   replaced by the loaded text.
5. Given the fetch fails (backend unavailable or non-success status), when the
   Guest loads or reloads, then an error state is shown, not a blank screen.
6. Given the fetch failed and an error state is showing, when the Guest triggers
   retry and the fetch succeeds, then the error state is replaced by the loaded text.
7. Given the database has no row for the text, when the Guest loads the site,
   then an error state with a retry control is shown, not a blank screen.
8. Given the stored text is an empty string, when the Guest loads the site, then
   an error state with a retry control is shown; the empty string is not rendered
   as if it were the loaded text.

## 6. Dependencies

- Go backend and PostgreSQL, for reading and serving the text; API contract in
  `docs/architecture/services.md` (TL).
- Seed migration inserting the single `Hello Word` row (TL, backend).
- `code/frontend/app/page.tsx` composition root and `code/frontend/app/globals.css`
  tokens (TL scaffold) — this story mounts one component into `page.tsx`.

## 7. Notes for implementation

- The frontend is built and reviewed on mock data first; `lib/mock/` is deleted
  when the real API replaces it.
- Backend must never return an empty string as success; treat no-row and empty
  value as an error the frontend maps to the error state.

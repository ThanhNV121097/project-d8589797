# Story — Display Hello Word from database

Module: `content`
Plan item: Display Hello Word from database
Requirement ids: CONTENT-001, CONTENT-002, CONTENT-003

## User story

As a Guest, I want to load the site and see the text `Hello Word` fetched from
the database, so that I receive the product's only message, which originates
from the database and is never hardcoded in the frontend.

## In scope

- A Go backend endpoint that reads the single content row from PostgreSQL and
  returns it. Contract per `docs/architecture/services.md`: `GET /v1/content`
  with the §2.3 error envelope (`NOT_FOUND`, `UNAVAILABLE`, `INTERNAL`).
- A Next.js frontend component that fetches the value from the backend and
  renders it centred on a white background with dark text and no animation.
- Three UI states: loading (dots pulse), loaded (h1 with the fetched value),
  and error (message plus a retry control).
- The seed data: a migration inserts the single row whose value is exactly
  `Hello Word`.

## Out of scope

- Editing or writing the text — the value is fixed to `Hello Word`.
- Authentication or accounts — the only actor is a Guest; no sign-in.
- Any second screen, navigation, or routing — the product is a single screen.
- Caching or offline persistence of the value.
- The preview-only `Simulate error` / demo pills from the mockup — build the
  four states, do not ship the controls.

## UI scope

Touches the single content screen from the approved design: the status panel
with its three shipped states (loading, loaded, error). The error state adds a
retry control per CONTENT-003. The `Empty` demo state is not shipped — an empty
stored value renders as the error state.

## Acceptance criteria

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | Database holds the value `Hello Word` | Guest loads the site | Page shows exactly `Hello Word`, centred, dark text on white, no animation |
| AC-2 | Database holds the value `Hello Word` | Frontend source is inspected | The string `Hello Word` is not hardcoded in the frontend; it arrives from the API |
| AC-3 | Fetch is in progress | Guest loads the site | Loading state shown until the fetch completes |
| AC-4 | Fetch succeeds | Fetch completes | Loading state replaced by the loaded text |
| AC-5 | Fetch fails (backend unavailable or returns an error) | Guest loads or reloads | Error state shown, not a blank screen |
| AC-6 | Fetch failed, error state showing | Guest triggers retry and fetch succeeds | Error state replaced by the loaded text |
| AC-7 | Database has no row for the text | Guest loads the site | Error state with retry shown, not a blank screen |
| AC-8 | Stored text is an empty string | Guest loads the site | Error state with retry shown; empty string not rendered as success |

## Dependencies

- Backend endpoint `GET /v1/content` and `/healthz` — defined in
  `docs/architecture/services.md`, implemented in the BE PR.
- PostgreSQL seeded via the migration that inserts the single `Hello Word` row
  (matches `docs/architecture/erd.md` §3.1: table `contents`, column `value`,
  singleton index `uq_contents_singleton`, constraint `ck_contents_value_not_blank`).
- Frontend reads the API base from `NEXT_PUBLIC_API_URL`, falling back to
  `/api` (Next proxy strips the `/api` prefix before the backend).

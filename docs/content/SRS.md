# SRS — Content

Module: `content`
Last updated: 2026-05-27
Design: [View the approved design](http://localhost:8080/design/d8589797-aee6-42c3-94c6-db252a05e886)
Design system: `design/design-system.md`

> One file per module, at `docs/{module}/SRS.md`. It covers only the functions
> that belong to this module. Never write `docs/SRS.md`.

## 1. Purpose

This module renders the site's single screen of content. It displays the exact
text `Hello Word` fetched from a PostgreSQL database. A guest visiting the site
sees the text; the product delivers nothing without it, because this screen is
the entire product.

## 2. Actors

| Actor | Who they are | What they may do in this module |
|---|---|---|
| Guest | Anyone visiting the site, not signed in | View the displayed text; trigger a reload of the text |

There is no sign-in, no account, and no write action in this module.

## 3. Scope

**In scope** — the functions specified below, by their plan titles:

- Display Hello Word from database

**Out of scope** — name what a reader would reasonably expect here and say
where it lives instead.

- Editing the displayed text — deliberately not built; the text is fixed to `Hello Word`.
- Authentication and accounts — no module needs them; the only actor is a Guest.
- Any second screen or navigation — the product is a single screen.

## 4. Functional requirements

### 4.1 Display Hello Word from database

**Requirement CONTENT-001 — render the database value**

*As a* Guest, *I want to* load the site and see the text `Hello Word`, *so that*
I receive the product's only message, which originates from the database.

Behaviour:

1. The Guest opens the site. The screen renders a loading state while the text
   is being fetched.
2. The text is read from the PostgreSQL database and served to the frontend
   through an API.
3. The frontend displays the fetched value centred on a white background with
   dark text and no animation.
4. The displayed value must equal exactly `Hello Word` — the value comes from the
   database, never hardcoded in the frontend.

**Requirement CONTENT-002 — loading state**

*As a* Guest, *I want to* see a loading state while the text is fetched, *so
that* I know the page is working rather than blank.

Behaviour:

1. On first load, before the text arrives, the screen shows a loading state.
2. The loading state is replaced by the loaded text once the fetch completes.
3. The loading state never ships a control to trigger it — it appears only
   during an actual fetch.

**Requirement CONTENT-003 — error state and retry**

*As a* Guest, *I want to* see an error state and be able to retry when the text
cannot be fetched, *so that* a failed fetch is not a blank or broken screen.

Behaviour:

1. If the text cannot be fetched, the screen shows an error state.
2. The error state offers a way to retry the fetch.
3. A retry that succeeds replaces the error state with the loaded text.

**Acceptance criteria** — each maps one-to-one onto a test case in
`docs/content/test-cases/display-hello-word.md`.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | Database holds the value `Hello Word` | The Guest loads the site | The page shows exactly `Hello Word`, centred, dark text on white, no animation |
| AC-2 | Database holds the value `Hello Word` | The frontend is inspected | The string `Hello Word` is not hardcoded in the frontend source; it arrives from the API |
| AC-3 | The fetch is in progress | The Guest loads the site | A loading state is shown until the fetch completes |
| AC-4 | The fetch succeeds | The fetch completes | The loading state is replaced by the loaded text |
| AC-5 | The fetch fails (backend unavailable or returns an error) | The Guest loads or reloads | An error state is shown, not a blank screen |
| AC-6 | The fetch failed and an error state is showing | The Guest triggers retry and the fetch succeeds | The error state is replaced by the loaded text |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Not found | The database has no row for the text | Error state with a retry; not a blank screen |
| Empty value | The stored text is empty | Error state with a retry; the screen must not render an empty string as if it were success |
| Upstream failure | Database or backend is unavailable, or the API returns a non-success status | Error state with a retry; no partial or cached text is shown |
| Slow response | The fetch is slow but eventually succeeds | Loading state persists until the fetch completes; then the loaded text |
| Permission | Guest triggers a retry | Always allowed; there is no permission check in this module |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| text | text | yes | Non-empty; value is exactly `Hello Word` |

## 5. Screens

The design is the source of truth for appearance; this section maps functions
onto it so nothing in the design is unaccounted for and nothing specified here
is missing from the design.

| Screen | Section in the design | Functions it serves | States that must exist |
|---|---|---|---|
| Single content screen | Centred "Hello Word" | CONTENT-001, CONTENT-002, CONTENT-003 | loading, loaded, empty, error |

The design's `Loading` / `Loaded` / `Empty` / `Error` demo pills and caption are
preview-only controls: the four states are built, the controls are not shipped.

## 6. Non-functional requirements

| Area | Requirement |
|---|---|
| Performance | The loaded text renders within 2s on a typical connection after the fetch completes |
| Accessibility | Text is real text, not an image; visible focus on any retry control; contrast ≥ 4.5:1 |
| Responsive | Works at 320px and up; no horizontal page scroll |
| Localisation | Copy is in English; the displayed value is exactly `Hello Word` |

## 7. Dependencies and assumptions

- **Depends on:** the Go backend and PostgreSQL, for reading and serving the
  text. TL owns their contracts in `docs/architecture/services.md`.
- **Assumption:** the database is seeded with a single row whose text value is
  `Hello Word`. If it turns out empty or absent, the error state in §4 covers
  the screen; the seed itself is TL's deployment concern.

| Open question | Proposed default | Who decides |
|---|---|---|
| — none — | — | — |

## 8. Traceability

| Plan item | Requirement ids | Test cases |
|---|---|---|
| Display Hello Word from database | CONTENT-001, CONTENT-002, CONTENT-003 | `test-cases/display-hello-word.md` |

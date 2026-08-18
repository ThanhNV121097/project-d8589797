# Architecture Overview — Let's build a website

> Foundation document. It fixes the stack, folder layout, conventions and the
> CI gate every later PR must pass. The data model and API contracts are not
> here — they live in `docs/architecture/erd.md` and
> `docs/architecture/services.md`.

Shape: **fullstack** — Next.js frontend, Go backend, PostgreSQL.

## 1. Tech stack

| Layer | Choice | Version | Why |
|---|---|---|---|
| Frontend | Next.js App Router | 15.x | Default; server/client boundary is explicit and the committed Dockerfile assumes `output: "standalone"` |
| Frontend language | TypeScript | 5.x | Default; strict mode on |
| Styling | Tailwind | 3.x | Default; tokens are CSS custom properties in `globals.css` |
| Backend | Go | 1.22+ | Default; single `main` package as the Dockerfile expects |
| Database | PostgreSQL | 16 | Default; only persisted value is the single content row |
| DB driver | pgx | v5 stdlib | `database/sql`-compatible; the one non-stdlib dependency |

## 2. Folder structure

```
code/
  backend/
    cmd/api/main.go          # entry point: migrate → listen → serve
    migrations/              # timestamped .up.sql/.down.sql, embedded
    go.mod / go.sum          # module, dependency hashes
    .env.example             # DATABASE_URL, PORT, APP_PORT
    Dockerfile               # committed, untouched
  frontend/
    app/
      layout.tsx             # root layout, html/body
      page.tsx               # composition root only — stories mount here
      globals.css            # finished; all six token categories
    components/              # one component per story, added later
    lib/mock/                # per-story mocks, deleted when API lands
    package.json             # pinned deps; scripts dev/build/lint
    next.config.js           # output: "standalone"
    tailwind.config.ts       # colours map to var(--color-*)
    Dockerfile               # committed, untouched
docs/
  architecture/              # this document + erd.md + services.md
  content/SRS.md             # requirements
design/                      # approved mockup + design system
docker-compose.yml           # committed, untouched — boots the whole stack
deploy.yml                   # committed — deployment overlay, untouched
.env.example                 # shared compose keys
```

## 3. Key design decisions

1. **Self-migrating backend.** The runtime creates an empty database and
   nothing else applies the schema, so `cmd/api/main.go` applies every pending
   `*.up.sql` from the embedded `migrations/` directory on boot, in filename
   order, tracked in a `schema_migrations` table. Re-running is a no-op.
   *Rejected:* a separate migration runner — adds a step nobody runs and leaves
   a server that starts cleanly then fails on its first query.
2. **Health check proves the database.** `/healthz` returns 200 only after
   migrations succeeded **and** `SELECT 1` works. *Rejected:* a process-liveness
   check — reports a broken app as healthy.
3. **Embedded migrations live beside the directive.** `//go:embed` resolves
   relative to the declaring file, so the directive sits in
   `migrations/embed.go`, not in `cmd/api` where it would look for
   `cmd/api/migrations/`. This is a Go constraint, not a preference.
4. **`page.tsx` is a composition root, not a page.** Every story adds one import
   and one element. Filling it with markup forces the first story to rewrite it.
5. **Design tokens live in `globals.css`, which is finished now.** All six
   categories (colour, spacing, typography, radius, shadow, motion) are defined.
   Story authors are forbidden from adding tokens and must express every value
   through them; CI rejects hardcoded values in `*.module.css`.
6. **Committed container/CI files are fixed.** `docker-compose.yml`, both
   `Dockerfile`s, and `.github/workflows/container.yml` + `publish.yml` were
   committed before this task and are not edited. The scaffold holds up the two
   conventions they assume: frontend at `code/frontend/` with
   `output: "standalone"`, backend at `code/backend/` with exactly one `main`
   package (`cmd/api`).

## 4. Naming conventions

- One component per story: `components/{Component}.tsx`, PascalCase from the
  story title, styles in `components/{Component}.module.css`.
- Every React component file uses `export default function ComponentName()` —
  never `const X = () =>` or a bare `function X()`.
- Client components start with the literal first line `"use client"`. Server
  components do not; they cannot use event handlers, hooks or browser APIs.
- Story mock data lives in `lib/mock/{story-slug}.ts`, deleted when the API
  replaces it.
- Go packages: `cmd/api` for the entry point, `migrations` for embedded SQL,
  `internal/…` for anything later. One `main` package only.

## 5. Environment variables

Backend (`code/backend/.env.example`):

- `DATABASE_URL` — full PostgreSQL DSN; the runtime injects this, never assemble
  from parts.
- `PORT` — listen port; falls back to `APP_PORT`, then `8080`.
- `APP_PORT` — legacy fallback, read only when `PORT` is unset.

Frontend (`code/frontend/.env.example`):

- `NEXT_PUBLIC_API_URL` — base URL the browser uses to reach the API; baked into
  the bundle at build time.

Compose (`.env.example`): `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`,
`BACKEND_PORT`, `FRONTEND_PORT`, `NEXT_PUBLIC_API_URL`, `IMAGE_REPO`,
`IMAGE_TAG`.

## 6. How to run

```sh
docker compose --profile local up -d --wait
```

The `local` profile brings up the project's own Postgres; the default profile
excludes it because a deployment uses a shared instance (see `deploy.yml`).
The backend waits for the DB health check, self-migrates, then serves; the
frontend waits for the backend. Browse to `http://localhost:3000`.

## 7. CI gates

- `.github/workflows/ci.yml` — backend (`go build`/`vet`/`test`), frontend
  (`npm ci`, `lint`, `build`, `test`), compose config validation, and a design
  token grep that rejects hardcoded colours and `var(--x, fallback)` in
  `*.module.css`.
- `.github/workflows/container.yml` — builds and boots the whole stack, asserting
  every healthcheck passes.
- `.github/workflows/publish.yml` — builds and pushes images to GHCR on `main`.

## 8. Risks and unknowns

- **Seeding is a migration, not a runtime default.** The single `Hello Word` row
  is inserted by `0001_create_content.up.sql`. If a deployment wants a different
  value, that is a data change, not a code change.
- **Empty/absent value renders as error**, per SRS §4; the frontend must never
  treat an empty string as success.
- **`prefers-reduced-motion`** guard is included in `globals.css`; the design
  system listed it as a known deviation.

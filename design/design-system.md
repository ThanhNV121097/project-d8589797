# Design System — Let's build a website

> Source of truth: the approved `index.html` (preview: http://localhost:8080/design/d8589797-aee6-42c3-94c6-db252a05e886).
> Every value below is extracted from it. Changing a value here without
> changing the approved design is a defect.

Last updated: 2026-05-27

## 1. Foundations

### 1.1 Color

Semantic tokens. Name by job, never by hue.

| Token | Value | Used for |
|---|---|---|
| `--color-bg` | `#FFFFFF` | Page background |
| `--color-surface` | `#F9FAFB` | Status panel background |
| `--color-border` | `#E5E7EB` | Panel border, button border |
| `--color-text` | `#111827` | Body text, heading |
| `--color-text-muted` | `#6B7280` | Loading text, caption |
| `--color-primary` | `#2563EB` | Primary button background, focus ring |
| `--color-primary-text` | `#FFFFFF` | Text on primary button |
| `--color-danger` | `#B91C1C` | Error message text |
| `--color-focus` | `#2563EB` | Focus ring |

`--color-focus` and `--color-primary` share a value (`#2563EB`) but name different jobs; keep them as separate tokens.

#### Contrast audit

Every text-on-background pair actually used. Body text ≥ 4.5:1, large text (≥ 18.66px bold or ≥ 24px) ≥ 3:1, UI borders ≥ 3:1.

| Foreground | Background | Ratio | Passes |
|---|---|---|---|
| `--color-text` `#111827` | `--color-bg` `#FFFFFF` | `17.7:1` | AA |
| `--color-text-muted` `#6B7280` | `--color-surface` `#F9FAFB` | `4.6:1` | AA (narrow) |
| `--color-text-muted` `#6B7280` | `--color-bg` `#FFFFFF` | `4.8:1` | AA |
| `--color-danger` `#B91C1C` | `--color-surface` `#F9FAFB` | `6.2:1` | AA |
| `--color-primary-text` `#FFFFFF` | `--color-primary` `#2563EB` | `5.2:1` | AA |
| `--color-border` `#E5E7EB` | `--color-bg` `#FFFFFF` | `1.2:1` | FAIL (border < 3:1) |
| `--color-border` `#E5E7EB` | `--color-surface` `#F9FAFB` | `1.2:1` | FAIL (border < 3:1) |

### 1.2 Spacing

Base unit: `4px`. Every margin, padding, and gap in the product uses one of these.

| Token | Value | Used for |
|---|---|---|
| `--space-3` | `12px` | Button group gap |
| `--space-4` | `16px` | Button horizontal padding |
| `--space-5` | `20px` | Caption top margin |
| `--space-6` | `24px` | Page padding, controls top margin |
| `--space-8` | `32px` | Status panel padding |

### 1.3 Typography

Font families (system stack, no external font loaded):

- Body: `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif`
- Headings: same stack, weight 700
- Mono: not used

| Token | Size | Line height | Weight | Used for |
|---|---|---|---|---|
| `--text-xs` | `0.85rem` (13.6px) | normal | 400 | Caption note |
| `--text-sm` | `0.9rem` (14.4px) | normal | 400 | Button label |
| `--text-base` | `1rem` (16px) | normal | 400 | Body, loading, error |
| `--text-display` | `clamp(2rem, 8vw, 3.5rem)` (32–56px) | normal | 700 | h1 "Hello Word" |

The display heading also applies `letter-spacing: -0.02em`. No explicit line-heights are set; browser defaults apply. Heading levels are used in order (single `h1`, no skipped levels).

### 1.4 Radius, border, shadow, motion

| Token | Value | Used for |
|---|---|---|
| `--radius-md` | `8px` | Button |
| `--radius-lg` | `12px` | Status panel |
| `--border-width` | `1px` | Panel border, button border |
| `--shadow-*` | none | No shadows used — flat design |
| `--duration-pulse` | `1.2s` | Loading dot pulse |
| `--easing` | `ease-in-out` | Loading dot pulse |

Button hover uses an instant `filter: brightness(1.05)` — no transition duration.

Motion does **not** yet respect `prefers-reduced-motion: reduce` (see Known deviations).

### 1.5 Layout and breakpoints

| Property | Value |
|---|---|
| Content container | `main`, `width: 100%`, `max-width: 560px`, centered |
| Alignment | Body is a flex column centered both axes; `text-align: center` |
| Breakpoints | None — no media queries. Responsive via `clamp()` on the heading and `flex-wrap` on the button group |

Z-index scale — no z-index values are used; only the implicit base layer (`0`) exists.

## 2. Components

One subsection per reusable component. Every component lists **all** states.

### 2.1 Status panel

**Purpose** — the single region that shows the DB-fetched text; use it wherever a remote value is displayed and may still be loading or may fail. Not for static content.

**Anatomy** — `[panel] → [state content]`; panel is `1px` border, `#F9FAFB` surface, `12px` radius, `32px` padding, `min-height: 120px`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Default | `--color-surface`, `--color-border` | Always |

**Sizes** — one size only.

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default (Loading) | "Loading" + three pulsing dots | `--color-text-muted`, `--duration-pulse` |
| Loaded | `h1` shows the fetched text ("Hello Word") | `--color-text`, `--text-display` |
| Focus | n/a — not an interactive element | |
| Active / pressed | n/a | |
| Disabled | n/a | |
| Loading | dots pulse `opacity 0.2 → 1` | `--color-text-muted` |
| Error | "Could not load the text from the database." | `--color-danger` |
| Empty | **Not designed** — see Known deviations | |

**Accessibility** — panel carries `role="status"` and `aria-live="polite"`; the error message carries `role="alert"`. State switch toggles a `.hidden` class on the three state nodes.

### 2.2 Button

**Purpose** — trigger an action (Reload, Simulate error). Not for navigation.

**Anatomy** — `[label]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Primary | bg `--color-primary`, border `--color-primary`, text `--color-primary-text` | The main action |
| Secondary | bg `--color-bg`, border `--color-border`, text `--color-text` | Any secondary action |

**Sizes** — one size only.

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | ~38px (content-driven) | `10px 16px` | `--text-sm` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Variant fill/border/text | per variant |
| Hover | `filter: brightness(1.05)` | — |
| Focus (keyboard) | `outline: 2px solid`, `outline-offset: 2px` | `--color-focus` |
| Active / pressed | **Not designed** | |
| Disabled | **Not designed** | |
| Loading | **Not designed** | |
| Error | n/a | |
| Empty | n/a | |

**Accessibility** — `type="button"`, native keyboard focus, visible focus ring (never `outline: none`). Minimum hit target is not met: measured height ~38px vs the 44×44px requirement (see Known deviations).

## 3. Content and formatting

- Voice and tone: plain, factual, one short line per state.
- Date, time, number, currency formats: none used — no formatted values on screen.
- Capitalization: sentence case for buttons and messages; the heading "Hello Word" is product copy and is reproduced verbatim.
- Error-message pattern: "Could not load the text from the database." Empty-state wording pattern: not defined (see Known deviations).

## 4. Known deviations

Places where the approved design does not follow its own rules or the
anti-patterns in `references/ai-defaults.md`. Record, do not silently fix.

| Where | Deviation | Why it stands | Follow-up |
|---|---|---|---|
| Loading dots | Pulse animation has no `prefers-reduced-motion` guard | Motion shipped without a reduced-motion fallback | Add the guard when motion is implemented |
| Status panel | No empty state designed — only loading / loaded / error | A DB returning no row has no defined screen | Design the empty state before the API returns empty |
| Button | Active/pressed, disabled, and loading states not designed | Only default/hover/focus specified | Add when the button gains those behaviors |
| Button | Vertical padding `10px` breaks the 4px spacing scale | Off-scale value | Normalize to `8px` or `12px` |
| Button | Height ~38px, below the 44px minimum hit target | Touch users get a small target | Increase padding to reach 44px |
| Panel border | `--color-border` at `1.2:1` fails the 3:1 non-text contrast bar | Light border is decorative; the panel is already distinguished by its surface color | Confirm the boundary is decorative, or darken the border |

Avoided AI defaults (no action needed): no purple/indigo palette (`#2563EB` blue chosen for a single clear accent), no gradients, no `9999px` pill rounding (8/12px scale), no oversized padding (single centered card), no shadows, no emoji icons, realistic copy throughout, and a visible focus ring is present.

## 5. Change log

| Date | Change | Design PR |
|---|---|---|
| 2026-05-27 | Initial extraction from approved `index.html` | |

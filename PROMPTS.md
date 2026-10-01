# AI Usage Log — Prompts & Human Decisions

This file documents how AI was used to build this project: every prompt executed,
the configuration behind it (model, mode, skills, subagents, tools) and what I did
by hand — reviews, corrections, decisions and their rationale.

## Workflow

| Role | Responsibility |
|---|---|
| **Me (author)** | Requirements, plan review, approval gates, design decisions, manual verification, final review |
| **Claude Code — Claude Opus 5.5** | Implementation, tests, tooling and commits inside this repository |
| **Claude chat — Claude Opus 5.5** | Separate assistant used to draft and refine prompts and to review plans before I sent them. Drafting conversations are not logged; the prompts actually executed are. |

Persistent rules live in [`CLAUDE.md`](CLAUDE.md): phased work with a stop for approval
after each phase, run tests before claiming success, Conventional Commits, no AI
attribution in commits, never push without asking, and a fixed end-of-phase report
(the source for the fields below).

**Environment:** WSL2 (Ubuntu) · Go 1.23 · Node 20 · Docker

---

## 01 — Master prompt: architecture plan + Phase 0 scaffold

| Field | Value |
|---|---|
| **Model** | Claude Opus 5.5 |
| **Mode** | Plan mode → Auto mode (implementation) |
| **Skills** | — |
| **Subagents** | `Explore` — surveyed the repo state · `Plan` — designed API contract, types, folder trees and phases |
| **Tools** | Bash (go build/vet/run, npm create vite, npm install/build/lint/test, curl, git), Read, Edit, Write |
| **Phase** | Planning + Phase 0 — Scaffold |

**Prompt:**
```text
Read CLAUDE.md first; its rules apply to the whole project.

## Assignment (Sezzle take-home)
Build a full-stack calculator: React + TypeScript frontend consuming a Go REST microservice.
Operations: add, subtract, multiply, divide. Also: power, square root, percentage.

## Backend (Go, standard library only)
- Layout: backend/cmd/server, backend/internal/calculator (pure domain, no HTTP),
  backend/internal/api (handlers, DTOs, validation, middleware).
- Operations behind a small interface + registry (strategy pattern): adding an
  operation = one new type + registration, no handler changes.
- Proposed contract (challenge it if you see a better design):
  POST /api/v1/calculate {"operation":"add","operands":[2,3]} -> 200 {"result":5}
  GET  /api/v1/operations -> supported operations with arity
  GET  /healthz
  Errors: {"error":{"code":"DIVISION_BY_ZERO","message":"..."}}
  400 invalid input, 422 math domain errors, 404/405 where appropriate.
- Edge cases, each with an explicit test: division by zero, sqrt of negative,
  wrong operand count, unknown operation, malformed JSON, unknown fields,
  missing fields, NaN/±Inf/overflow results, oversized body, wrong Content-Type.
- Define precisely what "percentage" means and justify it.
- float64; document the precision trade-off vs a decimal library.
- slog request logging, panic recovery, CORS for local dev, graceful shutdown, PORT env var.
- Table-driven tests + httptest. Target >90% coverage on internal/.

## Frontend (React + TypeScript strict, Vite)
- Structure: api/ (typed client, the only place that knows the HTTP contract),
  hooks/ (useCalculator), components/ (presentational), types/.
- Calculator keypad + display, clear/backspace, keyboard support, loading state,
  friendly messages mapped from backend error codes, client-side validation
  (backend remains the source of truth).
- Mobile-first, accessible (labels, aria-live result, visible focus).
- Vite dev proxy to the backend; base URL via env.
- Vitest + React Testing Library; mock the API (justify MSW vs mocked client).

## Delivery
- Multi-stage Dockerfile: Go server also serves the built frontend; docker-compose.yml.
- Makefile shortcuts, but README must also show the plain commands.
- GitHub Actions: go vet + go test, tsc + eslint + vitest.
- README: overview, Mermaid architecture diagram, setup, run backend/frontend/Docker,
  tests + coverage reports, curl examples for every operation and error case,
  design decisions, assumptions, what I'd do with more time.

## What I need now (plan only, no code)
1. Architecture and final API contract (with any changes you recommend and why).
2. Full folder tree.
3. Testing strategy per layer.
4. List of assumptions and open questions for me.
5. Phase-by-phase plan with the intended commit sequence.
```

**Human actions & decisions:**
- Before this prompt: set up the repository with my git identity, wrote the project rules
  in `CLAUDE.md`, disabled AI commit attribution in `.claude/settings.json` and forced LF
  line endings with `.gitattributes` (Windows host + WSL + Linux containers).
- Defined the backend layering and the strategy-pattern registry myself in the prompt, so
  the AI designed within my architecture instead of choosing one.
- Left the API contract explicitly open to challenge and asked for assumptions and open
  questions, to surface ambiguity before any code was written.
- The plan was accepted in auto mode right away, so Phase 0 (scaffold only) ran before
  my detailed review. I let it finish because scaffolding doesn't depend on the open
  design points, and did the full review before any business logic was written (entry 02).
- Verified the Phase 0 commits: all authored by me, no AI trailers.

**Outcome:**
- Plan: percentage defined as `a% of b` (`percentage(20, 50) = 10`); flat
  `{operation, operands}` contract with per-operation arity; error split 400 (malformed
  request) vs 422 (undefined math); 11 error codes; hand-written API-client mocks instead of MSW.
- Phase 0: Go module serving `/healthz`, Vite + React + TS strict scaffold with ESLint
  flat config and Vitest/RTL, root Makefile. Build, vet, lint and typecheck pass.

**Commits:** `1a1bb08` chore(backend): init go module and scaffold entrypoint ·
`0f14900` chore(frontend): scaffold vite react-ts app with eslint and vitest ·
`2462978` chore: add root makefile skeleton

---

## 02 — Plan review and corrections

| Field | Value |
|---|---|
| **Model** | Claude Opus 5.5 |
| **Mode** | Auto mode |
| **Skills** | — |
| **Subagents** | — |
| **Tools** | Edit (plan file), git (amend), task list |
| **Phase** | Planning — revision before Phase 1 |

**Prompt:**
```text
Phase 0 looks good.

0. PROMPTS.md now has content. Amend the last commit (90eee8d, not pushed yet)
   so it includes it: git add PROMPTS.md && git commit --amend --no-edit
   Do not edit PROMPTS.md.

Before starting Phase 1, apply these changes to the plan
(update the saved plan file so it stays the source of truth):

1. HTTP status codes: oversized body -> 413 REQUEST_TOO_LARGE,
   wrong Content-Type -> 415 UNSUPPORTED_MEDIA_TYPE. Everything else as planned.
2. GET /api/v1/operations returns an object, not a bare array:
   {"operations":[{"name":"add","arity":{"min":2,"max":2}}, ...]}
   so the contract can grow without breaking clients.
3. Frontend input model: pocket-calculator "immediate execution".
   Each binary operator applies to the running result (2 + 3 × 4 = 20, not 14);
   sqrt and percentage act on the current value/operands as a real calculator would —
   propose the exact UX for those two before Phase 3. The frontend never evaluates
   math itself; every result comes from the API. Document this in the README.
4. Phase 3 must use the frontend-design skill for the visual direction.
   Report in the phase summary which aesthetic direction you chose and why.
5. Add Phase 6 — Independent review: spawn a subagent that has not seen the
   implementation to review the whole repo as a strict Sezzle interviewer
   (correctness, edge cases, idiomatic Go/React, test gaps, README clarity).
   Show me the findings; I decide which ones get fixed.
6. CI: backend job fails if coverage on internal/ drops below 90%.

Confirm the updated plan, then start Phase 1 (calculator domain, tests first).
Show me the test table for the edge cases before writing the implementation.
```

**Human actions & decisions:**
- Reviewed the approved plan and the Phase 0 report before any business logic existed.
- **413 / 415 instead of a generic 400:** the plan returned 400 for everything malformed;
  HTTP has specific codes for an oversized body and an unsupported media type, and
  clients can react to them differently.
- **`/operations` wrapped in an object:** a bare JSON array can't gain fields
  (version, metadata) without a breaking change.
- **Immediate-execution input model:** the plan never said how a keypad sequence like
  `2 + 3 × 4` maps to a binary API. I chose pocket-calculator semantics (each operator
  resolves the pending operation with one API call) because it keeps all math in the
  backend and the frontend a simple state machine. Operator precedence would require an
  expression parser — listed as future work instead.
- **`frontend-design` skill for Phase 3:** explicitly chose it so the UI has a deliberate
  visual direction instead of a default template, and asked for the rationale in the report.
- **Phase 6 independent review:** added a review by a fresh subagent with no context of the
  implementation, so the code isn't graded by the same context that wrote it. I keep the
  decision on which findings get fixed.
- **Coverage as a CI gate:** the plan had it as a "soft gate"; made it a hard failure under 90%.
- Kept TDD as a checkpoint: asked to see the edge-case test table before implementation.

**Outcome:** Plan file updated with the six changes; Phase 6 added to the task list.
Claude proposed one refinement: a `Registry.Calculate` method that centralizes arity
validation and NaN/±Inf classification in the domain layer, leaving the API layer only
transport concerns (JSON, Content-Type, body size). It then presented the Phase 1 test
table for review (entry 03).

**Commits:** none (plan revision only)

---

## 03 — Test table review + Phase 1: calculator domain

| Field | Value |
|---|---|
| **Model** | Claude Opus 5.5 |
| **Mode** | Auto mode |
| **Skills** | — |
| **Subagents** | — |
| **Tools** | Bash (go build/vet/test, go tool cover, git), Read, Write, Edit |
| **Phase** | 1 — Calculator domain (TDD) |

**Prompt:**
```text
Test table approved, and I accept the Registry.Calculate refinement:
the domain owns math and contract rules, the API layer owns transport only.

0. First: PROMPTS.md was rewritten after your last amend. Amend HEAD again so the
   commit contains the current file: git add PROMPTS.md && git commit --amend --no-edit
   Do not edit PROMPTS.md.

Add these cases before implementing:
1. add: 0.1 + 0.2 -> 0.30000000000000004. This test documents the float64
   precision trade-off on purpose; reference it from the README later.
2. Negative zero: Calculate normalizes -0 to 0 (e.g. multiply -4 × 0 -> 0, not -0),
   so the API never returns "-0". Test it.
3. Registry:
   - registering the same operation name twice panics (programmer error);
   - List() returns operations sorted by name (stable /operations output).
4. NewDefaultRegistry() contains exactly the 7 expected operations.

Then implement Phase 1: write the tests, run them and show them failing,
then implement until green. Run go vet and go test -cover on internal/calculator.
Commit following the plan and give me the Phase 1 report.
```

**Human actions & decisions:**
- Reviewed the proposed edge-case test table before any implementation existed (TDD gate).
- **Accepted `Registry.Calculate`:** one domain entry point owns arity validation and
  NaN/±Inf classification, so the API layer only deals with transport concerns.
- **Added `0.1 + 0.2` as an explicit test:** turns the float64 precision trade-off into
  executable documentation instead of a README sentence.
- **Added negative-zero normalization:** `-4 × 0` is `-0` in IEEE-754 and would serialize
  as `"-0"` in JSON — technically valid, confusing on a calculator display.
- **Added registry tests:** duplicate registration must panic (fail fast on a programmer
  error) and `List()` must be sorted, because Go map iteration order is random and
  `/operations` would otherwise change order between calls.
- **Added a default-registry test:** guarantees all 7 operations stay registered.
- Accepted Claude's change from 4 planned commits to 3: a tests-only commit would have
  left the build broken in between, and every commit should build and pass on its own.
- Verified the history: the `PROMPTS.md` fix was amended into the docs commit, and all
  commits are authored by me with no AI trailers.

**Outcome:** Operation interface, sorted registry with centralized `Calculate`, 5 sentinel
errors and 7 operations, each in its own file with table-driven tests. Generic registry
behavior is tested with stub operations, independent of the real ones.
`go test ./internal/calculator/... -cover` → **100% coverage**, `go vet` clean, and each
of the 3 commits builds and passes on its own.

**Commits:** `e90698d` feat(calculator): add operation interface and registry ·
`3a110f4` feat(calculator): implement arithmetic operations (add, subtract, multiply, divide) ·
`969e108` feat(calculator): implement power, sqrt, percentage

---

## 04 — Percentage fix + Phase 2: REST API layer

| Field | Value |
|---|---|
| **Model** | Claude Opus 5.5 |
| **Mode** | Auto mode |
| **Skills** | — (no specialized skill applies to a Go HTTP layer) |
| **Subagents** | `general-purpose` — black-box contract tester: received only the API contract, not the source, and wrote `contract_test.go` (37 subtests) |
| **Tools** | Bash (go build/vet/test, go tool cover, curl, git), Read, Write, Edit, Agent |
| **Phase** | 2 — API layer |

**Prompt:**
```text
Phase 1 approved. Before Phase 2:

0. Add this rule to CLAUDE.md under "Working rules" and commit it (chore: ...):
   "Use skills and subagents whenever they add independent value (e.g. a subagent
   that tests or reviews without seeing the implementation, a skill for a specialized
   domain). Never use them just for show. Justify each one in the phase report."

1. Bug found in review: percentage computes (a/100)*b, which gives
   percentage(7,100) = 7.000000000000001 and percentage(29,100) = 28.999999999999996.
   Change it to a*b/100 (exact for typical inputs; extreme overflow is still caught as
   RESULT_OUT_OF_RANGE). Add both cases as regression tests, show them failing first,
   then fix. Commit as fix(calculator): ...

Then Phase 2 — API layer, as planned (413/415, /operations wrapped in an object).
Also handle and test these transport edge cases:
- null inside operands: [1, null] must be 400, not silently decoded as 0
  (Go decodes null into float64 as zero — that would turn it into a division by zero).
- numbers sent as strings ("2") -> 400.
- trailing data after the JSON object ({...}{...}) -> 400.
- Content-Type "application/json; charset=utf-8" must be accepted (parse the media type).
- operands: [] and operands missing are both MISSING_FIELD or INVALID_OPERAND_COUNT —
  pick one, be consistent, document it.

Subagent: once the handlers compile, spawn a general-purpose subagent that receives
ONLY the API contract (endpoints, status codes, error codes — not the source code)
and writes black-box contract tests in backend/internal/api/contract_test.go using
httptest against the real router. Fix any failure it finds and report what it caught.

Finish with go vet, go test -cover ./..., and a curl smoke test of every operation
and every error code against the running server. Include the curl commands in the
report so I can rerun some of them myself. Give me the Phase 2 report.
```

**Human actions & decisions:**
- **Found a precision bug in code review** (the Phase 1 tests passed): `percentage` computed
  `(a/100)*b`, so `7% of 100` returned `7.000000000000001`. I reproduced it, then required
  `a*b/100`, which is exact for typical inputs; the only cost is overflow slightly earlier for
  values near 1e306, already caught as `RESULT_OUT_OF_RANGE`. Required regression tests
  red-first.
- **Added transport edge cases the plan missed:**
  - `null` inside `operands`: Go silently decodes `null` into a `float64` as `0`, so
    `divide [1, null]` would surface as a misleading division by zero instead of bad input.
  - Numbers sent as strings.
  - Trailing data after the JSON object.
  - `application/json; charset=utf-8` must be accepted, because real clients send it.
- **Made skill/subagent usage a project rule** in `CLAUDE.md`: use them only when they add
  independent value, and justify each one in the phase report.
- **Manually verified the running server** with curl from a separate terminal: a normal sum,
  `7% of 100` → `7` (the fix), division by zero → 422 `DIVISION_BY_ZERO`, `[1, null]` →
  400 instead of a false division by zero, and malformed JSON → 400 `INVALID_JSON`.
  All responses matched the contract.
- **Designed the subagent's role:** a black-box tester that never sees the implementation,
  so its tests check the contract instead of mirroring the code.
- Reviewed and accepted Claude's decisions:
  - `operands` missing or `[]` → `INVALID_OPERAND_COUNT`, one code for "wrong number of operands".
  - `Run` refactored into `Serve(ctx)` so graceful shutdown is testable without OS signals.
  - Feature commits merged where splitting them would break the build.

**Outcome:**
- Handlers for `/calculate`, `/operations` and `/healthz`; middleware for panic recovery,
  `slog` logging and CORS; 413 / 415; graceful shutdown via `signal.NotifyContext`;
  `PORT` environment variable.
- The blind contract tests found no violations (37/37 pass); Claude re-ran them itself
  instead of trusting the subagent's summary.
- `go vet` clean; coverage on `internal/` is **97.6%** (api 96.3%, calculator 100%).

**Commits:** `81c8df3` chore: require justifying skill and subagent use in phase reports ·
`d76b1d8` fix(calculator): compute percentage as a*b/100 to avoid rounding error ·
`6d2320d` feat(api): add DTOs and error mapping ·
`0fbeb68` feat(api): implement handlers, middleware, and server wiring with graceful shutdown ·
`747278f` test(api): add black-box contract tests from the published API spec

---

<!-- Next entries are appended below -->

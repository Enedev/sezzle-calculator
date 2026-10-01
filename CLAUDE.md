# CLAUDE.md — Project rules for Claude Code

## Project
Full-stack calculator for the Sezzle take-home assignment.
- `backend/`  — Go REST microservice (standard library only).
- `frontend/` — React + TypeScript (Vite).
Evaluation priorities: correctness, clarity, maintainability, testability. Not feature count.
Time budget: ~2–4 hours of equivalent work. Do not over-engineer.

## Environment
- Development happens in WSL2 (Ubuntu). Repo path: `/mnt/c/dev/sezzle-calculator`.
- Toolchain: Go 1.23, Node 20, npm 10, GNU Make. Docker is available from Phase 4.
- Run every command from WSL. Never call Windows executables (`*.exe`).

## Git rules (mandatory)
- Conventional Commits: `feat:`, `fix:`, `test:`, `refactor:`, `docs:`, `chore:`, `ci:`, `build:`.
  Use scopes when useful, e.g. `feat(api): ...`, `test(calculator): ...`.
- Small, meaningful commits. Commit after each milestone, only with tests passing.
- NEVER add AI attribution of any kind to commits or PRs: no `Co-Authored-By: Claude`,
  no "Generated with Claude Code", no session links, no emojis signatures. The author is the human.
- NEVER push without asking first.
- Do not change git config (user, email, remote).

## PROMPTS.md
- `PROMPTS.md` is maintained by the author outside of this session. NEVER edit it.
- If it has uncommitted changes when you commit, include it in a separate commit:
  `docs(prompts): update AI usage log`.

## Working rules
- Work in phases. At the end of each phase STOP and wait for approval.
- Never claim something works without running it. Run tests and linters after every change.
- If a requirement is ambiguous, state the assumption explicitly and record it for the README.
- Prefer simple, idiomatic code over clever code. Comment only non-obvious decisions.
- Use skills and subagents whenever they add independent value (e.g. a subagent that tests
  or reviews without seeing the implementation, a skill for a specialized domain). Never use
  them just for show. Justify each one in the phase report.

## End-of-phase report (always use this format)
```
### Phase <N> report
- Summary: <2–3 lines>
- Skills used: <name — why> | none
- Subagents used: <type — what for> | none
- Tools used: <e.g. Bash (go test, npm run test), Read/Write/Edit, git>
- Assumptions made: <list> | none
- Tests: <command> → <result, coverage %>
- Commits: <hash short — message>
- Open questions for the author: <list> | none
```

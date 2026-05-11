---
description: Implement a specific feature for the wget (01-edu style) clone project.
---

# Implementation instruction: {Ticket}

You are tasked with implementing a specific slice of the **wget**-style downloader described in this repository. Produce correct, maintainable code in the **compiled language chosen for the project** (for example Go, Rust, or C), aligned with the subject brief and verification material.

## 1. Primary directives

1. **Read the sources of truth** before writing code:
   - **`docs/requirements.md`**: required behavior, flags, examples, and mirror semantics.
   - **`docs/audit.md`**: acceptance checks grouped under Functional, Mirror, and Bonus.
   - **`docs/ticket-tracker.md`**: ticket description, dependencies, coverage tags (`RQ-*`, `AQ-*`), and sprint ordering.
2. **Match the public contract**: CLI behavior, flag forms (`-O=…`, `-P=…`, `-i=…`, `--rate-limit=…`, `-B`, `--mirror`, and mirror companions), log layout for background mode, and console output shape should follow **`docs/requirements.md`** unless **`docs/audit.md`** demands a stricter check—then **`docs/audit.md`** wins for verification wording.
3. **Verification gate**: Work is not complete until the ticket’s stated checks are satisfied and you can point to tests or manual steps that prove them.

## 2. Technical standards

- **HTTP(S) first**: Implement reliable single-file retrieval before advanced mirror logic; handle non-200 responses as specified (message and clean exit).
- **Concurrency**: Batch downloads from `-i` must be **asynchronous** relative to each other (overlap), not strictly sequential, per **`docs/requirements.md`**.
- **Mirror mode**: Recursive crawl with domain-based output folder, following link discovery rules in **`docs/requirements.md`** (`a`, `link`, `img` with `href` / `src`); optional flags (`--reject` / `-R`, `--exclude` / `-X`, `--convert-links`) apply only in conjunction with `--mirror` as described there.
- **Observability**: Timestamps `yyyy-mm-dd hh:mm:ss`, status line, content size (raw and human-rounded), destination path, and progress presentation should match the examples in **`docs/requirements.md`** where applicable.
- **Small, testable units**: Prefer focused functions or modules (argument parsing, HTTP fetch, progress reporting, mirror crawl, filesystem layout) over monolithic `main`.

## 3. Workflow

### Step 1 — Analysis
- Locate the ticket in **`docs/ticket-tracker.md`** and note **dependencies** and **`RQ` / `AQ`** tags.
- Re-read the relevant sections of **`docs/requirements.md`** and **`docs/audit.md`** for that slice (e.g. flags, mirror, or logging).

### Step 2 — Implementation
- Implement the smallest change that fulfills the ticket; avoid unrelated refactors.
- Keep flag parsing consistent with examples (`-O=name`, `-P=path`, `-i=file`, `--rate-limit=300k`, etc.).

### Step 3 — Testing and verification
- Add or extend automated tests where practical (redirects, error status, path handling, rate limit behavior, mirror rules).
- Run the project’s test and build commands; manually exercise any **`docs/audit.md`** items that are not fully automated.
- If you find a defect: add a minimal failing test (or reproducible command), fix it, then re-run the suite.

### Step 4 — Documentation and handover
- Write the PR narrative using **`docs/pr-messages/pr-template.md`** as the outline.
- Save the completed PR text as a new file under **`docs/pr-messages/`**, named from the ticket id and a short slug (for example `docs/pr-messages/W03-progress-bar-pr.md`).
- Update **`docs/ticket-tracker.md`**: set the ticket status to **Done** when merged or when the branch fully satisfies the gate.

---

## 4. Ticket context: {Ticket}

> [!IMPORTANT]
> Paste the full ticket row (or description + verification bullets) from **`docs/ticket-tracker.md`** here before executing.

**Begin implementation.** Prefer clear errors, bounded resource usage on large downloads, and behavior that matches **`docs/requirements.md`** and **`docs/audit.md`**.

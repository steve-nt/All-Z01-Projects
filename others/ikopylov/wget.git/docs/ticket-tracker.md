> Legend: 🔴 Blocked · 🟡 Ready · 🟢 In Progress · ✅ Done · ⬜ Not Started
>
> **Requirements legend**
> - **RQ-***: Requirements from **`docs/requirements.md`**.
> - **AQ-***: Acceptance checks from **`docs/audit.md`**.
>
> **Precedence**: If **`docs/audit.md`** and **`docs/requirements.md`** disagree on a verification detail, treat **`docs/audit.md`** as the grading checklist and **`docs/requirements.md`** as the primary design brief.

---

# wget Ticket Tracker

Last refreshed: 2026-04-21

## 1) Scope contract

This tracker is **requirements-first** and **audit-first**.

Execution order:

1. Satisfy **`docs/requirements.md`** for the feature slice.
2. Map completed work to **`docs/audit.md`** Functional, Mirror, and Bonus sections.
3. Use **`docs/AGENT_WORKFLOW.md`** for the agent implementation loop and PR message rules.

There is no separate PRD/SDS in this repo; **`docs/requirements.md`** plus **`docs/audit.md`** define the product.

---

## 2) Tracks

| Track | Focus |
|-------|--------|
| **D** | Scaffolding, build, tests, README, traceability |
| **W** | Downloader core: HTTP, CLI output, flags, mirror |

---

## 3) Requirement IDs (`docs/requirements.md`)

| ID | Requirement summary |
|----|------------------------|
| RQ1 | Implement as a **compiled** program; usage matches subject style (`./wget` or documented equivalent). |
| RQ2 | Download a file from a **single URL**; save under the expected default name. |
| RQ3 | Print **start time** in `yyyy-mm-dd hh:mm:ss`. |
| RQ4 | Show **HTTP status**; proceed only on success path as specified; otherwise report status and exit with error. |
| RQ5 | Show **content length** as raw bytes and rounded MB/GB (or MB as in examples). |
| RQ6 | Show **destination path/name** before saving. |
| RQ7 | **Progress bar**: KiB/MiB amounts, percentage, **ETA**, smooth updates on large files; **finish time** in required format. |
| RQ8 | **`-O=`** save under a different filename. |
| RQ9 | **`-P=`** save under a chosen directory (compose with `-O` as in examples). |
| RQ10 | **`--rate-limit=`** with `k` and `M` semantics. |
| RQ11 | **`-B`**: background download, user-facing line about log file, **silent** terminal, **`wget-log`** structure as specified. |
| RQ12 | **`-i=`** read many URLs from a file; downloads run **asynchronously** (overlapping). |
| RQ13 | **`--mirror`**: mirror site into a **host-named folder**; recurse using **`a`**, **`link`**, **`img`** with **`href`** / **`src`**. |
| RQ14 | Mirror **reject** list: **`--reject` / `-R`** suffixes not downloaded. |
| RQ15 | Mirror **exclude** paths: **`--exclude` / `-X`** paths not followed. |
| RQ16 | **`--convert-links`** for offline-local resource references. |
| RQ17 | **Bonus**: speed / avoiding unnecessary work (audit “quickly and effectively”). |
| RQ18 | **Bonus**: [good practices](https://public.01-edu.org/subjects/good-practices/README.md) (audit item). |

---

## 4) Audit IDs (`docs/audit.md`)

| ID | Audit summary |
|----|----------------|
| AQ1 | Twitter image URL downloads `EMtmPFLWkAA8CIS.jpg`. |
| AQ2 | Arbitrary chosen URL downloads expected file. |
| AQ3 | `go1.16.3.linux-amd64.tar.gz` downloads with correct filename. |
| AQ4 | Start time is shown. |
| AQ5 | Start and end times match `yyyy-mm-dd hh:mm:ss`. |
| AQ6 | Response status shown (e.g. 200 OK). |
| AQ7 | Content length of download shown. |
| AQ8 | Content length shown raw and rounded (MB/GB). |
| AQ9 | Saved file **name and path** shown. |
| AQ10 | Large sample (`Sample.zip`) downloads as expected. |
| AQ11 | Progress shows downloaded amount (KiB/MiB). |
| AQ12 | Progress shows **percentage**. |
| AQ13 | Progress shows **time remaining**. |
| AQ14 | Progress advances smoothly for the duration. |
| AQ15 | `-O=test_20MB.zip` saves under that name. |
| AQ16 | `-O` with `-P=~/Downloads/` places file in Downloads. |
| AQ17 | `--rate-limit=300k` keeps speed below cap. |
| AQ18 | `--rate-limit=700k` keeps speed below cap. |
| AQ19 | `--rate-limit=2M` keeps speed below cap. |
| AQ20 | `-i=downloads.txt` downloads all listed assets. |
| AQ21 | `-i` downloads are **asynchronous** (order tip in audit). |
| AQ22 | `-B` prints log redirection message as specified. |
| AQ23 | `-B` run is **silent** in terminal during download. |
| AQ24 | `wget-log` structure matches audit template. |
| AQ25 | With `-B`, file is actually present after download. |
| AQ26 | `--mirror --convert-links` on sample site; offline `index.html` works in browser. |
| AQ27 | `--mirror` on sample site; offline site works. |
| AQ28 | `--mirror --reject=gif` excludes GIFs; site still usable. |
| AQ29 | Mirrored `trypap.com` layout matches expected listing (`css`, `img`, `index.html`). |
| AQ30 | `--mirror -X=/img` on `trypap.com` excludes `img` directory from result. |
| AQ31 | Mirror `theuselessweb.com` works offline. |
| AQ32 | Mirror arbitrary chosen site succeeds. |
| AQ33 | Bonus: project runs quickly/effectively (no needless requests, etc.). |
| AQ34 | Bonus: good practices link satisfied. |

---

## Sprint 0 — Bootstrap

| ID | Ticket | Size | Status | Deps | Coverage |
|----|--------|------|--------|------|----------|
| D01 | **Project scaffolding**: Choose language; layout sources; produce binary or documented `./wget` entrypoint; minimal “hello” or no-op CLI. | S | ✅ | - | RQ1 |
| D02 | **README runbook**: Install, build, run, and test commands; link to **`docs/requirements.md`** / **`docs/audit.md`**. | S | ✅ | D01 | RQ1 |
| D03 | **Test harness**: Add baseline automated tests (e.g. argument parsing, small local HTTP fixture if used). | M | ⬜ | D01 | RQ1 |

---

## Sprint 1 — Single-file HTTP and console contract

| ID | Ticket | Size | Status | Deps | Coverage |
|----|--------|------|--------|------|----------|
| W01 | **HTTP(S) GET to disk**: Download URL to default filename; handle errors. | M | ⬜ | D01 | RQ2, AQ1–AQ3 |
| W02 | **Console preamble**: Start time, status line, content size (raw + rounded), save path — match **`docs/requirements.md`** examples. | M | ⬜ | W01 | RQ3–RQ6, AQ4–AQ9 |
| W03 | **Progress and finish**: Progress bar fields (KiB/MiB, %, ETA); finish timestamp; stable on large file. | L | ⬜ | W02 | RQ7, AQ10–AQ14 |

---

## Sprint 2 — Output location and rate

| ID | Ticket | Size | Status | Deps | Coverage |
|----|--------|------|--------|------|----------|
| W04 | **`-O=` output name**: Override local filename. | S | ⬜ | W01 | RQ8, AQ15 |
| W05 | **`-P=` directory**: Resolve save directory with `-O` / default name. | M | ⬜ | W04 | RQ9, AQ16 |
| W06 | **`--rate-limit`**: Throttle with `k`/`M` suffixes; verified on large sample. | L | ⬜ | W03 | RQ10, AQ17–AQ19 |

---

## Sprint 3 — Background and batch

| ID | Ticket | Size | Status | Deps | Coverage |
|----|--------|------|--------|------|----------|
| W07 | **`-B` background mode**: Detach/silent UI, stdout message, `wget-log` format, successful file. | L | ⬜ | W02 | RQ11, AQ22–AQ25 |
| W08 | **`-i=` batch**: Parse file; concurrent downloads; completion behavior per brief. | L | ⬜ | W01 | RQ12, AQ20–AQ21 |

---

## Sprint 4 — Mirror core

| ID | Ticket | Size | Status | Deps | Coverage |
|----|--------|------|--------|------|----------|
| W09 | **`--mirror` crawl**: Host folder; fetch HTML/CSS assets; follow **`a`/`link`/`img`** rules. | L | ⬜ | W01 | RQ13, AQ27, AQ29, AQ31–AQ32 |
| W10 | **`--reject` / `-R`**: Skip suffixes during mirror. | M | ⬜ | W09 | RQ14, AQ28 |
| W11 | **`--exclude` / `-X`**: Skip path prefixes during mirror. | M | ⬜ | W09 | RQ15, AQ30 |
| W12 | **`--convert-links`**: Rewrite links for offline viewing. | L | ⬜ | W09 | RQ16, AQ26 |

---

## Sprint 5 — Gates and bonus

| ID | Ticket | Size | Status | Deps | Coverage |
|----|--------|------|--------|------|----------|
| D04 | **Traceability pass**: Update matrices below; ensure each RQ/AQ has an owning ticket or explicit deferral. | S | ⬜ | W01–W12 | RQ1–RQ18 |
| D05 | **Audit sweep**: Run **`docs/audit.md`** checklist; document gaps or N/A. | M | ⬜ | D04 | AQ1–AQ34 |
| D06 | **Bonus evidence** (optional): Performance notes; good-practices checklist. | S | ⬜ | D05 | RQ17–RQ18, AQ33–AQ34 |

---

## 5) Verification gates

### Gate G1 — Core download

- Single-URL downloads and console contract behave per **`docs/requirements.md`**.
- Evidence: tests + manual runs covering AQ1–AQ9 (as applicable).

### Gate G2 — Progress, paths, throttle

- Progress bar and `-O`/`-P`/rate limit satisfied.
- Evidence: AQ10–AQ19.

### Gate G3 — Background and batch

- `-B` and `-i` behavior satisfied.
- Evidence: AQ20–AQ25.

### Gate G4 — Mirror

- Mirror recursion and optional mirror flags satisfied.
- Evidence: AQ26–AQ32.

### Gate G5 — Bonus (optional)

- AQ33–AQ34 addressed or explicitly scoped out.

---

## 6) Requirements coverage matrix

| RQ | Tickets | Gate |
|----|---------|------|
| RQ1 | D01, D02 | G1 |
| RQ2 | W01 | G1 |
| RQ3–RQ6 | W02 | G1 |
| RQ7 | W03 | G2 |
| RQ8 | W04 | G2 |
| RQ9 | W05 | G2 |
| RQ10 | W06 | G2 |
| RQ11 | W07 | G3 |
| RQ12 | W08 | G3 |
| RQ13 | W09 | G4 |
| RQ14 | W10 | G4 |
| RQ15 | W11 | G4 |
| RQ16 | W12 | G4 |
| RQ17–RQ18 | D06 | G5 |

---

## 7) Audit coverage matrix

| AQ | Tickets | Gate |
|----|---------|------|
| AQ1–AQ3 | W01 | G1 |
| AQ4–AQ9 | W02 | G1 |
| AQ10–AQ14 | W03 | G2 |
| AQ15 | W04 | G2 |
| AQ16 | W05 | G2 |
| AQ17–AQ19 | W06 | G2 |
| AQ22–AQ25 | W07 | G3 |
| AQ20–AQ21 | W08 | G3 |
| AQ26–AQ32 | W09–W12 | G4 |
| AQ33–AQ34 | D06 | G5 |

---

## 8) Next work queue (suggested)

1. Close **D01**–**D03** so builds and tests exist.
2. Deliver **W01**–**W03** for a credible single-file experience.
3. Layer **W04**–**W08**, then **W09**–**W12**.
4. Run **D04**–**D06** for traceability and audit sign-off.

---

## Summary counts

| Track | Ticket count (approx.) |
|-------|-------------------------|
| D | 6 |
| W | 12 |
| **Total** | **18** |

# go-reloaded

A Go CLI tool: `go run main.go input.txt output.txt`

Reads input, applies transformation functions, writes output. Standard library only.

## Presentation

Open `deck.html` in your browser. Navigate with arrow keys.

---

## The 7 Transformation Functions

### Numbers
| Function | What it does | Example |
|----------|-------------|---------|
| `(hex)` | Previous word: hex to decimal | `1E (hex)` -> `30` |
| `(bin)` | Previous word: binary to decimal | `10 (bin)` -> `2` |

### Casing
| Function | What it does | Example |
|----------|-------------|---------|
| `(up)` | Previous word to UPPERCASE | `go (up)` -> `GO` |
| `(low)` | Previous word to lowercase | `SHOUTING (low)` -> `shouting` |
| `(cap)` | Previous word to Capitalized | `bridge (cap)` -> `Bridge` |
| `(up, N)` | Previous N words to UPPERCASE | `so exciting (up, 2)` -> `SO EXCITING` |
| `(low, N)` | Previous N words to lowercase | `BREAKFAST IN BED (low, 3)` -> `breakfast in bed` |
| `(cap, N)` | Previous N words to Capitalized | `harold wilson (cap, 2)` -> `Harold Wilson` |

### Punctuation
- `. , ! ? : ;` attach to the previous word, space after (unless end of text)
- Groups like `...` or `!?` stay together
- `I was thinking ... You` -> `I was thinking... You`
- `there ,and then BAMM !!` -> `there, and then BAMM!!`

### Quotes
- `'` always comes in pairs
- Remove spaces inside: `' awesome '` -> `'awesome'`
- Multiple words: `' I am the best '` -> `'I am the best'`

### Articles
- `a` becomes `an` before a vowel (`a, e, i, o, u`) or `h`
- `a amazing` -> `an amazing`
- `a honest` -> `an honest`

---

## The 4 Audit Cases

Your auditor runs these exact inputs. Output must match character for character.

**#1**
```
In:  If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?
Out: If I make you breakfast in bed just say thank you instead of: How did you get in MY HOUSE?
```

**#2**
```
In:  I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure
Out: I have to pack 5 outfits. Packed 26 just to be sure
```

**#3**
```
In:  Don not be sad ,because sad backwards is das . And das not good
Out: Don not be sad, because sad backwards is das. And das not good
```

**#4**
```
In:  harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '
Out: Harold Wilson: 'I am an optimist, but an optimist who carries a raincoat.'
```

---

## Week 1: Plan

No coding yet. Think, plan, analyse.

### Deliverables

- [ ] **PRD** — requirements, non-goals, acceptance criteria (use `templates/prd-template.md`)
- [ ] **Golden tests** — the 4 audit cases + 2-3 extra edge cases you design
- [ ] **Architecture choice** — Pipeline or FSM, with your rationale
- [ ] **Task decomposition** — 5-8 small testable tasks (use `templates/task-card-template.md`)

### Pace

| Days | Focus |
|------|-------|
| 1-2 | Read the project, draft PRD, list edge cases |
| 3-4 | Write golden tests, draft architecture comparison |
| 5-7 | Finalize PRD, task breakdown, prepare for peer audit |

---

## Phase 2 Gate: Peer Audit

You don't start coding until your plan is peer-audited and accepted.

### How it works

1. Pair up. One auditor, one auditee, then swap.
2. The auditee shows their Week 1 deliverables in the repo.
3. The auditor checks clarity and completeness, then writes the report.

### Audit report

Copy `templates/audit-report-template.txt` into the auditee's repo:

```
audit/peer-audit-YYYY-MM-DD-<auditor>.txt
```

You do NOT start Phase 2 until you have `Outcome: Accept` (recommended >= 8/10).

---

## Task Decomposition

Break the project into 5-8 small, testable tasks.

Each task must include: goal, acceptance criteria, and the test(s) you'll use to prove it works.

Put them in your repo:
```
tasks/TASK-01.txt
tasks/TASK-02.txt
...
```

---

## AI Usage

You completed the AI Piscine. Now apply it.

Starting from task decomposition, keep a running log in your repo:

```
ai/task-decomposition-index.txt
```

For each AI interaction, log:
- Tool/model used
- Date
- What you asked
- What you kept vs what you changed
- Which task card it affected

You must be able to explain every line of your code.

---

## Templates

Copy these into your project repo:

| File | Purpose |
|------|---------|
| `templates/prd-template.md` | Product Requirements Document |
| `templates/task-card-template.md` | One per task in your `tasks/` folder |
| `templates/audit-report-template.txt` | For peer audits |

# Pipeline vs. FSM (Finite State Machine)

## What Are These?

They are **programming models** — ways to **organize the logic of a program**.  
They are *not* languages, libraries, or frameworks.  
Think of them as different ways to **structure how your code thinks**.

---

## The Pipeline Model

A **pipeline** processes text in a sequence of steps —  
each step receives input, transforms it, and passes the result to the next step.

> Imagine a factory conveyor belt, except instead of assembling cars,  
> we're assembling *slightly less chaotic text*.

### Example Transformation
`"Hello world !"` → `"Hello World!"`

### Possible Pipeline Steps:
1. Remove extra spaces → `Hello world!`
2. Capitalize words → `Hello World!`

**Great for:**
- `(hex)`
- `(bin)`
- `(up)`
- `(low)`
- `(cap)`
- Article correction (`a` → `an`)

**Why?**  
These transformations are **local** — they don't require the program to remember previous context.

---

## The FSM Model

An **FSM (Finite State Machine)** behaves differently depending on **its current state**.  
It maintains **context**.

It includes:

- **States** — e.g., `OutsideQuotes`, `InsideQuotes`
- **Transitions** — rules defining when to switch states

### Example Need
Remove extra spaces *only* outside quoted text.

### FSM Logic Example

| State          | Action                            |
|----------------|-----------------------------------|
| OutsideQuotes  | Normalize spaces and punctuation   |
| InsideQuotes   | Preserve text exactly as-is        |

Switch state when encountering ' or ".
Yes, it’s basically a tiny mood-driven robot.

FSMs are useful when the program must **remember where it is** in the text.

---

**Best For:**  
Punctuation, grouped punctuation (`...`, `?!`), and quote handling —  
because these require **remembering context** in the text.

---

## Why Choose the Pipeline Here?

While both models are valid and educational, for this exercise I'm choosing **Pipeline** because:

### 1. Easier Testing & Debugging
Each step is **isolated** — perfect for TDD and incremental development.

### 2. Simpler Mental Model
No need to track states. Each function just transforms text → done.

### 3. Fully Sufficient for the Assignment
Even tricky cases like quotes can be handled with careful string logic.

---

## Pipeline Drawbacks

A pipeline **doesn't remember context**, so:

- Nested or deeply complicated quote/punctuation scenarios require extra caution.
- If your input starts behaving like Shakespeare or legal documents,  
  an FSM might eventually be the better therapist.

---
# Claude Instructions

Use the Rust skills installed under `~/.claude/skills/`.

For Rust work, start with:
- `rust-router` for routing questions and implementation work
- `coding-guidelines` for style and conventions
- `unsafe-checker` for any `unsafe` or FFI review
- `m01-ownership`, `m06-error-handling`, and `m07-concurrency` when those topics appear

Default Rust project settings:

```toml
[package]
edition = "2024"
rust-version = "1.85"

[lints.rust]
unsafe_code = "warn"

[lints.clippy]
all = "warn"
pedantic = "warn"
```

General rules:
- Prefer `Result` and `?` over `unwrap()` in library code.
- Every `unsafe` block must include a `// SAFETY:` comment.
- Keep Rust naming idiomatic: `snake_case` for functions and variables, `PascalCase` for types, `SCREAMING_SNAKE_CASE` for constants.

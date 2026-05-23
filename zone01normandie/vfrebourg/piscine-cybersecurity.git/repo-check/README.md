### secret-leak

#### Goal

Detect hardcoded secrets or sensitive data in a codebase.

#### Instructions

- Clone or simulate a repository with embedded secrets (API keys, tokens, passwords)
- Use tools like `gitleaks`, `truffleHog`, or `git-secrets`
- Document findings and provide cleaned version

#### Deliverables

- Scan report
- Cleaned code (if possible)
- Summary of risks and remediation

### dep-alert

#### Goal

Perform a security audit of project dependencies.

#### Instructions

- Use `npm audit`, `pip-audit`, `yarn audit`, or `safety`
- Identify vulnerable packages
- Explain their impact and how to fix or upgrade them

#### Deliverables

- Audit report
- Fix plan
- Updated `package.json`, `requirements.txt` or similar

### lint-sec

#### Goal

Add security-focused linting to a codebase.

#### Instructions

- Integrate tools like `bandit` (Python), `eslint-plugin-security`, or `semgrep`
- Run on existing code and interpret results
- Adjust rules or CI if necessary

#### Deliverables

- Linter configuration
- Output of scan
- Explanation of rules chosen and results
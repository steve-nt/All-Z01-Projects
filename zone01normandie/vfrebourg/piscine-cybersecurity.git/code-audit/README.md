### xss-find

#### Goal

Analyze a web application source code snippet and identify possible Cross-Site Scripting (XSS) vulnerabilities.

#### Instructions

You will be provided with a sample vulnerable HTML/JS/PHP/React file.

- Identify where user input is rendered without proper sanitization
- Describe how the vulnerability could be exploited
- Suggest how to fix it (e.g., escaping, CSP, sanitizers)

#### Deliverables

- Annotated source showing flaws
- Explanation of risks and how to fix


### sql-spot

#### Goal

Detect SQL injection vulnerabilities in a source code snippet.

#### Instructions

You are given insecure database queries in PHP, Python, or Node.js.

- Locate the insecure query
- Explain how the input could be injected
- Recommend parameterized queries or ORM alternatives

#### Deliverables

- Annotated code showing vulnerable queries
- A sample payload demonstrating injection (if safe)
- Suggested secure version

### auth-bypass

#### Goal

Identify logic flaws in authentication or session handling code.

#### Instructions

You are provided with backend snippets with potential logic bugs (e.g., flawed password check, missing token verification, weak session logic).

- Identify flawed conditions or checks
- Explain how an attacker might bypass login/session protection
- Suggest improvements

#### Deliverables

- Annotated code with issues explained
- Exploitation scenario
- Secure recommendation

### code-fix

#### Goal

Apply secure patches to flawed code provided.

#### Instructions

You are given insecure code (mix of flaws: XSS, SQLi, logic bugs). Your task:

- Patch the code to mitigate the vulnerabilities
- Explain what you changed and why
- Ensure the code still runs

#### Deliverables

- Fixed code version
- Commented summary of changes
- Optional: test output

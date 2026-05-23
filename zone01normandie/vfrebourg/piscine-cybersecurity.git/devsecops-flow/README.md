### pipeline-flaws

#### Goal

Analyze a given CI/CD pipeline and identify security flaws or missing controls.

#### Instructions

You'll be given a simplified CI/CD YAML config (GitHub Actions, GitLab, or similar). Your task is to:

- Identify common weaknesses (e.g., no secrets handling, no image scanning)
- List what could be exploited by attackers
- Suggest mitigation measures for each issue

#### Deliverables

- A markdown report listing:
  - The flaws found
  - Their severity
  - Suggested fixes

#### Resources

- OWASP CICD Top 10
- DevSecOps best practices

### scanner-add

#### Goal

Integrate at least one security scanner into a CI/CD workflow.

#### Instructions

Choose a tool such as:

- `Trivy`, `Bandit`, `Semgrep`, `Checkov`, `Snyk`, etc.

Set it up inside a pipeline config (GitHub Actions or GitLab preferred), and show that it runs correctly on a sample repo.

You can mock a repo with a vulnerable script or Dockerfile.

#### Deliverables

- Pipeline YAML with scanner integrated
- Screenshot or log output of the scan running
- Short explanation of the tool's purpose

### workflow-sec

#### Goal

Redesign a DevOps workflow with security in mind.

#### Instructions

Based on a typical pipeline (code > build > test > deploy), describe how you would embed security at each phase:

- Secrets detection
- Dependency scanning
- Linting and static analysis
- Container security
- Access controls
- Environment protections

#### Deliverables

- A markdown diagram or list of steps with associated tools and controls
- Short explanations of each security measure
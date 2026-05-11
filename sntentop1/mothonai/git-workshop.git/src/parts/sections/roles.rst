.. _roles:

#####
Roles
#####

There are mostly two roles on git repository:

1. maintainer and
2. contributor.

**********
Maintainer
**********

A maintainer is a role for someone that is in charge of which branches will make
it to the main branch. Assuming that the "main" branch is the one that keeps
track of the released version of a software or other type of product, a
maintainer's role is to make sure the repository is accessible for contributors.

More extensively:

- Review pull requests (PRs) or merge requests (MRs) for quality and consistency.
- Merge approved contributions while maintaining a clean commit history.
- Ensure contributions follow project coding standards and guidelines.
- Enforce coding standards, linters, and formatting rules.
- Ensure tests pass before merging changes.
- Monitor code coverage and refactor when necessary.
- Remove or archive outdated code.
- Structure the repository logically (folders, modules, docs).
- Maintain important files: `README.md`, `CONTRIBUTING.md`, `LICENSE`.
- Tag releases and maintain versioning (semantic versioning recommended).
- Triage issues: label, prioritize, and assign.
- Respond to questions, bug reports, or feature requests.
- Close stale or resolved issues.
- Guide contributors and provide constructive feedback.
- Welcome new contributors and enforce a code of conduct.
- Handle conflicts or disagreements diplomatically.
- Prepare and publish new releases.
- Update changelogs and release notes.
- Ensure proper Git tagging for releases.
- Monitor and fix security vulnerabilities.
- Keep dependencies up to date.
- Ensure backups and data integrity.
- Maintain user and developer guides, API references.
- Keep documentation up to date with code changes.

Many of those, in a project setting won't be a solo responsibility. For example,
it's considered a good practise to have a whole team review someone's commits
before proceeding to merge the branch onto the main one.

***********
Contributor
***********

On the other hand, contributor can do any kind of work they want.

- Create a personal copy of the repository to work on.
- Adhere to the repository’s rules, coding standards, and style guides.
- Work on isolated branches rather than directly on `main` or `master`.
- Write clear, descriptive commit messages and make atomic commits.
- Ensure your changes don’t break existing functionality.
- Keep your branch updated with the main repository to avoid conflicts.
- Submit changes for review, providing context and explanations.
- Review others’ PRs and provide constructive feedback.
- Update your PRs based on reviewer comments.
- Update README, docs, or comments as necessary.


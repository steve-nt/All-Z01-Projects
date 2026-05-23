# Cybersecurity Piscine — Quest & Exercise Overview

---

## Branch: `sec-ops`

### Quest: `intro-soc`

* `soc-roles`: SOC roles and responsibilities.
* `soc-tools`: Tools used by SOC teams (SIEM, SOAR, etc.).
* `incident-types`: Common types of incidents.
* `layered-defense`: Multi-layered defense principles.

### Quest: `phishing-mail`

* `spot-bait`: Identify phishing characteristics.
* `header-hunt`: Analyze email headers.
* `safe-open`: Open attachments securely.
* `write-alert`: Draft an internal incident alert.

### Quest: `ransomware-log`

* `log-signs`: Detect ransomware in logs.
* `ransom-behaviors`: Typical attack behavior.
* `host-isolate`: Isolate a compromised host.
* `ioc-notes`: Document IOCs.
* `recover-plan`: Simulate incident response.

### Quest: `data-leak`

* `leak-signs`: Detect leak symptoms.
* `exfil-methods`: Analyze exfiltration tactics.
* `dns-scan`: Investigate DNS logs.
* `conf-fix`: Correct insecure settings.

### Quest: `mini-playbook`

* `template-kit`: Create IR playbook template.
* `fill-steps`: Populate steps for a case.
* `share-doc`: Internal documentation sharing.

---

## Branch: `osint`

### Quest: `osint-basics`

* `what-osint`: Define OSINT in cyber.
* `osint-legal`: Legal limits overview.
* `toolbox`: Key OSINT tools.
* `osint-ethics`: Ethical boundaries.

### Quest: `fake-profile`

* `profile-check`: Spot fake profile traits.
* `bot-signs`: Detect bot-like behavior.
* `deep-dive`: Investigate account history.
* `timeline-build`: Build target timeline.

### Quest: `domain-check`

* `whois-scan`: Basic domain investigation.
* `history-dig`: Check domain history.
* `reputation-check`: Check IP/domain reputation.
* `url-hunt`: Investigate linked URLs.

### Quest: `photo-trace`

* `exif-pull`: Extract image metadata.
* `img-rev`: Reverse image search.
* `geo-clues`: Spot geo indicators.
* `source-hunt`: Trace back origin.

### Quest: `mini-investigation`

* `find-target`: OSINT on fictional target.
* `data-timeline`: Organize findings.
* `cross-ref`: Validate sources.
* `report-out`: Draft investigation summary.

---

## Branch: `legal`

### Quest: `gdpr-quiz`

* `gdpr-terms`: Key GDPR definitions.
* `rights-check`: Data subject rights.
* `roles-id`: Controller vs processor.
* `legal-bases`: Lawful processing grounds.

### Quest: `leak-case`

* `breach-story`: Analyze case scenario.
* `notify-dpa`: Simulate DPA report.
* `user-info`: Write breach email.
* `fix-plan`: Recommend remediations.

### Quest: `cookie-check`

* `banner-audit`: Cookie banner review.
* `consent-check`: Consent validation.
* `cookie-scan`: Site cookie inspection.

### Quest: `legal-risk`

* `who-liable`: Assign responsibility.
* `fine-types`: Understand GDPR penalties.
* `case-study`: Review known cases.

### Quest: `compliance-kit`

* `checklist`: Build a legal checklist.
* `doc-pack`: Prepare legal templates.
* `kit-ready`: Assemble compliance starter kit.

---

## Branch: `appsec`

### Quest: `devsecops-flow`

* `pipeline-flaws`: Audit insecure CI/CD.
* `scanner-add`: Add security tools.
* `workflow-sec`: Harden DevOps flow.

### Quest: `code-audit`

* `xss-find`: Detect XSS flaws.
* `sql-spot`: Spot SQL injections.
* `auth-bypass`: Authentication flaws.
* `code-fix`: Secure patching.

### Quest: `repo-check`

* `secret-leak`: Scan secrets in code.
* `dep-alert`: Dependency audit.
* `lint-sec`: Enforce secure linting.

### Quest: `owasp-top10`

* `owasp-match`: Map risks to OWASP Top 10.
* `owasp-quiz`: OWASP concept quiz.
* `real-case`: Analyze real attacks.

### Quest: `mini-ctf`

* `web-pwn`: Solve a web challenge.
* `crypto-fun`: Break a crypto task.
* `file-trace`: Investigate file metadata.
* `dataset-clues`: Hunt clues in open datasets.
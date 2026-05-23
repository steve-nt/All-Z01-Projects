## conf-fix

#### Functional

##### Has the student identified all insecure or risky settings in the configuration files?

##### Has the student provided corrected configuration snippets that address the security issues?

##### Are the proposed fixes compatible with service functionality (no overly restrictive changes)?

##### Has the student explained why each change improves security?

#### Bonus

###### + Does the student provide automated scripts or configuration management snippets (e.g., Ansible, Puppet) for applying the fixes?

###### + Does the student test the configuration changes and document the results?

###### + Are best practices and current security guidelines referenced in the corrections?

---

## dns-scan

#### Functional

##### Has the student identified suspicious DNS queries or patterns from the provided logs?

##### Has the student documented the indicators of compromise found in the DNS data?

##### Has the student proposed practical remediation or monitoring strategies for DNS traffic?

##### Did the student differentiate between normal and malicious DNS behavior?

#### Bonus

###### + Does the student use tools or commands to support their analysis (e.g., grep, Wireshark, custom scripts)?

###### + Is the analysis report clear, detailed, and actionable?

###### + Does the student explain DNS tunneling techniques and how to detect them?

---

## exfil-methods

# exfil-methods - Audit

- Student completed **at least 3** exercises from the malware-traffic-analysis.net training platform.
- Each exercise involved data exfiltration or suspicious outbound traffic.
- The student provided a short report with:
  - Links to exercises
  - Identified techniques
  - Proof of completion (e.g., screenshots or Wireshark output)

---

## leak-signs

#### Functional

##### Check `leak_findings.md` presence and naming.

##### Verify at least 3 suspicious log entries are listed.

##### For each finding:

- Filename and line number referenced?
- Clear explanation of suspicious behavior?
- Logical rationale provided?

##### Summary of findings included?

##### Are detections realistic (volumes, times, patterns)?

##### Correlation between logs done when relevant?

#### Bonus

###### +Possible causes suggested?

###### +Visual aids (tables, timelines)?

###### +Professional and concise English used?
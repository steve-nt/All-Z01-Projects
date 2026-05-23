### log-signs

#### Context

Ransomware attacks often leave traces in system and application logs before and during the encryption process. Being able to recognize those indicators early is critical to initiate a proper incident response.

In this exercise, you’ll analyze a set of log files and identify suspicious patterns or anomalies that may suggest the presence of ransomware.

#### Instructions

- Analyze the provided log files.
- Identify at least 3 potential indicators of ransomware activity (e.g. file renaming, suspicious processes, encryption tools).
- Document each finding with the line reference and a short justification.

#### Deliverable

- Submit a file named `findings.md` that contains:
  - A list of indicators found.
  - The line number where each was found.
  - A 1-2 sentence explanation for each.

- Example:
```
* Line 120: Process "encryptor.exe" started by unknown user — possible ransomware payload.
* Line 348: Mass file renaming with ".locked" extension — typical ransomware behavior.
* Line 812: User "admin" added to local group "Administrators" — suspicious privilege escalation.
```

#### 📂 Deliverables (logs to analyze)

Download the realistic log files below:
[👉 Download `ransomware_logs.zip`](https://zone01normandie.org/git/vfrebourg/piscine-cybersecurity/raw/branch/master/ransomware-logs/ransomware_logs.zip)

The ZIP archive includes:

* `system.log`
* `application.log`
* `security.log`

These logs contain the following suspicious activity:

* Execution of a suspicious file (`encryptor.exe`)
* Mass creation of `.locked` files
* Addition of a user to the administrators group
* Outbound network connections to an unidentified IP address

### ransom-behaviors

#### Context

Ransomware attacks often follow a recognizable pattern. Understanding these behaviors is essential to identify, prevent, and respond to ransomware threats.

As a SOC analyst in training, you're tasked with studying and identifying the typical behavior patterns of ransomware attacks. You'll document each phase and describe real-world indicators for each one.

#### Objectives

In this exercise, you will:

- Analyze the typical flow of a ransomware attack.
- Document each phase: Initial Access, Execution, Persistence, Lateral Movement, Encryption, and Exfiltration.
- Provide at least one real-world indicator or technique per phase.

#### Instructions

Create a markdown file named `ransomware_behaviors.md` with the following sections:

1. **Initial Access**  
   - Common vectors (phishing, RDP, etc.)
   - Example indicator

2. **Execution**  
   - Methods used to launch ransomware
   - Example tool or command

3. **Persistence**  
   - Techniques to maintain access
   - Example registry entry or scheduled task

4. **Lateral Movement**  
   - How ransomware spreads in a network
   - Example protocol or method

5. **Encryption**  
   - File targeting strategy
   - How ransom notes are deployed

6. **Exfiltration (if applicable)**  
   - What data is stolen and how
   - Example destination or tool

#### Deliverables

- A file named `ransomware_behaviors.md` with the 6 detailed sections.

#### Tips

- Refer to MITRE ATT&CK framework for accurate terminology.
- Use real ransomware case studies like WannaCry, LockBit, Maze, etc.

### host-isolate

#### Context

When a host is compromised, isolating it from the rest of the network is one of the first critical containment steps. This action prevents lateral movement, data exfiltration, or further spread of malware.

As a SOC analyst, you must understand and demonstrate how to perform host isolation using multiple tools and techniques, across different environments.

#### Objectives

In this exercise, you will:

- Explain the importance of host isolation during incident response.
- Describe and simulate how to isolate a host using different methods.
- Document the commands, interfaces, or tools used.

#### Instructions

Create a markdown file named `host_isolation.md` including the following sections:

1. **Theory**  
   - Why isolate a host?  
   - When should isolation be performed?  
   - Risks of improper isolation.

2. **Techniques**  
   - Describe at least two different methods to isolate a host:
     - One on a Windows environment.
     - One on a Linux environment.
     - Optionally: one using an EDR solution or firewall/network-based isolation.

3. **Commands and Tools**  
   - Include example commands or procedures to isolate a host (e.g., `netsh`, `iptables`, EDR GUI).
   - Document any response logs or verification commands.

4. **Simulation or Case Study**  
   - Briefly describe a realistic incident scenario where host isolation was critical.
   - What happened before/after isolation?

#### Deliverables

- A markdown file named `host_isolation.md` containing:
  - Theory section
  - Techniques section
  - Commands/tools section
  - Simulation or case study section

No script or executable is required for this exercise. Manual procedures and command documentation are sufficient.

#### Tips

- You can use real-world tools (CrowdStrike, SentinelOne, Windows Defender, etc.) if you're familiar.
- Avoid cutting the network stack in a way that prevents management or monitoring.

### ioc-notes

#### Context

Indicators of Compromise (IOCs) are forensic artifacts left behind by attackers during or after a cyberattack. They are essential for identifying, containing, and remediating incidents. Documenting and sharing IOCs correctly helps defenders recognize malicious activity more quickly and collaborate efficiently.

As a SOC analyst, it is crucial to identify and properly format IOCs discovered during investigations.

#### Objectives

In this project, you will:

- Collect IOCs related to a simulated attack scenario.
- Document these IOCs in a structured format.
- Understand the different types of IOCs and their role in threat intelligence.

#### Instructions

1. Read the following simulated incident summary:

> "A ransomware campaign was detected in the organization. The malware was delivered via email attachment. Once executed, it attempted to reach out to `http://malicious-domain.biz/update` and downloaded a second-stage payload. The executable was named `encryptor.exe`, located in `C:\Users\Public\Tools\`. A scheduled task named `WindowsUpdateChecker` was created to maintain persistence. The file had a SHA256 hash of `a3b5c9f3e29a...` and was signed with an untrusted certificate. The C2 IP observed in outbound traffic was `192.168.66.22`."

2. Based on this scenario, create a file named `ioc_report.md` and document at least 6 relevant IOCs, covering at least 3 different IOC types (e.g., file hash, domain, IP, file path, mutex, process name, etc.).

3. Structure your report with the following format:

```

Indicator Name: [e.g., Encryptor Executable]

* Type: [File hash / IP address / Domain name / etc.]
* Value: [IOC value]
* Description: [What this IOC represents or how it was found]
* Relevance: [Why this IOC is important]

```

4. Include a short introductory paragraph explaining what IOCs are and why they matter in incident response.

#### Deliverables

- A markdown file named `ioc_report.md` containing:
  - An introductory paragraph
  - At least 6 well-structured IOC entries

#### Tips

- Be concise but informative.
- Choose a variety of indicators to demonstrate a broad understanding.
- Use a consistent format to improve readability.

### recover-plan

#### Context

Following a ransomware incident, your organization needs to carry out proper incident response and recovery actions. A well-documented recovery plan helps teams coordinate effectively, reduce downtime, and prevent reinfection. This exercise aims to simulate your organization's response to a ransomware event.

#### Objectives

In this project, you will:

- Simulate a full incident response plan based on a ransomware infection.
- Document each phase of the response using the NIST framework (Preparation, Detection & Analysis, Containment, Eradication, Recovery, and Post-Incident Activity).
- Reflect on how to improve future preparedness.

#### Instructions

1. You are provided with the following simulated scenario:

> "An employee opened a malicious attachment named `invoice_details.docm`. The file contained a macro that executed `encryptor.exe`, encrypting files across the user's system and spreading via mapped drives. Several `.locked` files were found, and a ransom note was placed in each folder."

2. Based on this scenario, create a file named `incident_response_plan.md` and document your recovery actions across the following phases:
   - Preparation
   - Detection & Analysis
   - Containment
   - Eradication
   - Recovery
   - Post-Incident Activity

3. For each phase, answer:
   - What actions are taken?
   - Who is involved?
   - What tools or procedures are used?
   - What challenges or risks are identified?

4. Include a conclusion summarizing:
   - What went well
   - What needs improvement
   - Key lessons learned

#### Deliverables

- A markdown file named `incident_response_plan.md` containing:
  - 6 structured sections (one per NIST phase)
  - A final conclusion section

#### Tips

- Be realistic and operational — think like a SOC or IT team.
- Be specific when listing tools (EDR, backup systems, SIEM alerts, etc.)
- Don't forget user communication, documentation, and learning process.
### soc-roles

Welcome to the SOC. This project introduces the basic concepts of Security Operations Centers and the key roles involved.

#### Objectives

The goal of this project is to understand and identify the main responsibilities of a SOC team. You will read, analyze, and classify various security-related tasks into the appropriate SOC roles.

You will submit a classification of tasks and match each one with the correct role inside a SOC: Tier 1 Analyst, Tier 2 Analyst, Tier 3 (Threat Hunter), Incident Responder, or SOC Manager.

This project will not involve any code, but you will need to demonstrate precision, clarity, and a structured approach to handling incidents.

#### Introduction

- The 15 tasks you must classify are provided in the file: [tasks.txt](https://zone01normandie.org/git/vfrebourg/piscine-cybersecurity/raw/branch/master/intro-soc/tasks.txt)
- Each task represents a real-world responsibility typically handled inside a Security Operations Center.
- Your job is to assign the correct SOC role to each task.
- Each task must be matched to one and only one role.

#### Roles Overview

Here is a quick reminder of SOC roles:

- **Tier 1 Analyst**: Monitors alerts, escalates incidents, handles routine tasks.
- **Tier 2 Analyst**: Investigates escalated alerts, performs deeper analysis, correlates events.
- **Tier 3 / Threat Hunter**: Proactively hunts for threats, creates detection rules.
- **Incident Responder**: Coordinates response during a breach, containment and recovery.
- **SOC Manager**: Oversees SOC operations, team coordination, reporting.

#### Allowed resources

- Any public documentation (OWASP, CISA, ENISA, NIST, etc.)

#### Usage

```console
$ cat roles.txt
1. Tier 1 Analyst
2. Tier 3 / Threat Hunter
3. Incident Responder
...
15. SOC Manager
```

### soc-tools

Welcome to the SOC toolbox. In this project, you will explore the essential tools used by SOC teams and understand their purposes.

#### Objectives

You will map each SOC tool to its category and function. This exercise is designed to help you understand the core components of a SOC's technical stack: SIEM, SOAR, EDR, threat intelligence platforms, ticketing tools, etc.

Your deliverable will be a text file containing a classification table with two columns: tool name and associated category.

#### Introduction

- Your submission must be a text file named `tools.txt`.
- It must contain exactly 10 tools, each on its own line, following the usage format.
- Each tool must be correctly matched to one of the following categories:
  - SIEM
  - SOAR
  - EDR
  - Threat Intelligence Platform
  - Ticketing System

- No tool should appear more than once.

#### Allowed resources

- Official documentation of SOC tools
- Open-source intelligence (OSINT)
- Any reliable public documentation

#### Usage

```console
$ cat tools.txt
1. Splunk - SIEM
2. Cortex XSOAR - SOAR
3. CrowdStrike Falcon - EDR
4. MISP - Threat Intelligence Platform
5. TheHive - Ticketing System
...
10. Wazuh - SIEM
```

This project will help you learn about:

- SOC tooling ecosystem
- Tool categories and responsibilities
- How tools interact during incident workflows

### incident-types

Welcome to the frontline. In this exercise, you will explore the most common types of security incidents that SOC teams must handle.

#### Objectives

You will create a list of typical security incidents, along with a brief description and their potential impact. The goal is to get familiar with the incident taxonomy and learn to distinguish between different scenarios that a SOC analyst might face.

#### Introduction

- Your deliverable must be a text file named `incidents.txt`.
- It should list at least 8 different types of security incidents.
- For each incident, include:
  - A short name (one line)
  - A brief description (1-2 lines)
  - An example scenario (1 line)

- The format must follow the usage example exactly.

#### Allowed resources

- Public documentation (SANS, NIST)
- Security blogs
- SOC playbooks

#### Usage

```console
$ cat incidents.txt
1. Phishing
Description: A fraudulent attempt to obtain sensitive information via email.
Example: An employee receives an email mimicking Microsoft asking for password reset.

2. Malware Infection
Description: A system is compromised by malicious software.
Example: A user downloads a trojan disguised as a PDF invoice.

...

8. Misconfiguration
Description: Improper configuration of systems or services that exposes vulnerabilities.
Example: An S3 bucket with public read access contains internal documents.
```

This project will help you learn about:

- Real-world SOC incident categories
- Describing and identifying incidents
- Providing concrete, illustrative examples

### layered-defense

Layer upon layer, we build the fortress. In this exercise, you will explore the principles of multi-layered defense (also called "defense in depth") and understand how it applies to modern cybersecurity infrastructures.

#### Objectives

Your goal is to map out a basic multi-layered defense strategy, identifying key security controls at each layer. This will help you grasp the logic behind layered security and why it's critical for protecting systems.

#### Introduction

- Your deliverable must be a text file named `defense-layers.txt`.
- It should describe **at least 5 layers** of defense.
- For each layer, include:
  - The name of the layer
  - Its primary objective
  - 2 to 3 common security controls used at this layer

- Follow the usage format strictly.

#### Allowed resources

- OWASP
- NIST SP 800-53
- CIS Controls

#### Usage

```console
$ cat defense-layers.txt
1. Perimeter Layer
Objective: Prevent unauthorized access from external networks.
Controls: Firewall, IDS/IPS, Network ACLs

2. Network Layer
Objective: Segment and monitor internal traffic.
Controls: VLANs, Internal firewalls, Network monitoring tools

3. Host Layer
Objective: Secure individual devices.
Controls: Endpoint Protection, Patch Management, Device Hardening

...

5. Data Layer
Objective: Protect sensitive information.
Controls: Encryption, Access Control, DLP
```

This project will help you learn about:

- Defense in depth strategy
- Security control mapping
- Layer-specific objectives
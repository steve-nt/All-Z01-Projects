### spot-bait

#### Background

Phishing is one of the most common techniques used by attackers to deceive users into revealing sensitive information. Recognizing phishing emails is a fundamental skill in cybersecurity.

#### Objectives

In this exercise, you will analyze several email samples to identify common phishing indicators such as suspicious sender addresses, urgent language, fake links, and unusual attachments.

#### Instructions

- You will be provided with a set of **5 fake email samples** in `.eml` and/or `.txt` format (located in the `emails/` folder).
- For each email, identify **at least 3 phishing indicators** and justify why they are suspicious.
- Summarize your findings in a file named `phishing-analysis.txt`.

#### Deliverables

- A file named `phishing-analysis.txt` containing:
  - A numbered list from 1 to 5 (one for each email).
  - For each, at least three phishing red flags.
  - A short justification (1–2 sentences per flag).

#### Constraints

- Do not use generic lists or AI-generated answers.
- Stick to the actual email content: language, links, formatting, sender, etc.

#### Tips

- Look carefully at the sender’s domain (e.g., amaz0n.com vs amazon.com).
- Hover over links (or inspect in plaintext) for mismatches.
- Check if there's pressure language ("urgent", "act now").
- Misspellings or outdated logos are also common red flags.

#### Resources

All email samples are located in the `emails/` folder you can find [on this link](https://zone01normandie.org/git/vfrebourg/Custom-Cursus-Projects/src/branch/master/phishing-mail/emails)

### header-hunt

#### Instructions

In this exercise, you'll learn how to analyze email headers to detect potential phishing, spoofing, or suspicious behavior. Email headers contain vital information about the sender, server relays, and metadata that can help you identify forged or malicious messages.

You are provided with several raw email headers located in the `headers/` directory you can find [on this link](https://zone01normandie.org/git/vfrebourg/Custom-Cursus-Projects/src/branch/master/phishing-mail/headers)

Your task is to:

- Analyze each header and extract:
  - The **sender's real IP address**
  - The **originating domain**
  - Any **SPF/DKIM/DMARC** validation status
  - The **suspicious signs** (e.g., mismatched domains, forged return-paths)

You must fill in a report file named `analysis.md` at the root of your repo, formatted as follows:

```
# Header Analysis Report

## header1.txt
- IP: ...
- Domain: ...
- SPF: ...
- DKIM: ...
- DMARC: ...
- Suspicious signs: ...

## header2.txt
...

(etc.)
```

Use tools such as:

- [MXToolbox Email Header Analyzer](https://mxtoolbox.com/EmailHeaders.aspx)
- [Google Admin Toolbox Messageheader](https://toolbox.googleapps.com/apps/messageheader/)

You are allowed to copy/paste the headers into online tools but must **interpret the results** yourself.

#### Resources

- [How to read email headers](https://www.csoonline.com/article/3303361/how-to-read-email-headers.html)
- [Google Toolbox Messageheader tool](https://toolbox.googleapps.com/apps/messageheader/)
- [MXToolbox header tool](https://mxtoolbox.com/EmailHeaders.aspx)

### safe-open

#### Objective

Analyze the security of suspicious email attachments. Understand the risks of opening attachments directly and learn how to use secure analysis tools.

#### Context

You received an email with several attachments from an unknown sender. The files are of different types (executable, document with macros, script, archive).

#### Instructions

1. Analyze each file using antivirus scanners, sandboxes, or static analysis tools.
2. Identify potential risks associated with each file.
3. Prepare a concise report summarizing your findings and recommendations.
4. Do not open any files directly on your main machine.

#### Attachments to analyze

- file1.exe (Windows executable)
- file2.docm (Word document with macros)
- file3.ps1 (PowerShell script)
- file4.zip (archive containing a suspicious script)

### write-alert

#### Scenario

You are a SOC analyst who just detected a confirmed phishing attack targeting your company’s employees. Your task is to draft a clear and professional internal alert email to notify the relevant teams (IT, HR, Management) and give them initial instructions.

#### Instructions

- Write an email alert addressed to internal stakeholders.
- Include:
  - A concise description of the incident (what happened).
  - The potential impact on the company.
  - Immediate actions to take by the recipients.
  - Contact details for follow-up.

- The tone should be formal, precise, and clear.

#### Deliverable

- One email text file named `incident_alert.txt`.
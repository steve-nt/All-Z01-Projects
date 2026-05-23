## header-hunt

## header-hunt

### Expected deliverable

A file named `analysis.md` containing the analysis of each header located in the `headers/` directory.

### Evaluation checklist

- File `analysis.md` is present and correctly formatted
- Each provided header file (`header1.txt`, `header2.txt`, etc.) has its own section
- For each header:
  - The sender's IP address is identified
  - The originating domain is identified
  - SPF, DKIM, and DMARC results are mentioned (Pass / Fail / Absent)
  - Suspicious signs (if any) are listed clearly
- Signs of critical thinking are visible (e.g., explanation why a header is suspicious)

### Manual questions

1. What is the difference between SPF, DKIM, and DMARC? How do they complement each other?
2. How can an attacker spoof a sender domain despite SPF being "Pass"?
3. Why is the `Return-Path` field important in header analysis?
4. What are signs of a forged email header that tools might miss?

### Bonus (optional)

If the learner uses CLI tools (like `exim`, `swaks`, or Python scripts) to parse headers instead of online tools, award bonus points.

---

## safe-open

# Audit safe-open

#### Functional

##### Did you correctly identify the type of each attachment?

##### What methods did you use to analyze the security of each file? (e.g., antivirus scanner, sandbox, online tools)

##### For the executable and PowerShell script, did you detect any suspicious behavior (e.g., network connections, hidden commands)?

##### Explain why it is important not to open these files directly on your main machine.

##### Provide a screenshot or report from VirusTotal/Hybrid Analysis for at least one suspicious file.

##### What best practices would you recommend for handling unknown attachments in a professional environment?

#### Bonus

###### +Did you test multiple analysis tools (e.g., VirusTotal plus sandbox) to increase confidence in your analysis?

###### +Did you identify any evasion or obfuscation techniques in the suspicious files?

###### +Is your report clear, concise, and structured by analyzed file?

---

## spot-bait

## Files

- `phishing-analysis.txt` must be present.
- `emails/` folder must contain 5 files (email1.txt to email5.txt).

## Evaluation Steps

1. Confirm presence of `phishing-analysis.txt`.
2. Check each of the 5 email entries:
   - Are 3 distinct indicators mentioned per email?
   - Are the justifications relevant to the content?
   - Do indicators correspond to actual clues (e.g., suspicious link, wrong domain, emotional trigger)?

## Scoring

- ✅ 3 phishing signs + justification = +2 pts/email
- ✅ High accuracy and specificity = +1 pt/email
- ❌ Copy/paste or vague answers = -1 pt/email

## Bonus

- +1 point if learner adds a real phishing email example (anonymized) and analysis at the end of the file.

---

## write-alert

#### Functional

##### Does the email clearly describe the incident in the first paragraph?

##### Are the potential impacts to the company clearly stated?

##### Are the immediate recommended actions for recipients clear and actionable?

##### Is the tone professional and suitable for internal communication?

##### Does the alert include contact information for follow-up questions?

#### Quality

##### Is the email free of spelling and grammar mistakes?

##### Is the message concise, without unnecessary jargon or verbosity?

#### Bonus

##### + Does the alert anticipate possible questions or concerns the recipients might have?

##### + Is the alert structured for easy reading (paragraphs, bullet points)?

##### + Would you consider this alert effective in a real SOC environment?
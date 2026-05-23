### conf-fix

In this exercise, you will analyze a set of system configuration files to identify insecure or misconfigured settings that could lead to data leaks, unauthorized access, or other security vulnerabilities.

Your tasks are:

- Review the provided configuration files from various services (e.g., SSH, Apache, FTP).
- Identify insecure settings such as weak permissions, enabled insecure protocols, or default credentials.
- Propose and implement corrected configurations to harden the system.
- Document the changes and explain why they improve security.

#### Deliverables (configuration files)

You will find a ZIP archive with several configuration files to analyze and fix:  

👉 Download the file included [on this link](https://zone01normandie.org/git/vfrebourg/Custom-Cursus-Projects/src/branch/master/data-leak/conf-fix)
- `sshd_config` (OpenSSH server configuration)  
- `apache2.conf` (Apache HTTP server configuration)  
- `vsftpd.conf` (FTP server configuration)

The files contain intentionally insecure settings such as:  
- PermitRootLogin enabled  
- Weak SSL/TLS settings  
- Anonymous FTP access allowed  
- Default or empty passwords

### dns-scan

In this exercise, you will investigate DNS logs to detect suspicious or malicious activity that may indicate a security incident, such as data exfiltration, command-and-control (C2) communication, or reconnaissance.

You will receive a set of DNS log files simulating normal and abnormal DNS query patterns. Your tasks are:

- Analyze the DNS logs to identify unusual DNS queries or patterns that could indicate malicious behavior.
- Look for DNS tunneling attempts, suspicious domain names, or unusual query frequencies.
- Document your findings, including which entries look suspicious and why.
- Suggest possible remediation steps or detection strategies to monitor DNS traffic effectively.


#### 📂 Deliverables (DNS logs to analyze)

You will find a ZIP archive containing realistic DNS logs for analysis [on this link](https://zone01normandie.org/git/vfrebourg/Custom-Cursus-Projects/src/branch/master/data-leak/dns-scan)

The archive includes:  
- `dns_queries.log`  
- `dns_responses.log`  

The logs contain:  
- Regular DNS queries from internal users  
- Suspicious queries to newly registered or uncommon domains  
- DNS TXT records potentially used for data exfiltration  
- High-frequency DNS requests from a single host  

#### Notes

- Focus on understanding DNS query patterns and anomalies.  
- The logs simulate a realistic environment with both benign and malicious DNS activity.  
- Practical detection and response recommendations are important parts of this exercise.

### exfil-methods

#### Goal

Understand and analyze common data exfiltration techniques using real-world PCAP samples.

#### Instructions

1. Visit the following external platform for PCAP analysis training:  
   https://www.malware-traffic-analysis.net/training-exercises.html

2. Complete **at least 3** training exercises from the platform that focus on **data exfiltration or suspicious traffic**.

3. Refer to the platform's tutorials if needed:  
   https://www.malware-traffic-analysis.net/tutorials.html

4. Document the techniques identified in each exercise and reflect on how exfiltration was performed or attempted.

#### Deliverables

- A report (PDF or Markdown) summarizing the 3 exercises, with:
  - The links to each chosen exercise
  - The exfiltration method involved
  - Screenshots or notes proving completion

### leak-signs

#### Context

Data leaks can happen through cloud misconfigurations, insider threats, or stolen credentials. Detecting early signs from logs is essential to prevent major damage.

#### Objectives

- Analyze provided logs to identify potential data leak indicators.
- Document suspicious entries with explanations.
- Correlate information across logs.

#### Instructions

1. Download and analyze the provided logs [on this link](https://zone01normandie.org/git/vfrebourg/Custom-Cursus-Projects/src/branch/master/data-leak/leak-signs)

2. Look for:

   - Large or unusual data downloads.
   - Access during unusual hours.
   - Multiple downloads of sensitive files.
   - Data sent to unknown external IP addresses.
   - Abnormal VPN behavior.

3. Create `leak_findings.md` documenting:

   - Each suspicious entry found.
   - Explanation why it’s suspicious.
   - Filename and line number references.
   - Summary of conclusions.

#### Deliverables

- `leak_findings.md`

#### Tips

- Justify why an entry is suspicious.
- Correlate logs when possible.
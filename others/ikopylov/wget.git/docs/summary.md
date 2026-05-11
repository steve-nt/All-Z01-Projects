# GNU Wget – Summary

## 🧠 What Wget is
- A **command-line tool** used to download files from the internet.  
- Supports **HTTP, HTTPS, and FTP** protocols.  
- Designed to work **non-interactively** (runs in the background without user input).  

---

## ⚙️ Core Features
- **Background downloading** → start a download and disconnect; it continues.  
- **Recursive downloads** → can download entire websites by following links.  
- **Resume support** → continues interrupted downloads instead of restarting.  
- **Mirroring** → keeps local copies of sites or FTP directories.  
- **Offline browsing** → can rewrite links so downloaded pages work locally.  
- **Robustness** → retries automatically on network failures.  

---

## 🌐 Networking Capabilities
- Works with **proxies and firewalls**.  
- Supports **IPv6 and IPv4**.  
- Handles both **HTTP and FTP advanced options** (headers, authentication, SSL/TLS, etc.).  

---

## 🛠️ How You Use It

Basic command:

    wget [options] [URL]

- You can customize behavior with:
  - **command-line options**
  - or a config file (`.wgetrc`)  

---

## ⚡ Advanced Capabilities
- **Time-stamping** → downloads only updated files  
- **Filtering rules** → choose which files/links to download  
- **Automation-friendly** → ideal for scripts and cron jobs  
- **Highly configurable** → almost everything can be tuned  

---

## 📌 Big Picture

Think of Wget as:

> a **scriptable, reliable downloader** that can fetch anything from a single file to an entire website — automatically and without supervision.


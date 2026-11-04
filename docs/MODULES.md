# ROAR Modules Reference

This document outlines the exhaustive list of modules embedded directly within ROAR, including how they map explicitly to user interface panels and what artifact outputs they generate.

## 1. Mock Data / Test Targets
ROAR includes a hard-routed demo trap-door configured for local developer testing. Submitting any of the following targets intercepts actual system tools and initiates an immediate generation of synthetic mock payloads to validate your UI state natively across every artifact:
*   `keyword.com`
*   `0x88.roar`
*   `demo.roar.com`

### Generated Mock Payloads
The following items explicitly describe what is returned when the internal test suite fires.

| UI Tab | Module Name | Mock Target Output | Vulnerability Type / Finding | Severity |
| :--- | :--- | :--- | :--- | :--- |
| **Links** | **URL-Enum** | `https://0x88.roar/api`<br>`https://0x88.roar/admin`<br>`https://0x88.roar/test` | *URL-Enum* | **INFO** |
| **Links** | **JS-Enum** | `https://0x88.roar/main.js`<br>`https://0x88.roar/vendor.chunk.js` | *JS-Enum* | **INFO** |
| **Vulnerabilities** | **Nuclei** | `https://0x88.roar/api` | `cve-2023-1000` | **HIGH** |
| **Vulnerabilities** | **S3 Scan** | `0x88.roar-bucket` | `s3-scan-open` | **CRITICAL** |
| **Vulnerabilities** | **Backup Files** | `https://0x88.roar/backup.zip` | `backup-detection` | **HIGH** |
| **Vulnerabilities** | **Port Scan** | `0x88.roar` | `Open Port 8080 (http-alt)` | **INFO** |
| **Vulnerabilities** | **JS Analysis** | `https://0x88.roar/main.js` | `[AWS API Key]: AKIA1234567890` | **HIGH** |
| **Vulnerabilities** | **ZeroDays** | `https://0x88.roar` | `CVE-2024-XXXX Node.js Remote Code Execution` | **CRITICAL** |
| **Vulnerabilities** | **AEM Scan** | `https://0x88.roar/aem` | `AEM Default Credentials` | **CRITICAL** |
| **Vulnerabilities** | **Misconfig** | `https://0x88.roar/.git` | `Exposed Git Directory` | **MEDIUM** |
| **Vulnerabilities** | **Dep Confusion** | `https://0x88.roar/wp-content` | `WordPress Missing Theme` | **MEDIUM** |
| **Vulnerabilities** | **Dep Confusion** | `package.json` | `Dependency Confusion in 'internal-core'` | **HIGH** |
| **Vulnerabilities** | **FFUF Fuzzing** | `https://0x88.roar/secret-dir` | `Hidden Directory` | **MEDIUM** |
| **Vulnerabilities** | **XSS Detection** | `https://0x88.roar/?q=test` | `Reflected XSS in 'q' parameter` | **HIGH** |
| **Vulnerabilities** | **XSS Detection** | `https://0x88.roar/?id=test` | `Reflection: true (kxss)` | **HIGH** |
| **Vulnerabilities** | **SQL Injection** | `https://0x88.roar/?id=1` | `Time-based blind SQLi (MySQL)` | **CRITICAL** |
| **Vulnerabilities** | **GF Patterns** | `https://0x88.roar/?action=test` | `Potential XSS Endpoint` | **INFO** |
| **Vulnerabilities** | **DNS Takeover** | `https://test.0x88.roar` | `Dangling CNAME` | **CRITICAL** |
| **Vulnerabilities** | **Cloudflare DNS** | `cloud.0x88.roar` | `Cloudflare Error 1016 (Dangling Record)` | **HIGH** |
| **Vulnerabilities** | **Github Scan** | `https://github.com/org/repo` | `Leaked Stripe API Key in commit` | **CRITICAL** |

---

## 2. Active Production Modules
These modules execute the active payloads listed above dynamically against real targets dynamically mapping to exactly the identical panels displayed.

> *Maintained natively by ROAR's structural parser backend (`inferModuleFromFileName`).*

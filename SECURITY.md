# Security Policy

## Supported Versions

Security fixes are provided for the latest stable release only. Please update to
the newest version before reporting an issue.

| Version | Supported |
| ------- | --------- |
| latest  | Yes       |
| older   | No        |

## Reporting a Vulnerability

If you believe you found a security vulnerability, please report it through
GitHub Security Advisories:

https://github.com/serenityzn/ProjectZomboidAutoSave/security/advisories/new

If GitHub Security Advisories are not available, open a GitHub issue with a
short description and mark it as security-related. Please avoid posting exploit
details publicly if the issue could harm users.

Please include:

- A clear description of the issue
- Steps to reproduce it
- The affected operating system and app version
- Any relevant logs, screenshots, or sample files

I will try to review valid reports as soon as possible. If the report is
accepted, I will work on a fix and publish a patched release. If the report is
not considered a vulnerability, I will explain why.

## Scope

This project is a local backup utility. It does not intentionally collect data,
send telemetry, or make network requests. Security reports are most useful when
they relate to local file handling, archive extraction, dependency issues,
release binaries, or build/release automation.

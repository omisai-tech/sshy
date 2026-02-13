# Security Policy

We welcome and value responsible disclosure of security issues.

If you believe you have discovered a security vulnerability in any Omisai Technologies product, service, or repository, please report it to us privately.

## Security Best Practices

When using sshy, we recommend:

### SSH Key Management
- Use strong, unique SSH keys for each server
- Regularly rotate SSH keys
- Use SSH key passphrases
- Store keys securely (not in version control)

### Configuration Security
- Keep servers.yaml in a secure, access-controlled location
- Use local.yaml for personal/private server configurations
- Avoid storing sensitive information in shared configuration files
- Regularly audit your server configurations

### General Security
- Keep sshy updated to the latest version
- Use SSH features like StrictHostKeyChecking and UserKnownHostsFile
- Be cautious with SSH agent forwarding
- Use VPNs or bastion hosts for additional security layers

## Preferred reporting method

Please use our dedicated security reporting form:

[https://omisai.com/security](https://omisai.com/security)

This form is monitored by our security team and ensures your report is handled confidentially.

## Alternative reporting method

If you are unable to use the form, you may contact us by email:

[security@omisai.com](mailto:security@omisai.com)

## What to include

To help us investigate effectively, please include as much of the following as possible:

- Description of the vulnerability
- Steps to reproduce or proof-of-concept
- Affected systems, versions, or endpoints
- Any relevant logs, screenshots, or tool output
- Your environment (OS, browser, runtime, etc.)

## Our commitment

We strive to:

- Acknowledge reports in a timely manner
- Assess and validate reported issues
- Prioritize remediation based on impact
- Communicate status when appropriate
- Credit researchers for valid disclosures (if desired)

Please avoid public disclosure until the issue has been reviewed and resolved.

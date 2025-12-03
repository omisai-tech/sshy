
If you discover a security vulnerability in sshy, please report it to security@omisai-tech.com instead of creating a public issue.

We will acknowledge your report within 48 hours and provide a more detailed response within 7 days indicating our next steps.

We will keep you informed about our progress throughout the process of fixing the vulnerability.

## Scope

This security policy applies to the sshy CLI tool and its source code.

Third-party dependencies are not covered by this policy.If you discover a security vulnerability in sshy, please report it to security@omisai-tech.com instead of creating a public issue.

We will acknowledge your report within 48 hours and provide a more detailed response within 7 days indicating our next steps.

We will keep you informed about our progress throughout the process of fixing the vulnerability.

## Scope

This security policy applies to:

- The sshy CLI tool and its source code
- Official releases and packages
- Documentation and examples in this repository

### Out of Scope

This policy does not apply to:

- Third-party dependencies (please report to the respective maintainers)
- Configuration files or user data
- Issues in other SSH clients or servers
- General best practices or educational content

## Vulnerability Classification

We classify vulnerabilities using the following severity levels:

- Critical: Remote code execution, privilege escalation, data theft
- High: Significant security impact, bypass of security controls
- Medium: Limited security impact, information disclosure
- Low: Minor security issues, edge cases

## Disclosure Process

1. Report: You report the vulnerability to us
2. Acknowledgment: We acknowledge receipt within 48 hours
3. Investigation: We investigate and validate the vulnerability
4. Fix: We develop and test a fix
5. Release: We release the fix and security advisory
6. Public Disclosure: We coordinate public disclosure with you

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

## Security Updates

We will release security updates as needed. Security fixes will be:

- Documented in release notes
- Marked with appropriate security advisories
- Coordinated with CVE assignments when applicable

## Recognition

We appreciate security researchers who help keep our users safe. With your permission, we will acknowledge your contribution in our security advisory.

## Contact

For security-related questions or concerns:
- Email: security@omisai-tech.com
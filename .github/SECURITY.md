
If you discover a security vulnerability in sshy, please report it to security@omisai.com instead of creating a public issue.

We will acknowledge your report within 48 hours and provide a more detailed response within 7 days indicating our next steps.

We will keep you informed about our progress throughout the process of fixing the vulnerability.

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

## Contact

For security-related questions or concerns:
- Email: security@omisai.com
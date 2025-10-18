# Security Policy

## Supported Versions

The PhD Dissertation Planner project takes security seriously. We actively support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 5.1.x   | :white_check_mark: |
| 5.0.x   | :white_check_mark: |
| < 5.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in the PhD Dissertation Planner, please help us by reporting it responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities by emailing:
- **Email**: [security@phd-planner.dev](mailto:security@phd-planner.dev)
- **Subject**: `[SECURITY] Vulnerability Report - PhD Dissertation Planner`

### What to Include

When reporting a security vulnerability, please include:

1. **Description**: A clear description of the vulnerability
2. **Steps to Reproduce**: Detailed steps to reproduce the issue
3. **Impact**: Potential impact and severity of the vulnerability
4. **Affected Versions**: Which versions are affected
5. **Environment**: Your operating system, Go version, and other relevant details
6. **Contact Information**: How we can reach you for follow-up questions

### Response Timeline

We will acknowledge your report within 48 hours and provide a more detailed response within 7 days indicating our next steps.

We will keep you informed about our progress throughout the process of fixing the vulnerability.

### Disclosure Policy

- We follow a 90-day disclosure timeline from the initial report
- We will credit you (if desired) in our security advisory
- We will not disclose vulnerability details until a fix is available
- We may delay disclosure for critical infrastructure vulnerabilities

## Security Considerations

### Data Handling

The PhD Dissertation Planner processes timeline data that may contain sensitive academic information. Consider the following security practices:

- **Input Validation**: All CSV and configuration inputs are validated
- **No Network Transmission**: The tool operates locally and does not transmit data over networks
- **File Permissions**: Generated files inherit appropriate permissions from the working directory
- **Temporary Files**: Sensitive data in temporary files is properly cleaned up

### Dependencies

We regularly update our dependencies to address security vulnerabilities:

- **Go Modules**: Dependencies are managed through Go modules with regular updates
- **LaTeX Distribution**: Users should keep their TeX distribution updated
- **Python Dependencies**: PDF processing libraries should be kept current

### Best Practices for Users

1. **Keep Dependencies Updated**:
   ```bash
   go mod tidy
   go get -u ./...
   ```

2. **Use Secure File Permissions**:
   ```bash
   # Set appropriate permissions on input files
   chmod 600 your_timeline.csv
   ```

3. **Regular Backups**: Backup your timeline data regularly

4. **Environment Isolation**: Consider running in a containerized environment

### Known Security Considerations

#### LaTeX Processing
- LaTeX processing can execute arbitrary commands through `\write` and other primitives
- The tool generates LaTeX code but does not include user-controlled content that could lead to code injection
- Users should be cautious with LaTeX distributions that include shell escape enabled

#### File System Access
- The tool reads from and writes to the local file system
- Generated files are written to the `generated/` directory
- No elevation of privileges occurs during normal operation

#### CSV Processing
- CSV files are parsed using Go's standard `encoding/csv` package
- Input validation prevents common CSV injection attacks
- Large CSV files are processed in memory (consider file size limits)

## Security Updates

Security updates will be:
- Released as patch versions (e.g., 5.1.1, 5.1.2)
- Documented in the CHANGELOG.md with appropriate severity indicators
- Announced through GitHub Security Advisories
- Tagged with appropriate CVSS scores when applicable

## Contact

For security-related questions or concerns:
- **Security Issues**: [security@phd-planner.dev](mailto:security@phd-planner.dev)
- **General Support**: [support@phd-planner.dev](mailto:support@phd-planner.dev)
- **GitHub Issues**: For non-security related issues

Thank you for helping keep the PhD Dissertation Planner secure!

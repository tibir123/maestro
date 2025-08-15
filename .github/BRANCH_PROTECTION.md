# Branch Protection Rules

This document outlines the branch protection rules configured for the Maestro repository to ensure code quality and prevent issues from reaching the main branch.

## Main Branch Protection

The following rules are enforced on the `main` branch:

### Required Checks
- **CI Complete** - All CI workflow jobs must pass
- **Lint** - Code must pass linting checks
- **Test** - All tests must pass with >80% coverage
- **Build** - All binaries must build successfully
- **Security** - Vulnerability checks must pass

### Pull Request Requirements
- **Require pull request reviews**: 1 approval required
- **Dismiss stale reviews**: Enabled (when new commits are pushed)
- **Require review from code owners**: Enabled (when CODEOWNERS file exists)
- **Restrict dismissal of reviews**: Repository admins only

### Status Check Requirements
- **Require status checks to pass**: Enabled
- **Require up-to-date branches**: Enabled
- **Status checks that must pass**:
  - `ci-complete` (final check that all other jobs passed)
  - Individual jobs are implicit dependencies

### Additional Restrictions
- **Require conversation resolution**: All review conversations must be resolved
- **Require signed commits**: Recommended for enhanced security
- **Require linear history**: Prevents merge commits
- **Restrict pushes that create files**: Files can only be added via PR

### Administrative Settings
- **Include administrators**: Repository admins are subject to these rules
- **Allow force pushes**: Disabled
- **Allow deletions**: Disabled

## Repository Settings

### Merge Configuration
- **Allow squash merging**: Enabled (default and recommended)
- **Allow merge commits**: Disabled (to maintain linear history)
- **Allow rebase merging**: Enabled
- **Auto-merge**: Enabled for approved PRs
- **Automatically delete head branches**: Enabled

## Setting Up Branch Protection

### Via GitHub UI
1. Go to repository Settings → Branches
2. Click "Add rule" for the `main` branch
3. Configure the following settings:

#### General Settings
- [x] Restrict pushes that create files
- [x] Require a pull request before merging
  - [x] Require approvals: 1
  - [x] Dismiss stale pull request approvals when new commits are pushed
  - [x] Require review from code owners
- [x] Require status checks to pass before merging
  - [x] Require branches to be up to date before merging
  - Required status checks:
    - `ci-complete`
- [x] Require conversation resolution before merging
- [x] Require signed commits (recommended)
- [x] Require linear history
- [x] Do not allow bypassing the above settings
- [x] Restrict pushes that create files

### Via GitHub CLI
```bash
# Set up branch protection rules
gh api repos/madstone-tech/maestro/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["ci-complete"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}' \
  --field restrictions=null \
  --field required_linear_history=true \
  --field allow_force_pushes=false \
  --field allow_deletions=false \
  --field required_conversation_resolution=true
```

## Testing Branch Protection

### Verify Protection Works
1. **Test direct push to main** (should fail):
   ```bash
   git checkout main
   echo "test" > test.txt
   git add test.txt
   git commit -m "test direct push"
   git push origin main  # Should be rejected
   ```

2. **Test PR without required checks** (should block merge):
   - Create PR without passing CI
   - Merge button should be disabled with explanation

3. **Test PR with failing tests** (should block merge):
   - Create PR that breaks tests
   - CI should fail and prevent merge

### Status Check Verification
The CI workflow includes a final `ci-complete` job that verifies all required jobs passed:
- Lint job
- Test job (with coverage threshold)
- Build job
- Security job

## Quality Gates

### Code Coverage
- Minimum 80% test coverage required
- Coverage is calculated and enforced in CI
- Coverage reports uploaded to Codecov

### Code Quality
- All code must pass `golangci-lint` checks
- `go vet` static analysis must pass
- `go mod tidy` check ensures clean dependencies

### Security
- `govulncheck` scans for known vulnerabilities
- No known security issues allowed

### Build Verification
- All binaries must compile successfully
- Basic CLI functionality must work

## Troubleshooting

### Common Issues
1. **Tests failing**: Check test output in CI, fix failing tests
2. **Coverage too low**: Add tests to increase coverage above 80%
3. **Lint failures**: Run `make lint` locally and fix issues
4. **Build failures**: Run `make build` locally and resolve compilation errors

### Emergency Procedures
Repository administrators can temporarily disable protection rules in case of emergency, but this should be:
1. Documented with reason
2. Re-enabled as soon as possible
3. Followed by immediate review of any changes made

## Future Enhancements

### Potential Additions
- Code owners file (`.github/CODEOWNERS`)
- Additional status checks for specific components
- Integration test requirements
- Documentation checks
- Dependency vulnerability scanning

### Monitoring
- Regular review of failed PR attempts
- Analysis of common protection rule violations
- Updates to rules based on project evolution
# Pull Request

## Description
<!-- Provide a brief description of the changes in this PR -->

## Type of Change
<!-- Mark the type of change with an [x] -->
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Test improvement

## Related Issues
<!-- Link to any related issues using #issue_number -->
Closes #
Related to #

## Implementation Details
<!-- Describe the technical implementation and any architectural decisions -->

### Domain Layer Changes
- [ ] No domain layer changes
- [ ] Added new entities or value objects
- [ ] Modified repository interfaces
- [ ] Updated domain errors or business rules

### Infrastructure Changes
- [ ] No infrastructure changes
- [ ] AppleScript template changes
- [ ] gRPC protocol updates
- [ ] Database/cache modifications

### Application Layer Changes
- [ ] No application layer changes
- [ ] New command/query handlers
- [ ] Service modifications
- [ ] Session management updates

### Presentation Layer Changes
- [ ] No presentation changes
- [ ] CLI command updates
- [ ] TUI view modifications
- [ ] API endpoint changes

## Testing
<!-- Describe the testing strategy and coverage -->

### Test Coverage
- [ ] Unit tests added/updated for new functionality
- [ ] Integration tests added/updated
- [ ] Test coverage maintains >80% threshold
- [ ] All existing tests pass

### Manual Testing
- [ ] Tested with Music.app locally
- [ ] Verified CLI commands work correctly
- [ ] Tested error scenarios
- [ ] Performance tested (if applicable)

### Test Commands
```bash
# Commands used to test the changes
make test
make build
./bin/maestro --help
# Add specific test commands here
```

## Performance Impact
<!-- Describe any performance implications -->
- [ ] No performance impact
- [ ] Performance improvement
- [ ] Potential performance impact (explain below)

<!-- If there's a performance impact, describe it -->

## Breaking Changes
<!-- List any breaking changes and migration steps -->
- [ ] No breaking changes
- [ ] Breaking changes (list below)

<!-- If there are breaking changes, list them and provide migration guidance -->

## Checklist
<!-- Ensure all items are completed before submitting -->

### Code Quality
- [ ] Code follows project conventions and architecture patterns
- [ ] Code is well-documented with clear comments
- [ ] Error handling is comprehensive
- [ ] Logging is appropriate and structured
- [ ] No hardcoded values or magic numbers

### Domain-Driven Design
- [ ] Changes respect domain boundaries
- [ ] Business logic is in the domain layer
- [ ] Infrastructure dependencies are properly abstracted
- [ ] Repository pattern is followed correctly

### Security
- [ ] No sensitive information exposed
- [ ] Input validation is implemented
- [ ] Certificate handling is secure (if applicable)
- [ ] No SQL injection or similar vulnerabilities

### Documentation
- [ ] README updated if needed
- [ ] CLAUDE.md updated if architecture changes
- [ ] Inline documentation added for complex logic
- [ ] API documentation updated (if applicable)

### Deployment
- [ ] Changes are backward compatible
- [ ] Migration scripts provided (if needed)
- [ ] Configuration changes documented
- [ ] No hardcoded environment-specific values

## Additional Notes
<!-- Any additional context, concerns, or notes for reviewers -->

## Screenshots/Demos
<!-- Add screenshots or demo outputs if applicable, especially for UI changes -->

## Reviewer Notes
<!-- Specific areas you'd like reviewers to focus on -->
- [ ] Please review the domain logic in [specific file]
- [ ] Check the error handling in [specific scenario]
- [ ] Verify the performance of [specific operation]
- [ ] Review the test coverage for [specific component]

---

**For Reviewers:**
- Ensure all CI checks pass before approving
- Verify test coverage meets the 80% threshold
- Check that domain-driven design principles are followed
- Confirm error handling and logging are appropriate
# Jira Process Workflow

This document outlines the agile workflow processes used to track developments in the SETA DAM backend.

## Issue States
- **Backlog**: Features, improvements, or bugs waiting to be prioritized.
- **Ready for Development**: Prioritized tasks ready for assignment.
- **In Progress**: Active work is being done. Branch naming convention: `feature/JIRA-ID-description` or `bugfix/JIRA-ID-description`.
- **In Code Review**: Pull requests submitted. Requires passing CI/CD tests and approval.
- **In QA / Testing**: Deployed to staging environment, ready for manual check.
- **Done**: Merged to main and deployed to production.

## Git Commit Guidelines
Commits should reference the Jira task:
```text
[DAM-123] Implement folder cycle checks and CTE queries
- Add IsDescendantOf query to folder repository
- Handle ErrCycleDetected error responses in HTTP handler
```

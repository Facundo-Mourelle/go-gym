# Git Practices Skill - Examples

## Example 1: Feature Completion

**User:** "I've completed the authentication feature"

**Assistant:**
```
I've staged the authentication feature files:

- internal/service/auth_service.go
- internal/api/handler/auth_handler.go
- internal/api/middleware/auth.go
- internal/service/auth_service_test.go

Proposed commit message:
feat: Implement authentication service

- User registration with password hashing
- Login with JWT token generation
- Token validation middleware
- 6 new tests covering auth flow

Exact command to run:
git commit -m "feat: Implement authentication service" -m "- User registration with password hashing
- Login with JWT token generation
- Token validation middleware
- 6 new tests covering auth flow"

Ready to commit? (yes/no)
```

**User:** "yes"

**Assistant:** [Executes the commit]

---

## Example 2: Bug Fix

**User:** "Fixed the login bug"

**Assistant:**
```
I've staged the bug fix files:

- internal/service/auth_service.go
- internal/api/handler/auth_handler.go

Proposed commit message:
fix: Resolve login validation issue

- Incorrect email validation regex
- Added proper RFC 5322 validation
- Added test for edge cases

Exact command to run:
git commit -m "fix: Resolve login validation issue" -m "- Incorrect email validation regex
- Added proper RFC 5322 validation
- Added test for edge cases"

Ready to commit? (yes/no)
```

---

## Example 3: Documentation Update

**User:** "Updated the README"

**Assistant:**
```
I've staged the documentation files:

- README.md
- docs/installation.md

Proposed commit message:
docs: Update installation instructions

- Added Docker deployment steps
- Updated environment variables
- Added troubleshooting section

Exact command to run:
git commit -m "docs: Update installation instructions" -m "- Added Docker deployment steps
- Updated environment variables
- Added troubleshooting section"

Ready to commit? (yes/no)
```

---

## Example 4: Test Addition

**User:** "Added tests for the calculator"

**Assistant:**
```
I've staged the test files:

- internal/domain/calculator/calculator_test.go
- internal/domain/calculator/volume_test.go

Proposed commit message:
test: Add calculator unit tests

- 1RM formula calculations
- Volume calculations
- Edge case handling
- 85% test coverage achieved

Exact command to run:
git commit -m "test: Add calculator unit tests" -m "- 1RM formula calculations
- Volume calculations
- Edge case handling
- 85% test coverage achieved"

Ready to commit? (yes/no)
```

---

## Example 5: Push Request

**User:** "Push the changes"

**Assistant:**
```
Current branch: main
Commits to push: 3

Commits:
- feat: Add security middleware (11ed207)
- fix: Resolve login validation issue (b664ddf)
- docs: Update installation instructions (9abc182)

Exact command to run:
git push origin main

Ready to push? (yes/no)
```

---

## Example 6: Force Push (Requires Extra Caution)

**User:** "I need to force push after rebase"

**Assistant:**
```
⚠️  WARNING: Force push requested

Current branch: feature/auth
Remote branch: origin/feature/auth
Commits to force push: 5

This will overwrite remote history. Are you sure?

Exact command to run:
git push --force-with-lease origin feature/auth

Confirm force push? (yes/no)
```

---

## Example 7: Error Handling

**User:** "Commit the changes"

**Assistant:**
```
I've staged the files:

- src/app.js
- src/app.test.js

Proposed commit message:
feat: Add new dashboard component

- Real-time metrics display
- Chart integration
- Responsive design

Exact command to run:
git commit -m "feat: Add new dashboard component" -m "- Real-time metrics display
- Chart integration
- Responsive design"

Ready to commit? (yes/no)
```

**User:** "yes"

**Assistant:** [Executes commit, but linting fails]
```
❌ Commit failed: ESLint errors in src/app.js

Fixed the issues:
- Added missing semicolons
- Fixed variable naming
- Updated import order

Updated commit message:
feat: Add new dashboard component

- Real-time metrics display
- Chart integration
- Responsive design
- Fixed ESLint issues

Exact command to run:
git commit -m "feat: Add new dashboard component" -m "- Real-time metrics display
- Chart integration
- Responsive design
- Fixed ESLint issues"

Ready to retry commit? (yes/no)
```

---

## Example 8: Multiple Files

**User:** "I've refactored the API layer"

**Assistant:**
```
I've staged 12 files:

- internal/api/handler/auth_handler.go
- internal/api/handler/exercise_handler.go
- internal/api/handler/session_handler.go
- internal/api/handler/workout_handler.go
- internal/api/middleware/auth.go
- internal/api/middleware/cors.go
- internal/api/middleware/ratelimit.go
- internal/api/middleware/validation.go
- internal/api/router.go
- internal/api/handler/auth_handler_test.go
- internal/api/handler/exercise_handler_test.go
- internal/api/handler/session_handler_test.go

Proposed commit message:
refactor: Standardize API handler patterns

- Consistent error handling
- Unified response formats
- Standardized middleware usage
- Added handler tests

Exact command to run:
git commit -m "refactor: Standardize API handler patterns" -m "- Consistent error handling
- Unified response formats
- Standardized middleware usage
- Added handler tests"

Ready to commit? (yes/no)
```

---

## Best Practices Demonstrated

1. **Auto-staging**: Files are automatically staged
2. **Transparency**: Shows exactly what will be committed
3. **Command preview**: Shows exact command to run
4. **Approval required**: Waits for explicit approval
5. **Error handling**: Fixes issues and re-asks for approval
6. **Conventional commits**: Uses standard commit types
7. **Detailed messages**: Includes context in commit body

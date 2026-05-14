# Git Practices Skill

## Overview
This skill enforces strict git practices to prevent accidental data loss and ensure proper workflow.

## Core Principles

### 1. Authorization Required
**ALWAYS ask for explicit authorization before performing:**
- Creating commits
- Pushing to remote
- Rebasing branches
- Force operations (push --force, reset --hard, etc.)
- Deleting branches
- Merging branches
- Creating/deleting tags

### 2. Auto-Stage Workflow
**For completed work:**
1. Automatically stage files with `git add`
2. Show staged files and proposed commit message
3. Wait for approval
4. Only then execute the commit

### 3. Commit Message Format
**Use Conventional Commits:**
```
<type>: <description>

[optional body]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes
- `perf`: Performance improvements

## Authorization Procedure

### For Commits
When you complete work, I will:

1. **Stage files automatically:**
   ```bash
   git add <files>
   ```

2. **Show what will be committed:**
   ```
   Files to commit:
   - internal/api/middleware/cors.go
   - internal/api/middleware/ratelimit.go
   - internal/api/middleware/validation.go
   
   Proposed commit message:
   feat: Add security middleware
   
   - CORS middleware with whitelist support
   - Rate limiting (60 req/min per IP)
   - Input validation (email, password, XSS protection)
   ```

3. **Show the exact command to run:**
   ```bash
   git commit -m "feat: Add security middleware" -m "- CORS middleware with whitelist support
   - Rate limiting (60 req/min per IP)
   - Input validation (email, password, XSS protection)"
   ```

4. **Wait for your approval** before executing

### For Pushes
When you want to push:

1. **Show current branch and commits:**
   ```
   Branch: main
   Commits to push: 12
   ```

2. **Show the exact command:**
   ```bash
   git push origin main
   ```

3. **Wait for approval**

### For Force Operations
**NEVER perform without explicit authorization:**
- `git push --force`
- `git push --force-with-lease`
- `git reset --hard`
- `git rebase -i`
- `git clean -f`
- `git branch -D`

**Always show:**
- What will be affected
- Why it's needed
- The exact command
- Wait for approval

## Error Handling

### If Commit Fails
1. **Fix the issue** (linting, tests, etc.)
2. **Re-stage files**
3. **Show updated commit message**
4. **Ask for re-approval**
5. **Only then retry commit**

### If Push Fails
1. **Analyze the error**
2. **Determine if rebase is needed**
3. **Ask for authorization before proceeding**
4. **Follow rebase procedure if needed**

## Commit Templates

### Feature
```
feat: <short description>

- What this feature does
- How it works
- Any important notes
```

### Fix
```
fix: <short description>

- What was broken
- How it's fixed
- Testing notes
```

### Test
```
test: <short description>

- What tests were added
- Coverage areas
- Test commands
```

### Documentation
```
docs: <short description>

- What was added/updated
- Location of changes
- Related documentation
```

### Refactor
```
refactor: <short description>

- What was refactored
- Why it was needed
- No functional changes
```

### Chore
```
chore: <short description>

- What maintenance task
- Why it's needed
- Any side effects
```

## Workflow Checklist

### For New Features
- [ ] Code complete
- [ ] Tests added
- [ ] Linting passes
- [ ] Build successful
- [ ] Auto-staged files
- [ ] Commit message approved
- [ ] Commit executed
- [ ] Push approved

### For Bug Fixes
- [ ] Issue identified
- [ ] Root cause found
- [ ] Fix implemented
- [ ] Test added
- [ ] Regression tested
- [ ] Auto-staged files
- [ ] Commit message approved
- [ ] Commit executed
- [ ] Push approved

### For Documentation
- [ ] Content complete
- [ ] Review complete
- [ ] Auto-staged files
- [ ] Commit message approved
- [ ] Commit executed
- [ ] Push approved

## Forbidden Operations

**NEVER perform these without explicit authorization:**
1. Force push
2. Force push with lease
3. Hard reset
4. Interactive rebase
5. Clean with force
6. Delete branches
7. Delete tags
8. Amend commits (unless requested)

## Quick Reference

### Always Ask For:
```
git commit
git push
git rebase
git push --force
git reset --hard
git branch -D
git tag
```

### Can Do Automatically:
```
git add
git status
git diff
git log
git show
git branch (list)
git checkout
```

## Example Interaction

**User:** "I'm done with the auth feature"

**Assistant:** 
```
I've staged the auth feature files:

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

**Created:** May 14, 2026  
**Location:** `.opencode/skills/git-practices/`  
**Version:** 1.0

# Git Practices Skill

A strict git workflow skill that enforces authorization before any git operations.

## Installation

This skill is already installed in your project at:
```
.opencode/skills/git-practices/
```

## What This Skill Does

### 🔒 Enforces Authorization
- **ALWAYS asks** before commits, pushes, rebases, force operations, merges, branch deletions, and tag operations
- Shows exactly what will be done before doing it
- Provides the exact command that will be executed

### 📝 Auto-Stage Workflow
1. When work is complete, automatically stages relevant files
2. Shows all staged files
3. Proposes a conventional commit message
4. Shows the exact git command
5. Waits for your approval
6. Only then executes the commit

### ✅ Conventional Commits
Enforces standard commit message format:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation
- `test:` - Tests
- `refactor:` - Code refactoring
- `chore:` - Maintenance
- `ci:` - CI/CD changes
- `perf:` - Performance

### 🛡️ Safety Features
- Never performs destructive operations without explicit approval
- Fixes errors and re-asks for approval
- Shows full transparency of what will happen

## Files

- `SKILL.md` - Complete skill documentation
- `EXAMPLES.md` - Real-world usage examples
- `QUICK_REFERENCE.md` - Quick reference card
- `package.json` - Skill metadata
- `validate.sh` - Validation script

## Usage

The skill is automatically active. When you complete work:

**You say:** "I've completed the authentication feature"

**Assistant will:**
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

**You approve:** "yes"

**Assistant executes the commit**

## Operations Requiring Authorization

✋ **ALWAYS requires approval:**
- `git commit`
- `git push`
- `git rebase`
- `git push --force`
- `git reset --hard`
- `git branch -D`
- `git merge`
- `git tag` operations
- `git commit --amend`

✅ **Can do automatically:**
- `git add` (auto-staging)
- `git status`
- `git diff`
- `git log`
- `git show`
- `git branch` (list only)

## Error Handling

If a commit fails (linting, tests, pre-commit hooks):
1. Fix the issue
2. Re-stage files
3. Show updated commit message
4. Ask for re-approval
5. Retry commit

## Quick Reference

See `QUICK_REFERENCE.md` for a one-page reference card.

## Examples

See `EXAMPLES.md` for detailed usage examples.

## Version

**Version:** 1.0.0  
**Created:** 2026-05-14  
**Location:** `.opencode/skills/git-practices/`

## Testing

To test the validation script:
```bash
cd .opencode/skills/git-practices
./validate.sh "feat: Test commit message"
```

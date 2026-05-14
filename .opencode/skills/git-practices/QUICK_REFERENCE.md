# Git Practices Skill - Quick Reference Card

## 🚨 ALWAYS ASK FOR AUTHORIZATION

### Commands Requiring Approval
```
git commit
git push
git rebase
git push --force
git push --force-with-lease
git reset --hard
git rebase -i
git clean -f
git branch -D
git tag -d
git merge
git commit --amend
```

### Auto-Stage Workflow
1. **Complete work** → Auto-stage files
2. **Show staged files** + commit message
3. **Show exact command** to run
4. **Wait for approval** → Execute

### Conventional Commits
```
feat: Add new feature
fix: Fix bug
docs: Update documentation
style: Code style changes
refactor: Code refactoring
test: Add tests
chore: Maintenance tasks
ci: CI/CD changes
perf: Performance improvements
```

## 📋 Example Approval Flow

**When work is complete:**
```
Files staged:
- src/feature.js
- test/feature.test.js

Commit message:
feat: Add user profile feature

- Profile page with avatar upload
- Settings management
- 3 new tests

Command to run:
git commit -m "feat: Add user profile feature" -m "- Profile page with avatar upload
- Settings management
- 3 new tests"

Approve? (yes/no)
```

## ⚠️ NEVER DO WITHOUT ASKING

- Force operations
- Hard resets
- Branch deletions
- Interactive rebases
- Amending commits
- Tag operations

## 🔧 Error Handling

If commit fails:
1. Fix the issue
2. Re-stage files
3. Show updated message
4. Ask for re-approval
5. Retry commit

## 📁 Location
`.opencode/skills/git-practices/`

**Version:** 1.0  
**Created:** 2026-05-14

#!/bin/bash
# Git Practices Skill - Validation Script

echo "=== Git Practices Validation ==="

# Check if in git repository
if [ ! -d .git ]; then
    echo "❌ Not a git repository"
    exit 1
fi

echo "✅ In git repository"

# Check staged files
STAGED_FILES=$(git diff --cached --name-only)
if [ -z "$STAGED_FILES" ]; then
    echo "❌ No files staged for commit"
    exit 1
fi

echo "✅ Files staged:"
echo "$STAGED_FILES" | sed 's/^/  - /'

# Check commit message format
if [ -z "$1" ]; then
    echo "❌ No commit message provided"
    exit 1
fi

COMMIT_MSG="$1"
echo "✅ Commit message: $COMMIT_MSG"

# Validate conventional commit format
if [[ ! "$COMMIT_MSG" =~ ^(feat|fix|docs|style|refactor|test|chore|ci|perf): ]]; then
    echo "⚠️  Warning: Commit message doesn't follow conventional format"
    echo "   Expected: <type>: <description>"
    echo "   Types: feat, fix, docs, style, refactor, test, chore, ci, perf"
fi

# Check message length
if [ ${#COMMIT_MSG} -gt 72 ]; then
    echo "⚠️  Warning: Commit message exceeds 72 characters"
fi

echo ""
echo "=== Ready for Approval ==="
echo "Files staged: $(echo "$STAGED_FILES" | wc -l | tr -d ' ')"
echo "Commit message: $COMMIT_MSG"
echo ""
echo "To commit:"
echo "git commit -m \"$COMMIT_MSG\""
echo ""
echo "Approve? (yes/no)"

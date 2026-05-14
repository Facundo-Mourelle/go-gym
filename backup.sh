#!/bin/bash
set -e

BACKUP_DIR=${BACKUP_DIR:-./backups}
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/gymtracker_$TIMESTAMP.sql"

echo "=== Database Backup ==="

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Get database container name
DB_CONTAINER=$(docker-compose -f docker-compose.prod.yml ps -q db)

if [ -z "$DB_CONTAINER" ]; then
    echo "ERROR: Database container not found"
    exit 1
fi

# Backup database
echo "Backing up database to $BACKUP_FILE..."
docker exec "$DB_CONTAINER" pg_dump -U gym gymtracker > "$BACKUP_FILE"

# Compress backup
gzip -f "$BACKUP_FILE"
BACKUP_FILE="$BACKUP_FILE.gz"

echo "Backup complete: $BACKUP_FILE"
echo "Size: $(du -h "$BACKUP_FILE" | cut -f1)"

# Keep only last 7 backups
echo "Cleaning old backups (keeping last 7)..."
ls -t "$BACKUP_DIR"/gymtracker_*.sql.gz 2>/dev/null | tail -n +8 | xargs -r rm -f

echo "Backup rotation complete"

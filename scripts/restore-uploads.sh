#!/bin/bash

# Restore uploaded files to Docker volume
echo "🔄 Restoring uploaded files..."

# Check if backup file is provided
if [ -z "$1" ]; then
  echo "❌ Please provide backup file path"
  echo "Usage: ./restore-uploads.sh <backup_file.tar.gz>"
  echo ""
  echo "Available backups:"
  ls -la ./backups/uploads/ 2>/dev/null || echo "No backups found"
  exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
  echo "❌ Backup file not found: $BACKUP_FILE"
  exit 1
fi

echo "📂 Restoring from: $BACKUP_FILE"

# Restore backup to Docker volume
docker run --rm \
  -v ecom-golang-clean-architecture_uploads_data:/target \
  -v $(pwd):/backup \
  alpine \
  sh -c "cd /target && tar xzf /backup/$BACKUP_FILE"

if [ $? -eq 0 ]; then
  echo "✅ Restore completed successfully!"
  echo "🔍 Verifying restored files..."
  
  # List restored files
  docker run --rm \
    -v ecom-golang-clean-architecture_uploads_data:/uploads \
    alpine \
    find /uploads -type f | head -10
else
  echo "❌ Restore failed!"
  exit 1
fi

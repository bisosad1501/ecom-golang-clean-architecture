#!/bin/bash

# Backup uploaded files from Docker volume
echo "🔄 Backing up uploaded files..."

# Create backup directory
mkdir -p ./backups/uploads

# Get the current date for backup naming
BACKUP_DATE=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="./backups/uploads/uploads_backup_${BACKUP_DATE}.tar.gz"

# Create backup from Docker volume
docker run --rm \
  -v ecom-golang-clean-architecture_uploads_data:/source \
  -v $(pwd)/backups/uploads:/backup \
  alpine \
  tar czf /backup/uploads_backup_${BACKUP_DATE}.tar.gz -C /source .

if [ $? -eq 0 ]; then
  echo "✅ Backup created successfully: ${BACKUP_FILE}"
  echo "📁 Backup location: $(pwd)/backups/uploads/"
  ls -lh ./backups/uploads/uploads_backup_${BACKUP_DATE}.tar.gz
else
  echo "❌ Backup failed!"
  exit 1
fi

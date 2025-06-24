# File Upload Management Guide

## 🗂️ Upload Storage Solutions

### Development Environment (Current)
- **Method**: Bind mount to local directory
- **Location**: `./uploads` (host) → `/app/uploads` (container)
- **Persistence**: ✅ Files persist across container rebuilds
- **Backup**: Not needed (files are on host filesystem)

### Production Environment
- **Method**: Named Docker volumes with host bind mounts
- **Location**: `/var/ecom/uploads` (host) → `/app/uploads` (container)
- **Persistence**: ✅ Files persist across container rebuilds
- **Backup**: Required using backup scripts

## 🔄 Backup & Restore

### For Development
```bash
# Files are already on host filesystem at ./uploads
# Just commit to Git (if needed) or copy to backup location
cp -r ./uploads ./backups/uploads_$(date +%Y%m%d_%H%M%S)
```

### For Production (Docker Volumes)
```bash
# Backup uploads
./scripts/backup-uploads.sh

# Restore uploads
./scripts/restore-uploads.sh backups/uploads/uploads_backup_20250624_143022.tar.gz
```

## 🚀 Deployment Commands

### Development
```bash
# Start with bind mount (current setup)
docker-compose up -d

# Files are stored in ./uploads directory
ls -la ./uploads/images/
```

### Production
```bash
# Create production directories
sudo mkdir -p /var/ecom/{postgres_data,redis_data,uploads}
sudo chown -R 1000:1000 /var/ecom/

# Start production stack
docker-compose -f docker-compose.prod.yml up -d
```

## 📁 Directory Structure
```
project/
├── uploads/                  # Development uploads
│   ├── .gitkeep             # Keep directory in Git
│   └── images/              # Uploaded product images
├── backups/                 # Backup storage
│   └── uploads/             # Upload backups
├── scripts/
│   ├── backup-uploads.sh    # Backup script
│   └── restore-uploads.sh   # Restore script
├── docker-compose.yml       # Development config
└── docker-compose.prod.yml  # Production config
```

## 🔧 Configuration

### Environment Variables
```env
# Upload configuration
UPLOAD_PATH=/app/uploads
MAX_UPLOAD_SIZE=10485760  # 10MB
```

### File Access URLs
- **Development**: `http://localhost:8080/uploads/images/{filename}`
- **Production**: `https://yourdomain.com/uploads/images/{filename}`

## ⚠️ Important Notes

1. **Development**: Files in `./uploads` are ignored by Git (except `.gitkeep`)
2. **Production**: Always backup before major deployments
3. **Security**: Ensure proper file permissions and validation
4. **Performance**: Consider CDN for production file serving

## 🛡️ Security Best Practices

1. **File validation**: Only allow specific image formats
2. **Size limits**: Enforce maximum file size (currently 5MB)
3. **Path traversal**: Validate file paths and names
4. **Virus scanning**: Consider antivirus scanning for uploads
5. **Content-Type**: Validate MIME types

## 📊 Monitoring

### Check Upload Storage
```bash
# Development
du -sh ./uploads

# Production (Docker volume)
docker run --rm -v ecom-golang-clean-architecture_uploads_data:/uploads alpine du -sh /uploads
```

### List Recent Uploads
```bash
# Development
find ./uploads -type f -mtime -1 -ls

# Production
docker run --rm -v ecom-golang-clean-architecture_uploads_data:/uploads alpine find /uploads -type f -mtime -1 -ls
```

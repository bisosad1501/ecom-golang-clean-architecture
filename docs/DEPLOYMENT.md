# Deployment Guide

This guide covers different deployment options for the E-commerce API.

## 🐳 Docker Deployment

### Prerequisites
- Docker and Docker Compose installed
- Git (to clone the repository)

### Quick Start

1. **Clone the repository**
```bash
git clone https://github.com/bisosad1501/ecom-golang-clean-architecture.git
cd ecom-golang-clean-architecture
```

2. **Start with Docker Compose**
```bash
docker-compose up -d
```

This will start:
- PostgreSQL database (port 5432)
- Redis cache (port 6379)
- E-commerce API (port 8080)
- pgAdmin (port 5050) - optional

3. **Verify deployment**
```bash
curl http://localhost:8080/health
```

### Environment Configuration

Create a `.env` file or modify `docker-compose.yml`:

```env
# Application
APP_NAME=ecom-api
APP_ENV=production
APP_PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=ecommerce_db

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT
JWT_SECRET=your_super_secret_jwt_key_change_this_in_production
JWT_EXPIRE_HOURS=24

# Email (optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
FROM_EMAIL=noreply@yourcompany.com

# Payment (optional)
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
```

## ☁️ Cloud Deployment

### AWS ECS Deployment

1. **Build and push Docker image**
```bash
# Build image
docker build -t ecom-api .

# Tag for ECR
docker tag ecom-api:latest 123456789012.dkr.ecr.us-west-2.amazonaws.com/ecom-api:latest

# Push to ECR
docker push 123456789012.dkr.ecr.us-west-2.amazonaws.com/ecom-api:latest
```

2. **Create ECS Task Definition**
```json
{
  "family": "ecom-api",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "ecom-api",
      "image": "123456789012.dkr.ecr.us-west-2.amazonaws.com/ecom-api:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {"name": "APP_ENV", "value": "production"},
        {"name": "DB_HOST", "value": "your-rds-endpoint"},
        {"name": "REDIS_HOST", "value": "your-elasticache-endpoint"}
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/ecom-api",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Google Cloud Run Deployment

1. **Build and deploy**
```bash
# Build and deploy to Cloud Run
gcloud run deploy ecom-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="APP_ENV=production,DB_HOST=your-cloud-sql-ip"
```

### Heroku Deployment

1. **Create Heroku app**
```bash
heroku create your-ecom-api
```

2. **Add PostgreSQL and Redis**
```bash
heroku addons:create heroku-postgresql:hobby-dev
heroku addons:create heroku-redis:hobby-dev
```

3. **Set environment variables**
```bash
heroku config:set APP_ENV=production
heroku config:set JWT_SECRET=your_super_secret_key
```

4. **Deploy**
```bash
git push heroku main
```

## 🔧 Production Configuration

### Security Checklist

- [ ] Change default JWT secret
- [ ] Use strong database passwords
- [ ] Enable SSL/TLS
- [ ] Configure CORS properly
- [ ] Set up rate limiting
- [ ] Enable request logging
- [ ] Configure monitoring

### Database Setup

1. **PostgreSQL Configuration**
```sql
-- Create database
CREATE DATABASE ecommerce_db;

-- Create user
CREATE USER ecom_user WITH PASSWORD 'secure_password';

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE ecommerce_db TO ecom_user;
```

2. **Redis Configuration**
```redis
# redis.conf
requirepass your_redis_password
maxmemory 256mb
maxmemory-policy allkeys-lru
```

### Monitoring

1. **Health Checks**
```bash
# Application health
curl http://your-domain/health

# Database connectivity
curl http://your-domain/api/v1/categories/root
```

2. **Logging**
```bash
# View application logs
docker logs ecom_api

# View database logs
docker logs ecom_postgres
```

### Backup Strategy

1. **Database Backup**
```bash
# Create backup
docker exec ecom_postgres pg_dump -U postgres ecommerce_db > backup.sql

# Restore backup
docker exec -i ecom_postgres psql -U postgres ecommerce_db < backup.sql
```

2. **Automated Backups**
```bash
# Add to crontab
0 2 * * * /path/to/backup-script.sh
```

## 🚀 Scaling

### Horizontal Scaling

1. **Load Balancer Configuration**
```nginx
upstream ecom_api {
    server api1:8080;
    server api2:8080;
    server api3:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://ecom_api;
    }
}
```

2. **Database Read Replicas**
```go
// Configure read/write splitting
masterDB := connectToMaster()
replicaDB := connectToReplica()
```

### Vertical Scaling

1. **Resource Limits**
```yaml
# docker-compose.yml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

## 🔍 Troubleshooting

### Common Issues

1. **Database Connection Failed**
```bash
# Check database status
docker logs ecom_postgres

# Test connection
docker exec -it ecom_postgres psql -U postgres -d ecommerce_db
```

2. **API Not Responding**
```bash
# Check API logs
docker logs ecom_api

# Check port binding
netstat -tulpn | grep 8080
```

3. **Memory Issues**
```bash
# Monitor memory usage
docker stats

# Check application metrics
curl http://localhost:8080/metrics
```

### Performance Optimization

1. **Database Indexing**
```sql
-- Add indexes for frequently queried fields
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_created_at ON orders(created_at);
```

2. **Caching Strategy**
```go
// Implement Redis caching for frequently accessed data
cache.Set("categories", categories, 1*time.Hour)
```

3. **Connection Pooling**
```go
// Configure database connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

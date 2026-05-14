# Go Gym - Deployment Guide

## Prerequisites

- Docker and Docker Compose installed
- PostgreSQL 15+ (or use Docker Compose)
- Go 1.25+ (for local development)
- Node.js 20+ (for frontend development)

## Quick Start (Development)

```bash
# 1. Clone repository
git clone <repository-url>
cd go-gym

# 2. Copy environment template
cp .env.example .env

# 3. Edit .env with your values
# IMPORTANT: Change JWT_SECRET to a secure random value
# Generate with: openssl rand -base64 32

# 4. Start all services
./start.sh

# Services will be available at:
# - Backend: http://localhost:8080
# - Frontend: http://localhost:5173
# - Database: localhost:5432
```

## Production Deployment

### Option 1: Docker Compose (Recommended)

```bash
# 1. Set required environment variables
export JWT_SECRET="$(openssl rand -base64 32)"
export DATABASE_URL="postgres://gym:gympassword@db:5432/gymtracker?sslmode=disable"
export ENVIRONMENT="production"
export CORS_ALLOWED_ORIGINS="https://yourdomain.com"

# 2. Run deployment script
./deploy.sh

# 3. Verify deployment
curl http://localhost:8080/health
```

### Option 2: Manual Docker Compose

```bash
# 1. Create .env.production file
cat > .env.production << EOF
DB_USER=gym
DB_PASSWORD=<secure-password>
DB_NAME=gymtracker
DB_PORT=5432

JWT_SECRET=<generated-secret>
ENVIRONMENT=production
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=60
CORS_ALLOWED_ORIGINS=https://yourdomain.com

BACKEND_PORT=8080
FRONTEND_PORT=3000
EOF

# 2. Start services
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# 3. Check logs
docker-compose -f docker-compose.prod.yml logs -f
```

### Option 3: Kubernetes (Advanced)

See `k8s/` directory for Kubernetes manifests (coming soon).

## Environment Variables

### Required

| Variable | Description | Example |
|----------|-------------|---------|
| `JWT_SECRET` | Secret key for JWT signing | `openssl rand -base64 32` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@host:5432/db` |

### Optional

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Backend server port |
| `ENVIRONMENT` | `development` | Environment mode |
| `RATE_LIMIT_ENABLED` | `true` | Enable rate limiting |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | `60` | Rate limit threshold |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | Allowed CORS origins (comma-separated) |

## Database Setup

### Using Docker Compose (Included)

The production docker-compose file includes PostgreSQL. No additional setup needed.

### Using External Database

```bash
# 1. Create database
createdb gymtracker

# 2. Run migrations (if available)
psql -U gym -d gymtracker -f scripts/migrations/001_initial.sql

# 3. (Optional) Seed data
psql -U gym -d gymtracker -f scripts/seed.sql
```

## Health Checks

```bash
# Backend health
curl http://localhost:8080/health

# Expected response:
# {"status":"ok","service":"gym-tracker-api"}

# Run full health check
./healthcheck.sh
```

## Backup & Restore

### Backup

```bash
# Automated backup (creates timestamped backup in ./backups/)
./backup.sh

# Manual backup
docker exec gym-db pg_dump -U gym gymtracker > backup.sql
```

### Restore

```bash
# From backup file
docker exec -i gym-db psql -U gym gymtracker < backup.sql

# From compressed backup
gunzip -c backup.sql.gz | docker exec -i gym-db psql -U gym gymtracker
```

## Monitoring

### Logs

```bash
# All services
docker-compose -f docker-compose.prod.yml logs -f

# Specific service
docker-compose -f docker-compose.prod.yml logs -f backend

# Last 100 lines
docker-compose -f docker-compose.prod.yml logs --tail=100
```

### Metrics

```bash
# Container stats
docker stats

# Resource usage
docker-compose -f docker-compose.prod.yml ps
```

## Scaling

### Horizontal Scaling (Multiple Backend Instances)

```bash
# Scale backend to 3 instances
docker-compose -f docker-compose.prod.yml up -d --scale backend=3

# Add load balancer (nginx/traefik) in front
```

### Database Connection Pooling

Configure in `DATABASE_URL`:
```
postgres://user:pass@host:5432/db?pool_max_conns=20&pool_min_conns=5
```

## SSL/TLS Configuration

### Using Nginx Reverse Proxy

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/ssl/certs/cert.pem;
    ssl_certificate_key /etc/ssl/private/key.pem;

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
    }
}
```

### Using Let's Encrypt

```bash
# Install certbot
apt-get install certbot python3-certbot-nginx

# Obtain certificate
certbot --nginx -d yourdomain.com

# Auto-renewal is configured automatically
```

## Troubleshooting

### Backend won't start

```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs backend

# Common issues:
# 1. JWT_SECRET not set or using default value
# 2. Database connection failed
# 3. Port already in use
```

### Database connection errors

```bash
# Check database is running
docker-compose -f docker-compose.prod.yml ps db

# Check database logs
docker-compose -f docker-compose.prod.yml logs db

# Test connection
docker exec gym-db psql -U gym -d gymtracker -c "SELECT 1"
```

### Frontend not loading

```bash
# Check frontend logs
docker-compose -f docker-compose.prod.yml logs frontend

# Verify API URL is correct
# Check VITE_API_URL in frontend container
```

### Rate limiting too aggressive

```bash
# Increase rate limit
export RATE_LIMIT_REQUESTS_PER_MINUTE=120

# Or disable temporarily
export RATE_LIMIT_ENABLED=false

# Restart backend
docker-compose -f docker-compose.prod.yml restart backend
```

## Security Checklist

- [ ] JWT_SECRET is a secure random value (not default)
- [ ] Database password is strong and unique
- [ ] CORS_ALLOWED_ORIGINS is set to your domain (not *)
- [ ] SSL/TLS is configured for production
- [ ] Database is not exposed to public internet
- [ ] Regular backups are scheduled
- [ ] Security updates are applied regularly
- [ ] Rate limiting is enabled
- [ ] Logs are monitored for suspicious activity

## Performance Tuning

### Backend

```bash
# Increase Go max procs
export GOMAXPROCS=4

# Database connection pool
export DATABASE_URL="postgres://...?pool_max_conns=50"
```

### Frontend

- Enable CDN for static assets
- Configure browser caching
- Enable gzip compression (already configured in nginx)

### Database

```sql
-- Add indexes for common queries
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_exercises_user_id ON exercises(user_id);
CREATE INDEX idx_performed_sets_session_id ON performed_sets(session_id);
```

## Updating

```bash
# Pull latest changes
git pull origin main

# Rebuild images
docker-compose -f docker-compose.prod.yml build

# Restart services (zero-downtime with multiple instances)
docker-compose -f docker-compose.prod.yml up -d

# Run migrations if needed
docker exec gym-backend ./migrate
```

## Rollback

```bash
# Stop current version
docker-compose -f docker-compose.prod.yml down

# Restore database backup
gunzip -c backups/gymtracker_YYYYMMDD_HHMMSS.sql.gz | \
  docker exec -i gym-db psql -U gym gymtracker

# Start previous version
git checkout <previous-commit>
docker-compose -f docker-compose.prod.yml up -d
```

## Support

- GitHub Issues: <repository-url>/issues
- Documentation: See README.md and AGENTS.md
- Security Issues: Report privately to maintainers

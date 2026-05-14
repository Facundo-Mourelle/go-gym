#!/bin/bash
set -e

echo "=== Go Gym Deployment Script ==="

# Check for required environment variables
if [ -z "$JWT_SECRET" ]; then
    echo "ERROR: JWT_SECRET environment variable is required"
    echo "Generate one with: openssl rand -base64 32"
    exit 1
fi

if [ -z "$DATABASE_URL" ]; then
    echo "ERROR: DATABASE_URL environment variable is required"
    echo "Example: postgres://user:password@host:5432/database"
    exit 1
fi

# Set defaults
ENVIRONMENT=${ENVIRONMENT:-production}
CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS:-"http://localhost:3000"}
RATE_LIMIT_REQUESTS_PER_MINUTE=${RATE_LIMIT_REQUESTS_PER_MINUTE:-60}

echo "Environment: $ENVIRONMENT"
echo "CORS Origins: $CORS_ALLOWED_ORIGINS"
echo "Rate Limit: $RATE_LIMIT_REQUESTS_PER_MINUTE per minute"

# Create .env file for docker-compose
cat > .env.production << EOF
# Database
DB_USER=$(echo $DATABASE_URL | sed -n 's|.*://\([^:]*\):.*|\1|p')
DB_PASSWORD=$(echo $DATABASE_URL | sed -n 's|.*://[^:]*:\([^@]*\)@.*|\1|p')
DB_NAME=$(echo $DATABASE_URL | sed -n 's|.*://[^/]*/\([^?]*\).*|\1|p')
DB_PORT=5432

# Backend
JWT_SECRET=$JWT_SECRET
ENVIRONMENT=$ENVIRONMENT
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=$RATE_LIMIT_REQUESTS_PER_MINUTE
CORS_ALLOWED_ORIGINS=$CORS_ALLOWED_ORIGINS

# Ports
BACKEND_PORT=${BACKEND_PORT:-8080}
FRONTEND_PORT=${FRONTEND_PORT:-3000}

#!/bin/bash
set -e

echo "=== Testing Go Gym Locally (No Commit Required) ==="

# Check if .env exists
if [ ! -f .env ]; then
    echo "Creating .env from template..."
    cp .env.example .env
    # Set a test JWT secret
    sed -i 's/JWT_SECRET="your-secret-key-change-in-production"/JWT_SECRET="test-secret-for-local-development-12345"/' .env
fi

echo ""
echo "1. Testing Backend Build..."
go build -o /tmp/gym-api cmd/api/main.go
if [ $? -eq 0 ]; then
    echo "   ✅ Backend builds successfully"
    rm /tmp/gym-api
else
    echo "   ❌ Backend build failed"
    exit 1
fi

echo ""
echo "2. Running Backend Tests..."
go test ./... -short
if [ $? -eq 0 ]; then
    echo "   ✅ All tests passing"
else
    echo "   ❌ Tests failed"
    exit 1
fi

echo ""
echo "3. Testing Frontend Build..."
cd gym-tracker-frontend
npm run build > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✅ Frontend builds successfully"
else
    echo "   ❌ Frontend build failed"
    exit 1
fi
cd ..

echo ""
echo "4. Starting Services..."
echo "   - Database: Starting PostgreSQL..."
docker-compose up -d

echo "   - Waiting for database..."
sleep 5

echo "   - Backend: Starting on port 8080..."
./start_backend.sh

echo "   - Frontend: Starting on port 5173..."
./start_frontend.sh

echo ""
echo "=== Services Started ==="
echo "Backend:  http://localhost:8080"
echo "Frontend: http://localhost:5173"
echo "Health:   http://localhost:8080/health"
echo ""
echo "Test the application in your browser!"
echo ""
echo "To stop: ./stop.sh"

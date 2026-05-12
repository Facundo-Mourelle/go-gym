#!/bin/bash

echo "=== Shutting Down Go Gym ==="

echo "[1/3] Stopping backend..."
BACKEND_PID=$(lsof -ti :8080 2>/dev/null || true)
if [ -n "$BACKEND_PID" ]; then
  kill "$BACKEND_PID" 2>/dev/null || true
  echo "Backend stopped"
else
  echo "Backend not running"
fi

echo "[2/3] Stopping frontend..."
FRONTEND_PID=$(lsof -ti :5173 2>/dev/null || true)
[ -n "$FRONTEND_PID" ] && kill "$FRONTEND_PID" 2>/dev/null || true
pkill -f "gym-tracker-frontend" 2>/dev/null || true
echo "Frontend stopped"

echo "[3/3] Stopping PostgreSQL..."
docker-compose down

echo ""
echo "=== All services stopped ==="

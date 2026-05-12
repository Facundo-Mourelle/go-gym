#!/bin/bash
set -e

RUN_TESTS=false

while [[ $# -gt 0 ]]; do
  case $1 in
    -t|--tests)
      RUN_TESTS=true
      shift
      ;;
    -h|--help)
      echo "Usage: $0 [-t|--tests]"
      echo "  -t, --tests    Run tests after initialization"
      exit 0
      ;;
    *)
      shift
      ;;
  esac
done

echo "=== Initializing Go Gym ==="

# 1. Start DB (docker-compose)
echo "[1/3] Starting PostgreSQL..."
docker-compose up -d

echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
  if docker-compose exec -T db pg_isready -U gym > /dev/null 2>&1; then
    echo "PostgreSQL is ready!"
    break
  fi
  echo "  Waiting... ($i/30)"
  sleep 1
done

if ! docker-compose exec -T db pg_isready -U gym > /dev/null 2>&1; then
  echo "WARNING: PostgreSQL not confirmed ready after 30s. Backend may fail."
fi

# 2. Start backend
echo "[2/3] Starting Go API..."
./start_backend.sh
sleep 2

# 3. Start frontend
echo "[3/3] Starting React frontend..."
./start_frontend.sh

echo ""
echo "=== Services running ==="
echo "DB:      localhost:5432"
echo "Backend: localhost:8080"
echo "Frontend: http://localhost:5173"
echo ""

if [ "$RUN_TESTS" = true ]; then
  echo "=== Running tests ==="

  echo "[Backend] Running Go tests..."
  go test ./... -v

  echo ""
  echo "[Frontend] Running React tests..."
  cd gym-tracker-frontend
  npm test -- --run

  echo ""
  echo "=== All tests passed ==="
fi
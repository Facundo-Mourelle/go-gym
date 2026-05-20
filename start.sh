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

echo "=== Starting Go Gym (Docker Dev) ==="

# Build and start all services
echo "[1/1] Starting all services..."
docker compose up -d --build

echo "Waiting for services to be ready..."
for i in {1..30}; do
  if docker compose exec -T db pg_isready -U gym > /dev/null 2>&1; then
    echo "All services ready!"
    break
  fi
  echo "  Waiting... ($i/30)"
  sleep 1
done

echo ""
echo "=== Services running ==="
echo "Caddy:   http://localhost"
echo "Backend: http://localhost/api (proxied via Caddy)"
echo "Frontend: http://localhost (proxied via Caddy)"
echo "DB:      localhost:5432"
echo ""
echo "View logs: docker compose logs -f"
echo "Stop:      docker compose down"
echo ""

if [ "$RUN_TESTS" = true ]; then
  echo "=== Running tests ==="

  echo "[Backend] Running Go tests..."
  docker compose exec backend go test ./... -v

  echo ""
  echo "[Frontend] Running React tests..."
  docker compose exec frontend pnpm test -- --run

  echo ""
  echo "=== All tests passed ==="
fi
#!/bin/bash

echo "=== Go Gym Health Check ==="

# Check backend
echo "Backend health:"
curl -s -f http://localhost:8080/health && echo " ✅" || echo " ❌"

# Check database connection via backend
echo "Database connection:"
curl -s http://localhost:8080/health | grep -q "ok" && echo " ✅" || echo " ❌"

# Check frontend
echo "Frontend health:"
curl -s -f http://localhost:3000/health && echo " ✅" || echo " ❌"

# Check Docker containers
echo ""
echo "=== Docker Containers ==="
docker-compose -f docker-compose.prod.yml ps

# Check container logs for errors
echo ""
echo "=== Recent Errors ==="
docker-compose -f docker-compose.prod.yml logs --tail=20 | grep -i error | head -5

# Check resource usage
echo ""
echo "=== Resource Usage ==="
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | head -5

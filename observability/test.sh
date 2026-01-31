#!/bin/bash

echo "Stress testing with 1000 requests..."
echo "==================================="

echo "Starting stress test..."
echo ""

for i in {1..100}; do
  curl -s "http://localhost:8080/api/v1/data?request=$i" > /dev/null
done

echo ""
echo "Test completed!"
echo ""
echo "Final metrics:"
echo "--------------"
curl -s http://localhost:8080/metrics | grep -E "(http_requests_total|http_requests_by_status)" | grep -v "#"

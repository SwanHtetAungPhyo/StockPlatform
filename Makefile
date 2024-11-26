
docker:
	@echo "Starting Docker Compose..."
	@docker-compose up --build &

load_test:
	@echo "Running Load Test..."
	@k6 run ./test/loadTest.js &


main: docker load_test
	@echo "Waiting for both Docker and Load Test to finish..."
	@wait


FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp

FRONT_BINARY=frontEndApp

BROKER_VERSION=1.0.0



## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"


## build_dockerfiles: builds all dockerfile images
build_dockerfiles: build_auth build_broker build_listener build_logger build_mailer build_front_linux
	@echo "Building dockerfiles..."
	docker build -f caddy.dockerfile -t effectone/micro-caddy:1.0.0 .
	docker build -f ../front-end/front-end.dockerfile -t effectone/front-end:1.0.0 .
	docker build -f ../authentication-service/authentication-service.dockerfile -t effectone/authentication-service:${AUTH_VERSION} .
	docker build -f ../broker-service/broker-service.dockerfile -t effectone/broker-service:1.0.0 .
	docker build -f ../listener-service/listener-service.dockerfile -t effectone/listener-service:1.0.2 .
	docker build -f ../mail-service/mail-service.dockerfile -t effectone/mail-service:1.0.0 .
	docker build -f ../logger-service/logger-service.dockerfile -t effectone/logger-service:1.0.0 .

## push_dockerfiles: pushes tagged versions to docker hub
push_dockerfiles:
	docker push effectone/micro-caddy:1.0.0
	docker push effectone/front-end:1.0.0


## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker 
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"


## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the front end binary
build_front:
	@echo "Building front end binary..."
	cd front-end && env GOOS=linux && CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

swarm_up:
	@echo "Starting swarm..."
	docker stack deploy -c swarm.yml myapp

## swarm_down: stops the swarm
swarm_down:
	@echo "Stopping swarm..."
	docker stack rm myapp

## test: runs all tests
test:
	@echo "Testing..."
	go test -v ./...

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@cd broker-service && rm -f ${BROKER_BINARY}
	@cd broker-service && go clean
	@cd front-end && go clean
	@cd front-end && rm -f ${FRONT_END_BINARY}
	@echo "Cleaned!"

## help: displays help
help: Makefile
	@echo " Choose a command:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

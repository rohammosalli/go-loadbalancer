
## build: builds all binaries
build: clean build_lb
	@printf "All binaries built!\n"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"
	
## build_back: builds the LB
build_lb:
	@echo "Building LB..."
	@go build -o dist/main main.go
	@echo "LB built!"
	
## start_lb: starts the LB
start_lb: build_lb
	@echo "Starting the LB..."
	@./dist/main &
	@echo "LBrunning!"

## stop: stops the front and LB
stop:  stop_lb
	@echo "All applications stopped"

stop_lb:
	@echo "Stopping the LB..."
	@-pkill -SIGTERM -f "./dist/main"
	@echo "Stopped LB"

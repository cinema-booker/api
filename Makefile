# build the api binary
build:
	@go build -o bin/api cmd/main.go

# remove the api binary
clean:
	@rm -rf bin/api

# build and run the api binary
run: clean build
	@./bin/api

# install all dependencies
install:
	go get -u ./...

# run the tests
test:
	@go test -v ./...

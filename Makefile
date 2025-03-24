build:
	go build -o bin/app

run: build
	./bin/app

clean:
	rm -f bin/app

test:
	go test -v ./... -count=1


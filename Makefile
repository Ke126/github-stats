.PHONY: build
build:
	go build .

.PHONY: clean
clean:
	rm -f ./github-stats

.PHONY: run
run:
	go run .

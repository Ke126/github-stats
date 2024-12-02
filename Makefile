.PHONY: build
build:
	go build .

.PHONY: clean
clean:
	rm -f ./github-stats
	rm -rf ./_site

.PHONY: run
run:
	go run .

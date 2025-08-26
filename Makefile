.PHONY: test clean build

build:
	go mod tidy
	go build -o generator generator.go

test: build
	./run-test.sh 600

clean:
	rm -f generator
	pkill -f "otelcol\|generator" || true

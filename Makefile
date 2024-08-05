all: test fmt vet vuln staticcheck build

test:
	go test ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l
	test -z $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l)

vet:
	go vet ./...

vuln:
	govulncheck ./...

staticcheck:
	staticcheck -f stylish ./...

build:
	go build -o bin/survey-worker ./cmd/survey-worker
	go build -o bin/survey-curator-worker ./cmd/survey-curator-worker

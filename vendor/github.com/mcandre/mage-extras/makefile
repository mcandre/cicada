.POSIX:
.SILENT:
.PHONY: all

all:
	cargo install --force unmake@0.0.17

	go install github.com/alexkohler/nakedret@v1.0.1
	go install github.com/kisielk/errcheck@v1.6.3
	go install github.com/magefile/mage@v1.14.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@2023.1.3
	go mod tidy

	npm install -g snyk

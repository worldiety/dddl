# Build binaries to the out directory.
# The filenames should be chosen so that Node's OS API (https://nodejs.org/api/os.html#os_os_arch)
# can be used to select a binary.
build: test
	GOOS=linux GOARCH=amd64 go build -o vsc-plugin/out/bin/dddlsp-linux-x64 github.com/worldiety/dddl/cmd/ddd-lsp
	GOOS=darwin GOARCH=amd64 go build -o vsc-plugin/out/bin/dddlsp-darwin-x64 github.com/worldiety/dddl/cmd/ddd-lsp
	# Skip building of darwin-arm64 until github actions can do that.
	# GOOS=darwin GOARCH=arm64 go build -o ../out/bin/wdyspec-darwin-arm64 cmd/dyml.go

test:
	golangci-lint run || true

# Download the newest LSP types from https://github.com/golang/tools.
# License and code will be copied to the protocol directory.
protocol:
	rm -rf /tmp/dddl/go-tools
	git clone https://github.com/golang/tools.git /tmp/dddl/go-tools
	mkdir protocol || true
	cp /tmp/dddl/go-tools/LICENSE lsp/protocol
	cp /tmp/dddl/go-tools/internal/lsp/protocol/tsprotocol.go lsp/protocol

.PHONY: protocol ebnf

ebnf:
	go run github.com/worldiety/dddl/cmd/dddc -format=grammar > ebnf/grammar.ebnf
	go run github.com/alecthomas/participle/v2/cmd/railroad@latest < ebnf/grammar.ebnf > ebnf/grammar.html
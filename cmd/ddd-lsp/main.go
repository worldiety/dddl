// Package main provides the ddd lsp server.
package main

import (
	"bufio"
	"context"
	"github.com/worldiety/dddl/lsp"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	reader := bufio.NewReader(os.Stdin)
	server := lsp.NewServer()

	lsp.HandleRequests(ctx, server, reader)

}

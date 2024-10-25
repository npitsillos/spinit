package main

import (
	"context"

	"github.com/npitsillos/spinit/cmd/root"
)

func main() {
	rootCmd := root.NewRootCommand()
	rootCmd.ExecuteContext(context.Background())
}

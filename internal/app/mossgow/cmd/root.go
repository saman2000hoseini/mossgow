package cmd

import (
	"github.com/spf13/cobra"

	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/cmd/detect"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/config"
)

// NewRootCommand creates a new mossgow root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "mossgow",
	}

	cfg := config.Init()

	detect.Register(root, cfg)

	return root
}

package cmd

import (
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/cmd/detect"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/config"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new virtual-box root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "mossgow",
	}

	cfg := config.Init()

	detect.Register(root, cfg)

	return root
}

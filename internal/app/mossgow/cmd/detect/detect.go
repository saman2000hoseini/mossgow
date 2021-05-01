package detect

import (
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/cleaner"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/config"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/detect"
	"github.com/saman2000hoseini/mossgow/pkg/unzip"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func main(input, moss string, layers int, cfg config.Config) {
	if err := os.Mkdir(cfg.OutputDir, 0777); err != nil {
		if err := os.RemoveAll(cfg.OutputDir); err != nil {
			logrus.Fatal("output dir already exists")
		}

		if err := os.Mkdir(cfg.OutputDir, 0777); err != nil {
			logrus.Fatalf("could not create output dir: %s", err.Error())
		}
	}

	files, err := unzip.Unzip(input, cfg.OutputDir)
	if err != nil {
		logrus.Infof("could not extract: %s", err.Error())
	}

	cl := cleaner.New(layers, cfg)
	cl.Cleanup(files)

	if err := detect.Detect(moss, cfg); err != nil {
		logrus.Fatalf("Error executing moss: %s", err.Error())
	}
}

// Register registers server command for virtual-box binary.
func Register(root *cobra.Command, cfg config.Config) {
	var input, moss string
	var layers int

	root.AddCommand(
		&cobra.Command{
			Use:   "detect",
			Short: "detect software similarity",
			Run: func(cmd *cobra.Command, args []string) {
				main(input, moss, layers, cfg)
			},
		},
	)

	root.Flags().StringVarP(&input, "input", "I", cfg.InputDir, "To define input zip file")
	root.Flags().StringVarP(&moss, "moss", "M", cfg.MossDir, "To define path to moss")
	root.Flags().IntVarP(&layers, "layers", "L", cfg.PathLayers, "To define path layers")
}

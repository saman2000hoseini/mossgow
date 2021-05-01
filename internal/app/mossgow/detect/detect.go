package detect

import (
	"fmt"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/config"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

const DIR = "-d"

func Detect(path string, cfg config.Config) error {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	fmt.Println(dir)

	output := fmt.Sprintf("%s/%s", dir, cfg.OutputDir)
	for i := 0; i < cfg.PathLayers; i++ {
		output = fmt.Sprintf("%s/*", output)
	}

	cmd := &exec.Cmd{
		Path:   path,
		Args:   []string{path, DIR, output, output},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	logrus.Infof("Executing %s ==>", cmd.String())

	return cmd.Run()
}

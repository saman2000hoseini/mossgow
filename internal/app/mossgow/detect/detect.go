package detect

import (
	"fmt"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/config"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

const (
	DIR  = "-d"
	BASE = "-b"
)

func Detect(path, baseFile string, cfg config.Config) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	output := fmt.Sprintf("%s/%s", pwd, cfg.OutputDir)
	for i := 0; i < cfg.PathLayers; i++ {
		output = fmt.Sprintf("%s/*", output)
	}

	path = fmt.Sprintf("%s/%s", pwd, path)

	var cmd *exec.Cmd

	if baseFile != "" {
		cmd = exec.Command(path, DIR, output, output, BASE, baseFile)
	} else {
		cmd = exec.Command(path, DIR, output, output)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	logrus.Infof("Executing %s ==>", cmd.String())

	return cmd.Run()
}

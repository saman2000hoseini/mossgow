package detect

import (
	"fmt"
	"github.com/gen2brain/go-unarr"
	"github.com/mholt/archiver/v3"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/config"
	"github.com/saman2000hoseini/mossgow/pkg/unzip"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

const minFiles = 3

func main(input string, cfg config.Config) {
	if err := os.Mkdir(cfg.OutputDir, 0755); err != nil {
		if err := os.RemoveAll(cfg.OutputDir); err != nil {
			logrus.Fatal("output dir already exists")
		}

		if err := os.Mkdir(cfg.OutputDir, 0755); err != nil {
			logrus.Fatalf("could not create output dir: %s", err.Error())
		}
	}

	files, err := unzip.Unzip(input, cfg.OutputDir)
	if err != nil {
		logrus.Infof("could not extract: %s", err.Error())
	}

	extract(files, cfg.ExtraFiles)
	removeAllExtra(cfg.OutputDir, cfg.ExtraFiles)
	removeAllExtra(cfg.OutputDir, cfg.ExtraFiles)
}

func removeAllExtra(output string, extra []string) {
	items, _ := ioutil.ReadDir(output)
	for _, item := range items {
		path := fmt.Sprintf("%s/%s", output, item.Name())
		if item.IsDir() {
			subitems, _ := ioutil.ReadDir(path)
			if len(subitems) == 0 || strings.Contains(path, "venv") {
				os.RemoveAll(path)
				paths := strings.Split(path, "/")
				path = strings.Join(paths[:len(paths)-2], "/")
				if path != "" && path != output {
					removeAllExtra(path, extra)
				}
			} else {
				removeAllExtra(path, extra)
			}
		} else {
			if !strings.Contains(item.Name(), ".") {
				os.RemoveAll(path)
				paths := strings.Split(path, "/")
				path = strings.Join(paths[:len(paths)-2], "/")
				if path != "" && path != output {
					removeAllExtra(path, extra)
				}
				continue
			}

			for _, extraFile := range extra {
				if strings.Contains(item.Name(), extraFile) {
					os.RemoveAll(path)
					paths := strings.Split(path, "/")
					path = strings.Join(paths[:len(paths)-2], "/")
					if path != "" && path != output {
						removeAllExtra(path, extra)
					}
					continue
				}
			}
		}
	}
}

func extract(files, extra []string) {
	for i := range files {
		for _, extraFile := range extra {
			if strings.Contains(files[i], extraFile) {
				os.RemoveAll(files[i])
				continue
			}
		}

		fileInfo, err := os.Stat(files[i])
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			continue
		}

		parts := strings.Split(files[i], "/")
		src := parts[0]

		for j := 1; j < len(parts); j++ {
			dirFiles, err := ioutil.ReadDir(src)
			if err != nil {
				logrus.Infof("could not fetch files in dir %s: %s", src, err.Error())
			}

			src = fmt.Sprintf("%s/%s", src, parts[j])
			if len(dirFiles) > minFiles {
				break
			}
		}

		src = strings.Replace(src, "_assignsubmission_file_", "", -1)

		if strings.Contains(files[i], ".zip") || strings.Contains(files[i], ".rar") {
			extracted, err := unzip.Unzip(files[i], src)
			if err != nil {
				a, err := unarr.NewArchive(files[i])
				if err != nil {
					if err := archiver.DecompressFile("test.tar.gz", "test"); err != nil {
						logrus.Infof("could not extract: %s", err.Error())
					}
					continue
				}

				out, err := a.Extract(src)
				if err != nil {
					logrus.Infof("could not extract: %s", err.Error())
				}

				a.Close()
				extract(out, extra)
			}

			os.RemoveAll(files[i])
			extract(extracted, extra)
		} else {
			if err := move(files[i], src); err != nil {
				logrus.Infof("could not move: %s", err.Error())
			}
		}
	}
}

func move(src, dst string) error {
	parts := strings.Split(src, "/")

	dst = fmt.Sprintf("%s/%s", dst, parts[len(parts)-1])

	return os.Rename(src, dst)
}

// Register registers server command for virtual-box binary.
func Register(root *cobra.Command, cfg config.Config) {
	var input string

	root.AddCommand(
		&cobra.Command{
			Use:   "detect",
			Short: "detect software similarity",
			Run: func(cmd *cobra.Command, args []string) {
				main(input, cfg)
			},
		},
	)

	root.Flags().StringVarP(&input, "input", "I", cfg.InputDir, "To define input zip file")
}

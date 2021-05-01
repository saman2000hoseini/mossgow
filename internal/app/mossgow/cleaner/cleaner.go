package cleaner

import (
	"fmt"
	"github.com/gen2brain/go-unarr"
	"github.com/mholt/archiver/v3"
	"github.com/saman2000hoseini/mossgow/internal/app/mossgow/config"
	"github.com/saman2000hoseini/mossgow/pkg/unzip"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Cleaner struct {
	Layers             int
	SupportedLanguages []string
	Output             string
	VisitedDirs        map[string]bool
}

func New(layers int, languages []string, cfg config.Config) *Cleaner {
	return &Cleaner{
		Layers:             layers,
		SupportedLanguages: languages,
		Output:             cfg.OutputDir,
		VisitedDirs:        map[string]bool{},
	}
}

func (c *Cleaner) Cleanup(files []string) {
	c.extract(files)
	c.removeExtra(c.Output)
}

func (c *Cleaner) extract(files []string) {
	for i := 0; i < len(files); i++ {
		if strings.Contains(files[i], "venv") ||
			strings.Contains(files[i], "__") || strings.Contains(files[i], ".DS_Store") {
			if err := os.RemoveAll(files[i]); err != nil {
				logrus.Infof("could not remove: %s", err.Error())
			}
			continue
		}

		fileInfo, err := os.Stat(files[i])
		if err != nil || fileInfo.IsDir() {
			continue
		}

		parts := strings.Split(files[i], "/")
		src := strings.Join(parts[:c.Layers], "/")
		src = strings.Replace(src, "_assignsubmission_file_", "", -1)

		if strings.Contains(files[i], ".zip") || strings.Contains(files[i], ".rar") {
			extracted, err := unzip.Unzip(files[i], src)
			if err != nil {
				a, err := unarr.NewArchive(files[i])
				if err != nil {
					if err := archiver.DecompressFile(files[i], src); err != nil {
						if err := archiver.Unarchive(files[i], src); err != nil {
							logrus.Infof("could not extract: %s", err.Error())
						}
					}

					newFiles := findAllFiles(src)
					files = append(files, newFiles...)

					if err := os.RemoveAll(files[i]); err != nil {
						logrus.Infof("could not remove: %s", err.Error())
					}
					continue
				}

				out, err := a.Extract(src)
				if err != nil {
					if err := archiver.DecompressFile(files[i], src); err != nil {
						if err := archiver.Unarchive(files[i], src); err != nil {
							logrus.Infof("could not extract: %s", err.Error())
						}
					}

					newFiles := findAllFiles(src)
					files = append(files, newFiles...)

					if err := os.RemoveAll(files[i]); err != nil {
						logrus.Infof("could not remove: %s", err.Error())
					}
					continue
				}

				if err := a.Close(); err != nil {
					logrus.Infof("could not close: %s", err.Error())
				}
				c.extract(out)

				if err := os.RemoveAll(files[i]); err != nil {
					logrus.Infof("could not remove: %s", err.Error())
				}
				continue
			}

			if err := os.RemoveAll(files[i]); err != nil {
				logrus.Infof("could not remove: %s", err.Error())
			}

			c.extract(extracted)
		} else {
			if err := c.move(files[i]); err != nil {
				logrus.Infof("could not move: %s", err.Error())
			}
		}
	}
}

func (c *Cleaner) removeChild(path string) {
	if err := os.RemoveAll(path); err == nil {
		paths := strings.Split(path, "/")
		path = strings.Join(paths[:len(paths)-1], "/")

		if path != "" && path != "." && path != c.Output {
			c.removeExtra(path)
		}
	}
}

func (c *Cleaner) removeExtra(dir string) {
	items, _ := ioutil.ReadDir(dir)
	if len(items) == 0 || strings.Contains(dir, "venv") ||
		strings.Contains(dir, "__") || strings.Contains(dir, ".DS_Store") {
		c.removeChild(dir)
		return
	}

	for _, item := range items {
		path := fmt.Sprintf("%s/%s", dir, item.Name())
		if item.IsDir() {
			c.removeExtra(path)
		} else {
			if !strings.Contains(item.Name(), ".") {
				c.removeChild(path)
				continue
			}

			isSupported := false
			for _, supported := range c.SupportedLanguages {
				if strings.HasSuffix(item.Name(), supported) {
					isSupported = true
					break
				}
			}

			if !isSupported {
				c.removeChild(path)
			} else {
				parts := strings.Split(path, "/")
				if len(parts) > c.Layers {
					if err := c.move(path); err != nil {
						logrus.Infof("Error moving file: %s", err)
					}
				}
			}
		}
	}
}

func findAllFiles(root string) []string {
	files := []string{}

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			files = append(files, path)
			return nil
		})
	if err != nil {
		logrus.Info(err)
	}

	return files
}

func (c *Cleaner) move(src string) error {
	parts := strings.Split(src, "/")

	dst := strings.Join(parts[:c.Layers], "/")
	dst = fmt.Sprintf("%s/%s", dst, parts[len(parts)-1])

	return os.Rename(src, dst)
}

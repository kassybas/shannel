package snlloader

import (
	"os"

	"github.com/kassybas/shannel/internal/snlapi"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func readFile(filePath string) ([]byte, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Load reads the shannel file and returns an unmarshalled instance of it
func Load(path string) (snlapi.SnlFile, error) {
	f, err := readFile(path)
	if err != nil {
		logrus.WithField("error", err).Fatalf("could not read file '%s'", path)
	}

	snlfile := snlapi.SnlFile{}

	err = yaml.UnmarshalStrict(f, &snlfile)
	if err != nil {
		return snlapi.SnlFile{}, err
	}
	return snlfile, nil
}

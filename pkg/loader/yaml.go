package loader

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type YamlLoader struct{}

func (l *YamlLoader) LoadFromFile(location string) (interface{}, error) {
	byt, err := os.ReadFile(location)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read yaml file")
	}

	var out interface{}
	err = yaml.Unmarshal(byt, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}

	return out, nil
}

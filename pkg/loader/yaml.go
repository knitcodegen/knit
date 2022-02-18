package loader

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type YamlLoader struct{}

func (l *YamlLoader) LoadFromData(data []byte) (interface{}, error) {
	var out interface{}
	err := yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}
	return out, nil
}

func (l *YamlLoader) LoadFromFile(location string) (interface{}, error) {
	byt, err := os.ReadFile(location)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read yaml file")
	}
	return l.LoadFromData(byt)
}

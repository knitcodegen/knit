package loader

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type JsonLoader struct{}

func (l *JsonLoader) LoadFromData(data []byte) (interface{}, error) {
	var out interface{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}
	return out, nil
}

func (l *JsonLoader) LoadFromFile(location string) (interface{}, error) {
	byt, err := os.ReadFile(location)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read yaml file")
	}
	return l.LoadFromData(byt)
}

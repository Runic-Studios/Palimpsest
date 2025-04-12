package merger

import (
	"os"

	"github.com/goccy/go-yaml"
)

type YAMLLoader struct{}

func NewYAMLLoader() ConfigLoader {
	return &YAMLLoader{}
}

func (m *YAMLLoader) Load(path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &data)
	return data, err
}

func (m *YAMLLoader) Write(path string, data map[string]interface{}) error {
	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

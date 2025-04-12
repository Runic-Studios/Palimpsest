package merger

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type TOMLLoader struct{}

func NewTOMLLoader() ConfigLoader {
	return &TOMLLoader{}
}

func (m *TOMLLoader) Load(path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(content, &data)
	return data, err
}

func (m *TOMLLoader) Write(path string, data map[string]interface{}) error {
	out, err := toml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

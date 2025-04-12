package merger

import (
	"encoding/json"
	"os"
)

type JSONLoader struct{}

func NewJSONLoader() ConfigLoader {
	return &JSONLoader{}
}

func (m *JSONLoader) Load(path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &data)
	return data, err
}

func (m *JSONLoader) Write(path string, data map[string]interface{}) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

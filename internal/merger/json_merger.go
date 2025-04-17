package merger

import (
	"encoding/json"
	"os"
)

type JSONLoader struct{}

func NewJSONLoader() ConfigLoader {
	return &JSONLoader{}
}

func (m *JSONLoader) Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw interface{}
	if err := json.Unmarshal(content, &raw); err != nil {
		return nil, err
	}

	if data, ok := raw.(map[string]interface{}); ok {
		return data, nil
	}

	return map[string]interface{}{"$root": raw}, nil
}

func (m *JSONLoader) Write(path string, data Config) error {
	var toWrite interface{} = data
	if len(data) == 1 {
		if val, ok := data["$root"]; ok {
			toWrite = val
		}
	}

	out, err := json.MarshalIndent(toWrite, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

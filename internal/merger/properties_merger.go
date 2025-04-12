package merger

import (
	"errors"
	"os"

	"github.com/magiconair/properties"
)

type PropertiesLoader struct{}

func NewPropertiesLoader() ConfigLoader {
	return &PropertiesLoader{}
}

func (m *PropertiesLoader) Load(path string) (map[string]interface{}, error) {
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{}, len(p.Map()))
	for k, v := range p.Map() {
		data[k] = v
	}
	return data, nil
}

func (m *PropertiesLoader) Write(path string, data map[string]interface{}) (err error) {
	p := properties.NewProperties()
	for k, v := range data {
		if str, ok := v.(string); ok {
			if err = p.SetValue(k, str); err != nil {
				return err
			}
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	_, err = p.Write(file, properties.UTF8)
	return
}

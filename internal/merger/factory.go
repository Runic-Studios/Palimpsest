package merger

import (
	"fmt"
)

func ForExtension(ext string) (ConfigLoader, error) {
	switch ext {
	case ".yaml", ".yml":
		return NewYAMLLoader(), nil
	case ".json":
		return NewJSONLoader(), nil
	case ".toml":
		return NewTOMLLoader(), nil
	case ".properties":
		return NewPropertiesLoader(), nil
	default:
		return nil, fmt.Errorf("unsupported extension: %s", ext)
	}
}

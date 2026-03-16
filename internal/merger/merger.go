package merger

type Config = map[string]interface{}

type ConfigLoader interface {
	Load(path string) (Config, error)
	Write(path string, data Config) error
}

// Merge merges 0 or more overlays together.
// Overlays are applied in order: later overlays override earlier ones.
func Merge(overlays []Config) Config {
	result := make(Config)

	for _, overlay := range overlays {
		for k, v := range overlay {
			if existing, ok := result[k]; ok {
				switch existingTyped := existing.(type) {
				case Config:
					if vTyped, ok := v.(Config); ok {
						result[k] = Merge([]Config{existingTyped, vTyped})
						continue
					}
				case []interface{}:
					if vList, ok := v.([]interface{}); ok {
						result[k] = append(existingTyped, vList...)
						continue
					}
				}
			}
			result[k] = v
		}
	}

	return result
}

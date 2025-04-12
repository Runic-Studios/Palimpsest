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
				// If both existing and new values are Configs, merge them recursively
				if existingMap, ok := existing.(Config); ok {
					if vMap, ok := v.(Config); ok {
						result[k] = Merge([]Config{existingMap, vMap})
						continue
					}
				}
			}
			// Otherwise, overlay replaces or sets the value
			result[k] = v
		}
	}

	return result
}

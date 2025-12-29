package configuration

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	SearchPaths []string
	MaxDepth    int
	FzfFlags    string
	Verbose     bool
}

// Le fichier de configuration doit être conforme au format indiqué dans le README.md.
//
// La configuration:
//   - utilise le `configPath` fourni si donné
//   - si `configPath` est vide, utilise `$XDG_CONFIG_DIR/tmux-goway/config.toml`.
//   - si `XDG_CONFIG_DIR` est vide, utilise `$HOME/.config/tmux-goway/config.toml`.
//   - il est assumé que le fichier existe, il ne sera pas créé (et crashera le programme).
func GetConfig(configPath string) Config {
	var cfg Config

	configFile := configPath
	if configPath == "" {
		dir := os.Getenv("XDG_CONFIG_DIR")
		if dir == "" {
			home := os.Getenv("HOME")
			dir = home + "/.config/"
		}
		configFile = dir + "tmux-goway/config.toml"
	}

	content, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal([]byte(content), &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

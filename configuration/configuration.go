package configuration

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	SearchPaths    []string
	MaxDepth       int
	FzfFlags       string
	Verbose        bool
	IgnoreDotfiles bool
}

// Le fichier de configuration doit être conforme au format indiqué dans le README.md.
//
// La configuration:
//   - utilise le `configPath` si donnée
//   - utilise le config dir du système (XDG_CONFIG_HOME ou $HOME/.config sur
//     linux, %AppData% on windows, etc.) et joint "tmux-goway/config.toml" (ex:
//     `$HOME/.config/tmux-goway/config.toml`)
func GetConfig(configPath string) Config {
	var cfg Config

	configFile := configPath
	if configPath == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		configFile = filepath.Join(dir, "tmux-goway", "config.toml")
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

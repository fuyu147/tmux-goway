# Tmux-goway

Réadaptation de `tmux-sessionizer` de @ThePrimeagen en Go.

# Configuration

Le fichier par défault devrait être placé à `$XDG_CONFIG_DIR/tmux-goway/config.toml`,
un autre endroit peut être donné à l'aide du flag `-c`:
```
tmux-goway -c $HOME/.dotfiles/config/some-file.toml
```

Le fichier est considéré être en format TOML, exemple:

```toml
searchpaths = [ "~/", "~/desk" ]
maxdepth    = 2
fzfflags    = "--tmux=top,80%,50%"
```

> [!IMPORTANT]
> Si le fichier de configuration n'existe pas ou si la configuration TOML est
> invalide, alors le programme crashera (ceci est un comportement attendu et qui
> sera peut-être résolu plus tard).

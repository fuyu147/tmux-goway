# Tmux-goway

Un outil CLI pour switcher entre des sessions TMUX (ou en créé des nouvelles).
Grandement inspiré de `tmux-sessionizer` par [ThePrimeagen](https://github.com/ThePrimeagen) en
Go.

# Installation

Se référé au [Makefile](./Makefile)

# Utilisation

Par défault, `tmux-goway` va amené un menu avec `fzf`, puis va attacher à une
session `tmux` si elle existe ou la créé.
```
tmux-goway
```

En ajoutant `-p`, `tmux-goway` changera pour la précédente
session, ne fera rien si aucune session existe:
```
tmux-goway -p
```

# Configuration

Le fichier par défault devrait être placé à `$XDG_CONFIG_DIR/tmux-goway/config.toml`,
un autre endroit peut être donné en ajoutant le flag `-c`:
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

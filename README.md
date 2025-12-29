# Tmux-goway

Un outil CLI pour basculer entre des sessions TMUX (ou en créer de nouvelles).
Grandement inspiré de `tmux-sessionizer` par
[ThePrimeagen](https://github.com/ThePrimeagen) en Go.

# Installation

Référez-vous au [Makefile](./Makefile)

# Utilisation

Par défaut, `tmux-goway` affichera un menu `fzf`, puis s'attachera à une session
`tmux` existante ou en créera une nouvelle.
```
tmux-goway
```

En ajoutant `-p`, `tmux-goway` basculera vers la session précédente et ne fera
rien si aucune session de ce type n'existe:
```
tmux-goway -p
```

# Configuration

Le fichier par défaut devrait être placé à
`$XDG_CONFIG_DIR/tmux-goway/config.toml`, un autre endroit peut être donné en
ajoutant le flag `-c`:
```
tmux-goway -c $HOME/.dotfiles/config/some-file.toml
```

Le fichier doit être au format TOML, exemple:
```toml
searchpaths = [ "~/", "~/desk" ]
maxdepth    = 2
fzfflags    = "--tmux=top,80%,50%"
```

> [!IMPORTANT]
> Si le fichier de configuration n'existe pas ou si la configuration TOML est
> invalide, le programme plantera (ceci est un comportement attendu et sera
> peut-être résolu ultérieurement).
> Aucun fichier ne sera créé.

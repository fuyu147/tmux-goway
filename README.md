# Tmux-goway

Un outil CLI pour basculer entre des sessions TMUX (ou en créer de nouvelles).
Grandement inspiré de [tmux-sessionizer](https://github.com/ThePrimeagen/tmux-sessionizer) par
[ThePrimeagen](https://github.com/ThePrimeagen).

# Installation

Référez-vous au [Makefile](./Makefile)

# Utilisation

Par défaut, `tmux-goway` affichera un menu `fzf`, puis s'attachera à une session
`tmux` existante ou en créera une nouvelle au chemin sélectionné par FZF.
```
tmux-goway
```

En ajoutant `-p`, `tmux-goway` basculera vers la session précédente et ne fera
rien si aucune session de ce type n'existe:
```
tmux-goway -p
```

Un chemin de configuration custom peut être donné avec `-c`:
```
tmux-goway -c $HOME/.dotfiles/config/some-file.toml
```

# Configuration
`tmux-goway` regardera à ces endroits en ordre de priorité:
- chemin donné par `-c`
- `$XDG_CONFIG_DIR/tmux-goway/config.toml`
- `$OME/.config/tmux-goway/config.toml`

Le fichier doit être au format TOML, par exemple:
```toml
searchpaths = [ "~/personal", "~/proj", "~/art", "/run/media/user/œstrogen" ] # paths looked
maxdepth    = 2 # max depth to look for subdirectories
fzfflags    = "--tmux=top,80%,50%"
```

> [!IMPORTANT]
> Si le fichier de configuration n'existe pas ou si la configuration TOML est
> invalide, le programme plantera (ceci est un comportement attendu et sera
> peut-être résolu ultérieurement).
> Aucun fichier ne sera créé.

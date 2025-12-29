build:
	go build

local-install: build
	mkdir ~/.local/bin -p
	mv tmux-goway ~/.local/bin

root-install: build
	mv tmux-goway /usr/bin/

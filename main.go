package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"tmux-goway/configuration"
	"tmux-goway/session"
)

func main() {
	if !session.Sanity() {
		return
	}

	previous := flag.Bool("p", false, "Switch to previous session")
	configPath := flag.String("c", "", "Custom config path")
	flag.Parse()

	args := flag.Args()

	cfg := configuration.GetConfig(*configPath)

	if cfg.Verbose {
		fmt.Println("Config.SearchPaths:", cfg.SearchPaths)
		fmt.Println("Config.MaxDepth:", cfg.MaxDepth)
		fmt.Println("Config.FzfFlags:", cfg.FzfFlags)
		fmt.Println("Config.Verbose:", cfg.Verbose)
		fmt.Println("Config.IgnoreDotfiles:", cfg.IgnoreDotfiles)
	}

	if *previous && len(args) == 0 {
		session.HandlePreviousSession()
		return
	}

	userSelected := ""
	if len(args) == 1 {
		selected, err := filepath.Abs(args[0])
		if err != nil {
			panic("main(): Can't convert path to an absolute path")
		}
		if cfg.Verbose {
			fmt.Printf("Select: ARG <%s>\n", selected)
		}
		userSelected = selected
	} else if len(args) == 0 {
		selected, err := session.GetSelection(cfg)
		if err != nil {
			panic(err)
		}

		userSelected = selected
		if cfg.Verbose {
			fmt.Printf("Select: FZF <%s>\n", userSelected)
		}
	} else {
		if cfg.Verbose {
			fmt.Println("Too many args passed.")
		}
		return
	}

	if userSelected == "" {
		if cfg.Verbose {
			fmt.Println("nothing selected")
		}
		return
	}

	sessionName, sessionPath := session.GetSessionName(cfg, userSelected)
	if cfg.Verbose {
		fmt.Printf("Session name: %s, Path: %s\n", sessionName, sessionPath)
	}

	err := session.SwitchTo(cfg, sessionName, sessionPath)
	if err != nil {
		panic(err)
	}
}

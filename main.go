package main

import (
	"flag"
	"fmt"
	"os/exec"
	"tmux-goway/configuration"
	"tmux-goway/session"
)

func main() {
	if !session.Sanity() {
		return
	}
	cfg := configuration.GetConfig("")

	if cfg.Verbose {
		fmt.Println("Config.SearchPaths:", cfg.SearchPaths)
		fmt.Println("Config.MaxDepth:", cfg.MaxDepth)
		fmt.Println("Config.FzfFlags:", cfg.FzfFlags)
		fmt.Println("Config.Verbose:", cfg.Verbose)
	}

	previous := flag.Bool("p", false, "Switch to previous session")
	flag.Parse()

	args := flag.Args()

	if *previous && len(args) == 0 {
		session.HandlePreviousSession()
		return
	}

	userSelected := ""
	if len(args) == 1 {
		userSelected = args[0]
		if cfg.Verbose {
			fmt.Printf("Select: ARG <%s>\n", userSelected)
		}
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
		fmt.Println("Too many args passed.")
		return
	}

	if userSelected == "" {
		fmt.Println("nothing selected")
		return
	}

	fmt.Println(userSelected)

	sessionName, path := session.GetSessionName(userSelected)
	if cfg.Verbose {
		fmt.Printf("Session name: %s, Path: %s\n", sessionName, path)
	}

	if !session.HasSession(sessionName) {
		cmd := exec.Command("tmux", "new-session", "-ds", sessionName)
		if path != "" {
			cmd.Args = append(cmd.Args, "-c", path)
		}
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}

	err := session.SwitchTo(cfg, sessionName)
	if err != nil {
		panic(err)
	}
}

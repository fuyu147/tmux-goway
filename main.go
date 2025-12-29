package main

import (
	"fmt"
	"os"
	"os/exec"
	"tmux-goway/configuration"
	"tmux-goway/session"
)

func main() {
	if !session.Sanity() {
		return
	}
	cfg := configuration.GetConfig("")

	fmt.Println("Config.SearchPaths:", cfg.SearchPaths)
	fmt.Println("Config.MaxDepth:", cfg.MaxDepth)
	fmt.Println("Config.FzfFlags:", cfg.FzfFlags)

	args := os.Args

	for i, a := range args {
		fmt.Printf("Arg: [%d] %s\n", i, a)
	}

	userSelected := ""
	if len(args) == 2 {
		userSelected = args[1]
		fmt.Printf("Select: ARG <%s>\n", userSelected)
	} else if len(args) == 1 {
		selected, err := session.GetSelection(cfg)
		if err != nil {
			panic(err)
		}

		userSelected = selected
		fmt.Printf("Select: FZF <%s>\n", userSelected)
	}

	if userSelected == "" {
		fmt.Println("nothing selected")
		return
	}

	fmt.Println(userSelected)

	sessionName, path := session.GetSessionName(userSelected)
	fmt.Printf("Session name: %s, Path: %s\n", sessionName, path)

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

	err := session.SwitchTo(sessionName)
	if err != nil {
		panic(err)
	}
}

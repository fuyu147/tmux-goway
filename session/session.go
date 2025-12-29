package session

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"tmux-goway/configuration"
)

func Sanity() bool {
	_, errTMUX := exec.LookPath("tmux")
	_, errFZF := exec.LookPath("fzf")

	TMUXInstalled := errTMUX == nil
	FZFInstalled := errFZF == nil

	if !TMUXInstalled {
		fmt.Println("TMUX is not installed, please install it.")
	}
	if !FZFInstalled {
		fmt.Println("FZF is not installed, please install it.")
	}

	return TMUXInstalled && FZFInstalled
}

func SwitchTo(cfg configuration.Config, session string, sessionPath string) error {
	var cmd *exec.Cmd

	tmuxServerRunning := isTmuxRunning()
	insideTmux := os.Getenv("TMUX") != ""

	if cfg.Verbose {
		fmt.Println("SwitchTo: Tmux running: ", tmuxServerRunning)
		fmt.Println("SwitchTo: Inside Tmux")
	}

	if !HasSession(session) {
		err := exec.Command("tmux", "new-session", "-ds", session, "-c", sessionPath).Run()
		if err != nil {
			return err
		}
	}

	if insideTmux {
		cmd = exec.Command("tmux", "switch-client", "-t", session)
	} else {
		cmd = exec.Command("tmux", "attach-session", "-t", session)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func GetSelection(cfg configuration.Config) (string, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", fmt.Errorf("fzf not found: %v", err)
	}

	var buf bytes.Buffer
	findDirs(cfg, &buf)

	cmd := exec.Command("fzf", cfg.FzfFlags)
	cmd.Stdin = &buf
	out, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}

func HandlePreviousSession() {
	if !isTmuxRunning() {
		fmt.Println("TMUX is not running, cannot switch to previous session.")
		return
	}

	if os.Getenv("TMUX") != "" {
		// on est dans TMUX, switch utilisant `tmux switch-client -l`
		err := exec.Command("tmux", "switch-client", "-l").Run()
		if err != nil {
			fmt.Println("Failed to switch to previous session")
		}
	} else {
		// on est en dehors de TMUX, switch pour la session la plus récente

		cmd := exec.Command(
			"sh",
			"-c",
			`tmux list-sessions -F '#{session_activity} #{session_name}' | sort -rn | head -n 1 | cut -d' ' -f2-`,
		)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to list tmux sessions:", err)
			os.Exit(1)
		}

		lastSession := strings.TrimSpace(out.String())
		if lastSession == "" {
			fmt.Println("No previous session to attach to.")
			os.Exit(1)
		}

		attach := exec.Command("tmux", "attach-session", "-t", lastSession)
		attach.Stdin = os.Stdin
		attach.Stdout = os.Stdout
		attach.Stderr = os.Stderr

		if err := attach.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to attach session:", err)
			os.Exit(1)
		}

	}
}

func HasSession(session string) bool {
	err := exec.Command("tmux", "has-session", "-t", session).Run()
	return err == nil
}

func GetSessionName(cfg configuration.Config, selected string) (string, string) {
	if strings.HasPrefix(selected, "[TMUX] ") {
		sessionName := selected[7:]
		return sessionName, ""
	} else {
		path := filepath.Clean(selected)

		parent := filepath.Base(filepath.Dir(path))
		child := filepath.Base(path)

		if cfg.Verbose {
			fmt.Printf("GetSessionName: parent %s, child %s\n", parent, child)
		}

		parent = strings.ReplaceAll(parent, ".", "_")
		child = strings.ReplaceAll(child, ".", "_")

		sessionName := parent + "_" + child

		return sessionName, path
	}
}

func isTmuxRunning() bool {
	isInTmuxSession := os.Getenv("TMUX") != ""
	err := exec.Command("pgrep", "tmux").Run()
	tmuxIsRunning := err == nil
	return isInTmuxSession || tmuxIsRunning
}

func findDirs(cfg configuration.Config, w io.Writer) {
	printTmuxSessions(w)

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get home directory: %v", err)
	}

	for _, entry := range cfg.SearchPaths {
		path, depth := parsePathDepth(entry, cfg.MaxDepth)
		if home != "" && strings.HasPrefix(path, "~") {
			path = strings.Replace(path, "~", home, 1)
		}

		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			continue
		}

		walkWithDepth(w, path, depth)
	}
}

func printTmuxSessions(w io.Writer) {
	if _, err := exec.LookPath("tmux"); err != nil {
		return
	}

	var current string
	if os.Getenv("TMUX") != "" {
		out, err := exec.Command("tmux", "display-message", "-p", "#S").Output()
		if err == nil {
			current = strings.TrimSpace(string(out))
		}
	}

	cmd := exec.Command("tmux", "list-sessions", "-F", "[TMUX] #{session_name}")
	out, err := cmd.Output()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if current != "" && line == "[TMUX] "+current {
			continue
		}
		fmt.Fprintln(w, line)
	}
}

func walkWithDepth(w io.Writer, root string, maxDepth int) {
	rootDepth := strings.Count(filepath.Clean(root), string(os.PathSeparator))

	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking path %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == ".cache" {
				return filepath.SkipDir
			}

			depth := strings.Count(filepath.Clean(path), string(os.PathSeparator)) - rootDepth
			if depth > maxDepth {
				return filepath.SkipDir
			}
			fmt.Fprintln(w, path)
		}
		return nil
	})
}

func parsePathDepth(entry string, defaultDepth int) (string, int) {
	if i := strings.LastIndex(entry, ":"); i != -1 {
		if d, err := strconv.Atoi(entry[i+1:]); err == nil {
			return entry[:i], d
		}
	}
	return entry, defaultDepth
}

func parseSearchPaths(env string) []string {
	if env == "" {
		return nil
	}
	return strings.FieldsFunc(env, func(r rune) bool {
		return r == ':' || r == '\n'
	})
}

func parseIntDefault(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}

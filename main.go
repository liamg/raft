package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/liamg/aminal/gui"
	"github.com/liamg/aminal/platform"
	"github.com/liamg/aminal/terminal"
	"github.com/riywo/loginshell"
)

func main() {

	runtime.LockOSThread()

	conf := getConfig()
	logger, err := getLogger(conf)
	if err != nil {
		fmt.Printf("Failed to create logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Infof("Allocating pty...")

	pty, err := platform.NewPty(80, 25)
	if err != nil {
		logger.Fatalf("Failed to allocate pty: %s", err)
	}

	shellStr := conf.Shell
	if shellStr == "" {
		loginShell, err := loginshell.Shell()
		if err != nil {
			logger.Fatalf("Failed to ascertain your shell: %s", err)
		}
		shellStr = loginShell
	}

	os.Setenv("TERM", "xterm-256color") // controversial! easier than installing terminfo everywhere, but obviously going to be slightly different to xterm functionality, so we'll see...
	os.Setenv("COLORTERM", "truecolor")

	_, err = pty.CreateGuestProcess(shellStr)
	if err != nil {
		pty.Close()
		logger.Fatalf("Failed to start your shell: %s", err)
	}

	logger.Infof("Creating terminal...")
	terminal := terminal.New(pty, logger, conf)

	g, err := gui.New(conf, terminal, logger)
	if err != nil {
		logger.Fatalf("Cannot start: %s", err)
	}
	if err := g.Render(); err != nil {
		logger.Fatalf("Render error: %s", err)
	}

}

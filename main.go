// Package main as wedit entry point
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sofmon/wedit/builder"
	"github.com/sofmon/wedit/service"
)

const version = "1"

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

// Action requested when running wedit
type Action string

const (
	actionHelp    Action = "help"
	actionVersion Action = "version"
	actionInit    Action = "init"
	actionEdit    Action = "edit"
	actionBuild   Action = "build"
)

func determineAction() Action {

	if len(os.Args) < 2 {
		return actionEdit
	}

	return Action(strings.ToLower(os.Args[1]))
}

func main() {

	action := determineAction()

	switch action {

	case actionInit:
		initialize()

	case actionEdit:
		edit()

	case actionBuild:
		build()

	case actionVersion:
		log.Printf("★ wedit version: %s\n", version)

	case actionHelp:
		printHelp()

	default:
		log.Printf("✘ unknown action '%s' ... printing 'help' text\n", action)
		printHelp()
		os.Exit(1)
	}

	os.Exit(0)
}

func initialize() {
	log.Fatalln("✘ TODO: init action not implemented")
}

func edit() {

	err := LoadConfig()
	if err != nil {
		log.Fatalf("✘ unable to read config file. Did you forget `wedit init`? Error: %v\n", err)
	}
	err = builder.LoadConfig("wedit.json")
	if err != nil {
		log.Fatalf("✘ unable to read config file. Did you forget `wedit init`? Error: %v\n", err)
	}
	err = service.LoadConfig("wedit.json")
	if err != nil {
		log.Fatalf("✘ unable to read config file. Did you forget `wedit init`? Error: %v\n", err)
	}

	done := make(chan error)

	go func() {
		done <- service.ListenAndServe()
	}()

	if cfg.OpenBrowser {
		time.Sleep(1 * time.Second) // wait 1s for web server to start
		openBrowser(fmt.Sprintf("http://%s:%d/", cfg.Host, cfg.Port))
	}

	log.Printf("✔ running local http server on http://%s:%d/\n", cfg.Host, cfg.Port)

	err = <-done
	if err != nil {
		log.Fatalf("✘ error in serving wedit http requests doe to an error: %v\n", err)
	}
}

func build() {

	err := builder.LoadConfig("wedit.json")
	if err != nil {
		log.Fatalf("✘ unable to read config file. Did you forget `wedit init`? Error: %v\n", err)
	}

	err = builder.RebuildAll()
	if err != nil {
		log.Fatalf("✘ error in building public folder doe to an error: %v\n", err)
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {

	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command(
			filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe"),
			"url.dll,FileProtocolHandler",
			url,
		)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	}

	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Fatalf("✘ unable to start browser due to an error: %v\n", err)
		}
	}
}

func printHelp() {
	log.Print(
		"\n┊ web edit tool from https://wedit.io/\n┊\n" +
			"┊ usage: wedit [action]\n┊\n" +
			"┊ actions:\n" +
			"┊   edit    - edit current website using the default web browser (default action)\n" +
			"┊   init    - prepare the current folder as a wedit project\n" +
			"┊   build   - [re]generate the static website, for example after template change\n" +
			"┊   version - prints wedit version\n" +
			"┊   coffee  - support wedit maintainers by donating through PayPal\n" +
			"┊   help    - prints this message\n\n",
	)
}

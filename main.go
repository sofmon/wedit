// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"runtime"

	"github.com/sofmon/wedit/generator"
)

// Version is initialized in compilation time by go build.
var Version string

// Name is initialized in compilation time by go build.
var Name string

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

// Action requested when running wedit
type Action int

const (
	actionUnknown Action = iota
	actionHelp
	actionVersion
	actionInit
	actionEdit
	actionBuild
	//ActionTry
)

func determineAction() Action {
	if len(os.Args) < 2 {
		return actionEdit
	}

	switch os.Args[1] {
	case "help":
		return actionHelp
	case "version":
		return actionVersion
	case "init":
		return actionInit
	case "edit":
		return actionEdit
	case "build":
		return actionBuild
	}

	return actionUnknown
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Web edit tool from wedit.io provided by sofmon.com\n\n")
		fmt.Fprintf(os.Stderr, "Usage: wedit\n  Edit the current wedit project or shows the help text\n\n")
		fmt.Fprintf(os.Stderr, "Usage: wedit action\n")
		fmt.Fprintf(os.Stderr, "  actions:\n")
		fmt.Fprintf(os.Stderr, "    help    - Prints this message.\n")
		fmt.Fprintf(os.Stderr, "    version - Prints wedit version.\n")
		fmt.Fprintf(os.Stderr, "    init    - Prepare the current folder as a wedit project.\n")
		fmt.Fprintf(os.Stderr, "    edit    - Edit the website using web browser.\n")
		fmt.Fprintf(os.Stderr, "    build   - (Re)Generate the static website. Usualy after template change.\n")
		//fmt.Fprintf(os.Stderr, "    try - Prepare the current folder for wedit\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	action := determineAction()

	if action == actionUnknown {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	switch action {
	case actionVersion:
		log.Println(Version)
		os.Exit(0)
		break
	case actionInit:
		initialize()
		os.Exit(0)
		break
	case actionEdit:
		edit()
		os.Exit(0)
		break
	case actionBuild:
		build()
		os.Exit(0)
		break
	}
}

func initialize() {

}

func edit() {
	generator, err := generator.NewGenerator("wedit.json")
	if err != nil {
		log.Fatalf("There was an error initializing wedit. Did you forget to do 'wedit init'?")
	}
	go generator.Serve()

	var cmd *exec.Cmd

	switch runtime.GOOS {

	case "darwin":
		cmd = exec.Command("open", "http://localhost:5000/!/")
		break

	case "windows":
		cmd = exec.Command(
			filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe"),
			"url.dll,FileProtocolHandler",
			"http://localhost:5000/!/",
		)
		break

	case "linux":
		cmd = exec.Command("xdg-open", "http://localhost:5000/!/")
		break

	}

	if cmd != nil {
		cmd.Start()
	}

	<-generator.Done // Wait for end of serve
}

func build() {

}

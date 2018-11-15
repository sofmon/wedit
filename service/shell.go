// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func shellCommandHandler(w http.ResponseWriter, r *http.Request, c string) {
	fmt.Printf("executing shell command: %s\n", c)
	shCmd, ok := cfg.ShellCommands[c]
	if !ok {
		log.Printf("could not execute unknown shell command %v\n", c)
		http.Error(w, "could not execute shell command", http.StatusInternalServerError)
		return
	}

	lastQuote := rune(0)
	fieldsCheck := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	for _, script := range shCmd.Script {
		if script == "" {
			continue
		}
		fmt.Printf("$ %s\n", script)
		split := strings.FieldsFunc(script, fieldsCheck)
		cmd := exec.Command(split[0], split[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("could not execute shell command %v; error: %v;\n", c, err)
			http.Error(w, "could not execute shell command", http.StatusInternalServerError)
			return
		}
	}
}

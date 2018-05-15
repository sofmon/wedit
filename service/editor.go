// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"net/http"

	"github.com/sofmon/wedit/editor"
)

const (
	pageDataTemplate = `{"h":"","s":"","p":"","t":"","e":[],"r":[],"s":{"e":"","r":""}}`
)

/*
var editorJsPart1 *string
var editorJsPart2 *string

func getEditorJs() (*string, *string) {

	if editorJsPart1 != nil && editorJsPart2 != nil {
		return editorJsPart1, editorJsPart2
	}

	weditJSBuff := []byte(editor.EditorJSCode)

	split := strings.Split(string(weditJSBuff), pageDataTemplate)

	editorJsPart1 = &split[0]
	editorJsPart2 = &split[1]

	return editorJsPart1, editorJsPart2
}

func getPreparedJs(page model.PageWithSettings) (string, error) {

	var buffer bytes.Buffer

	varJsPart1, varJsPart2 := getEditorJs()

	_, err := buffer.WriteString(*varJsPart1)
	if err != nil {
		return "", err
	}

	pageJSON, err := json.Marshal(page)
	if err != nil {
		return "", err
	}

	_, err = buffer.WriteString(string(pageJSON))
	if err != nil {
		return "", err
	}

	_, err = buffer.WriteString(*varJsPart2)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
*/
func (s *service) editorHandler(w http.ResponseWriter, r *http.Request) {
	/*
		path := getPathWithoutAction(r)

		page, err := s.bld.ReadPageData(path)
		if err != nil {
			log.Printf("could not create page data; path: %v; error: %v;\n", path, err)
			http.Error(w, "could not create page data", http.StatusInternalServerError)
			return
		}

		pageWithSettings := model.PageWithSettings{
			Page: page,
			Settings: model.Settings{
				EditAttribute:   s.cfg.EditAttr,
				RepeatAttribute: s.cfg.RepeatAttr,
				MenuTextColor:   s.cfg.MenuTextColor,
			},
		}

		for k, v := range s.cfg.ShellCommands {
			pageWithSettings.Settings.Commands = append(
				pageWithSettings.Settings.Commands,
				model.Command{
					Name:  k,
					Color: v.Color,
				},
			)
		}

		preparedJs, err := getPreparedJs(pageWithSettings)
		if err != nil {
			log.Printf("could not prepare editor.js; path: %v; error: %v;\n", path, err)
			http.Error(w, "could not prepare editor.js", http.StatusInternalServerError)
			return
		}
	*/

	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Write([]byte(editor.EditorJSCode))
}

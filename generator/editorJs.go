// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

import (
	"bytes"
	"strings"

	"github.com/sofmon/wedit/editor"
)

const (
	page_data_template = "{\"h\":\"\",\"s\":\"\",\"p\":\"\",\"t\":\"\",\"e\":[],\"r\":[]}"
)

var editorJsPart1 *string
var editorJsPart2 *string

func getEditorJs() (*string, *string, error) {

	if editorJsPart1 != nil && editorJsPart2 != nil {
		return editorJsPart1, editorJsPart2, nil
	}

	varcharJSBuff := []byte(editor.EditorJSCode)

	split := strings.Split(string(varcharJSBuff), page_data_template)

	editorJsPart1 = &split[0]
	editorJsPart2 = &split[1]

	return editorJsPart1, editorJsPart2, nil
}

func getPreparedJs(subject Subject, page *Page) (*string, error) {

	var buffer bytes.Buffer

	varJsPart1, varJsPart2, err := getEditorJs()
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(*varJsPart1)
	if err != nil {
		return nil, err
	}

	err = WritePageAsEscapedJson(subject.Path(), subject.Host, page, &buffer)
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(*varJsPart2)
	if err != nil {
		return nil, err
	}

	var preparedJS = buffer.String()

	return &preparedJS, nil
}

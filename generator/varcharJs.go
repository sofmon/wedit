// Copyright (c) 2016, Haralampi Staykov. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
package generator

import (
	"bytes"
	"io/ioutil"
	"strings"
)

const (
	page_data_template = "{\"h\":\"\",\"s\":\"\",\"p\":\"\",\"t\":\"\",\"e\":[],\"r\":[]}"
)

var varcharJsPart1 *string
var varcharJsPart2 *string

func getVarcharJs() (*string, *string, error) {

	if varcharJsPart1 != nil && varcharJsPart2 != nil {
		return varcharJsPart1, varcharJsPart2, nil
	}

	varcharJSBuff, err := ioutil.ReadFile("varchar.js")
	if err != nil {
		return nil, nil, err
	}

	split := strings.Split(string(varcharJSBuff), page_data_template)

	varcharJsPart1 = &split[0]
	varcharJsPart2 = &split[1]

	return varcharJsPart1, varcharJsPart2, nil
}

func getPreparedJs(subject Subject, page *Page) (*string, error) {

	var buffer bytes.Buffer

	varJsPart1, varJsPart2, err := getVarcharJs()
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(*varJsPart1)
	if err != nil {
		return nil, err
	}

	err = WritePageAsEscapedJson(subject.Path, subject.Host, page, &buffer)
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

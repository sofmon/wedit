#!/bin/bash
echo "=== Start building dart applicaiton === "
dart2js ./editor/editor.dart -m -o ./editor/editor.js
rm ./editor/editor.go
echo "package editor" >> ./editor/editor.go
echo "const EditorJSCode = \`" >> ./editor/editor.go
cat ./editor/editor.js >> ./editor/editor.go
echo "\`" >> ./editor/editor.go
echo "=== Start building go applicaiton === "
go build

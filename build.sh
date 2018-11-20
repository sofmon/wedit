#!/bin/bash
echo "=== Start building dart applicaiton === "
cd ./editor; pub get; cd ..
dart2js ./editor/editor.dart -m -o ./editor/editor.js
rm ./editor/editor.go
echo "package editor" >> ./editor/editor.go
echo "const EditorJSCode = \`" >> ./editor/editor.go
sed 's/`/`+"`"+`/g' ./editor/editor.js >> ./editor/editor.go # fix bug in golang with long strings
echo "\`" >> ./editor/editor.go
echo "=== Start building go applicaiton === "
go build -o wedit .
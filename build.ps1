Write-Host "=== Start building dart applicaiton === "
cd ./editor; pub get; cd ..
dart compile js ./editor/editor.dart -m -o ./editor/editor.js
If (Test-Path ./editor/editor.go) {
    Remove-Item ./editor/editor.go 
}
$editorJS = Get-Content ./editor/editor.js
"package editor" | Out-File ./editor/editor.go -Encoding UTF8 
"const EditorJSCode = ``" | Out-File ./editor/editor.go -Encoding UTF8 -Append
$jsContent = Get-Content ./editor/editor.js
$jsContent = $jsContent.Replace("``","``+""``""+``")  #fix bug in golang with long strings
$jsContent | Out-File ./editor/editor.go -Encoding UTF8 -Append
"``" | Out-File ./editor/editor.go -Encoding UTF8 -Append
Write-Host "=== Start building go applicaiton === "
go build -o wedit

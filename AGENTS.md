# AGENTS.md - Wedit Implementation Details

This document provides comprehensive technical implementation details for AI agents working with the wedit codebase.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Package Details](#package-details)
3. [Data Models](#data-models)
4. [Core Algorithms](#core-algorithms)
5. [HTTP API](#http-api)
6. [Frontend Implementation](#frontend-implementation)
7. [Build System](#build-system)
8. [File Formats](#file-formats)
9. [Key Implementation Patterns](#key-implementation-patterns)

---

## Architecture Overview

Wedit is a static website WYSIWYG editor built with:
- **Backend**: Go (HTTP server, HTML processing, content management)
- **Frontend**: Dart compiled to JavaScript (WYSIWYG editor UI)
- **Architecture**: Single binary with embedded frontend

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  HTML Page + Injected Editor Script                   │  │
│  │  - CTRL key detection                                 │  │
│  │  - Element/Image/Repeat management                    │  │
│  │  - Auto-save on changes                               │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            ↕ HTTP API
┌─────────────────────────────────────────────────────────────┐
│                      Service Package                        │
│  - HTTP Server (localhost:5000)                             │
│  - Static file serving with editor injection                │
│  - API endpoints: load, save, image upload                  │
└─────────────────────────────────────────────────────────────┘
                            ↕
┌─────────────────────────────────────────────────────────────┐
│                      Builder Package                        │
│  - HTML parsing and rendering                               │
│  - Template resolution (variation support)                  │
│  - Content merging (local + root)                           │
│  - Markdown to HTML conversion                              │
│  - Image resizing (responsive srcset)                       │
│  - Static site generation                                   │
└─────────────────────────────────────────────────────────────┘
                            ↕
┌─────────────────────────────────────────────────────────────┐
│                      File System                            │
│  template/  - Static HTML templates with wedit attributes  │
│  content/   - JSON content files + images                   │
│  public/    - Generated static website                      │
└─────────────────────────────────────────────────────────────┘
```

---

## Package Details

### 1. Main Package

**Files**: `main.go`, `config.go`

**Responsibilities**:
- CLI command routing
- Project initialization
- Browser launching (cross-platform)

**Commands**:

```go
// main.go
func main() {
    action := "edit" // default
    if len(os.Args) > 1 {
        action = os.Args[1]
    }

    switch action {
    case "init":     initialize()
    case "edit":     edit()
    case "build":    build()
    case "clean":    clean()
    case "version":  version()
    case "help":     help()
    }
}
```

**Key Functions**:

- `initialize()`: Creates project structure
  - Creates `template/`, `content/`, `public/` folders
  - Generates default `wedit.json`
  - Copies default template if available

- `edit()`: Starts editing server
  - Loads config
  - Starts HTTP server via `service.StartEditing()`
  - Opens default browser to `http://localhost:5000/`

- `build()`: Regenerates static site
  - Loads config
  - Calls `builder.RebuildAll()`

- `clean()`: Removes orphaned content
  - Loads config
  - Calls `builder.CleanUp()`

- `openBrowser(url string)`: Cross-platform browser launching
  - macOS: `open {url}`
  - Windows: `cmd /c start {url}`
  - Linux: `xdg-open {url}`

### 2. Model Package

**Files**: `model/*.go`

**Purpose**: Shared data structures with JSON serialization

#### Page (`page.go`)

```go
type Page struct {
    Title    string
    Repeats  PageRepeats   // map[Key]Repeat
    Elements PageElements  // map[Key]Element
    Images   PageImages    // map[Key]Image
}
```

**Custom Serialization**: Maps are serialized as sorted arrays for deterministic JSON output.

```go
// MarshalJSON converts map to sorted array
func (p PageElements) MarshalJSON() ([]byte, error) {
    keys := make([]string, 0, len(p))
    for k := range p {
        keys = append(keys, string(k))
    }
    sort.Strings(keys)

    arr := make([]Element, len(keys))
    for i, k := range keys {
        arr[i] = p[Key(k)]
    }
    return json.Marshal(arr)
}
```

#### Element (`element.go`)

```go
type Element struct {
    Key  Key    `json:"k"`
    Text string `json:"t"` // Markdown content
}
```

**Usage**: Represents editable text content. Text is stored as markdown for:
- Version control friendliness (plain text diffs)
- SEO (readable content in source)
- Flexibility (supports HTML fallback)

#### Image (`image.go`)

```go
type Image struct {
    Key          Key       `json:"k"`
    Type         string    `json:"t"` // "jpg", "png", "gif"
    Width        []int     `json:"w"` // For srcset width descriptors [320, 640, 1024]
    PixelDensity []float64 `json:"x"` // For pixel density [1.0, 2.0, 3.0]
}
```

**Srcset Parsing**: Extracts descriptors from HTML attribute

```go
// NewImage parses srcset attribute
func NewImage(key Key, img *html.Node) (*Image, error) {
    srcset := getAttr(img, "srcset")
    // Parses: "img-320w.jpg 320w, img-640w.jpg 640w, img-1024w.jpg 1024w"
    // Extracts width descriptors: [320, 640, 1024]
    // Also supports pixel density: "img.jpg 1x, img@2x.jpg 2x"
}
```

**File Naming Methods**:

```go
// FileName() returns base filename: "key.jpg"
// FileNameW(320) returns: "key-320w.jpg"
// FileNameX(2.0) returns: "key-2x.jpg"
// FileNameMin() returns smallest variant
// Srcset() generates complete srcset attribute value
```

#### Repeat (`repeat.go`)

```go
type Repeat struct {
    Key      Key     `json:"k"`
    CopyKeys []Key   `json:"c"` // Ordered list of instance keys
}
```

**Copy Key Strategy**: Each repeat instance gets a unique suffix:
- Original template: `"item"`
- Copy keys: `["item", "1234567890123", "9876543210987"]`
- Element keys become: `"text"`, `"text1234567890123"`, `"text9876543210987"`

#### Key (`key.go`)

```go
type Key string

func (k Key) IsGlobal() bool {
    return strings.HasPrefix(string(k), "!")
}
```

**Global Keys**: Keys starting with `!` are shared across all pages.

#### Settings (`settings.go`)

```go
type Settings struct {
    EditAttribute    string    `json:"ea"`
    RepeatAttribute  string    `json:"ra"`
    DarkMode         bool      `json:"dm"`
    Commands         []Command `json:"cm,omitempty"`
}
```

Passed to frontend for configuration.

#### Command (`command.go`)

```go
type Command struct {
    Name  string `json:"n"`
    Color string `json:"c"`
}
```

Custom shell commands accessible from page menu.

### 3. Service Package

**Files**: `service/*.go`

**Responsibilities**:
- HTTP server
- Static file serving
- Editor script injection
- REST API for editor

#### HTTP Server (`service.go`)

```go
type Config struct {
    Port        int
    Host        string
    OpenBrowser bool
}

func StartEditing(config Config) error {
    http.HandleFunc("/~.js", editorHandler)
    http.HandleFunc("/~", apiHandler)
    http.HandleFunc("/", pageHandler)

    addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
    return http.ListenAndServe(addr, nil)
}
```

**Route Priority**:
1. `/~.js` - Editor JavaScript
2. `/~` - API endpoints (query param determines action)
3. `/` - Static files or HTML with editor injection

#### Page Handler (`page.go`)

**Purpose**: Serves static files with editor injection for HTML pages

```go
func pageHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path

    // Serve from public folder
    publicPath := filepath.Join(publicFolder, path)

    // If HTML file, inject editor script
    if strings.HasSuffix(path, ".html") {
        html, err := ioutil.ReadFile(publicPath)
        if err != nil {
            // File doesn't exist in public, serve from template
            templateHTML, _ := builder.ReadPageTemplate(path)
            html = []byte(templateHTML)
        }

        // Inject editor script before </body>
        htmlStr := string(html)
        bodyIndex := strings.LastIndex(htmlStr, "</body>")
        if bodyIndex != -1 {
            htmlStr = htmlStr[:bodyIndex] +
                      `<script src="/~.js"></script>` +
                      htmlStr[bodyIndex:]
        }

        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(htmlStr))
        return
    }

    // Serve static file as-is
    http.ServeFile(w, r, publicPath)
}
```

#### Editor Handler (`editor.go`)

```go
// editor/editor.go (generated during build)
package editor

const EditorJSCode = `...` // Entire compiled Dart JS

// service/editor.go
func editorHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/javascript")
    w.Write([]byte(editor.EditorJSCode))
}
```

#### Load Handler (`load.go`)

**Endpoint**: `GET /~?p=/path/to/page`

**Response**: Page data + editor settings

```go
type PageWithSettings struct {
    Page     model.Page     `json:"p"`
    Settings model.Settings `json:"s"`
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Query().Get("p")

    // Load page data (local + root merged)
    page, err := builder.ReadPageData(path)

    // Build settings from config
    settings := model.Settings{
        EditAttribute:   config.EditAttribute,
        RepeatAttribute: config.RepeatAttribute,
        DarkMode:        config.DarkMode,
        Commands:        config.Commands,
    }

    response := PageWithSettings{
        Page:     page,
        Settings: settings,
    }

    json.NewEncoder(w).Encode(response)
}
```

#### Save Handler (`save.go`)

**Endpoint**: `PUT /~?p=/path/to/page`

**Request Body**: Page JSON

```go
func saveHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Query().Get("p")

    // Parse incoming page
    var newPage model.Page
    json.NewDecoder(r.Body).Decode(&newPage)

    // Load existing page (for merge)
    existingPage, _ := builder.ReadPageData(path)

    // Merge: keep existing elements not in newPage
    for k, v := range existingPage.Elements {
        if _, exists := newPage.Elements[k]; !exists {
            newPage.Elements[k] = v
        }
    }
    // Same for Images and Repeats

    // Write page (triggers HTML generation)
    builder.WritePage(path, &newPage)

    w.WriteHeader(http.StatusOK)
}
```

#### Image Handler (`image.go`)

**Endpoint**: `POST /~?k={key}&n={filename}&p=/path`

**Request Body**: Image binary data

```go
func imageHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("k")
    name := r.URL.Query().Get("n")
    path := r.URL.Query().Get("p")

    // Read image data
    data, _ := ioutil.ReadAll(r.Body)

    // Save and process image (generates srcset variants)
    image, err := builder.WriteImage(path, model.Key(key), name, data)

    // Return updated image metadata
    json.NewEncoder(w).Encode(image)
}
```

### 4. Builder Package

**Files**: `builder/*.go`

**Responsibilities**:
- Template resolution
- Content management
- HTML processing
- Static site generation

#### Config (`config.go`)

```go
type Config struct {
    TemplateFolder      string
    ContentFolder       string
    PublicFolder        string
    VariationFolderName string
    RootJsonFile        string
    RootKeyPrefix       string
    EditAttribute       string
    RepeatAttribute     string
    IncludeAttribute    string
    KeepAttributes      bool
    DefaultPage         string
    AllowedPageExt      []string
    SitemapHost         string
}
```

#### Template Resolution (`template.go`)

**Purpose**: Maps URL paths to template files, handles variations

```go
func findTemplatePath(path string) (string, error) {
    // Split path: /blogs/my-article/index.html
    // → ["blogs", "my-article", "index.html"]

    subFolders := strings.Split(path, "/")
    currentPath := templateFolder

    for _, folder := range subFolders {
        // Try exact match first
        exactPath := filepath.Join(currentPath, folder)
        if exists(exactPath) {
            currentPath = exactPath
            continue
        }

        // Try variation folder
        variationPath := filepath.Join(currentPath, variationFolderName)
        if exists(variationPath) {
            currentPath = variationPath
            continue
        }

        return "", fmt.Errorf("template not found")
    }

    return currentPath, nil
}
```

**Example**:
```
Template: template/blogs/-/index.html
URLs:     /blogs/first-post/index.html  → matches
          /blogs/second-post/index.html → matches
          /blogs/anything/index.html    → matches
```

#### Page Data (`page.go`)

**ReadPageData**: Loads page content with root data merged

```go
func ReadPageData(path string) (*model.Page, error) {
    // Read page JSON
    contentPath := filepath.Join(contentFolder, path+".json")
    pageJSON, _ := ioutil.ReadFile(contentPath)

    var page model.Page
    json.Unmarshal(pageJSON, &page)

    // Merge root data
    addRootData(&page)

    // Update image srcset info from template
    templatePath, _ := findTemplatePath(path)
    templateHTML, _ := ioutil.ReadFile(templatePath)
    doc, _ := html.Parse(strings.NewReader(templateHTML))

    // Scan template for image elements, update srcset descriptors
    updateImageInfo(doc, &page)

    return &page, nil
}
```

**WritePage**: Saves content and regenerates HTML

```go
func WritePage(path string, page *model.Page) error {
    // Split root vs local data
    localPage, rootData := splitRootData(page)

    // Save local content
    contentPath := filepath.Join(contentFolder, path+".json")
    os.MkdirAll(filepath.Dir(contentPath), 0755)
    json, _ := json.MarshalIndent(localPage, "", "  ")
    ioutil.WriteFile(contentPath, json, 0644)

    // Update root data if changed
    wasUpdated := updateRootData(rootData)

    // Render HTML
    templatePath, _ := findTemplatePath(path)
    publicPath := filepath.Join(publicFolder, path)
    renderHTML(publicPath, templatePath, page)

    // Copy template assets (CSS, JS, images)
    copyTemplateAssets(templatePath, publicPath)

    // Update sitemap
    smap := loadSitemap()
    smap.Update(path)
    smap.save()

    // If root data changed, rebuild all pages
    if wasUpdated {
        RebuildAll()
    }

    return nil
}
```

#### HTML Processing (`html.go`)

**renderHTML**: Main HTML generation pipeline

```go
func renderHTML(targetPath, templatePath string, page *model.Page) error {
    // Read template
    templateHTML, _ := ioutil.ReadFile(templatePath)

    // Parse HTML
    doc, _ := html.Parse(strings.NewReader(templateHTML))

    // Process includes (recursively)
    includesProcessNode(doc, filepath.Dir(templatePath))

    // Process wedit attributes (elements, images, repeats)
    renderProcessNode(doc, page)

    // Render to file
    os.MkdirAll(filepath.Dir(targetPath), 0755)
    f, _ := os.Create(targetPath)
    defer f.Close()
    html.Render(f, doc)

    return nil
}
```

**renderProcessNode**: Recursive DOM traversal

```go
func renderProcessNode(n *html.Node, page *model.Page) {
    if n.Type == html.ElementNode {
        // Check for wedit attribute
        if editKey := getAttr(n, config.EditAttribute); editKey != "" {
            removeAttr(n, config.EditAttribute)

            if n.Data == "img" {
                processImage(model.Key(editKey), n, page)
            } else {
                processElement(model.Key(editKey), n, page)
            }
        }

        // Check for repeat attribute
        if repeatKey := getAttr(n, config.RepeatAttribute); repeatKey != "" {
            removeAttr(n, config.RepeatAttribute)
            processRepeat(model.Key(repeatKey), n, page)
            return // Don't process children (repeat does it)
        }
    }

    // Recurse children
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        renderProcessNode(c, page)
    }
}
```

**includesProcessNode**: HTML fragment inclusion

```go
func includesProcessNode(n *html.Node, templateDir string) {
    if n.Type == html.ElementNode {
        if includePath := getAttr(n, config.IncludeAttribute); includePath != "" {
            removeAttr(n, config.IncludeAttribute)

            // Load include file
            fullPath := filepath.Join(templateDir, includePath)
            includeHTML, _ := ioutil.ReadFile(fullPath)

            // Parse include
            includeDoc, _ := html.Parse(strings.NewReader(includeHTML))

            // Recursively process includes in the fragment
            includesProcessNode(includeDoc, filepath.Dir(fullPath))

            // Append children to current node
            for c := includeDoc.FirstChild; c != nil; c = c.NextSibling {
                n.AppendChild(c)
            }
        }
    }

    // Recurse
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        includesProcessNode(c, templateDir)
    }
}
```

#### Element Processing (`element.go`)

**processElement**: Converts markdown to HTML and replaces node content

```go
func processElement(key model.Key, n *html.Node, page *model.Page) {
    element, exists := page.Elements[key]
    if !exists {
        return // No content, leave template as-is
    }

    // Convert markdown to HTML
    htmlBytes := blackfriday.Run([]byte(element.Text))
    htmlStr := string(htmlBytes)

    // Remove single <p> wrapper if only one paragraph
    htmlStr = strings.TrimSpace(htmlStr)
    if strings.HasPrefix(htmlStr, "<p>") &&
       strings.HasSuffix(htmlStr, "</p>") &&
       strings.Count(htmlStr, "<p>") == 1 {
        htmlStr = htmlStr[3 : len(htmlStr)-4]
    }

    // Parse HTML
    contentDoc, _ := html.ParseFragment(strings.NewReader(htmlStr), n)

    // Remove existing children
    for c := n.FirstChild; c != nil; {
        next := c.NextSibling
        n.RemoveChild(c)
        c = next
    }

    // Append new content
    for _, c := range contentDoc {
        n.AppendChild(c)
    }
}
```

**Markdown Options**: Uses blackfriday v2 with default extensions:
- Autolink
- Strikethrough
- Tables
- Fenced code blocks
- Heading IDs

#### Image Processing (`image.go`)

**processImage**: Updates src and srcset attributes

```go
func processImage(key model.Key, n *html.Node, page *model.Page) {
    img, exists := page.Images[key]
    if !exists {
        return
    }

    // Set src to minimum size variant
    setAttr(n, "src", "./"+img.FileNameMin())

    // Set srcset with all variants
    setAttr(n, "srcset", img.Srcset())
}
```

**WriteImage**: Saves and processes uploaded image

```go
func WriteImage(path string, key model.Key, name string, data []byte) (*model.Image, error) {
    // Determine type from extension
    ext := strings.ToLower(filepath.Ext(name))
    imgType := strings.TrimPrefix(ext, ".")

    // Save original to content folder
    contentDir := filepath.Join(contentFolder, filepath.Dir(path))
    os.MkdirAll(contentDir, 0755)
    originalPath := filepath.Join(contentDir, string(key)+"."+imgType)
    ioutil.WriteFile(originalPath, data, 0644)

    // Load page to get srcset requirements
    page, _ := ReadPageData(path)
    img := page.Images[key]

    // Update type
    img.Type = imgType

    // Resize and save variants
    publicDir := filepath.Join(publicFolder, filepath.Dir(path))
    resizeImage(publicDir, img, data)

    return img, nil
}
```

**resizeImage**: Generates responsive image variants

```go
func resizeImage(folder string, img *model.Image, data []byte) error {
    // Decode original
    srcImage, _, _ := image.Decode(bytes.NewReader(data))

    os.MkdirAll(folder, 0755)

    // Generate width variants
    for _, w := range img.Width {
        // Resize using Lanczos3 (high quality)
        resized := resize.Resize(uint(w), 0, srcImage, resize.Lanczos3)

        // Save
        filename := filepath.Join(folder, img.FileNameW(w))
        f, _ := os.Create(filename)
        defer f.Close()

        switch img.Type {
        case "jpg", "jpeg":
            jpeg.Encode(f, resized, &jpeg.Options{Quality: 90})
        case "png":
            png.Encode(f, resized)
        case "gif":
            gif.Encode(f, resized, nil)
        }
    }

    // Generate pixel density variants
    for _, pd := range img.PixelDensity {
        // Similar process with different scaling
    }

    return nil
}
```

**Image Quality**: JPEG quality set to 90 for good balance between quality and size.

#### Repeat Processing (`repeat.go`)

**processRepeat**: Clones node for each repeat instance

```go
func processRepeat(key model.Key, n *html.Node, page *model.Page) {
    repeat, exists := page.Repeats[key]
    if !exists {
        return
    }

    parent := n.Parent

    // Process each copy key
    for i, copyKey := range repeat.CopyKeys {
        var clone *html.Node

        if i == 0 {
            // First copy uses original node
            clone = n
        } else {
            // Clone node
            clone = cloneNode(n)
            parent.InsertBefore(clone, n)
        }

        // Update wedit attributes with suffix
        suffix := strings.TrimPrefix(string(copyKey), string(key))
        updateWeditAttributes(clone, suffix)

        // Process clone's children
        renderProcessNode(clone, page)
    }
}
```

**updateWeditAttributes**: Adds suffix to all nested wedit attributes

```go
func updateWeditAttributes(n *html.Node, suffix string) {
    if n.Type == html.ElementNode {
        // Update wedit attribute
        if editKey := getAttr(n, config.EditAttribute); editKey != "" {
            setAttr(n, config.EditAttribute, editKey+suffix)
        }

        // Update repeat attribute
        if repeatKey := getAttr(n, config.RepeatAttribute); repeatKey != "" {
            setAttr(n, config.RepeatAttribute, repeatKey+suffix)
        }
    }

    // Recurse
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        updateWeditAttributes(c, suffix)
    }
}
```

**cloneNode**: Deep copies DOM node

```go
func cloneNode(n *html.Node) *html.Node {
    clone := &html.Node{
        Type:     n.Type,
        Data:     n.Data,
        DataAtom: n.DataAtom,
        Attr:     make([]html.Attribute, len(n.Attr)),
    }

    copy(clone.Attr, n.Attr)

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        clone.AppendChild(cloneNode(c))
    }

    return clone
}
```

#### Root Data Management (`root.go`)

**Purpose**: Manages global content shared across all pages

```go
var (
    rootData  *model.Page
    rootMutex sync.Mutex
)

func getRootData() *model.Page {
    rootMutex.Lock()
    defer rootMutex.Unlock()

    if rootData == nil {
        rootPath := filepath.Join(contentFolder, config.RootJsonFile)
        data, err := ioutil.ReadFile(rootPath)
        if err != nil {
            rootData = &model.Page{}
        } else {
            json.Unmarshal(data, &rootData)
        }
    }

    return rootData
}
```

**splitRootData**: Separates local vs global keys

```go
func splitRootData(page *model.Page) (*model.Page, *model.Page) {
    localPage := &model.Page{}
    rootPage := &model.Page{}

    // Split elements
    for k, v := range page.Elements {
        if k.IsGlobal() {
            rootPage.Elements[k] = v
        } else {
            localPage.Elements[k] = v
        }
    }
    // Same for Images and Repeats

    return localPage, rootPage
}
```

**addRootData**: Merges global content into page

```go
func addRootData(page *model.Page) {
    root := getRootData()

    for k, v := range root.Elements {
        page.Elements[k] = v
    }
    // Same for Images and Repeats
}
```

**updateRootData**: Saves root data if changed

```go
func updateRootData(newRoot *model.Page) bool {
    oldRoot := getRootData()

    // Compare JSON to detect changes
    oldJSON, _ := json.Marshal(oldRoot)
    newJSON, _ := json.Marshal(newRoot)

    if bytes.Equal(oldJSON, newJSON) {
        return false // No change
    }

    // Save new root data
    rootPath := filepath.Join(contentFolder, config.RootJsonFile)
    json, _ := json.MarshalIndent(newRoot, "", "  ")
    ioutil.WriteFile(rootPath, json, 0644)

    // Update cache
    rootMutex.Lock()
    rootData = newRoot
    rootMutex.Unlock()

    return true // Changed
}
```

#### Full Rebuild (`rebuild.go`)

**RebuildAll**: Regenerates entire static site

```go
func RebuildAll() error {
    // Copy all template assets to public
    copyDirectory(templateFolder, publicFolder)

    // Load sitemap
    smap := loadSitemap()

    // Walk content folder
    filepath.Walk(contentFolder, func(path string, info os.FileInfo, err error) error {
        if !strings.HasSuffix(path, ".json") {
            return nil
        }

        // Skip root.json
        if strings.HasSuffix(path, config.RootJsonFile) {
            return nil
        }

        // Get page path (remove content folder prefix and .json suffix)
        pagePath := strings.TrimPrefix(path, contentFolder)
        pagePath = strings.TrimSuffix(pagePath, ".json")

        // Load page data
        page, _ := ReadPageData(pagePath)

        // Render HTML
        templatePath, _ := findTemplatePath(pagePath)
        publicPath := filepath.Join(publicFolder, pagePath)
        renderHTML(publicPath, templatePath, page)

        // Update sitemap
        smap.Update(pagePath)

        return nil
    })

    // Save sitemap
    smap.save()

    return nil
}
```

#### Cleanup (`clean.go`)

**CleanUp**: Removes content files without matching templates

```go
func CleanUp() error {
    var toDelete []string

    // Walk content folder
    filepath.Walk(contentFolder, func(path string, info os.FileInfo, err error) error {
        if !strings.HasSuffix(path, ".json") {
            return nil
        }

        // Skip root.json
        if strings.HasSuffix(path, config.RootJsonFile) {
            return nil
        }

        // Get page path
        pagePath := strings.TrimPrefix(path, contentFolder)
        pagePath = strings.TrimSuffix(pagePath, ".json")

        // Check if template exists
        _, err = findTemplatePath(pagePath)
        if err != nil {
            toDelete = append(toDelete, path)
        }

        return nil
    })

    // Delete orphaned files
    for _, path := range toDelete {
        os.Remove(path)
    }

    // Remove empty directories
    removeEmptyDirs(contentFolder)

    return nil
}
```

#### Sitemap (`sitemap.go`)

```go
type sitemap struct {
    XMLName xml.Name       `xml:"urlset"`
    XMLNS   string         `xml:"xmlns,attr"`
    Urlset  []sitemapURL   `xml:"url"`
}

type sitemapURL struct {
    Loc     string `xml:"loc"`
    Lastmod string `xml:"lastmod"` // YYYY-MM-DD
}

func loadSitemap() *sitemap {
    smap := &sitemap{
        XMLNS:  "http://www.sitemaps.org/schemas/sitemap/0.9",
        Urlset: []sitemapURL{},
    }

    sitemapPath := filepath.Join(contentFolder, "sitemap.xml")
    data, err := ioutil.ReadFile(sitemapPath)
    if err == nil {
        xml.Unmarshal(data, smap)
    }

    return smap
}

func (s *sitemap) Update(pagePath string) {
    url := config.SitemapHost + pagePath
    today := time.Now().Format("2006-01-02")

    // Find existing entry
    for i, u := range s.Urlset {
        if u.Loc == url {
            s.Urlset[i].Lastmod = today
            return
        }
    }

    // Add new entry
    s.Urlset = append(s.Urlset, sitemapURL{
        Loc:     url,
        Lastmod: today,
    })
}

func (s *sitemap) save() {
    xml, _ := xml.MarshalIndent(s, "", "  ")

    // Save to both content and public
    ioutil.WriteFile(
        filepath.Join(contentFolder, "sitemap.xml"),
        xml, 0644)
    ioutil.WriteFile(
        filepath.Join(publicFolder, "sitemap.xml"),
        xml, 0644)
}
```

### 5. Editor Package

**Files**: `editor/*.dart`

**Responsibilities**:
- WYSIWYG editing interface
- CTRL key interaction
- Content editing (markdown)
- Image upload
- Repeat management
- Auto-save

#### Build Process

Dart code is compiled to JavaScript and embedded in Go:

```bash
# In editor/ directory
dart2js editor.dart -m -o editor.js

# Generate Go file with embedded JS
echo 'package editor' > editor.go
echo 'const EditorJSCode = `' >> editor.go
cat editor.js >> editor.go
echo '`' >> editor.go
```

#### Main Entry (`editor.dart`)

```dart
void main() {
    // Load page data
    var url = "/~?p=" + window.location.pathname;
    HttpRequest.getString(url).then(onDataLoaded);
}

void onDataLoaded(String responseText) {
    Map map = json.decode(responseText);
    new Page.fromMap(map);
}
```

#### Page Class (`page.dart`)

**Responsibilities**:
- Page state management
- CTRL key detection
- DOM synchronization
- Auto-save

```dart
class Page {
    String _title;
    Map<String, Element> _elements = {};
    Map<String, ImageElement> _images = {};
    Map<String, Repeat> _repeats = {};

    bool _ctrlPressed = false;
    Settings _settings;

    Page.fromMap(Map map) {
        // Parse settings
        Map settingsMap = map['s'];
        _settings = Settings.fromMap(settingsMap);

        // Parse page data
        Map pageMap = map['p'];
        _title = pageMap['t'];

        // Convert arrays to maps
        List elementsList = pageMap['e'] ?? [];
        for (Map e in elementsList) {
            var elem = Element.fromMap(e, this);
            _elements[elem.key] = elem;
        }
        // Same for images and repeats

        // Sync with DOM
        _syncElements();

        // Bind keyboard events
        window.onKeyDown.listen(_onKeyDown);
        window.onKeyUp.listen(_onKeyUp);

        // Dispatch ready event
        window.dispatchEvent(CustomEvent('wedit-ready'));
    }

    void _syncElements() {
        // Query all elements with wedit attribute
        var selector = '[' + _settings.editAttribute + ']';
        var nodes = document.querySelectorAll(selector);

        for (var node in nodes) {
            var key = node.getAttribute(_settings.editAttribute);

            if (node is html.ImageElement) {
                // Create image wrapper
                var img = ImageElement.fromNode(node, key, this);
                _images[key] = img;
            } else {
                // Create text element wrapper
                var elem = Element.fromNode(node, key, this);
                _elements[key] = elem;
            }
        }

        // Process repeats
        var repeatSelector = '[' + _settings.repeatAttribute + ']';
        var repeatNodes = document.querySelectorAll(repeatSelector);

        for (var node in repeatNodes) {
            var key = node.getAttribute(_settings.repeatAttribute);
            var repeat = Repeat.fromNode(node, key, this);
            _repeats[key] = repeat;
        }
    }

    void _onKeyDown(KeyboardEvent e) {
        if (e.ctrlKey && !_ctrlPressed) {
            _ctrlPressed = true;
            _elements.values.forEach((e) => e.mark());
            _images.values.forEach((i) => i.show());
            _repeats.values.forEach((r) => r.show());
        }
    }

    void _onKeyUp(KeyboardEvent e) {
        if (!e.ctrlKey && _ctrlPressed) {
            _ctrlPressed = false;
            _elements.values.forEach((e) => e.unmark());
            _images.values.forEach((i) => i.hide());
            _repeats.values.forEach((r) => r.hide());
        }
    }

    void save() {
        var url = "/~?p=" + window.location.pathname;
        var json = jsonEncode(toMap());

        HttpRequest.request(url,
            method: 'PUT',
            sendData: json,
            requestHeaders: {'Content-Type': 'application/json'}
        );
    }

    Map toMap() {
        return {
            't': _title,
            'e': _elements.values.map((e) => e.toMap()).toList(),
            'i': _images.values.map((i) => i.toMap()).toList(),
            'r': _repeats.values.map((r) => r.toMap()).toList(),
        };
    }
}
```

#### Element Class (`element.dart`)

**Responsibilities**:
- Text content editing
- Markdown rendering
- Visual highlighting

```dart
class Element {
    html.Element _node;
    String _key;
    String _text;
    Page _page;
    bool _editing = false;

    Element.fromMap(Map map, Page page) {
        _key = map['k'];
        _text = map['t'];
        _page = page;
    }

    Element.fromNode(html.Element node, String key, Page page) {
        _node = node;
        _key = key;
        _page = page;

        // Get text from page data or extract from DOM
        var elem = page.getElement(key);
        _text = elem?._text ?? _extractMarkdown(node);

        // Render markdown
        render();

        // Bind events
        _node.onClick.listen(_onClick);
    }

    void render() {
        if (_editing) {
            // Show raw markdown
            _node.contentEditable = 'true';
            _node.innerHtml = htmlEscape.convert(_text);
        } else {
            // Render markdown to HTML
            var html = markdownToHtml(_text);

            // Remove single <p> wrapper
            html = html.trim();
            if (html.startsWith('<p>') && html.endsWith('</p>') &&
                html.split('<p>').length == 2) {
                html = html.substring(3, html.length - 4);
            }

            _node.contentEditable = 'false';
            _node.innerHtml = html;
        }
    }

    void _onClick(Event e) {
        if (_page.isCtrlPressed) {
            e.preventDefault();
            e.stopPropagation();

            _editing = true;
            render();
            _node.focus();
        }
    }

    void mark() {
        _node.style.boxShadow = '0 0 2vw 0 rgba(0, 0, 0, .5), ' +
                                 'inset 0 0 2vw 0 rgba(255, 255, 255, .5)';
        _node.style.cursor = 'pointer';
    }

    void unmark() {
        if (!_editing) {
            _node.style.boxShadow = '';
            _node.style.cursor = '';
        }
    }

    String _extractMarkdown(html.Element node) {
        // Convert HTML to markdown (basic)
        var text = node.innerHtml;
        text = text.replaceAll('<br>', '\n');
        text = text.replaceAll('<div>', '\n');
        text = text.replaceAll('</div>', '');
        // More conversions...
        return text;
    }

    Map toMap() {
        return {'k': _key, 't': _text};
    }
}
```

#### Image Class (`image.dart`)

**Responsibilities**:
- Image upload
- Srcset management
- Visual controls

```dart
class ImageElement {
    html.ImageElement _image;
    String _key;
    String _type;
    List<int> _width;
    List<double> _pixelDensity;
    Page _page;

    html.DivElement _uploadButton;
    html.DivElement _revertButton;

    ImageElement.fromMap(Map map, Page page) {
        _key = map['k'];
        _type = map['t'];
        _width = map['w']?.cast<int>() ?? [];
        _pixelDensity = map['x']?.cast<double>() ?? [];
        _page = page;
    }

    ImageElement.fromNode(html.ImageElement node, String key, Page page) {
        _image = node;
        _key = key;
        _page = page;

        // Get metadata from page data
        var img = page.getImage(key);
        _type = img._type;
        _width = img._width;
        _pixelDensity = img._pixelDensity;

        // Create upload button
        _uploadButton = html.DivElement();
        _uploadButton.style.position = 'absolute';
        _uploadButton.style.width = '50px';
        _uploadButton.style.height = '50px';
        _uploadButton.style.borderRadius = '50%';
        _uploadButton.style.backgroundColor = 'rgba(0, 255, 0, 0.5)';
        _uploadButton.style.cursor = 'pointer';
        _uploadButton.style.display = 'none';
        _uploadButton.innerHtml = '↑'; // Upload icon
        _uploadButton.onClick.listen(_onUploadClick);

        // Similar for revert button (red)

        // Position buttons over image
        _updateButtonPosition();

        // Append to DOM
        _image.parent.append(_uploadButton);
        _image.parent.append(_revertButton);

        // Render
        render();
    }

    void render() {
        var suffix = _getSuffix();

        // Set src to minimum variant
        _image.src = './' + _key + suffix + '.' + _type;

        // Build srcset
        var srcset = '';
        for (var w in _width) {
            if (srcset.isNotEmpty) srcset += ', ';
            srcset += './' + _key + '-' + w.toString() + 'w.' + _type + ' ' + w.toString() + 'w';
        }
        for (var pd in _pixelDensity) {
            if (srcset.isNotEmpty) srcset += ', ';
            srcset += './' + _key + '-' + pd.toString() + 'x.' + _type + ' ' + pd.toString() + 'x';
        }

        if (srcset.isNotEmpty) {
            _image.srcset = srcset;
        }
    }

    void _onUploadClick(Event e) async {
        e.preventDefault();
        e.stopPropagation();

        // Create file input
        var input = html.FileUploadInputElement();
        input.accept = 'image/*';

        input.onChange.listen((e) async {
            var file = input.files[0];
            var reader = html.FileReader();

            reader.onLoadEnd.listen((e) async {
                var data = reader.result as Uint8List;

                // Upload to server
                var url = '/~?k=' + _key +
                         '&n=' + file.name +
                         '&p=' + window.location.pathname;

                var response = await HttpRequest.request(url,
                    method: 'POST',
                    sendData: data,
                    mimeType: 'application/octet-stream'
                );

                // Update metadata
                Map responseMap = json.decode(response.responseText);
                _type = responseMap['t'];
                _width = responseMap['w']?.cast<int>() ?? [];
                _pixelDensity = responseMap['x']?.cast<double>() ?? [];

                // Re-render
                render();

                // Save page
                _page.save();
            });

            reader.readAsArrayBuffer(file);
        });

        input.click();
    }

    void show() {
        _uploadButton.style.display = 'block';
        _revertButton.style.display = 'block';
    }

    void hide() {
        _uploadButton.style.display = 'none';
        _revertButton.style.display = 'none';
    }

    Map toMap() {
        return {
            'k': _key,
            't': _type,
            'w': _width,
            'x': _pixelDensity,
        };
    }
}
```

#### Repeat Class (`repeat.dart`)

**Responsibilities**:
- Managing multiple instances
- Add/remove/reorder operations
- DOM cloning

```dart
class Repeat {
    html.Element _template;
    String _key;
    List<String> _keyOrder;
    Map<String, RepeatShadow> _shadows = {};
    Page _page;

    Repeat.fromMap(Map map, Page page) {
        _key = map['k'];
        _keyOrder = map['c']?.cast<String>() ?? [];
        _page = page;
    }

    Repeat.fromNode(html.Element node, String key, Page page) {
        _template = node;
        _key = key;
        _page = page;

        // Get key order from page data
        var repeat = page.getRepeat(key);
        _keyOrder = repeat._keyOrder;

        // Create shadow for template (first in order)
        var templateKey = _keyOrder[0];
        var templateShadow = RepeatShadow(_template, templateKey, true, this);
        _shadows[templateKey] = templateShadow;

        // Clone for additional copies
        for (var i = 1; i < _keyOrder.length; i++) {
            var copyKey = _keyOrder[i];
            var clone = _cloneNode(_template);
            var suffix = copyKey.substring(_key.length);

            // Update wedit attributes with suffix
            _updateWeditAttributes(clone, suffix);

            // Insert after template
            _template.parent.insertBefore(clone, _template.nextElementSibling);

            // Create shadow
            var shadow = RepeatShadow(clone, copyKey, false, this);
            _shadows[copyKey] = shadow;

            // Register elements with page
            _registerElements(clone, suffix);
        }
    }

    void addCopy(String afterKey) {
        var newKey = _key + DateTime.now().millisecondsSinceEpoch.toString();

        // Clone template
        var clone = _cloneNode(_template);
        var suffix = newKey.substring(_key.length);
        _updateWeditAttributes(clone, suffix);

        // Insert after specified key
        var afterShadow = _shadows[afterKey];
        afterShadow._node.parent.insertBefore(clone,
                                               afterShadow._node.nextElementSibling);

        // Update key order
        var index = _keyOrder.indexOf(afterKey);
        _keyOrder.insert(index + 1, newKey);

        // Create shadow
        var shadow = RepeatShadow(clone, newKey, false, this);
        _shadows[newKey] = shadow;

        // Register elements
        _registerElements(clone, suffix);

        // Save
        _page.save();
    }

    void removeCopy(String key) {
        if (_shadows.length <= 1) return; // Can't remove last

        var shadow = _shadows[key];

        // Unregister elements
        var suffix = key.substring(_key.length);
        _unregisterElements(shadow._node, suffix);

        // Remove from DOM
        shadow._node.remove();

        // Remove from maps
        _shadows.remove(key);
        _keyOrder.remove(key);

        // Save
        _page.save();
    }

    void moveCopyUp(String key) {
        var index = _keyOrder.indexOf(key);
        if (index <= 0) return; // Already first

        // Swap in order
        var temp = _keyOrder[index - 1];
        _keyOrder[index - 1] = _keyOrder[index];
        _keyOrder[index] = temp;

        // Swap in DOM
        var shadow = _shadows[key];
        var prevShadow = _shadows[_keyOrder[index]];
        shadow._node.parent.insertBefore(shadow._node, prevShadow._node);

        // Save
        _page.save();
    }

    void moveCopyDown(String key) {
        var index = _keyOrder.indexOf(key);
        if (index >= _keyOrder.length - 1) return; // Already last

        // Swap in order
        var temp = _keyOrder[index + 1];
        _keyOrder[index + 1] = _keyOrder[index];
        _keyOrder[index] = temp;

        // Swap in DOM
        var shadow = _shadows[key];
        var nextShadow = _shadows[_keyOrder[index]];
        nextShadow._node.parent.insertBefore(nextShadow._node, shadow._node);

        // Save
        _page.save();
    }

    html.Element _cloneNode(html.Element node) {
        return node.clone(true) as html.Element;
    }

    void _updateWeditAttributes(html.Element node, String suffix) {
        // Update current node
        var editAttr = _page._settings.editAttribute;
        if (node.attributes.containsKey(editAttr)) {
            var key = node.getAttribute(editAttr);
            node.setAttribute(editAttr, key + suffix);
        }

        var repeatAttr = _page._settings.repeatAttribute;
        if (node.attributes.containsKey(repeatAttr)) {
            var key = node.getAttribute(repeatAttr);
            node.setAttribute(repeatAttr, key + suffix);
        }

        // Recurse children
        for (var child in node.children) {
            _updateWeditAttributes(child, suffix);
        }
    }

    void _registerElements(html.Element node, String suffix) {
        // Register elements with page for editing
        var editAttr = _page._settings.editAttribute;
        if (node.attributes.containsKey(editAttr)) {
            var key = node.getAttribute(editAttr);
            _page.registerElement(node, key);
        }

        // Recurse
        for (var child in node.children) {
            _registerElements(child, suffix);
        }
    }

    void show() {
        _shadows.values.forEach((s) => s.show());
    }

    void hide() {
        _shadows.values.forEach((s) => s.hide());
    }

    Map toMap() {
        return {
            'k': _key,
            'c': _keyOrder,
        };
    }
}
```

#### RepeatShadow Class (`repeatShadow.dart`)

**Responsibilities**:
- UI controls for single repeat instance
- Add/remove/move buttons

```dart
class RepeatShadow {
    html.Element _node;
    String _key;
    bool _isTemplate;
    Repeat _repeat;

    html.DivElement _addButton;
    html.DivElement _removeButton;
    html.DivElement _moveUpButton;
    html.DivElement _moveDownButton;

    RepeatShadow(this._node, this._key, this._isTemplate, this._repeat) {
        _createButtons();
    }

    void _createButtons() {
        // Add button (green)
        _addButton = html.DivElement();
        _addButton.style.position = 'absolute';
        _addButton.style.width = '30px';
        _addButton.style.height = '30px';
        _addButton.style.borderRadius = '50%';
        _addButton.style.backgroundColor = 'rgba(0, 255, 0, 0.5)';
        _addButton.style.cursor = 'pointer';
        _addButton.style.display = 'none';
        _addButton.innerHtml = '+';
        _addButton.onClick.listen((_) => _repeat.addCopy(_key));

        // Remove button (red, not for template)
        if (!_isTemplate) {
            _removeButton = html.DivElement();
            _removeButton.style.position = 'absolute';
            _removeButton.style.width = '30px';
            _removeButton.style.height = '30px';
            _removeButton.style.borderRadius = '50%';
            _removeButton.style.backgroundColor = 'rgba(255, 0, 0, 0.5)';
            _removeButton.style.cursor = 'pointer';
            _removeButton.style.display = 'none';
            _removeButton.innerHtml = '-';
            _removeButton.onClick.listen((_) => _repeat.removeCopy(_key));
        }

        // Move up button (blue)
        _moveUpButton = html.DivElement();
        _moveUpButton.style.position = 'absolute';
        _moveUpButton.style.width = '30px';
        _moveUpButton.style.height = '30px';
        _moveUpButton.style.borderRadius = '50%';
        _moveUpButton.style.backgroundColor = 'rgba(0, 0, 255, 0.5)';
        _moveUpButton.style.cursor = 'pointer';
        _moveUpButton.style.display = 'none';
        _moveUpButton.innerHtml = '↑';
        _moveUpButton.onClick.listen((_) => _repeat.moveCopyUp(_key));

        // Move down button (blue)
        _moveDownButton = html.DivElement();
        // Similar to move up
        _moveDownButton.onClick.listen((_) => _repeat.moveCopyDown(_key));

        // Append to body
        document.body.append(_addButton);
        if (_removeButton != null) document.body.append(_removeButton);
        document.body.append(_moveUpButton);
        document.body.append(_moveDownButton);
    }

    void show() {
        _updateButtonPositions();
        _addButton.style.display = 'block';
        if (_removeButton != null) _removeButton.style.display = 'block';
        _moveUpButton.style.display = 'block';
        _moveDownButton.style.display = 'block';
    }

    void hide() {
        _addButton.style.display = 'none';
        if (_removeButton != null) _removeButton.style.display = 'none';
        _moveUpButton.style.display = 'none';
        _moveDownButton.style.display = 'none';
    }

    void _updateButtonPositions() {
        var rect = _node.getBoundingClientRect();

        _addButton.style.top = '${rect.top}px';
        _addButton.style.left = '${rect.right - 30}px';

        if (_removeButton != null) {
            _removeButton.style.top = '${rect.top}px';
            _removeButton.style.left = '${rect.right - 65}px';
        }

        _moveUpButton.style.top = '${rect.top - 35}px';
        _moveUpButton.style.left = '${rect.right - 65}px';

        _moveDownButton.style.top = '${rect.top + 35}px';
        _moveDownButton.style.left = '${rect.right - 65}px';
    }
}
```

---

## Core Algorithms

### 1. Template to Public Flow

```
User Request: /blogs/my-article/index.html

1. findTemplatePath("/blogs/my-article/index.html")
   - Check: template/blogs/my-article/index.html → doesn't exist
   - Check: template/blogs/-/index.html → exists!
   - Return: "template/blogs/-/index.html"

2. ReadPageData("/blogs/my-article/index.html")
   - Load: content/blogs/my-article/index.html.json
   - Merge root data from: content/root.json
   - Update image srcset from template
   - Return: Page object

3. renderHTML("public/blogs/my-article/index.html",
              "template/blogs/-/index.html",
              page)
   - Parse template HTML
   - Process includes
   - Find wedit attributes
   - Replace content (elements, images, repeats)
   - Write to public folder

4. Result: public/blogs/my-article/index.html
```

### 2. Element Rendering Pipeline

```
Template: <div wedit="intro">Default text</div>
Content:  {"k": "intro", "t": "Hello **world**!"}

1. renderProcessNode() finds wedit="intro"
2. processElement("intro", node, page)
3. Get element text: "Hello **world**!"
4. Convert markdown to HTML: "<p>Hello <strong>world</strong>!</p>"
5. Remove <p> wrapper: "Hello <strong>world</strong>!"
6. Parse HTML and replace node children
7. Result: <div>Hello <strong>world</strong>!</div>
```

### 3. Repeat Rendering Pipeline

```
Template:
<ul>
  <li wedit-repeat="item">
    <span wedit="name">Name</span>
  </li>
</ul>

Content:
{
  "r": [{"k": "item", "c": ["item", "1234", "5678"]}],
  "e": [
    {"k": "name", "t": "First"},
    {"k": "name1234", "t": "Second"},
    {"k": "name5678", "t": "Third"}
  ]
}

1. renderProcessNode() finds wedit-repeat="item"
2. processRepeat("item", node, page)
3. Get copy keys: ["item", "1234", "5678"]
4. Process first copy (use original node):
   - Suffix: "" (no change)
   - Process children with key "name"
   - Result: <li><span>First</span></li>
5. Clone node for second copy:
   - Update attributes: wedit="name" → wedit="name1234"
   - Insert before original
   - Process children with key "name1234"
   - Result: <li><span>Second</span></li>
6. Clone node for third copy:
   - Update attributes: wedit="name" → wedit="name5678"
   - Insert before original
   - Process children with key "name5678"
   - Result: <li><span>Third</span></li>

Final HTML:
<ul>
  <li><span>Second</span></li>
  <li><span>Third</span></li>
  <li><span>First</span></li>
</ul>
```

### 4. Image Srcset Generation

```
Template:
<img wedit="hero" srcset="x 320w, y 640w, z 1024w">

Upload: user-photo.jpg (2000x1500)

1. WriteImage("/index.html", "hero", "user-photo.jpg", data)
2. Parse srcset from template: width=[320, 640, 1024]
3. Save original: content/hero.jpg
4. resizeImage():
   - Decode image (2000x1500)
   - Resize to 320w: public/hero-320w.jpg (320x240)
   - Resize to 640w: public/hero-640w.jpg (640x480)
   - Resize to 1024w: public/hero-1024w.jpg (1024x768)
5. Update page data:
   {"k": "hero", "t": "jpg", "w": [320, 640, 1024]}
6. Render HTML:
   <img src="./hero-320w.jpg"
        srcset="./hero-320w.jpg 320w, ./hero-640w.jpg 640w, ./hero-1024w.jpg 1024w">
```

### 5. Root Data Propagation

```
User edits page with root key:

1. User edits element with key "!header" on /about.html
2. Frontend saves page data
3. WritePage("/about.html", page)
4. splitRootData(page)
   - Local: elements without "!"
   - Root: elements with "!"
5. Save local: content/about.html.json
6. updateRootData(root)
   - Load existing: content/root.json
   - Compare JSON
   - If changed:
     - Save new root data
     - Return wasUpdated=true
7. If wasUpdated:
   - RebuildAll()
   - Regenerate ALL pages with updated root data
   - Result: /index.html, /about.html, /contact.html all updated
```

---

## HTTP API

### GET /~.js
**Purpose**: Serve editor JavaScript

**Response**: `application/javascript`
```javascript
// Compiled Dart code (minified)
```

### GET /~?p={path}
**Purpose**: Load page data for editing

**Example**: `GET /~?p=/blogs/my-article/index.html`

**Response**: `application/json`
```json
{
  "p": {
    "t": "My Article",
    "e": [
      {"k": "intro", "t": "Hello **world**!"},
      {"k": "body", "t": "Article content..."}
    ],
    "i": [
      {"k": "hero", "t": "jpg", "w": [320, 640, 1024], "x": null}
    ],
    "r": [
      {"k": "item", "c": ["item", "1234567890", "9876543210"]}
    ]
  },
  "s": {
    "ea": "wedit",
    "ra": "wedit-repeat",
    "dm": true,
    "cm": [
      {"n": "Deploy", "c": "green"},
      {"n": "Test", "c": "blue"}
    ]
  }
}
```

### PUT /~?p={path}
**Purpose**: Save page changes

**Example**: `PUT /~?p=/blogs/my-article/index.html`

**Request**: `application/json`
```json
{
  "t": "My Updated Article",
  "e": [
    {"k": "intro", "t": "Updated intro"},
    {"k": "body", "t": "Updated content"}
  ],
  "i": [...],
  "r": [...]
}
```

**Response**: `200 OK` (empty body)

**Side Effects**:
- Saves content JSON
- Updates root data if needed
- Regenerates HTML
- Updates sitemap
- Triggers full rebuild if root changed

### POST /~?k={key}&n={filename}&p={path}
**Purpose**: Upload image

**Example**: `POST /~?k=hero&n=photo.jpg&p=/index.html`

**Request**: `application/octet-stream` (binary image data)

**Response**: `application/json`
```json
{
  "k": "hero",
  "t": "jpg",
  "w": [320, 640, 1024],
  "x": null
}
```

**Side Effects**:
- Saves original to content folder
- Generates resized variants
- Updates page data
- Regenerates HTML

---

## File Formats

### wedit.json (Configuration)

```json
{
  "templateFolder": "template",
  "contentFolder": "content",
  "publicFolder": "public",
  "variationFolderName": "-",
  "rootJsonFile": "root.json",
  "rootKeyPrefix": "!",
  "editAttribute": "wedit",
  "repeatAttribute": "wedit-repeat",
  "includeAttribute": "wedit-include",
  "keepAttributes": false,
  "defaultPage": "index.html",
  "allowedPageExt": [".html"],
  "host": "localhost",
  "port": 5000,
  "openBrowser": true,
  "darkMode": true,
  "sitemapHost": "https://example.com",
  "commands": [
    {"name": "Deploy", "color": "green"},
    {"name": "Test", "color": "blue"}
  ]
}
```

### Page JSON (content/*.json)

```json
{
  "t": "Page Title",
  "e": [
    {
      "k": "element-key",
      "t": "Markdown **text** with [links](http://example.com)"
    }
  ],
  "i": [
    {
      "k": "image-key",
      "t": "jpg",
      "w": [320, 640, 1024],
      "x": null
    }
  ],
  "r": [
    {
      "k": "repeat-key",
      "c": ["repeat-key", "1234567890123", "9876543210987"]
    }
  ]
}
```

**Field Meanings**:
- `t` - Title (string)
- `e` - Elements array
- `i` - Images array
- `r` - Repeats array
- `k` - Key (string)
- `w` - Width descriptors (int array)
- `x` - Pixel density descriptors (float array)
- `c` - Copy keys (string array)

### sitemap.xml

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/index.html</loc>
    <lastmod>2025-01-15</lastmod>
  </url>
  <url>
    <loc>https://example.com/about.html</loc>
    <lastmod>2025-01-14</lastmod>
  </url>
</urlset>
```

---

## Build System

### Full Build Process

```bash
#!/bin/bash

# 1. Install Dart dependencies
cd editor
pub get

# 2. Compile Dart to JavaScript
dart2js editor.dart -m -o editor.js
# -m flag: minify output
# Output: editor.js + editor.js.map

# 3. Embed JavaScript in Go
echo 'package editor' > editor.go
echo '' >> editor.go
echo 'const EditorJSCode = `' >> editor.go
cat editor.js >> editor.go
echo '`' >> editor.go

cd ..

# 4. Install Go dependencies (if using modules)
go mod download

# 5. Build Go binary
go build -o wedit .

# Result: Single executable ./wedit
```

### Build Artifacts

```
wedit                    # Executable (5-10MB)
editor/editor.js         # Compiled Dart (500KB-1MB minified)
editor/editor.js.map     # Source map for debugging
editor/editor.go         # Generated Go file with embedded JS
```

---

## Key Implementation Patterns

### 1. Embedded Frontend

**Problem**: Distribute single executable
**Solution**: Embed compiled JavaScript as Go string constant

```go
package editor

const EditorJSCode = `...entire JS file...`
```

**Benefits**:
- Single binary distribution
- No external dependencies
- Simple deployment

### 2. JSON Key Compression

**Problem**: Large JSON payloads
**Solution**: Single-character keys

```json
{"t": "Title", "e": [...], "i": [...], "r": [...]}
```

**Benefits**:
- Reduced payload size (50%+ smaller)
- Faster serialization
- Lower bandwidth usage

### 3. Map to Sorted Array Serialization

**Problem**: Non-deterministic JSON output from maps
**Solution**: Custom MarshalJSON that sorts keys

```go
func (p PageElements) MarshalJSON() ([]byte, error) {
    keys := make([]string, 0, len(p))
    for k := range p {
        keys = append(keys, string(k))
    }
    sort.Strings(keys)
    // ... convert to array
}
```

**Benefits**:
- Deterministic output (git-friendly)
- Consistent diffs
- Better compression

### 4. Lazy Loading with Mutex

**Problem**: Root data accessed frequently
**Solution**: Singleton with lazy loading

```go
var (
    rootData  *model.Page
    rootMutex sync.Mutex
)

func getRootData() *model.Page {
    rootMutex.Lock()
    defer rootMutex.Unlock()

    if rootData == nil {
        // Load from file
    }

    return rootData
}
```

**Benefits**:
- Load only when needed
- Thread-safe access
- Cached for performance

### 5. DOM Cloning for Repeats

**Problem**: Duplicate HTML structure with different content
**Solution**: Deep clone DOM nodes

```go
func cloneNode(n *html.Node) *html.Node {
    clone := &html.Node{
        Type:     n.Type,
        Data:     n.Data,
        DataAtom: n.DataAtom,
        Attr:     make([]html.Attribute, len(n.Attr)),
    }
    copy(clone.Attr, n.Attr)

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        clone.AppendChild(cloneNode(c))
    }

    return clone
}
```

**Benefits**:
- Preserves structure
- Efficient memory usage
- Easy manipulation

### 6. Timestamp-based Keys

**Problem**: Unique identifiers for repeat copies
**Solution**: Millisecond timestamps

```dart
var newKey = baseKey + DateTime.now().millisecondsSinceEpoch.toString();
// Result: "item1234567890123"
```

**Benefits**:
- Guaranteed uniqueness
- Sortable (chronological order)
- No server round-trip needed

### 7. Box Shadow Highlighting

**Problem**: Highlight elements without layout shift
**Solution**: Nested box shadows

```dart
element.style.boxShadow =
    '0 0 2vw 0 rgba(0, 0, 0, .5), ' +
    'inset 0 0 2vw 0 rgba(255, 255, 255, .5)';
```

**Benefits**:
- No layout changes
- Clear visual feedback
- Works on any element

### 8. Absolute Positioning for UI

**Problem**: Editor controls shouldn't affect layout
**Solution**: Position buttons absolutely over content

```dart
button.style.position = 'absolute';
button.style.top = rect.top + 'px';
button.style.left = rect.right - 30 + 'px';
```

**Benefits**:
- No layout disruption
- Easy positioning
- Clean separation

### 9. Auto-save on Blur

**Problem**: Prevent data loss
**Solution**: Save on every blur event

```dart
element.onBlur.listen((_) {
    _editing = false;
    _text = element.innerText;
    render();
    _page.save();
});
```

**Benefits**:
- No manual save needed
- Immediate persistence
- Better UX

### 10. Full Rebuild on Root Change

**Problem**: Root data changes affect all pages
**Solution**: Detect change and rebuild everything

```go
wasUpdated := updateRootData(rootPage)
if wasUpdated {
    RebuildAll()
}
```

**Benefits**:
- Consistency guaranteed
- Simple logic
- Predictable behavior

---

## Common Development Tasks

### Adding New Model Field

1. Add field to struct in `model/`
2. Add JSON tag with single-char key
3. Update custom MarshalJSON/UnmarshalJSON if needed
4. Update Dart model in `editor/`
5. Update dataMap.dart with new key constant
6. Rebuild

### Adding New API Endpoint

1. Add handler function in `service/`
2. Register route in `service.go`
3. Update frontend to call new endpoint
4. Update AGENTS.md documentation

### Adding New Wedit Attribute

1. Add config field in `model/settings.go`
2. Load from wedit.json in builder/service config
3. Add processing logic in `builder/html.go`
4. Update Dart to handle new attribute
5. Document in README.md

### Adding New Command

1. Add command to wedit.json
2. Commands automatically appear in page menu
3. Create shell script/binary to execute
4. Test with POST to `/!{command}/`

### Debugging

**Backend**:
```bash
go run . edit
# Server logs to stdout
```

**Frontend**:
```bash
# Use non-minified Dart compilation
dart2js editor.dart -o editor.js
# Source maps available for browser debugging
```

**Template Processing**:
```go
// In builder/html.go, add logging:
fmt.Printf("Processing node: %+v\n", n)
```

---

This completes the comprehensive implementation guide for AI agents working with the wedit codebase.

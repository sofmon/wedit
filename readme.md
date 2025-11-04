# Wedit - WYSIWYG Static Website Editor

[![Website](https://img.shields.io/badge/website-wedit.dev-blue)](https://wedit.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Wedit is a command-line WYSIWYG editor that enables content management for static HTML websites directly in the browser. Edit text, manage images, and organize repeating sections with an intuitive CTRL+click interface.

## Features

- **WYSIWYG Editing**: Edit content directly in your browser with live preview
- **Markdown Support**: Write content in markdown, rendered as HTML
- **Responsive Images**: Automatic image resizing with srcset support
- **Repeatable Sections**: Dynamically add, remove, and reorder content blocks
- **Global Content**: Share headers, footers, and other elements across all pages
- **Template Variations**: One template serves multiple similar pages (blogs, localized content)
- **HTML Includes**: Reuse HTML fragments across different pages
- **Static Output**: Generates pure static HTML for any hosting
- **Version Control Friendly**: Content stored as JSON for easy diffing and merging
- **No Database**: File-based storage for maximum portability

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
  - [Initialize Project](#1-initialize-project)
  - [Setup Template](#2-setup-template)
  - [Start Editing](#3-start-editing)
  - [Build Static Site](#4-build-static-site)
  - [Deploy](#5-deploy)
- [Template Attributes](#template-attributes)
  - [Text Editing](#text-editing)
  - [Image Editing](#image-editing)
  - [Repeatable Sections](#repeatable-sections)
  - [Global Content](#global-content)
  - [HTML Includes](#html-includes)
  - [Template Variations](#template-variations)
- [Configuration](#configuration)
- [Commands](#commands)
- [Development Workflow](#development-workflow)
- [Architecture](#architecture)
- [Building from Source](#building-from-source)
- [Contributing](#contributing)
- [License](#license)

## Quick Start

```bash
# Install wedit
go get github.com/sofmon/wedit

# Create new project
mkdir my-website
cd my-website
wedit init

# Add your HTML template to template/
# Mark editable elements with wedit attributes

# Start editing
wedit

# Press CTRL in browser to highlight editable elements
# Click to edit, changes auto-save

# Deploy the public/ folder
```

## Installation

### Option 1: Using Go Get

Requires [Go](https://golang.org/) 1.11 or later.

```bash
go get github.com/sofmon/wedit
```

### Option 2: Download Binary

Download pre-built binaries from the [releases page](https://github.com/sofmon/wedit/releases).

### Option 3: Build from Source

Requires [Go](https://golang.org/) and [Dart](https://dart.dev/).

```bash
git clone https://github.com/sofmon/wedit.git
cd wedit
./build.sh  # or build.ps1 on Windows
```

## Usage

### 1. Initialize Project

Create a new directory and initialize wedit:

```bash
mkdir my-website
cd my-website
wedit init
```

This creates the following structure:

```
my-website/
├── template/       # HTML templates with wedit attributes
├── content/        # JSON content files (auto-generated)
├── public/         # Generated static website (auto-generated)
└── wedit.json      # Configuration file
```

### 2. Setup Template

Copy your static HTML files to the `template/` folder and add wedit attributes to mark editable elements.

**Example template** (`template/index.html`):

```html
<!DOCTYPE html>
<html>
<head>
    <title>My Website</title>
</head>
<body>
    <!-- Editable text -->
    <h1 wedit="page-title">Welcome to My Website</h1>

    <!-- Editable paragraph -->
    <div wedit="intro">
        This is the introduction paragraph.
    </div>

    <!-- Editable image with responsive sizes -->
    <img wedit="hero-image"
         src="hero.jpg"
         srcset="hero-320w.jpg 320w, hero-640w.jpg 640w, hero-1024w.jpg 1024w">

    <!-- Repeatable section (e.g., blog posts, team members) -->
    <ul>
        <li wedit-repeat="blog-post">
            <h2 wedit="post-title">Post Title</h2>
            <p wedit="post-excerpt">Post excerpt goes here.</p>
        </li>
    </ul>

    <!-- Global footer (shared across all pages) -->
    <footer>
        <p wedit="!copyright">© 2025 My Company</p>
    </footer>
</body>
</html>
```

### 3. Start Editing

Start the wedit server:

```bash
wedit
# or explicitly:
wedit edit
```

This will:
1. Start a local HTTP server (default: `http://localhost:5566`)
2. Open your default browser
3. Load the website in edit mode

**Editing in the browser**:
- **Press and hold CTRL** - Highlights all editable elements
- **Click an element** - Opens editor for that element
- **Edit text** - Type markdown, it renders as HTML
- **Click outside** - Auto-saves changes
- **Upload images** - Hover over image (with CTRL), click upload button
- **Manage repeats** - Use +/- buttons to add/remove/reorder items

### 4. Build Static Site

To regenerate the entire static site from content:

```bash
wedit build
```

This processes all content files and generates the complete static website in the `public/` folder.

### 5. Deploy

The `public/` folder contains your complete static website. Deploy it to any static hosting service:

- **Netlify**: Drag and drop `public/` folder
- **GitHub Pages**: Push to gh-pages branch
- **AWS S3**: Sync with `aws s3 sync public/ s3://your-bucket/`
- **Vercel**: Deploy via CLI or web interface
- **Traditional hosting**: FTP/SFTP upload `public/` contents

## Template Attributes

### Text Editing

Mark any HTML element as editable with the `wedit` attribute:

```html
<div wedit="unique-key">Content goes here</div>
```

**Supported markdown**:
- Paragraphs (separated by blank lines)
- **Bold** (`**bold**`)
- *Italic* (`*italic*`)
- [Links](http://example.com) (`[text](url)`)
- `Inline code` (`` `code` ``)
- Headers (`# H1`, `## H2`, etc.)
- Lists (`- item` or `1. item`)
- Images (`![alt](url)`)
- Inline HTML (for custom styling)

**Example**:

```html
<article wedit="blog-post">
# My First Post

This is a paragraph with **bold** and *italic* text.

- Bullet point one
- Bullet point two

Visit [my website](https://example.com) for more!
</article>
```

### Image Editing

Mark image elements as editable:

```html
<img wedit="image-key" src="default.jpg">
```

**With responsive srcset**:

```html
<img wedit="hero"
     src="hero.jpg"
     srcset="hero-320w.jpg 320w, hero-640w.jpg 640w, hero-1024w.jpg 1024w">
```

When you upload a new image:
- Original is saved to `content/`
- Multiple sizes are automatically generated (320w, 640w, 1024w)
- Srcset attribute is updated with all variants
- High-quality Lanczos3 resampling is used

**Pixel density support**:

```html
<img wedit="logo"
     src="logo.png"
     srcset="logo.png 1x, logo@2x.png 2x, logo@3x.png 3x">
```

### Repeatable Sections

Mark elements that should be repeatable with `wedit-repeat`:

```html
<div class="team-members">
    <div class="member" wedit-repeat="team-member">
        <img wedit="member-photo" src="placeholder.jpg">
        <h3 wedit="member-name">John Doe</h3>
        <p wedit="member-bio">Bio goes here</p>
    </div>
</div>
```

**Features**:
- **Add copies**: Click the green + button
- **Remove copies**: Click the red - button
- **Reorder**: Use blue up/down arrows
- **Nested editables**: Each copy can have its own editable elements

**How it works**:
- First element is the template
- Each copy gets a unique timestamp suffix
- Elements inside get the same suffix (e.g., `member-name1234567890`)
- All operations are instant and auto-saved

### CSS Class Selection

Mark elements with `wedit-class` to enable dynamic CSS class switching:

```html
<div wedit-class="hero:theme-light,theme-dark,theme-blue" class="theme-light">
    <h2 wedit="section-title">Section Title</h2>
    <p wedit="section-content">Section content goes here.</p>
</div>
```

**Usage**:
- Format: `wedit-class="key:class1,class2,class3"`
- The part before `:` is the unique identifier (key)
- The part after `:` is comma-separated available class options
- Press and hold **CTRL** to show the class dropdown
- Click a class name to apply it to the element
- The selected class replaces the element's entire class attribute
- Changes auto-save immediately

**Example use cases**:
- Theme switchers (light/dark mode)
- Layout variations (grid/list view)
- Color schemes (primary/secondary/accent)
- Size variations (small/medium/large)
- State indicators (active/inactive/disabled)

**JSON Storage**:
```json
{
  "c": [
    {"k": "hero-section", "v": "theme-dark"}
  ]
}
```

### Global Content

Prefix keys with `!` to make content global across all pages:

```html
<header>
    <h1 wedit="!site-title">My Website</h1>
    <nav wedit="!main-menu">
        <a href="/">Home</a>
        <a href="/about">About</a>
    </nav>
</header>

<footer>
    <p wedit="!copyright">© 2025 My Company</p>
</footer>
```

**Behavior**:
- Global content is stored in `content/root.json`
- Changes to global keys trigger a full site rebuild
- All pages receive the updated content
- Perfect for headers, footers, navigation, contact info

### HTML Includes

Include HTML fragments from separate files:

```html
<!-- In template/index.html -->
<header wedit-include="partials/header.html"></header>

<main>
    <div wedit="content">Page content</div>
</main>

<footer wedit-include="partials/footer.html"></footer>
```

**Included file** (`template/partials/footer.html`):

```html
<div class="footer">
    <p wedit="!company-name">ACME Corp</p>
    <p wedit="!address">123 Main St, City, State</p>
</div>
```

**Features**:
- Includes are processed at build time
- Can contain wedit attributes (local or global)
- Nested includes are supported
- Paths are relative to template folder

### Template Variations

Use a special folder name (default: `-`) as a wildcard for dynamic URLs:

**Template structure**:
```
template/
├── index.html
├── about.html
└── blogs/
    └── -/
        └── index.html   ← Template for all blog posts
```

**URL mapping**:
- `/blogs/my-first-post/` → Uses `template/blogs/-/index.html`
- `/blogs/another-article/` → Uses `template/blogs/-/index.html`
- `/blogs/anything/` → Uses `template/blogs/-/index.html`

**Content storage**:
- `/blogs/my-first-post/` → `content/blogs/my-first-post/index.html.json`
- `/blogs/another-article/` → `content/blogs/another-article/index.html.json`

**Use cases**:
- Blog posts with consistent layout
- Product pages
- Localized pages (`/services/en/`, `/services/de/`, etc.)
- Team member profiles
- Any collection of similar pages

**Creating new variations**:
1. Start wedit in edit mode
2. Navigate to the desired URL (e.g., `/blogs/new-post/`)
3. Edit content
4. Wedit creates the content file automatically

## Configuration

The `wedit.json` file controls all wedit behavior:

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
    "classAttribute": "wedit-class",
    "includeAttribute": "wedit-include",
    "keepAttributes": false,
    "defaultPage": "index.html",
    "allowedPageExt": [".html"],
    "host": "localhost",
    "port": 5566,
    "openBrowser": true,
    "darkMode": true,
    "sitemapHost": "https://example.com",
    "commands": [
        {
            "name": "Deploy",
            "color": "green"
        },
        {
            "name": "Test",
            "color": "blue"
        }
    ]
}
```

**Configuration options**:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `templateFolder` | string | `"template"` | Folder containing HTML templates |
| `contentFolder` | string | `"content"` | Folder for JSON content storage |
| `publicFolder` | string | `"public"` | Output folder for static site |
| `variationFolderName` | string | `"-"` | Folder name for URL wildcards |
| `rootJsonFile` | string | `"root.json"` | File for global content |
| `rootKeyPrefix` | string | `"!"` | Prefix for global keys |
| `editAttribute` | string | `"wedit"` | Attribute name for editables |
| `repeatAttribute` | string | `"wedit-repeat"` | Attribute name for repeats |
| `classAttribute` | string | `"wedit-class"` | Attribute name for class selection |
| `includeAttribute` | string | `"wedit-include"` | Attribute name for includes |
| `keepAttributes` | boolean | `false` | Keep wedit attributes in output |
| `defaultPage` | string | `"index.html"` | Default page for directories |
| `allowedPageExt` | array | `[".html"]` | Allowed page extensions |
| `host` | string | `"localhost"` | Server hostname |
| `port` | number | `5566` | Server port |
| `openBrowser` | boolean | `true` | Auto-open browser on start |
| `darkMode` | boolean | `true` | Dark mode for editor UI |
| `sitemapHost` | string | `""` | Base URL for sitemap.xml |
| `commands` | array | `[]` | Custom shell commands menu |

## Commands

Wedit provides several commands:

### `wedit` or `wedit edit`
Starts the editing server and opens the browser.

```bash
wedit
wedit edit
```

### `wedit init`
Initializes a new wedit project in the current directory.

```bash
wedit init
```

### `wedit build`
Rebuilds the entire static site from content files.

```bash
wedit build
```

### `wedit clean`
Removes content files that don't have matching templates (cleanup orphaned content).

```bash
wedit clean
```

### `wedit version`
Displays the wedit version.

```bash
wedit version
```

### `wedit help`
Shows help information.

```bash
wedit help
```

## Development Workflow

### Typical Workflow

1. **Design Phase**
   - Create HTML/CSS templates in `template/`
   - Add `wedit` attributes to editable elements
   - Use `wedit-repeat` for dynamic sections
   - Use `!` prefix for global elements

2. **Content Phase**
   - Run `wedit` to start editing
   - Add/edit content in the browser
   - Changes save automatically to `content/` folder
   - Static site generates in `public/` folder

3. **Review Phase**
   - Review generated site in `public/`
   - Test on different devices
   - Commit changes to version control

4. **Deploy Phase**
   - Run `wedit build` for final generation
   - Deploy `public/` folder to hosting
   - Monitor and iterate

### Version Control Best Practices

**Commit**:
- `template/` - Your HTML/CSS templates
- `content/` - All JSON content files
- `wedit.json` - Configuration
- `README.md` - Project documentation

**Ignore**:
- `public/` - Generated files (can be rebuilt)
- `.DS_Store`, `Thumbs.db` - OS files

**Example `.gitignore`**:
```
public/
.DS_Store
Thumbs.db
*.swp
```

### Team Collaboration

**Scenario**: Multiple content editors

1. One person designs templates
2. Commit templates to git
3. Each editor runs `wedit` locally
4. Editors create/modify content
5. Commit content JSON files
6. Merge via git (JSON files merge well)
7. Deploy merged result

**Scenario**: Separate design and content teams

1. Designers maintain `template/` folder
2. Content team maintains `content/` folder
3. Clear separation of concerns
4. Independent workflows
5. Merge via pull requests

## Architecture

Wedit consists of three main components:

### 1. Backend (Go)

- **HTTP Server**: Serves templates and handles API requests
- **Builder**: Processes templates and generates static HTML
- **Content Management**: Reads/writes JSON content files
- **Image Processing**: Resizes images and generates srcset variants

**Key packages**:
- `main` - CLI and command routing
- `model` - Data structures (Page, Element, Image, Repeat)
- `service` - HTTP server and API handlers
- `builder` - HTML processing and site generation
- `editor` - Embedded frontend JavaScript

### 2. Frontend (Dart → JavaScript)

- **Editor UI**: WYSIWYG editing interface
- **CTRL Detection**: Keyboard interaction for edit mode
- **DOM Manipulation**: Real-time content updates
- **Auto-save**: Persistent changes without manual save

**Compiled to JavaScript** and embedded in the Go binary for single-file distribution.

### 3. Content Storage

**JSON Format**: Human-readable, version-control friendly

```json
{
  "t": "Page Title",
  "e": [
    {"k": "intro", "t": "Markdown **content**"}
  ],
  "i": [
    {"k": "hero", "t": "jpg", "w": [320, 640, 1024], "x": null}
  ],
  "r": [
    {"k": "posts", "c": ["posts", "1234567890", "9876543210"]}
  ]
}
```

### Data Flow

```
Browser ←→ HTTP Server ←→ Builder ←→ File System
                                      ├── template/
                                      ├── content/
                                      └── public/
```

1. Browser requests page
2. Server serves HTML with injected editor
3. Editor loads page data via API
4. User edits content
5. Auto-save sends changes to server
6. Builder updates content JSON
7. Builder regenerates HTML
8. Updated page available in public/

## Building from Source

### Prerequisites

- [Go](https://golang.org/) 1.11 or later
- [Dart](https://dart.dev/) 2.0 or later

### Build Steps

```bash
# Clone repository
git clone https://github.com/sofmon/wedit.git
cd wedit

# Install Dart dependencies
cd editor
pub get

# Compile Dart to JavaScript
dart2js editor.dart -m -o editor.js

# Embed JavaScript in Go
cd ..
./embed.sh  # or embed.ps1 on Windows

# Build Go binary
go build -o wedit .

# Test
./wedit version
```

### Build Scripts

**Linux/macOS** (`build.sh`):
```bash
#!/bin/bash
cd editor
pub get
dart2js editor.dart -m -o editor.js
cd ..
echo 'package editor' > editor/editor.go
echo '' >> editor/editor.go
echo 'const EditorJSCode = `' >> editor/editor.go
cat editor/editor.js >> editor/editor.go
echo '`' >> editor/editor.go
go build -o wedit .
```

**Windows** (`build.ps1`):
```powershell
cd editor
pub get
dart2js editor.dart -m -o editor.js
cd ..
"package editor" | Out-File -FilePath editor/editor.go -Encoding UTF8
"" | Out-File -FilePath editor/editor.go -Append -Encoding UTF8
"const EditorJSCode = ``" | Out-File -FilePath editor/editor.go -Append -Encoding UTF8
Get-Content editor/editor.js | Out-File -FilePath editor/editor.go -Append -Encoding UTF8
"``" | Out-File -FilePath editor/editor.go -Append -Encoding UTF8
go build -o wedit.exe .
```

## Advanced Features

### Custom Commands

Add custom shell commands to the page menu:

```json
{
  "commands": [
    {
      "name": "Deploy to Production",
      "color": "green"
    },
    {
      "name": "Run Tests",
      "color": "blue"
    },
    {
      "name": "Backup Content",
      "color": "orange"
    }
  ]
}
```

Commands appear as buttons in the top menu when CTRL is pressed. Clicking executes `POST /!{command-name}/` which you can handle with a webhook or script.

### Sitemap Generation

Wedit automatically maintains an XML sitemap:

1. Set `sitemapHost` in `wedit.json`
2. Sitemap updates on every page save
3. Available at `public/sitemap.xml`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/index.html</loc>
    <lastmod>2025-01-15</lastmod>
  </url>
</urlset>
```

### Dark Mode

Enable dark mode for the editor UI:

```json
{
  "darkMode": true
}
```

Affects box shadows and UI colors for better visibility in low-light environments.

### Keep Attributes in Output

For debugging or special use cases, keep wedit attributes in generated HTML:

```json
{
  "keepAttributes": true
}
```

Generated HTML will include `wedit`, `wedit-repeat`, and `wedit-include` attributes.

## Troubleshooting

### Port Already in Use

If port 5566 is already in use:

```json
{
  "port": 8080
}
```

### Browser Doesn't Open

If the browser doesn't open automatically:

```json
{
  "openBrowser": false
}
```

Then manually navigate to `http://localhost:5566`

### Changes Not Appearing

1. Check browser console for errors
2. Verify content files in `content/` folder
3. Run `wedit build` to force regeneration
4. Clear browser cache

### Images Not Resizing

1. Verify srcset format in template
2. Check image file format (JPG, PNG, GIF supported)
3. Ensure Go has image processing libraries

### Template Not Found

1. Verify folder structure matches URL
2. Check `variationFolderName` in config
3. Use correct folder separator for OS

## Performance Tips

- **Large sites**: Use `wedit build` infrequently, edit mode regenerates on save
- **Many images**: Limit srcset variants to necessary sizes
- **Deep nesting**: Avoid excessive HTML nesting in repeats
- **File watching**: Not implemented, manual refresh needed after template changes

## Security Considerations

- **Local only**: Server binds to localhost by default
- **No authentication**: Intended for local development only
- **Do not expose**: Never expose wedit server to the internet
- **Production**: Deploy only the `public/` folder, not the entire project

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/sofmon/wedit.git
cd wedit
go mod download
cd editor
dart pub get
```

### Running Tests

```bash
go test ./...
```

### Code Style

- Go: Follow standard Go formatting (`gofmt`)
- Dart: Follow Dart style guide
- Commit messages: Use conventional commits

## FAQ

**Q: Can I use wedit for large websites?**
A: Yes, wedit scales well. The file-based approach is efficient and version-control friendly.

**Q: Does wedit support databases?**
A: No, wedit is intentionally file-based for simplicity and portability.

**Q: Can multiple people edit simultaneously?**
A: No, wedit is designed for local editing. Use version control for collaboration.

**Q: What browsers are supported?**
A: Modern browsers (Chrome, Firefox, Safari, Edge). IE not supported.

**Q: Can I customize the editor UI?**
A: The UI is embedded in the binary. Fork and rebuild to customize.

**Q: Does wedit support i18n?**
A: Use template variations for localization (e.g., `/pages/en/`, `/pages/de/`).

**Q: Can I use wedit with SSGs like Hugo/Jekyll?**
A: Wedit is a complete solution. However, you could use it to manage content and export for other SSGs.

**Q: Does wedit minify HTML/CSS/JS?**
A: No, use standard build tools for optimization.

**Q: Can I preview before saving?**
A: Changes are live in edit mode. Use version control to revert if needed.

**Q: Does wedit support video?**
A: Video elements can be edited like any other HTML, but no special handling.

## Resources

- **Website**: [wedit.io](https://wedit.io)
- **GitHub**: [github.com/sofmon/wedit](https://github.com/sofmon/wedit)
- **Documentation**: [wedit.io/docs](https://wedit.io/docs)
- **Issues**: [github.com/sofmon/wedit/issues](https://github.com/sofmon/wedit/issues)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Maintainers

Wedit is currently maintained by:

**Haralampi Staykov**
- Website: [haralampi.com](https://haralampi.com)
- GitHub: [@sofmon](https://github.com/sofmon)
- LinkedIn: [linkedin.com/in/haralampi](https://linkedin.com/in/haralampi)

## Acknowledgments

- Built with [Go](https://golang.org/)
- Frontend in [Dart](https://dart.dev/)
- HTML parsing: [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html)
- Markdown: [github.com/russross/blackfriday](https://github.com/russross/blackfriday)
- Image resizing: [github.com/nfnt/resize](https://github.com/nfnt/resize)

---

Made with ❤️ by [Sofmon](https://www.sofmon.com/)

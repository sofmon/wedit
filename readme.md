# Web edit tool ([wedit.io](https://wedit.io))

`wedit` uses static HTML website as template and allows the content to be changed through WYSIWYG editing using the web browser.

All changes are immediate and align with the intended website design.

# How it works

`wedit` works by:

- Opening self served HTTP server on the localhost serving a `template` folder.
- Starting the OS default browser and point it to the local server in edit mode.
- Allowing to edit the website in the browser by using `CTRL` key to highlight and edit any element marked with `wedit` attributes.
- Generates resulting static HTML files in the `public` folder.

# Usage

## 1. Initialize

To start fresh, create a new folder on your machine.

Navigate to the new folder and run:
``` sh
$ wedit init
```

This will generate the wedit folder structure:

- `template` - folder containing static HTML template.
- `content` - folder containing entered content.
- `public` - folder containing the resulting static website.
- `wedit.json` - wedit configuration file.

## 2. Static template

Copy a static HTML website to the `template` folder.

## 3. Prepare template

Every editable HTML element within template html file can have one of the attributes:
- `wedit="{unique-key}"` - marks HTML element as editable.
- `wedit-repeat="{unique-key}"` - marks HTML element as repeatable.
- `wedit-include="{path}"` - includes target HTML template file as child element (useful for headers and footers).

You can't use `wedit`, `wedit-repeat` or `wedit-include` on the same HTML element.

If `{unique-key}` start with prefix `!`, its value will persist throughout all pages with same key (useful for headers and footers)

### Text editing

Edit through markdown using paragraphs, [links](http://wedit.io), headers, bullet points, `basic` **text** *decorations*, inline images or even <span style="text-decoration:underline">HTML</span>.

```
<div wedit="text-key-1">...</div>
```

### Repeat sections

Replicated same section multiple times with different content.

```
<div class="carousel">
<div class="item" wedit-repear="repeat-key-1">
...
<div>
</div>
```

### Site-wide content

Mark content as available for the whole website, so changes are applied to all pages.

```
<div class="footer">
<span class="mantra" wedit="!footer-mantra">...</span>
</div>
```

### Images

Supports `<img>` elements.

```
<img wedit="image-key-1" src="..." srcset="...">
```

If `srcset` is used with 'width' descriptor (like `srcset="image-s.jpg 320w, image-m.jpg 640w, image-l.jpg 1024w"`), the new uploaded image will be automatically resized and different files will be supplied for each descriptor.

### Includes

Append entire HTML fragments from separate files to include similar sections in different pages.

```
<div wedit-include="footer.html"></div>
```

###  Variations

`wedit` follows the `template` folder structure for the static website. The one exception is the `variation` folder, which will be mapped to any string while `wedit` is in edit mode. By default the variation folder name is a single hyphen (`‐`).

As an example, the file in `template` folder `blogs/-/index.html` will be used as for URLs like `/blogs/some-article/index.html` and `/blogs/some-other-article/index.html`.

To create new variations, while in `wedit` edit more, navigate to the desire URL and `wedit` will create the variation in the `public` folder.

## 4. Start editing

Run `wedit` tool within the project folder:

``` sh
$ wedit
```

`wedit` opens the default browser with the template site in WYSIWYG edit mode.

Once the page is loaded at `http://{localhost}:{port}/` you can press `CTRL` key to start editing.

Entered content is saved in the `content` folder and the resulting static website is saved in the `public` folder.

## 5. Deployment

The resulting `public` folder is a static generated website that can be deployed to any hosting service/server.

# Config file

Default `wedit.json` file:
``` JSON
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
    "defaultPage":"index.html",
    "allowedPageExt": [".html"],
    "host": "localhost",
    "port": 5000,
    "openBrowser": true
}
```

# Get wedit

Requires [Go](https://golang.org/) to be installed.

``` sh
$ go get github.com/sofmon/wedit
```

# Build

The build requires setup for [Go](https://golang.org/) and [Dart](https://www.dartlang.org/) environments.

``` sh
$ ./build.sh # or build.ps1 for windows
```

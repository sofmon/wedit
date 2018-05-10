# Web edit tool ([wedit.io](https://wedit.io))

Wedit is a static website edit tool using web browser as WYSIWYG editor.

# How to use the wedit command

## 1. Create new website
Create a new folder, navigate to it and run: 
``` s
$ wedit init
```

This will generate the wedit folder structure 
- `template` - folder containing static HTML template
- `content` - folder containing entered content
- `public` - folder containing the result static website
- `wedit.json` - wedit configuration file

## 2. HTML template

### 2.1 Template folder

To remove ambiguous URLs, `wedit` expects all path to ends with `/`.

Home page path: `http(s)://{domain}/`

Sub pages path: `http(s)://{domain}/{level 1}/[{level 2}/[{level 3}/[...]]]`

`wedit` uses HTML template closest to the requested path.

Example of `wedit` template folder structure and URL mapping:


- `template`
  - {static files}
  - `index.html` - maps to `http(s)://{domain}/{anything}/`
  - {folder A}
    - {static files}
    - `index.html` - maps to `http(s)://{domain}/{folder A}/{anything}/`
    - {folder A.A}
      - {static files}
      - `index.html` - maps to `http(s)://{domain}/{folder A}/{folder A.A}/{anything}/`
      - ...
    - ...
  - {folder B}
    - {static files}
    - `index.html` - maps to `http(s)://{domain}/{folder B}/{anything}/`
    - ...
  - ...


### 2.1 Template HTML files `index.html`

Every editable element within template `index.html` file can have one of the attributes:
- `wedit="{unique-key}"` - marks HTML element as editable
- `wedit-repeat="{unique-key}"` - marks HTML element as repeatable
- `wedit-include="{path}"` - includes target HTML template file as child element (useful for headers and footers)

You can not use `wedit`, `wedit-repeat` or `wedit-include` on the same HTML element.

## 3. Edit the website through the browser

Run `wedit` tool within the project folder:

``` s
$ wedit
```

`wedit` opens the default browser with the template site in WYSIWYG edit mode.

Once the page is loaded at `http://{localhost}:{port}/!/...` you can press `CTRL` key to start edit.

As long as you are in edit mode (under the `/!/` root path), you can generate new pages just by navigating to the desired URL.

Once done press `CTRL+S` to save the page.

The result is saved in the `public` folder.

# How to build

The build requires setup for [Go](https://golang.org/) and [Dart](https://www.dartlang.org/) environments.

- On Linux or MAC run **./build.sh**
- On Windows run **./build.ps1**

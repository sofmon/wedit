# Web edit tool ([wedit.io](https://wedit.io))

`wedit` is a static website edit tool using web browser as WYSIWYG editor.

# How it works

`wedit` works by:

- Opening self served HTTP server on the localhost serving the `template` folder with `wedit` editor script

- Starting the os default browser and point it to the local server in edit mode (under the `/!/` root path)

- Allowing to edit the website in the browser by using `CTRL` key to highlight and edit any element marked with `wedit` attributes

- Generating result static HTML files in the `public` folder

# How to use the wedit command

## 1. Create a new folder

To start fresh, create a new folder on your machine.

## 2. Initialize

Navigate to the new folder rand run:
``` s
$ wedit init
```

This will generate the wedit folder structure:

- `template` - folder containing static HTML template
- `content` - folder containing entered content
- `public` - folder containing the result static website
- `wedit.json` - wedit configuration file

## 2. Provide static template

Copy a static HTML website to the `template` folder.

## 3. Update the HTML template

`wedit` required every servable web folder to have `index.html` file.

### Mark element as editable

Every editable element within template `index.html` file can have one of the attributes:
- `wedit="{unique-key}"` - marks HTML element as editable
- `wedit-repeat="{unique-key}"` - marks HTML element as repeatable
- `wedit-include="{path}"` - includes target HTML template file as child element (useful for headers and footers)

You can not use `wedit`, `wedit-repeat` or `wedit-include` on the same HTML element.

If `{unique-key}` start with prefix `!`, its value will be the same side wide (useful for headers and footers)


### URL mapping

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

### Prepare for editing

Wedit serves all public assets from the `public` folder. To prepare for editing run:

``` s
$ wedit build
```

The `build` command can also regenerate the whole website if the `template` changes

## 3. Edit website through the browser

Run `wedit` tool within the project folder:

``` s
$ wedit
```

`wedit` opens the default browser with the template site in WYSIWYG edit mode.

Once the page is loaded at `http://{localhost}:{port}/!/...` you can press `CTRL` key to start edit.

As long as you are in edit mode (under the `/!/` root path), you can generate new pages just by navigating to the desired URL.

Once done go to the menu (`CTRL` + hover on top of the page) and press `save` to save the page.

The result is saved in the `public` folder.

## 4. Deploy result to hosting service

The resulting `public` folder is a static generated website that can be deployed to any hosting service/server.

# Wedit config file

Default `wedit.json` file with two example `shellCommands`:
``` JSON
{
    "builder": {
        "templateFolder": "./template",
        "contentFolder": "./content",
        "publicFolder": "./public",
        "templateHtmlFile": "index.html",
        "pageHtmlFile": "index.html",
        "pageJsonFile": "index.json",
        "rootJsonFile": "root.json",
        "rootKeyPrefix": "!",
        "editAttribute": "wedit",
        "repeatAttribute": "wedit-repeat",
        "includeAttribute": "wedit-include",
        "keepAttributes": false
    },
    "service": {
        "host": "localhost",
        "port": 5000,
        "openBrowser": true,
        "menuTextColor": "#ffffff",
        "shellCommands": {
            "reset": {
                "color": "#99d250",
                "script": [
                    "git reset --hard origin/master",
                    "git pull"
                ]
            }
            ,
            "publish": {
                "color": "#d40d22",
                "script": [
                    "git add content",
                    "git add public",
                    "git commit -m 'content update by wedit'",
                    "git push"
                ]
            }
        }
    }
}
```

The example `shellCommands` shows how `wedit` can be used as content editing tool connected to `git`

# How to build

The build requires setup for [Go](https://golang.org/) and [Dart](https://www.dartlang.org/) environments.

- On Linux or MAC run **./build.sh**
- On Windows run **./build.ps1**

# Web edit tool ([wedit.io](https://wedit.io))

Wedit is a static website generation tool, which uses the web browser as WYSIWYG editor.

# How to use the wedit command

## 1. Create new webiste
Create a new folder, navigate to it and run: 
> wedit init

This will generate the wedit folder structure
- **public** - folder containing the static generated website 
- **template** - folder containing the HTML template
- **wedit.json** - wedit configuration file

## 2. Create an HTML template

### 2.1 Template folder
Wedit considered that web path always ends with '/'. This is required to ensure non ambiguous URLs.

Home page path:
> http(s)://< domain >/

Sub pages path:
> http(s)://< domain >/< level 1 >/[< level 2 >/[< level 3 >/[...]]]

When generating the static website, wedit will use the template that is closest to the requested path.

Example of wedit **template** folder structure and how it is mapped to URLs:
> - template
>   - < static files >
>   - index.html - used for: _http(s)://< domain >/< anything >/_
>   - < folder 1 >
>     - < static files >
>     - index.html - used for: _http(s)://< domain >/< folder 1 >/< anything >/_
>     - < folder 1.1 >
>       - < static files >
>       - index.html - used for: _http(s)://< domain >/< folder 1 >/< folder 1.1 >/< anything >/_
>       - ...
>   - < folder 2 >
>     - < static files >
>     - index.html - used for: _http(s)://< domain >/< folder 2 >/< anything >/_
>     - ...
>   - ...

### 2.1 Template data-var='...' and data-var-repeat='...' attributes

Within any **index.html** file for every element that needs to be editable, add one of the attributes:
- **data-var='< unique key >'** - makes the content of the HTML element editable through markdown
- **data-var-repeat='< unique key >'** - makes the HTML element repeatable

You can not use **data-var** and **data-var-repeat** on the same HTML element.

The **< unique key >** us used to map specific content to the specific place in the HTML.

## 3. Edit the website through the browser
Just run the wedit tool within the project folder:
> wedit

Wedit will open your default browser and allow you to edit the page in a WYSIWYG way.

Once the page is loaded and you are in http://< localhost >:< port >**/!/**... you can press **CTRL** key to start edit.

As long as you are in edit mode (under the **/!/** root path), you can generate new pages just by navigating to the desired URL.

Once done press **CTRL+S** to save the page.

The result is saved in the **public** folder.

# How to build

The build requires setup for [Go](https://golang.org/) and [Dart](https://www.dartlang.org/) environments.

- On Linux or MAC run **./build.sh**
- Windows build is coming soon, but you can check the **./build.sh** file for clues
- Brew, Scoop and other packages are coming soon

# Hey, this thing doesn't work on Windows!

Please be patient. If you read this page, you most probably found it while Wedit is still finding its way to production readiness.

A full Windows support is coming in few days.

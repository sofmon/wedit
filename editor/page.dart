// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Page {
  String _host, _site, _path;

  String _title;
  String get title => _title;
  set title(String t) => _title = t;

  Map<String, Map> _mappedElementsData;
  Map<String, Map> _mappedRepeatsData;
  Map<String, Map> _mappedImagesData;

  Map<String, Element> _elements;
  Map<String, Image> _images;
  Map<String, Repeat> _repeats;

  bool _ctrlPressed = false;
  bool get ctrlPressed => _ctrlPressed;

  String _editAttribute;
  String get editAttribute => _editAttribute;
  String _repeatAttribute;
  String get repeatAttribute => _repeatAttribute;

  PageMenu _pageMenu;
  Map<String, String> _commands;

  Page.fromMap(Map map) {
    _host = map[PAGE_HOST];
    _site = map[PAGE_SITE];
    _path = map[PAGE_PATH];

    var settings = map[PAGE_SETTINGS];
    _editAttribute = settings[PAGE_SETTINGS_EDITATTR];
    _repeatAttribute = settings[PAGE_SETTINGS_REPEATATTR];

    _commands = new Map<String, String>();
    List shellCommands = settings[PAGE_SETTINGS_COMMANDS];
    if (shellCommands != null) {
      settings[PAGE_SETTINGS_COMMANDS].forEach((c) =>
          _commands[c[PAGE_SETTINGS_COMMANDS_NAME]] =
              c[PAGE_SETTINGS_COMMANDS_COLOR]);
    }

    _title = map[PAGE_TITLE];

    _elements = new Map<String, Element>();
    _images = new Map<String, Image>();
    _repeats = new Map<String, Repeat>();

    _mappedElementsData = new Map<String, Map>();
    map[PAGE_ELEMENTS].forEach((e) => _mappedElementsData[e[ELEMENT_KEY]] = e);

    _mappedRepeatsData = new Map<String, Map>();
    map[PAGE_REPEATS].forEach((r) => _mappedRepeatsData[r[REPEAT_KEY]] = r);

    _mappedImagesData = new Map<String, Map>();
    map[PAGE_IMAGES].forEach((r) => _mappedImagesData[r[IMAGE_KEY]] = r);

    _syncElements();
    _bindEvents();

    html.window.dispatchEvent(new html.Event("wedit-ready"));

    _pageMenu =
        new PageMenu(this, _commands, settings[PAGE_SETTINGS_MENUTEXTCOLOR]);
  }

  Map toMap() {
    Map map = new Map();

    map[PAGE_HOST] = _host;
    map[PAGE_SITE] = _site;
    map[PAGE_PATH] = _path;

    map[PAGE_TITLE] = _title;

    List eList = new List();
    _elements.values.forEach((e) => eList.add(e.toMap()));
    map[PAGE_ELEMENTS] = eList;

    List rList = new List();
    _repeats.values.forEach((r) => rList.add(r.toMap()));
    map[PAGE_REPEATS] = rList;

    List iList = new List();
    _images.values.forEach((i) => iList.add(i.toMap()));
    map[PAGE_IMAGES] = iList;

    return map;
  }

  void registerElement(html.Element domElement) {
    var key = domElement.getAttribute(_editAttribute);

    if (key == null || key.isEmpty) return;

    if (Image.SUPPORTED_IMAGE_TAGS.contains(domElement.tagName.toLowerCase())) {
      Map cmsData = _mappedImagesData[key];

      Image image = new Image.fromMap(this, key, domElement, cmsData);

      _images[key] = image;

      image.render();
      return;
    }

    Map cmsData = _mappedElementsData[key];

    Element element = new Element.fromMap(this, key, domElement, cmsData);

    _elements[key] = element;

    element.render();
  }

  void unregisterElement(html.Element domElement) {
    var key = domElement.getAttribute(_editAttribute);

    if (Image.SUPPORTED_IMAGE_TAGS.contains(domElement.tagName.toLowerCase())) {
      var image = _images[key];
      image.prepareDomForHtmlSave();
      _images.remove(key);
      return;
    }

    var element = _elements[key];
    element.prepareDomForHtmlSave();
    _elements.remove(key);
  }

  void _syncElements() {
    html.document.title = _title;

    List<html.Element> repeatedDomElements = new List<html.Element>();

    html.ElementList<html.Element> domElements = html.querySelectorAll(
        "[" + _editAttribute + "],[" + _repeatAttribute + "]");

    for (int i = 0; i < domElements.length; i++) {
      html.Element domElement = domElements[i];

      String key = domElement.getAttribute(_repeatAttribute);

      if (key != null && !key.isEmpty) {
        Map cmsData = _mappedRepeatsData[key];

        Repeat repeat = new Repeat.fromMap(this, key, domElement, cmsData);
        _repeats[key] = repeat;

        List<html.Element> repeatedElements = repeat.render();

        for (int ri = 0; ri < repeatedElements.length; ri++) {
          _processDomElement(repeatedElements[ri]);
        }

        continue;
      }

      _processDomElement(domElement);
    }
  }

  void _processDomElement(html.Element domElement) {
    String key = domElement.getAttribute(_editAttribute);

    if (key == null || key.isEmpty) return;

    String elementTag = domElement.tagName.toLowerCase();

    if (Image.SUPPORTED_IMAGE_TAGS.contains(elementTag)) {
      Map cmsData = _mappedImagesData[key];

      Image image = new Image.fromMap(this, key, domElement, cmsData);
      _images[key] = image;

      image.render();
      return;
    }

    Map cmsData = _mappedElementsData[key];

    Element element = new Element.fromMap(this, key, domElement, cmsData);
    _elements[key] = element;

    element.render();
  }

  void _bindEvents() {
    html.window.onKeyDown.listen(_windowKeyDown);
    html.window.onKeyUp.listen(_windowKeyUp);
  }

  void _windowKeyDown(html.KeyboardEvent e) {
    _ctrlPressed = e.ctrlKey;

    if (e.ctrlKey) {
      _elements.values.forEach((e) => e.highlight());
      _images.values.forEach((e) => e.highlight());
      _repeats.values.forEach((r) => r.highlight());
    }
  }

  void _windowKeyUp(html.KeyboardEvent e) {
    _ctrlPressed = e.ctrlKey;

    _elements.values.forEach((e) => e.normalise());
    _images.values.forEach((e) => e.normalise());
    _repeats.values.forEach((r) => r.normalise());
  }

  void save(Function onSuccess, Function onFailure) {
    Map pageData = toMap();

    String jsonData = convert.json.encode(pageData);

    html.HttpRequest request =
        new html.HttpRequest(); // create a new XHR./wedit

    // add an event handler that is called when the request finishes
    request.onReadyStateChange.listen((_) {
      if (request.readyState == html.HttpRequest.DONE &&
          (request.status == 200 || request.status == 0)) {
        // data saved OK.
        print(request.responseText);
        // output the response from the server
        if (onSuccess != null) {
          onSuccess();
        }
      } else {
        if (onFailure != null) {
          onFailure();
        }
      }
    });

    var url = html.window.location.href.replaceAll("/!/", "/!save/");

    request.open("POST", url, async: false);
    request.send(jsonData);
  }

  void command(String cmd, Function onSuccess, Function onFailure) {
    Map pageData = toMap();

    String jsonData = convert.json.encode(pageData);

    html.HttpRequest request =
        new html.HttpRequest(); // create a new XHR./wedit

    // add an event handler that is called when the request finishes
    request.onReadyStateChange.listen((_) {
      if (request.readyState == html.HttpRequest.DONE &&
          (request.status == 200 || request.status == 0)) {
        // data saved OK.
        print(request.responseText);
        // output the response from the server
        if (onSuccess != null) {
          onSuccess();
        }
      } else {
        if (onFailure != null) {
          onFailure();
        }
      }
    });

    var url = html.window.location.href.replaceAll("/!/", "/!" + cmd + "/");

    request.open("POST", url, async: false);
    request.send(jsonData);
  }

  void prepareDomForHtmlSave() {
    _elements.values.forEach((e) => e.normalise());
    _images.values.forEach((e) => e.normalise());
    _repeats.values.forEach((r) => r.normalise());

    _elements.values.forEach((e) => e.prepareDomForHtmlSave());
    _images.values.forEach((e) => e.prepareDomForHtmlSave());
    _repeats.values.forEach((r) => r.prepareDomForHtmlSave());
  }

  void restoreDomAfterHtmlSave() {
    _elements.values.forEach((e) => e.restoreDomAfterHtmlSave());
    _images.values.forEach((e) => e.restoreDomAfterHtmlSave());
    _repeats.values.forEach((r) => r.restoreDomAfterHtmlSave());
  }
}

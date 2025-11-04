// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Page {
  String _host = '';
  String _path = '';

  String _title = '';
  String get title => _title;
  set title(String t) => _title = t;

  Map<String, Map<String, dynamic>> _mappedElementsData = {};
  Map<String, Map<String, dynamic>> _mappedRepeatsData = {};
  Map<String, Map<String, dynamic>> _mappedImagesData = {};
  Map<String, Map<String, dynamic>> _mappedClassesData = {};

  Map<String, Element> _elements = {};
  Map<String, Image> _images = {};
  Map<String, Repeat> _repeats = {};
  Map<String, Class> _classes = {};

  bool _ctrlPressed = false;
  bool get ctrlPressed => _ctrlPressed;

  String _editAttribute = '';
  String get editAttribute => _editAttribute;
  String _repeatAttribute = '';
  String get repeatAttribute => _repeatAttribute;
  String _classAttribute = '';
  String get classAttribute => _classAttribute;
  bool _darkMode = false;
  bool get darkMode => _darkMode;

  PageMenu? _pageMenu;
  Map<String, String> _commands = {};

  Page.fromMap(Map<String, dynamic> map) {
    _host = map[PAGE_HOST] as String? ?? '';
    //_site = map[PAGE_SITE];
    _path = map[PAGE_PATH] as String? ?? '';


    final settings = map[PAGE_SETTINGS] as Map<String, dynamic>?;
    if (settings != null) {
      _editAttribute = settings[PAGE_SETTINGS_EDITATTR] as String? ?? '';
      _repeatAttribute = settings[PAGE_SETTINGS_REPEATATTR] as String? ?? '';
      _classAttribute = settings[PAGE_SETTINGS_CLASSATTR] as String? ?? '';
      _darkMode = settings[PAGE_SETTINGS_DARK_MODE] as bool? ?? false;

      _commands = <String, String>{};
      final shellCommands = settings[PAGE_SETTINGS_COMMANDS] as List?;
      if (shellCommands != null) {
        for (var c in shellCommands) {
          final cmd = c as Map<String, dynamic>;
          _commands[cmd[PAGE_SETTINGS_COMMANDS_NAME] as String] =
              cmd[PAGE_SETTINGS_COMMANDS_COLOR] as String;
        }
      }
    }

    _title = map[PAGE_TITLE] as String? ?? '';

    _elements = <String, Element>{};
    _images = <String, Image>{};
    _repeats = <String, Repeat>{};
    _classes = <String, Class>{};

    _mappedElementsData = <String, Map<String, dynamic>>{};
    (map[PAGE_ELEMENTS] as List?)?.forEach((e) => _mappedElementsData[(e as Map)[ELEMENT_KEY] as String] = e as Map<String, dynamic>);

    _mappedRepeatsData = <String, Map<String, dynamic>>{};
    (map[PAGE_REPEATS] as List?)?.forEach((r) => _mappedRepeatsData[(r as Map)[REPEAT_KEY] as String] = r as Map<String, dynamic>);

    _mappedImagesData = <String, Map<String, dynamic>>{};
    (map[PAGE_IMAGES] as List?)?.forEach((r) => _mappedImagesData[(r as Map)[IMAGE_KEY] as String] = r as Map<String, dynamic>);

    _mappedClassesData = <String, Map<String, dynamic>>{};
    (map[PAGE_CLASSES] as List?)?.forEach((c) => _mappedClassesData[(c as Map)[CLASS_KEY] as String] = c as Map<String, dynamic>);

    _syncElements();
    _bindEvents();

    web.window.dispatchEvent(web.Event('wedit-ready'));

    _pageMenu = PageMenu(this, _commands, (settings?[PAGE_SETTINGS_MENUTEXTCOLOR] as String?) ?? '');

    web.window.postMessage("wedit.loaded".toJS, "*".toJS);
  }

  Map<String, dynamic> toMap() {
    final map = <String, dynamic>{};

    map[PAGE_HOST] = _host;
    //map[PAGE_SITE] = _site;
    map[PAGE_PATH] = _path;

    map[PAGE_TITLE] = _title;

    final eList = <Map<String, dynamic>>[];
    _elements.values.forEach((e) => eList.add(e.toMap()));
    map[PAGE_ELEMENTS] = eList;

    final rList = <Map<String, dynamic>>[];
    _repeats.values.forEach((r) => rList.add(r.toMap()));
    map[PAGE_REPEATS] = rList;

    final iList = <Map<String, dynamic>>[];
    _images.values.forEach((i) => iList.add(i.toMap()));
    map[PAGE_IMAGES] = iList;

    final cList = <Map<String, dynamic>>[];
    _classes.values.forEach((c) => cList.add(c.toMap()));
    map[PAGE_CLASSES] = cList;

    return map;
  }

  void registerElement(web.HTMLElement domElement) {
    final key = domElement.getAttribute(_editAttribute);

    if (key == null || key.isEmpty) return;

    if (Image.SUPPORTED_IMAGE_TAGS.contains(domElement.tagName.toLowerCase())) {
      final cmsData = _mappedImagesData[key];

      final image = Image.fromMap(this, key, domElement, cmsData);

      _images[key] = image;

      image.render();
      return;
    }

    final cmsData = _mappedElementsData[key];

    final element = Element.fromMap(this, key, domElement, cmsData);

    _elements[key] = element;

    element.render();
  }

  void unregisterElement(web.HTMLElement domElement) {
    final key = domElement.getAttribute(_editAttribute);

    if (Image.SUPPORTED_IMAGE_TAGS.contains(domElement.tagName.toLowerCase())) {
      final image = _images[key];
      image?.prepareDomForHtmlSave();
      _images.remove(key);
      return;
    }

    final element = _elements[key];
    element?.prepareDomForHtmlSave();
    _elements.remove(key);
  }

  void _syncElements() {
    web.document.title = _title;

    final domElements = web.document.querySelectorAll(
        "[$_editAttribute],[$_repeatAttribute],[$_classAttribute]");

    for (int i = 0; i < domElements.length; i++) {
      final domElement = domElements.item(i) as web.HTMLElement?;
      if (domElement == null) continue;

      // Check for class attribute first
      final classKey = domElement.getAttribute(_classAttribute);
      if (classKey != null && classKey.isNotEmpty) {
        final cmsData = _mappedClassesData[classKey];
        final classObj = Class.fromMap(this, classKey, domElement, cmsData);
        _classes[classKey] = classObj;
        // Don't continue - element may also have edit attribute
      }

      final key = domElement.getAttribute(_repeatAttribute);

      if (key != null && key.isNotEmpty) {
        final cmsData = _mappedRepeatsData[key];

        final repeat = Repeat.fromMap(this, key, domElement, cmsData);
        _repeats[key] = repeat;

        final repeatedElements = repeat.render();

        for (int ri = 0; ri < repeatedElements.length; ri++) {
          _processDomElement(repeatedElements[ri]);
        }

        continue;
      }

      _processDomElement(domElement);
    }
  }

  void _processDomElement(web.HTMLElement domElement) {
    final key = domElement.getAttribute(_editAttribute);

    if (key == null || key.isEmpty) return;

    final elementTag = domElement.tagName.toLowerCase();

    if (Image.SUPPORTED_IMAGE_TAGS.contains(elementTag)) {
      final cmsData = _mappedImagesData[key];

      final image = Image.fromMap(this, key, domElement, cmsData);
      _images[key] = image;

      image.render();
      return;
    }

    final cmsData = _mappedElementsData[key];

    final element = Element.fromMap(this, key, domElement, cmsData);
    _elements[key] = element;

    element.render();
  }

  void _bindEvents() {
    web.window.addEventListener('keydown', _windowKeyDown.toJS);
    web.window.addEventListener('keyup', _windowKeyUp.toJS);
  }

  void _windowKeyDown(web.KeyboardEvent e) {
    _ctrlPressed = e.ctrlKey;

    if (e.ctrlKey) {
      _elements.values.forEach((e) => e.highlight());
      _images.values.forEach((e) => e.highlight());
      _repeats.values.forEach((r) => r.highlight());
      _classes.values.forEach((c) => c.highlight());
    }
  }

  void _windowKeyUp(web.KeyboardEvent e) {
    _ctrlPressed = e.ctrlKey;

    _elements.values.forEach((e) => e.normalise());
    _images.values.forEach((e) => e.normalise());
    _repeats.values.forEach((r) => r.normalise());
    _classes.values.forEach((c) => c.normalise());
  }

  void save(Function onSuccess, Function onFailure) {
    if (isDemo) {
      web.window.postMessage("wedit.saved".toJS, "*".toJS);
      return;
    }

    final pageData = toMap();
    final jsonData = convert.json.encode(pageData);

    final url = "${web.window.location.protocol}//${web.window.location.host}/~?p=${web.window.location.pathname}";

    _sendRequest(url, "PUT", jsonData, onSuccess, onFailure);

    web.window.postMessage("wedit.saved".toJS, "*".toJS);
  }

  void command(String cmd, Function onSuccess, Function onFailure) {
    final pageData = toMap();
    final jsonData = convert.json.encode(pageData);

    final url = web.window.location.href.replaceAll("/!/", "/!$cmd/");

    _sendRequest(url, "POST", jsonData, onSuccess, onFailure);
  }

  void _sendRequest(String url, String method, String data, Function onSuccess, Function onFailure) {
    final headers = web.Headers();
    headers.set('Content-Type', 'application/json');

    final options = web.RequestInit(
      method: method,
      body: data.toJS,
      headers: headers,
    );

    web.window.fetch(url.toJS, options).toDart.then((response) {
      response.text().toDart.then((text) {
        print((text as JSString).toDart);
        onSuccess();
      }).catchError((_) {
        onFailure();
      });
    }).catchError((_) {
      onFailure();
    });
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

// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

uesc.HtmlUnescape htmlUnescape = uesc.HtmlUnescape();
convert.HtmlEscape htmlEscape = convert.HtmlEscape();

class Element {
  Page _page;

  String _key;

  String _text = ''; // in Markdown
  String get _html => _getHtml(_text); // in HTML

  String _getHtml(String text) {
    final result = md.markdownToHtml(text).toString();
    if (result.indexOf("<p>") == result.lastIndexOf("<p>")) {
      return result.replaceAll("<p>", "").replaceAll("</p>", "");
    }
    return result;
  }

  web.HTMLElement _domElement;

  String _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  String _DEFAULT_EDIT_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, 1), inset 0 0 2vw 0 rgba(255, 255, 255, 1)";

  static const _DEFAULT_EDIT_CURSOR = "pointer";

  String _originalBoxShadow = '';
  String _originalCursor = '';
  bool _isHighlighted = false;
  bool _isEditing = false;

  bool _wasEmpty = false;
  static const _NON_BREAKING_SPACE_HTML = "&nbsp;";

  String _originalPointerEvents = '';

  Element.fromMap(this._page, this._key, this._domElement, Map<String, dynamic>? map) {
    if (map != null) {
      _text = map[ELEMENT_TEXT] as String? ?? '';
      _text = _text.replaceAll("<br>", "\n").replaceAll("<q>", "\"");
    }

    _syncElement();

    _bindControls();

    if (_page.ctrlPressed) {
      highlight();
    }

    _wasEmpty = _domElement.textContent == "";

    if (_page.darkMode) {
      _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(255, 255, 255, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";
      _DEFAULT_EDIT_BOX_SHADOW = "0 0 2vw 0 rgba(255, 255, 255, 1), inset 0 0 2vw 0 rgba(0, 0, 0, 1)";
    }
  }

  Map<String, dynamic> toMap() {
    final map = <String, dynamic>{};

    map[ELEMENT_KEY] = _key;

    //_text = _domElement.textContent;
    map[ELEMENT_TEXT] = _text;

    return map;
  }

  void _syncElement() {
    _originalBoxShadow = _domElement.style.boxShadow;
    _originalCursor = _domElement.style.cursor;
    _originalPointerEvents = _domElement.style.pointerEvents;

    // Bind text
    if (_text.isEmpty) {
      _text = _domElement.textContent ?? '';
    }
  }

  void _bindControls() {
    _isHighlighted = false;
    _isEditing = false;

    _domElement.addEventListener('click', _elementClick.toJS);
    _domElement.addEventListener('contextmenu', _elementClick.toJS);
    _domElement.addEventListener('blur', _elementBlur.toJS);
  }

  void highlight() {
    _domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW;
    _domElement.style.cursor = _DEFAULT_EDIT_CURSOR;
    _domElement.style.pointerEvents = "all";

    if (_wasEmpty) {
      _domElement.innerHTML = _NON_BREAKING_SPACE_HTML.toJS;
    }

    _isHighlighted = true;
  }

  void normalise() {
    if (_isEditing) {
      return;
    }

    _domElement.style.boxShadow = _originalBoxShadow;
    _domElement.style.cursor = _originalCursor;
    _domElement.style.pointerEvents = _originalPointerEvents;

    if (_wasEmpty && _domElement.innerHTML == _NON_BREAKING_SPACE_HTML) {
      _domElement.textContent = "";
    }

    _isHighlighted = false;
  }

  void _elementClick(web.MouseEvent e) {
    if (!_isHighlighted) return;

    _domElement.style.boxShadow = _DEFAULT_EDIT_BOX_SHADOW;
    _domElement.style.cursor = _originalCursor;
    _domElement.style.pointerEvents = _originalPointerEvents;
    _domElement.contentEditable = "true";
    _domElement.focus();

    if (_isEditing) {
      return;
    }

    _domElement.innerHTML = htmlEscape.convert(_text).replaceAll("\n", "<br>").toJS;
    _isEditing = true;

    // Ensure that no links will be triggered
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();
  }

  void _elementBlur(web.Event e) {
    if (!_isEditing) return;
    _domElement.style.boxShadow = _originalBoxShadow;
    _domElement.style.cursor = _originalCursor;
    _domElement.style.pointerEvents = _originalPointerEvents;
    _domElement.contentEditable = "false";

    _wasEmpty = _domElement.textContent == "";

    _isHighlighted = _isEditing = false;

    _text = _htmlToMd((_domElement.innerHTML as JSString).toDart);
    _domElement.innerHTML = _html.toJS;

    _page.save(() {}, () {});
  }

  void render() {
    _domElement.innerHTML = _html.toJS;
  }

  void prepareDomForHtmlSave() {}

  void restoreDomAfterHtmlSave() {}

  String _htmlToMd(String raw) {
    raw = raw.replaceAll("</p>", "\n").replaceAll("<br>", "\n").replaceAll("<p>", "").replaceAll("</div>", "\n").replaceAll("<div>", "");

    while (raw.lastIndexOf("\n\n\n") != -1) {
      raw = raw.replaceAll("\n\n\n", "\n\n");
    }

    raw = htmlUnescape.convert(raw);

    return raw;
  }
}

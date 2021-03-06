// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

uesc.HtmlUnescape htmlUnescape = new uesc.HtmlUnescape();
convert.HtmlEscape htmlEscape = new convert.HtmlEscape();

class Element {
  Page _page;

  String _key;

  String _text; // in Markdown
  String get _html => _getHtml(_text); // in HTML

  String _getHtml(String text) {
    var result = md.markdownToHtml(text).toString();
    if (result.indexOf("<p>") == result.lastIndexOf("<p>")) {
      result = result.replaceAll("<p>", "").replaceAll("</p>", "");
    }
    return result;
  }

  html.Element _domElement;

  String _DEFAULT_MARK_BOX_SHADOW =
      "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  String _DEFAULT_EDIT_BOX_SHADOW =
      "0 0 2vw 0 rgba(0, 0, 0, 1), inset 0 0 2vw 0 rgba(255, 255, 255, 1)";

  static const _DEFAULT_EDIT_CURSOR = "pointer";

  String _originalBoxShadow;
  String _originalCursor;
  bool _isHighlighted;
  bool _isEditing;

  bool _wasEmpty;
  static const _NON_BREAKING_SPACE_HTML = "&nbsp;";

  String _originalPointerEvents;

  Element.fromMap(this._page, this._key, this._domElement, Map map) {
    if (map != null) {
      _text = map[ELEMENT_TEXT];
      _text = _text.replaceAll("<br>", "\n").replaceAll("<q>", "\"");
    }

    _syncElement();

    _bindControls();

    if (_page.ctrlPressed) {
      highlight();
    }

    _wasEmpty = _domElement.text == "";

    if(_page.darkMode) {
      _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(255, 255, 255, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";
      _DEFAULT_EDIT_BOX_SHADOW = "0 0 2vw 0 rgba(255, 255, 255, 1), inset 0 0 2vw 0 rgba(0, 0, 0, 1)";
    }
  }

  Map toMap() {
    var map = new Map();

    map[ELEMENT_KEY] = _key;

    //_text = _domElement.text;
    map[ELEMENT_TEXT] = _text;

    return map;
  }

  void _syncElement() {
    _originalBoxShadow = _domElement.style.boxShadow;
    _originalCursor = _domElement.style.cursor;
    _originalPointerEvents = _domElement.style.pointerEvents;

    // Bind text
    if (_text == null || _text.isEmpty) {
      _text = _domElement.text;
    }
  }

  void _bindControls() {
    _isHighlighted = false;
    _isEditing = false;

    _domElement.onClick.listen(_elementClick);
    _domElement.onContextMenu.listen(_elementClick);

    _domElement.onBlur.listen(_elementBlur);
  }

  void highlight() {
    _domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW;
    _domElement.style.cursor = _DEFAULT_EDIT_CURSOR;
    _domElement.style.pointerEvents = "all";

    if (_wasEmpty) {
      _domElement.innerHtml = _NON_BREAKING_SPACE_HTML;
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

    if (_wasEmpty && _domElement.innerHtml == _NON_BREAKING_SPACE_HTML) {
      _domElement.text = "";
    }

    _isHighlighted = false;
  }

  void _elementClick(html.MouseEvent e) {
    if (!_isHighlighted) return;

    _domElement.style.boxShadow = _DEFAULT_EDIT_BOX_SHADOW;
    _domElement.style.cursor = _originalCursor;
    _domElement.style.pointerEvents = _originalPointerEvents;
    _domElement.contentEditable = "true";
    _domElement.focus();

    if (_isEditing) {
      return;
    }

    _domElement.innerHtml = htmlEscape.convert(_text).replaceAll("\n", "<br>");
    _isEditing = true;

    // Ensure that no links will be triggered
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();
  }

  void _elementBlur(html.Event e) {
    if (!_isEditing) return;
    _domElement.style.boxShadow = _originalBoxShadow;
    _domElement.style.cursor = _originalCursor;
    _domElement.style.pointerEvents = _originalPointerEvents;
    _domElement.contentEditable = "false";

    _wasEmpty = _domElement.text == "";

    _isHighlighted = _isEditing = false;

    _text = _htmlToMd(_domElement.innerHtml);
    _domElement.setInnerHtml(_html, treeSanitizer: html.NodeTreeSanitizer.trusted);

    _page.save(null, null);
  }

  void render() {
    _domElement.setInnerHtml(_html, treeSanitizer: html.NodeTreeSanitizer.trusted);
  }

  void prepareDomForHtmlSave() {}

  void restoreDomAfterHtmlSave() {}

  String _htmlToMd(String raw) {

    raw = raw.replaceAll("</p>", "\n")
      .replaceAll("<br>", "\n")
      .replaceAll("<p>", "")
      .replaceAll("</div>", "\n")
      .replaceAll("<div>", "");

    while (raw.lastIndexOf("\n\n\n") != -1) {
      raw = raw.replaceAll("\n\n\n", "\n\n");
    }

    raw = htmlUnescape.convert(raw);

    return raw;
  }
}

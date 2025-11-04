// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Class {
  final Page _page;
  final String _key;
  final web.HTMLElement _domElement;

  String _value = '';
  List<String> _availableClasses = [];

  bool _isShown = false;
  web.HTMLSelectElement? _selectElement;

  String _originalBoxShadow = '';

  static const _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  static const _DEFAULT_MARK_BOX_SHADOW_DARK = "0 0 2vw 0 rgba(255, 255, 255, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";

  static const _SELECT_WIDTH = 120;
  static const _SELECT_HEIGHT = 20;

  Class.fromMap(this._page, this._key, this._domElement, Map<dynamic, dynamic>? map) {
    if (map != null) {
      _value = (map[CLASS_VALUE] as String?) ?? '';
    }

    // Parse available classes from the attribute value
    // Format: "key:class1,class2,class3"
    final classAttr = _domElement.getAttribute(_page._classAttribute);
    if (classAttr != null && classAttr.isNotEmpty) {
      // Check if it contains a colon (key:classes format)
      if (classAttr.contains(':')) {
        final parts = classAttr.split(':');
        if (parts.length >= 2) {
          // Use part after colon for available classes
          _availableClasses = parts[1].split(',').map((c) => c.trim()).toList();
        }
      } else {
        // No colon, treat entire value as single class option
        _availableClasses = [classAttr.trim()];
      }
    }

    _originalBoxShadow = _domElement.style.boxShadow;
    _createSelect();
    _syncElement();
  }

  void _syncElement() {
    // Apply the stored class value to the element
    if (_value.isNotEmpty) {
      _domElement.className = _value;
    }
  }

  void render() {
    // Apply the stored class value to the element
    if (_value.isNotEmpty) {
      _domElement.className = _value;
    }
  }

  void _createSelect() {
    _selectElement = web.document.createElement('select') as web.HTMLSelectElement;
    _selectElement!.style
      ..display = "none"
      ..position = "absolute"
      ..width = "${_SELECT_WIDTH}px"
      ..height = "${_SELECT_HEIGHT}px"
      ..backgroundColor = _page.darkMode ? "#333" : "#fff"
      ..color = _page.darkMode ? "#fff" : "#000"
      ..border = "1px solid " + (_page.darkMode ? "#555" : "#ccc")
      ..borderRadius = "4px"
      ..padding = "8px"
      ..fontSize = "14px"
      ..zIndex = "10000"
      ..boxShadow = "0 2px 8px rgba(0,0,0,0.3)"
      ..cursor = "pointer";

    // Add class options
    for (var className in _availableClasses) {
      final optionElement = web.document.createElement('option') as web.HTMLOptionElement;
      optionElement.value = className;
      optionElement.textContent = className;

      // Select current value
      if (className == _value) {
        optionElement.selected = true;
      }

      _selectElement!.appendChild(optionElement);
    }

    // Change handler
    _selectElement!.onChange.listen((event) {
      final selectedValue = _selectElement!.value;
      _selectClass(selectedValue);
    });

    web.document.body!.appendChild(_selectElement!);
  }

  void _selectClass(String className) {
    _value = className;

    // Update the actual class attribute on the element
    _domElement.className = className;

    // Auto-save
    _page.save(() {}, () {});

    // Hide select
    normalise();
  }

  void highlight() {
    if (_isShown) return;
    _isShown = true;

    _domElement.style.boxShadow = _page.darkMode ? _DEFAULT_MARK_BOX_SHADOW_DARK : _DEFAULT_MARK_BOX_SHADOW;
    _domElement.style.cursor = "pointer";

    _showSelect();
  }

  void normalise() {
    if (!_isShown) return;
    _isShown = false;

    _domElement.style.boxShadow = _originalBoxShadow;
    _domElement.style.cursor = "";

    _hideSelect();
  }

  int _getOffsetTop(web.HTMLElement? el) {
    int res = 0;

    while (el != null) {
      res += el.offsetTop.round();
      el = el.offsetParent as web.HTMLElement?;
    }

    return res;
  }

  int _getOffsetLeft(web.HTMLElement? el) {
    int res = 0;

    while (el != null) {
      res += el.offsetLeft.round();
      el = el.offsetParent as web.HTMLElement?;
    }

    return res;
  }

  void _showSelect() {
    if (_selectElement == null || _availableClasses.isEmpty) return;

    // Position the select near the element
    final rect = _domElement.getBoundingClientRect();
    _selectElement!.style
      ..display = "block"
      ..top = "${rect.bottom + 5}px"
      ..left = "${rect.left}px";

    final selectStyle = _selectElement?.style;
    if (selectStyle != null) {
      selectStyle.left = (_getOffsetLeft(_domElement) + _domElement.offsetWidth.round() - _SELECT_WIDTH).toString() + "px";
      selectStyle.top = (_getOffsetTop(_domElement) - _SELECT_HEIGHT).toString() + "px";
      selectStyle.display = "block";
    }
  }

  void _hideSelect() {
    if (_selectElement == null) return;
    _selectElement!.style.display = "none";
  }

  Map<String, dynamic> toMap() {
    return {CLASS_KEY: _key, CLASS_VALUE: _value};
  }
}

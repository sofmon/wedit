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
  web.HTMLElement? _dropdownElement;

  String _originalBoxShadow = '';

  static const _DEFAULT_MARK_BOX_SHADOW =
      "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  static const _DEFAULT_MARK_BOX_SHADOW_DARK =
      "0 0 2vw 0 rgba(255, 255, 255, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";

  Class.fromMap(this._page, this._key, this._domElement, Map<dynamic, dynamic>? map) {
    if (map != null) {
      _value = map[CLASS_VALUE] ?? '';
    }

    // Parse available classes from the attribute value (comma-separated)
    final classAttr = _domElement.getAttribute(_page._classAttribute);
    if (classAttr != null && classAttr.isNotEmpty) {
      _availableClasses = classAttr.split(',').map((c) => c.trim()).toList();
    }

    _originalBoxShadow = _domElement.style.boxShadow;
    _createDropdown();
  }

  void _createDropdown() {
    _dropdownElement = web.document.createElement('div') as web.HTMLElement;
    _dropdownElement!.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _page.darkMode ? "#333" : "#fff"
      ..border = "1px solid " + (_page.darkMode ? "#555" : "#ccc")
      ..borderRadius = "4px"
      ..padding = "8px"
      ..zIndex = "10000"
      ..boxShadow = "0 2px 8px rgba(0,0,0,0.3)"
      ..maxHeight = "200px"
      ..overflowY = "auto"
      ..minWidth = "150px";

    // Add class options as list items
    for (var className in _availableClasses) {
      final optionElement = web.document.createElement('div') as web.HTMLElement;
      optionElement.style
        ..padding = "6px 12px"
        ..cursor = "pointer"
        ..color = _page.darkMode ? "#fff" : "#000"
        ..borderBottom = "1px solid " + (_page.darkMode ? "#444" : "#eee");

      optionElement.textContent = className;

      // Highlight current selection
      if (className == _value) {
        optionElement.style
          ..backgroundColor = _page.darkMode ? "#444" : "#f0f0f0"
          ..fontWeight = "bold";
      }

      // Hover effects
      optionElement.onMouseEnter.listen((_) {
        if (className != _value) {
          optionElement.style.backgroundColor = _page.darkMode ? "#555" : "#e8e8e8";
        }
      });

      optionElement.onMouseLeave.listen((_) {
        if (className != _value) {
          optionElement.style.backgroundColor = "transparent";
        }
      });

      // Click handler
      optionElement.onClick.listen((event) {
        event.preventDefault();
        event.stopPropagation();
        _selectClass(className);
      });

      _dropdownElement!.appendChild(optionElement);
    }

    web.document.body!.appendChild(_dropdownElement!);
  }

  void _selectClass(String className) {
    _value = className;

    // Update the actual class attribute on the element
    _domElement.className = className;

    // Auto-save
    _page.save(() {}, () {});

    // Hide dropdown
    normalise();
  }

  void highlight() {
    if (_isShown) return;
    _isShown = true;

    _domElement.style.boxShadow = _page.darkMode
        ? _DEFAULT_MARK_BOX_SHADOW_DARK
        : _DEFAULT_MARK_BOX_SHADOW;
    _domElement.style.cursor = "pointer";

    _showDropdown();
  }

  void normalise() {
    if (!_isShown) return;
    _isShown = false;

    _domElement.style.boxShadow = _originalBoxShadow;
    _domElement.style.cursor = "";

    _hideDropdown();
  }

  void _showDropdown() {
    if (_dropdownElement == null || _availableClasses.isEmpty) return;

    // Position the dropdown near the element
    final rect = _domElement.getBoundingClientRect();
    _dropdownElement!.style
      ..display = "block"
      ..top = "${rect.bottom + 5}px"
      ..left = "${rect.left}px";
  }

  void _hideDropdown() {
    if (_dropdownElement == null) return;
    _dropdownElement!.style.display = "none";
  }

  Map<String, dynamic> toMap() {
    return {
      CLASS_KEY: _key,
      CLASS_VALUE: _value,
    };
  }
}

// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Repeat {
  Page _page;

  String _key;
  String get key => _key;

  web.HTMLElement _domElement;

  web.HTMLElement? _domTemplate;

  Map<String, RepeatShadow> _shadows = {};
  List<String> _keyOrder = [];

  Repeat.fromMap(this._page, this._key, this._domElement, Map<dynamic, dynamic>? map) {
    _bindControls();

    _syncElements();

    _shadows = <String, RepeatShadow>{};
    _shadows[_key] = RepeatShadow.fromRepeat(this, _key, _domElement);

    if (map != null) {
      print(map);
      print(map[REPEAT_COPY_KEYS]);
      _keyOrder = (map[REPEAT_COPY_KEYS] as List?)?.cast<String>() ?? [];
      print(_keyOrder);
      _renderShadows(_keyOrder);
    } else {
      _keyOrder = <String>[];
      _keyOrder.add(_key);
    }
  }

  Map<String, dynamic> toMap() {
    final map = <String, dynamic>{};

    map[REPEAT_KEY] = _key;
    map[REPEAT_COPY_KEYS] = _keyOrder;

    return map;
  }

  void _bindControls() {
    // nothing here
  }

  void _syncElements() {
    _domTemplate = _domElement.cloneNode(true) as web.HTMLElement;
    _domTemplate?.removeAttribute(_page.repeatAttribute);
  }

  void _renderShadows(List<String>? copyKeys) {
    if (copyKeys == null) {
      return;
    }

    final before = <String>[];
    final after = <String>[];
    bool beforeTemplate = true;
    for (int i = 0; i < copyKeys.length; i++) {
      final key = copyKeys[i];
      if (key == _key) {
        beforeTemplate = false;
        continue;
      }
      if(beforeTemplate){
        before.add(copyKeys[i]);
      } else {
        after.add(copyKeys[i]);
      }
    }

    for (int i = 0; i < before.length; i++) {
      final key = before[i];
      final copyElement = _createCopyDomElement(key);
      _domElement.insertAdjacentElement('beforebegin', copyElement);
      _shadows[key] = RepeatShadow.fromRepeat(this, key, copyElement);
    }

    for(int i=after.length-1; i>=0; i--) {
      final key = after[i];
      final copyElement = _createCopyDomElement(key);
      _domElement.insertAdjacentElement('afterend', copyElement);
      _shadows[key] = RepeatShadow.fromRepeat(this, key, copyElement);
    }
  }

  void highlight() {
    _shadows.values.forEach((rs) => rs.show());
  }

  void normalise() {
    _shadows.values.forEach((rs) => rs.hide());
  }

  void addCopy(String copyKey, web.HTMLElement targetDomElement) {
    final key = DateTime.now().millisecondsSinceEpoch.toString();

    final copyElement = _createCopyDomElement(key);
    targetDomElement.insertAdjacentElement('afterend', copyElement);

    _shadows[key] = RepeatShadow.fromRepeat(this, key, copyElement);
    _keyOrder.insert(_keyOrder.indexOf(copyKey) + 1, key);

    if (_page.ctrlPressed) {
      _shadows.values.forEach((rs) => rs.show());
    }
  }

  void removeCopy(String copyKey, web.HTMLElement targetDomElement) {
    if (copyKey == _key) {
      return;
    }

    final elements = targetDomElement.querySelectorAll("[${_page.editAttribute}]");

    for (int i = 0; i < elements.length; i++) {
      _page.unregisterElement(elements.item(i) as web.HTMLElement);
    }

    targetDomElement.remove();

    _shadows.remove(copyKey);
    _keyOrder.remove(copyKey);

    // Fix the position of the other repeat elements shadows
    _shadows.values.forEach((s) => s.show());
  }

  void moveCopyUp(String copyKey) {
    final orderIndex = _keyOrder.indexOf(copyKey);
    if (orderIndex == 0) {
      return;
    }

    _keyOrder.remove(copyKey);
    _keyOrder.insert(orderIndex - 1, copyKey);

    final copyDomElement = _shadows[copyKey]?.domElement;
    final previousSibling = copyDomElement?.previousElementSibling;
    if (previousSibling == null || copyDomElement == null) {
      return;
    }
    copyDomElement.remove();
    previousSibling.insertAdjacentElement('beforebegin', copyDomElement);

    // Fix the position of the other repeat elements shadows
    _shadows.values.forEach((s) => s.show());
  }

  void moveCopyDown(String copyKey) {
    final orderIndex = _keyOrder.indexOf(copyKey);
    if (orderIndex >= _keyOrder.length - 1) {
      return;
    }

    _keyOrder.remove(copyKey);
    _keyOrder.insert(orderIndex + 1, copyKey);

    final copyDomElement = _shadows[copyKey]?.domElement;
    final nextSibling = copyDomElement?.nextElementSibling;
    if (nextSibling == null || copyDomElement == null) {
      return;
    }
    copyDomElement.remove();
    nextSibling.insertAdjacentElement('afterend', copyDomElement);

    // Fix the position of the other repeat elements shadows
    _shadows.values.forEach((s) => s.show());
  }

  web.HTMLElement _createCopyDomElement(String key) {
    final copyElement = _domTemplate!.cloneNode(true) as web.HTMLElement;

    copyElement.removeAttribute(_page.repeatAttribute);
    final elements = copyElement.querySelectorAll("[${_page.editAttribute}]");

    for (int i = 0; i < elements.length; i++) {
      final element = elements.item(i) as web.HTMLElement;
      final eKey = (element.getAttribute(_page.editAttribute) ?? '') + key;
      element.removeAttribute(_page.editAttribute);
      element.setAttribute(_page.editAttribute, eKey);
      _page.registerElement(element);
    }

    return copyElement;
  }

  List<web.HTMLElement> render() {
    return <web.HTMLElement>[];
  }

  void prepareDomForHtmlSave() {
    _shadows.values.forEach((s) => s.prepareDomForHtmlSave());
  }

  void restoreDomAfterHtmlSave() {
    _shadows.values.forEach((s) => s.restoreDomAfterHtmlSave());
  }
}

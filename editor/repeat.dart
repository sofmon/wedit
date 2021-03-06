// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Repeat {
  Page _page;

  String _key;
  String get key => _key;

  html.Element _domElement;

  html.Element _domTemplate;

  Map<String, RepeatShadow> _shadows;
  List<String> _keyOrder;

  Repeat.fromMap(this._page, this._key, this._domElement, Map map) {
    _bindControls();

    _syncElements();

    _shadows = new Map<String, RepeatShadow>();
    _shadows[_key] = new RepeatShadow.fromRepeat(this, _key, _domElement);

    if (map != null) {
      print(map);
      print(map[REPEAT_COPY_KEYS]);
      _keyOrder = map[REPEAT_COPY_KEYS].cast<String>();
      print(_keyOrder);
      _renderShadows(_keyOrder);
    } else {
      _keyOrder = new List<String>();
      _keyOrder.add(_key);
    }
  }

  Map toMap() {
    Map map = new Map();

    map[REPEAT_KEY] = _key;
    map[REPEAT_COPY_KEYS] = _keyOrder;

    return map;
  }

  void _bindControls() {
    // nothing here
  }

  void _syncElements() {
    _domTemplate = _domElement.clone(true);
    _domTemplate.attributes.remove(_page.repeatAttribute);
  }

  void _renderShadows(List<String> copyKeys) {
    if (copyKeys == null) {
      return;
    }

    List<String> before = new List<String>();
    List<String> after = new List<String>();
    bool beforeTemplate = true;
    for (int i = 0; i < copyKeys.length; i++) {
      String key = copyKeys[i];
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
      String key = before[i];
      html.Element copyElement = _createCopyDomElement(key);
      _domElement.insertAdjacentElement("beforeBegin", copyElement);
      _shadows[key] = new RepeatShadow.fromRepeat(this, key, copyElement);
    }

    for(int i=after.length-1; i>=0; i--) {
      String key = after[i];
      html.Element copyElement = _createCopyDomElement(key);
      _domElement.insertAdjacentElement("afterEnd", copyElement);
      _shadows[key] = new RepeatShadow.fromRepeat(this, key, copyElement);
    }
  }

  void highlight() {
    _shadows.values.forEach((rs) => rs.show());
  }

  void normalise() {
    _shadows.values.forEach((rs) => rs.hide());
  }

  void addCopy(String copyKey, html.Element targetDomElement) {
    String key = new DateTime.now().millisecondsSinceEpoch.toString();

    html.Element copyElement = _createCopyDomElement(key);
    targetDomElement.insertAdjacentElement("afterEnd", copyElement);

    _shadows[key] = new RepeatShadow.fromRepeat(this, key, copyElement);
    _keyOrder.insert(_keyOrder.indexOf(copyKey) + 1, key);

    if (_page.ctrlPressed) {
      _shadows.values.forEach((rs) => rs.show());
    }
  }

  void removeCopy(String copyKey, html.Element targetDomElement) {
    if (copyKey == _key) {
      return;
    }

    html.ElementList<html.Element> elements =
        targetDomElement.querySelectorAll("[" + _page.editAttribute + "]");

    for (int i = 0; i < elements.length; i++) {
      _page.unregisterElement(elements[i]);
    }

    targetDomElement.remove();

    _shadows.remove(copyKey);
    _keyOrder.remove(copyKey);

    // Fix the position of the other repeat elements shadows
    _shadows.values.forEach((s) => s.show());
  }

  void moveCopyUp(String copyKey) {
    int orderIndex = _keyOrder.indexOf(copyKey);
    if (orderIndex == 0) {
      return;
    }

    _keyOrder.remove(copyKey);
    _keyOrder.insert(orderIndex - 1, copyKey);

    html.Element copyDomElement = _shadows[copyKey].domElement;
    html.Element previousSibling = copyDomElement.previousElementSibling;
    if (previousSibling == null) {
      return;
    }
    copyDomElement.remove();
    previousSibling.insertAdjacentElement("beforeBegin", copyDomElement);

    // Fix the position of the other repeat elements shadows
    _shadows.values.forEach((s) => s.show());
  }

  void moveCopyDown(String copyKey) {
    int orderIndex = _keyOrder.indexOf(copyKey);
    if (orderIndex >= _keyOrder.length - 1) {
      return;
    }

    _keyOrder.remove(copyKey);
    _keyOrder.insert(orderIndex + 1, copyKey);

    html.Element copyDomElement = _shadows[copyKey].domElement;
    html.Element nextSibling = copyDomElement.nextElementSibling;
    if (nextSibling == null) {
      return;
    }
    copyDomElement.remove();
    nextSibling.insertAdjacentElement("afterEnd", copyDomElement);

    // Fix the position of the other repeat elements shadows
    _shadows.values.forEach((s) => s.show());
  }

  html.Element _createCopyDomElement(String key) {
    html.Element copyElement = _domTemplate.clone(true);

    copyElement.attributes.remove(_page.repeatAttribute);
    html.ElementList<html.Element> elements =
        copyElement.querySelectorAll("[" + _page.editAttribute + "]");

    for (int i = 0; i < elements.length; i++) {
      String eKey = elements[i].getAttribute(_page.editAttribute) + key;
      elements[i].attributes.remove(_page.editAttribute);
      elements[i].setAttribute(_page.editAttribute, eKey);
      _page.registerElement(elements[i]);
    }

    return copyElement;
  }

  List<html.Element> render() {
    return new List<html.Element>();
  }

  void prepareDomForHtmlSave() {
    _shadows.values.forEach((s) => s.prepareDomForHtmlSave());
  }

  void restoreDomAfterHtmlSave() {
    _shadows.values.forEach((s) => s.restoreDomAfterHtmlSave());
  }
}

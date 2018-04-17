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
			String keysString = map[REPEAT_COPY_KEYS];
			_keyOrder = keysString.split(",");
			_renderShadows(_keyOrder);
		} else {
			_keyOrder = new List<String>();
			_keyOrder.add(_key);
		}
	}

	Map toMap() {
		Map map = new Map();

		map[REPEAT_KEY] = _key;
		map[REPEAT_COPY_KEYS] = _keyOrder.join(",");

		return map;
	}

	void _bindControls() {
		// nothing here
	}

	void _syncElements() {
		_domTemplate = _domElement.clone(true);
		_domTemplate.dataset.remove("var-repeat");
	}

	void _renderShadows(List<String> copyKeys) {

		if(copyKeys == null) {
			return;
		}

		bool beforeTemplate = true;

		for(int i=0; i<copyKeys.length; i++) {

			String key = copyKeys[i];

			if(key == _key) {
				beforeTemplate = false;
				continue;
			}

			html.Element copyElement = _createCopyDomElement(key);
			_domElement.insertAdjacentElement(beforeTemplate ? "beforeBegin" : "afterEnd", copyElement);

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
		_keyOrder.insert(_keyOrder.indexOf(copyKey)+1, key);

		if(_page.ctrlPressed) {
			_shadows.values.forEach((rs) => rs.show());
		}
	}

	void removeCopy(String copyKey, html.Element targetDomElement) {

		if(copyKey == _key) {
			return;
		}

		html.ElementList<html.Element> elements = targetDomElement.querySelectorAll("[data-var]");

		for(int i=0; i<elements.length; i++) {

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
		if(orderIndex == 0) {
			return;
		}

		_keyOrder.remove(copyKey);
		_keyOrder.insert(orderIndex-1, copyKey);

		html.Element copyDomElement = _shadows[copyKey].domElement;
		html.Element previousSibling = copyDomElement.previousElementSibling;
		if(previousSibling == null) {
			return;
		}
		copyDomElement.remove();
		previousSibling.insertAdjacentElement("beforeBegin", copyDomElement);

		// Fix the position of the other repeat elements shadows
		_shadows.values.forEach((s) => s.show());
	}

	void moveCopyDown(String copyKey) {

		int orderIndex = _keyOrder.indexOf(copyKey);
		if(orderIndex >= _keyOrder.length-1) {
			return;
		}

		_keyOrder.remove(copyKey);
		_keyOrder.insert(orderIndex+1, copyKey);

		html.Element copyDomElement = _shadows[copyKey].domElement;
		html.Element nextSibling = copyDomElement.nextElementSibling;
		if(nextSibling == null) {
			return;
		}
		copyDomElement.remove();
		nextSibling.insertAdjacentElement("afterEnd", copyDomElement);

		// Fix the position of the other repeat elements shadows
		_shadows.values.forEach((s) => s.show());
	}

	html.Element _createCopyDomElement(String key) {

		html.Element copyElement = _domTemplate.clone(true);

		copyElement.dataset.remove("var-repeat");
		html.ElementList<html.Element> elements = copyElement.querySelectorAll("[data-var]");

		for(int i=0;i<elements.length;i++) {
			String eKey = elements[i].dataset["var"] + key;
			elements[i].dataset.remove("var");
			elements[i].dataset["var"] = eKey;
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
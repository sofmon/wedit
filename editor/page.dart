// Copyright (c) 2015, Haralampi Staykov. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
part of wedit;

class Page {

	String _host, _site, _path;

	String _title;
	String get title => _title;
	set title(String t) => _title = t;

	Map<String, Map> _mappedElementsData;
	Map<String, Map> _mappedRepeatsData;

	Map<String, Element> _elements;
	Map<String, Image> _images;
	Map<String, Repeat> _repeats;

	bool _userIsEditor = false;
	bool _userIsAdministrator = false;

	static const _ALLOW_EDIT_HASH = "#var#";
	bool _allowEditMode;
	bool _inEditMode;
	bool get inEditMode => _inEditMode;

	bool _ctrlPressed = false;
	bool get ctrlPressed => _ctrlPressed;

	html.IFrameElement _authFrame;
	static const _DEFAULT_AUTH_FRAME_SHADOW = "0 0 4vw 0 rgba(0, 1, 20, 1)";
	async.Timer _authTimer;

	//PageMenu _pageMenu;
	Message _message;

	static const _TRY_MODE_SITE_NAME = "__try__";
	bool _inTryMode;

	bool get inTryMode => _inTryMode;

	Page.fromMap(Map map) {
		_host = map[PAGE_HOST];
		_site = map[PAGE_SITE];
		_path = map[PAGE_PATH];

		_title = map[PAGE_TITLE];

		_inTryMode = _site == _TRY_MODE_SITE_NAME;
		_inEditMode = _inTryMode;

		_allowEditMode = html.window.location.hash == _ALLOW_EDIT_HASH;

		_elements = new Map<String, Element>();
		_images = new Map<String, Image>();
		_repeats = new Map<String, Repeat>();

		_mappedElementsData = new Map<String, Map>();
		map[PAGE_ELEMENTS].forEach((e) =>
		_mappedElementsData[e[ELEMENT_KEY]] = e);

		_mappedRepeatsData = new Map<String, Map>();
		map[PAGE_REPEATS].forEach((r) =>
		_mappedRepeatsData[r[REPEAT_KEY]] = r);

		_syncElements();
		_bindEvents();

		html.window.dispatchEvent(new html.Event("varReady"));

		_message = new Message();
		//_pageMenu = new PageMenu(this);
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

		return map;
	}

	void registerElement(html.Element domElement) {
		var key = domElement.dataset["var"];

		if (key == null || key.isEmpty) return;

		Map cmsData = _mappedElementsData[key];

		Element element = new Element.fromMap(
			this, key, domElement, cmsData);

		_elements[key] = element;

		element.render();
	}

	void unregisterElement(html.Element domElement) {
		var key = domElement.dataset["var"];
		_elements.remove(key);
	}

	void _syncElements() {
		html.document.title = _title;

		List<html.Element> repeatedDomElements = new List<
			html.Element>();

		html.ElementList<html.Element> domElements = html
			.querySelectorAll("[data-var],[data-var-repeat]");

		for (int i = 0; i < domElements.length; i++) {
			html.Element domElement = domElements[i];

			String key = domElement.dataset["var-repeat"];

			if (key != null && !key.isEmpty) {
				Map cmsData = _mappedRepeatsData[key];

				Repeat repeat = new Repeat.fromMap(
					this, key, domElement, cmsData);
				_repeats[key] = repeat;

				List<html.Element> repeatedElements = repeat
					.render();

				for (int ri = 0; ri < repeatedElements.length;
				ri++) {
					_processDomElement(
						repeatedElements[ri]);
				}

				continue;
			}

			_processDomElement(domElement);

		}
	}

	void _processDomElement(html.Element domElement) {
		String key = domElement.dataset["var"];

		if (key == null || key.isEmpty) return;

		Map cmsData = _mappedElementsData[key];

		String elementTag = domElement.tagName.toLowerCase();

		if(Image.SUPPORTED_IMAGE_TAGS.contains(elementTag)) {

			Image image = new Image.fromMap(
				this, key, domElement, cmsData
			);
			_images[key] = image;

			image.render();
			return;
		}

		Element element = new Element.fromMap(
			this, key, domElement, cmsData);
		_elements[key] = element;

		element.render();
	}

	void _bindEvents() {
		html.window.onKeyDown.listen(_windowKeyDown);
		html.window.onKeyUp.listen(_windowKeyUp);
		html.window.onHashChange.listen(_windowHashChange);
		//window.onKeyUp.listen(_windowKeyUp);

	}

	void _windowKeyDown(html.KeyboardEvent e) {

		if ( _allowEditMode && !_inEditMode && e.ctrlKey && e.keyCode == html.KeyCode.E ) {
			_enterEditMode();
		}

		if (e.ctrlKey && _inEditMode && e.keyCode == html.KeyCode.S) {
			save();
		}

		_ctrlPressed = e.ctrlKey;

		if ( _inEditMode && e.ctrlKey ) {
			_elements.values.forEach((e)=>e.highlight());
			_images.values.forEach((e)=>e.highlight());
			_repeats.values.forEach((r)=>r.highlight());
		}

	}

	void _windowKeyUp(html.KeyboardEvent e) {
		_ctrlPressed = e.ctrlKey;

		if ( _inEditMode ) {
			_elements.values.forEach((e)=>e.normalise());
			_images.values.forEach((e)=>e.normalise());
			_repeats.values.forEach((r)=>r.normalise());
		}

	}

	void _windowHashChange(html.Event e) {
		_allowEditMode = html.window.location.hash == _ALLOW_EDIT_HASH;
	}

	void _enterEditMode() {
		//_inEditMode = true;

		var url = html.window.location.protocol + "://" + _host + "/" +
			_site + "/auth-request";

		// Create iframe
		_authFrame = new html.IFrameElement();
		_authFrame.src = url;
		_authFrame.style
			..background = "#fff"
			..border = "none"
			..boxShadow = _DEFAULT_AUTH_FRAME_SHADOW
			..width = "80vw"
			..height = "80vh"
			..position = "fixed"
			..left = "50%"
			..top = "50%"
			..transform = "translateX(-50%) translateY(-50%)";

		html.document.body.append(_authFrame);

		_authTimer = new async.Timer.periodic(
			new Duration(seconds: 3), _checkForPermissions);
	}

	void _checkForPermissions(async.Timer t) {
		var url = html.window.location.protocol + "://" + _host + "/" +
			_site + "/auth-check";

		html.HttpRequest.getString(url)
			.then((String status) {
			Map auth = convert.JSON.decode(status);
			_userIsEditor = auth["E"];
			_userIsAdministrator = auth["A"];
			if (_userIsEditor || _userIsAdministrator) {
				_authFrame.remove();
				_authTimer.cancel();
				_inEditMode = true;
				_message.ShowText("EDIT MODE");
			}
		})
			.catchError((Error error) {
			print(error);
		});
	}

	void save() {
		if (_inTryMode) {
			return;
		}

		Map map = toMap();

		html.HttpRequest request = new html
			.HttpRequest(); // create a new XHR

		// add an event handler that is called when the request finishes
		request.onReadyStateChange.listen((_) {
			if (request.readyState == html.HttpRequest.DONE &&
				(request.status == 200 ||
					request.status == 0)) {
				// data saved OK.
				print(request.responseText);
				// output the response from the server
			}
		});

		var url = html.window.location.protocol + "://" + _host + "/" +
			_site + "/save-page";

		request.open("POST", url, async: false);
		String jsonData = convert.JSON.encode(map);
		request.onLoad.listen((e)=>_message.ShowText("SAVED"));
		request.send(jsonData);
	}

}
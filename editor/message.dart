// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Message {

	static const _DEFAULT_BACKGROUND_COLOR = "#aaa";
	static const _DEFAULT_TEXT_COLOR = "#000";
	static const _DEFAULT_BORDER_RADIUS = "1vw";
	static const _DEFAULT_PADDING = "1vw";
	static const _DEFAULT_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, 1), inset 0 0 2vw 0 rgba(255, 255, 255, 1)";

	html.Element _domElement;

	Message() {
		_domElement = new html.Element.div();
		_domElement.style
			..position = "fixed"
			..top = "50%"
			..left = "50%"
			..transform = "translateX(-50%) translateY(-50%)"
			..backgroundColor = _DEFAULT_BACKGROUND_COLOR
			..color = _DEFAULT_TEXT_COLOR
			..borderRadius = _DEFAULT_BORDER_RADIUS
			..padding = _DEFAULT_PADDING
			..boxShadow = _DEFAULT_BOX_SHADOW
			..opacity = "0";

		html.document.body.children.add(_domElement);

	}

	void ShowText(String text) {

		_domElement.text = text;

		_domElement.animate( [{"opacity": "0"}, {"opacity": "1"}, {"opacity": "0"}], 1000);
	}

  void prepareDomForHtmlSave() {
    	_domElement.remove();
  }

  void restoreDomAfterHtmlSave() {
    html.document.body.children.add(_domElement);
  }

}
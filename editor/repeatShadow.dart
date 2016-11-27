// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class RepeatShadow {

	Repeat _repeat;

	html.Element _domElement;
	html.Element get domElement => _domElement;

	bool _isShown;

	String _key;
	bool _canBeDeleted;

	RepeatShadow.fromRepeat(this._repeat, this._key, this._domElement) {

		_isShown = false;

		_canBeDeleted = _key != _repeat.key;

		_bindControls();
	}

	html.Element _addDomElement;
	html.Element _removeDomElement;
	html.Element _moveUpDomElement;
	html.Element _moveDownDomElement;

	static const _REPEAT_BUTTON_SIZE = 20;

	String _originalBoxShadow;

	static const _DEFAULT_BUTTON_COLOR_ADD = "#0a0";
	static const _DEFAULT_BUTTON_COLOR_REMOVE = "#f00";
	static const _DEFAULT_BUTTON_COLOR_MOVE_UP = "#06f";
	static const _DEFAULT_BUTTON_COLOR_MOVE_DOWN = "#00f";
	static const _DEFAULT_BUTTON_OPACITY = ".3";


	static const _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
	static const _DEFAULT_MARK_BOX_SHADOW_ADD = "0 0 2vw 0 rgba(0, 155, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
	static const _DEFAULT_MARK_BOX_SHADOW_REMOVE = "0 0 2vw 0 rgba(255, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
	static const _DEFAULT_MARK_BOX_SHADOW_MOVE_UP = "0 0 2vw 0 rgba(0, 50, 255, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
	static const _DEFAULT_MARK_BOX_SHADOW_MOVE_DOWN = "0 0 2vw 0 rgba(0, 50, 255, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";

	static const SVG_UP = '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M16 1l-15 15h9v16h12v-16h9z"></path></svg>';
	static const SVG_DOWN = '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"> <path d="M16 31l15-15h-9v-16h-12v16h-9z"></path></svg>';
	static const SVG_PLUS = '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M31 12h-11v-11c0-0.552-0.448-1-1-1h-6c-0.552 0-1 0.448-1 1v11h-11c-0.552 0-1 0.448-1 1v6c0 0.552 0.448 1 1 1h11v11c0 0.552 0.448 1 1 1h6c0.552 0 1-0.448 1-1v-11h11c0.552 0 1-0.448 1-1v-6c0-0.552-0.448-1-1-1z"></path></svg>';
	static const SVG_REMOVE = '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M31.708 25.708c-0-0-0-0-0-0l-9.708-9.708 9.708-9.708c0-0 0-0 0-0 0.105-0.105 0.18-0.227 0.229-0.357 0.133-0.356 0.057-0.771-0.229-1.057l-4.586-4.586c-0.286-0.286-0.702-0.361-1.057-0.229-0.13 0.048-0.252 0.124-0.357 0.228 0 0-0 0-0 0l-9.708 9.708-9.708-9.708c-0-0-0-0-0-0-0.105-0.104-0.227-0.18-0.357-0.228-0.356-0.133-0.771-0.057-1.057 0.229l-4.586 4.586c-0.286 0.286-0.361 0.702-0.229 1.057 0.049 0.13 0.124 0.252 0.229 0.357 0 0 0 0 0 0l9.708 9.708-9.708 9.708c-0 0-0 0-0 0-0.104 0.105-0.18 0.227-0.229 0.357-0.133 0.355-0.057 0.771 0.229 1.057l4.586 4.586c0.286 0.286 0.702 0.361 1.057 0.229 0.13-0.049 0.252-0.124 0.357-0.229 0-0 0-0 0-0l9.708-9.708 9.708 9.708c0 0 0 0 0 0 0.105 0.105 0.227 0.18 0.357 0.229 0.356 0.133 0.771 0.057 1.057-0.229l4.586-4.586c0.286-0.286 0.362-0.702 0.229-1.057-0.049-0.13-0.124-0.252-0.229-0.357z"></path></svg>';


	void _mark() {
		_domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW;
	}

	void _markForAdd() {
		_domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW_ADD;
	}

	void _markForRemove() {
		_domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW_REMOVE;
	}

	void _markForMoveUp() {
		_domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW_MOVE_UP;
	}

	void _markForMoveDown() {
		_domElement.style.boxShadow =
			_DEFAULT_MARK_BOX_SHADOW_MOVE_DOWN;
	}

	void _clearMark() {
		_domElement.style.boxShadow = _isShown ? _DEFAULT_MARK_BOX_SHADOW : _originalBoxShadow;
	}

	void _bindControls() {
		_originalBoxShadow = _domElement.style.boxShadow;

		_addDomElement = new html.Element.div();
		_addDomElement.style
			..display = "none"
			..position = "absolute"
			..backgroundColor = _DEFAULT_BUTTON_COLOR_ADD
			..width = _REPEAT_BUTTON_SIZE.toString() + "px"
			..height = _REPEAT_BUTTON_SIZE.toString() + "px"
			..borderRadius = (_REPEAT_BUTTON_SIZE).toString() + "px"
			..opacity = _DEFAULT_BUTTON_OPACITY
			..cursor = "pointer";
		_addDomElement
			..children.add(new svg.SvgElement.svg(SVG_PLUS))
			..onMouseOver.listen((m) => _markForAdd())
			..onMouseLeave.listen((m) => _clearMark())
			..onClick.listen(_addCopy)
			..onContextMenu.listen(_addCopy);
		html.document.body.children.add(_addDomElement);



		if(_canBeDeleted) {
			_removeDomElement = new html.Element.div();
			_removeDomElement.style
				..display = "none"
				..position = "absolute"
				..backgroundColor = _DEFAULT_BUTTON_COLOR_REMOVE
				..width = _REPEAT_BUTTON_SIZE.toString() + "px"
				..height = _REPEAT_BUTTON_SIZE.toString() + "px"
				..borderRadius = (_REPEAT_BUTTON_SIZE).toString() + "px"
				..opacity = _DEFAULT_BUTTON_OPACITY
				..cursor = "pointer";
			_removeDomElement
				..children.add(new svg.SvgElement.svg(SVG_REMOVE))
				..onMouseOver.listen((m) => _markForRemove())
				..onMouseLeave.listen((m) => _clearMark())
				..onClick.listen(_removeCopy)
				..onContextMenu.listen(_removeCopy);
			html.document.body.children.add(_removeDomElement);
		}

		_moveUpDomElement = new html.Element.div();
		_moveUpDomElement.style
			..display = "none"
			..position = "absolute"
			..backgroundColor = _DEFAULT_BUTTON_COLOR_MOVE_UP
			..width = _REPEAT_BUTTON_SIZE.toString() + "px"
			..height = _REPEAT_BUTTON_SIZE.toString() + "px"
			..borderRadius = _REPEAT_BUTTON_SIZE.toString() + "px"
			..opacity = _DEFAULT_BUTTON_OPACITY
			..cursor = "pointer";
		_moveUpDomElement
			..children.add(new svg.SvgElement.svg(SVG_UP))
			..onMouseOver.listen((m) => _markForMoveUp())
			..onMouseLeave.listen((m) => _clearMark())
			..onClick.listen(_moveUp)
			..onContextMenu.listen(_moveUp);
		html.document.body.children.add(_moveUpDomElement);

		_moveDownDomElement = new html.Element.div();
		_moveDownDomElement.style
			..display = "none"
			..position = "absolute"
			..backgroundColor = _DEFAULT_BUTTON_COLOR_MOVE_DOWN
			..width = _REPEAT_BUTTON_SIZE.toString() + "px"
			..height = _REPEAT_BUTTON_SIZE.toString() + "px"
			..borderRadius = _REPEAT_BUTTON_SIZE.toString() + "px"
			..opacity = _DEFAULT_BUTTON_OPACITY
			..cursor = "pointer";
		_moveDownDomElement
			..children.add(new svg.SvgElement.svg(SVG_DOWN))
			..onMouseOver.listen((m) => _markForMoveDown())
			..onMouseLeave.listen((m) => _clearMark())
			..onClick.listen(_moveDown)
			..onContextMenu.listen(_moveDown);
		html.document.body.children.add(_moveDownDomElement);
	}

	void show() {

		_isShown = true;

		_mark();

		_addDomElement.style
			..left = (_domElement.offsetLeft +
				_domElement.offsetWidth -
				_REPEAT_BUTTON_SIZE * 4).toString() + "px"
			..top = (_domElement.offsetTop -
				_REPEAT_BUTTON_SIZE / 2).toString() + "px"
			..display = "block";

		if(_canBeDeleted) {

			_removeDomElement.style
				..left = (_domElement.offsetLeft +
					_domElement.offsetWidth -
					_REPEAT_BUTTON_SIZE * 2.5).toString() +
					"px"
				..top = (_domElement.offsetTop -
					_REPEAT_BUTTON_SIZE / 2).toString() +
					"px"
				..display = "block";
		}

		_moveUpDomElement.style
			..left = (_domElement.offsetLeft +
				_domElement.offsetWidth -
				_REPEAT_BUTTON_SIZE / 2).toString() + "px"
			..top = (_domElement.offsetTop -
				_REPEAT_BUTTON_SIZE / 2).toString() + "px"
			..display = "block";

		_moveDownDomElement.style
			..left = (_domElement.offsetLeft +
				_domElement.offsetWidth -
				_REPEAT_BUTTON_SIZE / 2).toString() + "px"
			..top = (_domElement.offsetTop +
				_REPEAT_BUTTON_SIZE * 0.6).toString() + "px"
			..display = "block";

	}

	void hide() {

		_isShown = false;

		_clearMark();

		_addDomElement.style
			..display = "none";

		if(_canBeDeleted) {

			_removeDomElement.style
				..display = "none";
		}

		_moveUpDomElement.style
			..display = "none";

		_moveDownDomElement.style
			..display = "none";

	}

	void _addCopy(html.MouseEvent e) {

		hide();
		_repeat.addCopy(_key, _domElement);
		show();

		// Ensure that no links will be triggered
		e.stopPropagation();
		e.stopImmediatePropagation();
		e.preventDefault();
	}

	void _removeCopy(html.MouseEvent e) {
		_repeat.removeCopy(_key, _domElement);

		_addDomElement.remove();

		_moveUpDomElement.remove();
		_moveDownDomElement.remove();

		if(_canBeDeleted) {
			_removeDomElement.remove();
		}

		// Ensure that no links will be triggered
		e.stopPropagation();
		e.stopImmediatePropagation();
		e.preventDefault();
	}

	void _moveUp(html.MouseEvent e) {

		_repeat.moveCopyUp(_key);

		// Ensure that no links will be triggered
		e.stopPropagation();
		e.stopImmediatePropagation();
		e.preventDefault();
	}

	void _moveDown(html.MouseEvent e) {

		_repeat.moveCopyDown(_key);

		// Ensure that no links will be triggered
		e.stopPropagation();
		e.stopImmediatePropagation();
		e.preventDefault();
	}

  void prepareDomForHtmlSave() {
    if(_addDomElement != null) _addDomElement.remove();
    if(_removeDomElement != null) _removeDomElement.remove();
    if(_moveUpDomElement != null) _moveUpDomElement.remove();
    if(_moveDownDomElement != null) _moveDownDomElement.remove();
  }

  void restoreDomAfterHtmlSave() {
    if(_addDomElement != null) html.document.body.children.add(_addDomElement);
    if(_removeDomElement != null) html.document.body.children.add(_removeDomElement);
    if(_moveUpDomElement != null) html.document.body.children.add(_moveUpDomElement);
    if(_moveDownDomElement != null) html.document.body.children.add(_moveDownDomElement);
  }

}

part of wedit;

class PageMenu {

	static const _DEFAULT_MENU_COLOR = "#aaa";
	static const _DEFAULT_MENU_OPACITY = ".6";
	static const _DEFAULT_MENU_WIDTH = 400;
	static const _DEFAULT_MENU_HEIGHT = 100;
	static const _DEFAULT_MENU_TRIGGER_SIZE = 20;
	static const _DEFAULT_MENU_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
	static const _DEFAULT_MENU_HOVER_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, 1), inset 0 0 2vw 0 rgba(255, 255, 255, 1)";

	Page _page;

	html.Element _domElement;
	html.InputElement _titleInputElement;
	html.Element _saveDomElement;
	html.Element _publishDomElement;
	html.Element _deleteDomElement;

	html.Point _pointer;
	bool _holdPointer;

	PageMenu(this._page) {

		_bindControls();
	}

	void _bindControls() {

		_domElement = new html.Element.div();

		_domElement.style
			..display = "none"
			..zIndex = "999999"
			..position = "fixed"
			..backgroundColor = _DEFAULT_MENU_COLOR
			..width = _DEFAULT_MENU_WIDTH.toString() + "px"
			..height = _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px"
			..top = "0px"
			..left = "50%"
			..transform = "translateX(-50%)"
			..overflow = "hidden"
			..borderBottomLeftRadius = (_DEFAULT_MENU_TRIGGER_SIZE/2).toString() + "px"
			..borderBottomRightRadius = (_DEFAULT_MENU_TRIGGER_SIZE/2).toString() + "px"
			..opacity = _DEFAULT_MENU_OPACITY
			..boxShadow = _DEFAULT_MENU_BOX_SHADOW
			..cursor = "pointer";

		_domElement
			..onMouseEnter.listen(_mouseOver)
			..onMouseLeave.listen(_mouseLeave);

		_titleInputElement = new html.InputElement();
		_titleInputElement.style
			..padding = "1px"
			..borderWidth = "2px"
			..margin = _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px"
			..width = (_DEFAULT_MENU_WIDTH - 2*_DEFAULT_MENU_TRIGGER_SIZE - 6).toString() + "px";
		_titleInputElement.text = _page.title;
		_domElement.children.add(_titleInputElement);

		_saveDomElement = new html.Element.div();
		_saveDomElement.style
			..marginLeft = _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px"
			..width = (2*_DEFAULT_MENU_TRIGGER_SIZE).toString() + "px"
			..height = (2*_DEFAULT_MENU_TRIGGER_SIZE).toString() + "px"
			..backgroundColor = "red";
		_domElement.children.add(_saveDomElement);


		html.document.body.children.add(_domElement);

		//html.document.onMouseMove.listen(_mouseMove);
		html.window.onKeyDown.listen(_windowKeyDown);
		html.window.onKeyUp.listen(_windowKeyUp);

	}

	void _windowKeyDown(html.KeyboardEvent e) {
		if (e.ctrlKey && _page.inEditMode) {
			show();
		}
	}

	void _windowKeyUp(html.KeyboardEvent e) {
		hide();
	}

	void _mouseOver(html.MouseEvent event) {

		_domElement.animate([
			{"height": _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px"},
			{"height": _DEFAULT_MENU_HEIGHT.toString() + "px"}
		], 100);

		_domElement.style
			..animationDelay = "2s"
			..height = _DEFAULT_MENU_HEIGHT.toString() + "px"
			..boxShadow = _DEFAULT_MENU_HOVER_BOX_SHADOW;
	}

	void _mouseLeave(html.MouseEvent event) {

		_domElement.animate([
			{"height": _DEFAULT_MENU_HEIGHT.toString() + "px"},
			{"height": _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px"}
		], 100);

		_domElement.style
			..animationDelay = "2s"
			..height = _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px"
			..boxShadow = _DEFAULT_MENU_BOX_SHADOW;
	}

	void show() {

		_domElement.style
			..display = "block";

		_holdPointer = true;

	}

	void hide() {

		_domElement.style
			..display = "none";

		_holdPointer = false;
	}

}
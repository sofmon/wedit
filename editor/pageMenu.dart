// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class PageMenu {
  static const _DEFAULT_MENU_COLOR = "#aaa";
  static const _DEFAULT_MENU_OPACITY = ".6";
  static const _DEFAULT_MENU_WIDTH = 400;
  static const _DEFAULT_MENU_TRIGGER_SIZE = 20;
  static const _DEFAULT_MENU_BUTTON_HEIGHT = 40;
  static const _DEFAULT_MENU_BOX_SHADOW =
      "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  static const _DEFAULT_MENU_HOVER_BOX_SHADOW =
      "0 0 2vw 0 rgba(0, 0, 0, 1), inset 0 0 2vw 0 rgba(255, 255, 255, 1)";
  static const _DEFAULT_COMMAND_SUCCESS_TEXT = "ok";
  static const _DEFAULT_COMMAND_FAILURE_TEXT = "ERROR";

  Page _page;
  Map<String, String> _commands;
  String _menuTextColor;

  html.Element _domElement;

  Map<String, html.Element> _commandDomElements;

  bool _lockMenu;

  PageMenu(this._page, this._commands, this._menuTextColor) {
    if(_commands.length <= 0) {
      return;
    }

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
      ..overflow = "hidden"
      ..borderBottomLeftRadius =
          (_DEFAULT_MENU_TRIGGER_SIZE / 2).toString() + "px"
      ..borderBottomRightRadius =
          (_DEFAULT_MENU_TRIGGER_SIZE / 2).toString() + "px"
      ..opacity = _DEFAULT_MENU_OPACITY
      ..boxShadow = _DEFAULT_MENU_BOX_SHADOW
      ..zIndex = "1000000"
      ..transform = "translateX(-50%) translateZ(1000000em)"
      ..cursor = "pointer";

    _domElement
      ..onMouseEnter.listen(_mouseOver)
      ..onMouseLeave.listen(_mouseLeave);

    html.document.body.children.add(_domElement);

    _commandDomElements = new Map<String, html.Element>();
    _commands.forEach((k, v) => _createCommandElement(k, v));

    html.window.onKeyDown.listen(_windowKeyDown);
    html.window.onKeyUp.listen(_windowKeyUp);
  }

  void _createCommandElement(String command, String color) {
    html.Element el = new html.Element.div();
    el.style
      ..marginTop = (0.5 * _DEFAULT_MENU_TRIGGER_SIZE).toString() + "px"
      ..marginBottom = (0.5 * _DEFAULT_MENU_TRIGGER_SIZE).toString() + "px"
      ..marginLeft = (0.05 * _DEFAULT_MENU_WIDTH).toString() + "px"
      ..width = (0.90 * _DEFAULT_MENU_WIDTH).toString() + "px"
      ..height = _DEFAULT_MENU_BUTTON_HEIGHT.toString() + "px"
      ..borderRadius = (_DEFAULT_MENU_BUTTON_HEIGHT / 2).toString() + "px"
      ..fontSize = _DEFAULT_MENU_BUTTON_HEIGHT.toString() + "px"
      ..textAlign = "center"
      ..color = _menuTextColor
      ..backgroundColor = color;
    el.text = command;
    el.onClick.listen(_commandClick);
    _commandDomElements[command] = el;
    _domElement.children.add(el);
  }

  void _windowKeyDown(html.KeyboardEvent e) {
    if (e.ctrlKey) {
      show();
    }
  }

  void _windowKeyUp(html.KeyboardEvent e) {
    if (!_lockMenu) {
      hide();
    }
  }

  void _mouseOver(html.MouseEvent event) {
    _domElement.style
      ..animationDelay = "2s"
      ..height = ""
      ..boxShadow = _DEFAULT_MENU_HOVER_BOX_SHADOW;

    _lockMenu = true;
  }

  void _mouseLeave(html.MouseEvent event) {
    _domElement.style
      ..animationDelay = "2s"
      ..height = _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px"
      ..boxShadow = _DEFAULT_MENU_BOX_SHADOW;

    _lockMenu = false;

    hide();
  }

  void _commandClick(html.MouseEvent event) {
    html.Element el = event.target;
    String cmd = el.text;
    if (cmd == _DEFAULT_COMMAND_SUCCESS_TEXT ||
        cmd == _DEFAULT_COMMAND_FAILURE_TEXT) {
      return;
    }

    _page.command(cmd, () => el.text = _DEFAULT_COMMAND_SUCCESS_TEXT,
        () => el.text = _DEFAULT_COMMAND_FAILURE_TEXT);
  }

  void show() {
    _domElement.style.display = "block";

    _commandDomElements.forEach((k, v) => v.text = k);
  }

  void hide() {
    _domElement.style..display = "none";
  }
}

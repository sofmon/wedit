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
  Map<String, String>? _commands;
  String _menuTextColor;

  web.HTMLElement? _domElement;

  Map<String, web.HTMLElement> _commandDomElements = {};

  bool _lockMenu = false;

  PageMenu(this._page, this._commands, this._menuTextColor) {
    if(_commands == null || _commands!.length <= 0) {
      return;
    }

    _bindControls();
  }

  void _bindControls() {
    _domElement = web.document.createElement('div') as web.HTMLElement;

    _domElement!.style
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

    _domElement!.addEventListener('mouseenter', _mouseOver.toJS);
    _domElement!.addEventListener('mouseleave', _mouseLeave.toJS);

    web.document.body?.appendChild(_domElement!);

    _commandDomElements = <String, web.HTMLElement>{};
    _commands?.forEach((k, v) => _createCommandElement(k, v));

    web.window.addEventListener('keydown', _windowKeyDown.toJS);
    web.window.addEventListener('keyup', _windowKeyUp.toJS);
  }

  void _createCommandElement(String command, String color) {
    final el = web.document.createElement('div') as web.HTMLElement;
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
    el.textContent = command;
    el.addEventListener('click', _commandClick.toJS);
    _commandDomElements[command] = el;
    _domElement?.appendChild(el);
  }

  void _windowKeyDown(web.KeyboardEvent e) {
    if (e.ctrlKey) {
      show();
    }
  }

  void _windowKeyUp(web.KeyboardEvent e) {
    if (!_lockMenu) {
      hide();
    }
  }

  void _mouseOver(web.MouseEvent event) {
    final style = _domElement?.style;
    if (style != null) {
      style.animationDelay = "2s";
      style.height = "";
      style.boxShadow = _DEFAULT_MENU_HOVER_BOX_SHADOW;
    }

    _lockMenu = true;
  }

  void _mouseLeave(web.MouseEvent event) {
    final style = _domElement?.style;
    if (style != null) {
      style.animationDelay = "2s";
      style.height = _DEFAULT_MENU_TRIGGER_SIZE.toString() + "px";
      style.boxShadow = _DEFAULT_MENU_BOX_SHADOW;
    }

    _lockMenu = false;

    hide();
  }

  void _commandClick(web.MouseEvent event) {
    final el = event.target as web.HTMLElement?;
    if (el == null) return;

    final cmd = el.textContent ?? '';
    if (cmd == _DEFAULT_COMMAND_SUCCESS_TEXT ||
        cmd == _DEFAULT_COMMAND_FAILURE_TEXT) {
      return;
    }

    _page.command(cmd, () => el.textContent = _DEFAULT_COMMAND_SUCCESS_TEXT,
        () => el.textContent = _DEFAULT_COMMAND_FAILURE_TEXT);
  }

  void show() {
    _domElement?.style.display = "block";

    _commandDomElements.forEach((k, v) => v.textContent = k);
  }

  void hide() {
    _domElement?.style.display = "none";
  }
}

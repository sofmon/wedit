// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class RepeatShadow {
  Repeat _repeat;

  final web.HTMLElement _domElement;
  web.HTMLElement get domElement => _domElement;

  bool _isShown = false;

  String _key;
  bool _canBeDeleted = false;

  RepeatShadow.fromRepeat(this._repeat, this._key, this._domElement) {
    _isShown = false;

    _canBeDeleted = _key != _repeat.key;

    _bindControls();

    if (_repeat._page.darkMode) {
      _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(255, 255, 255, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";
      _DEFAULT_MARK_BOX_SHADOW_ADD = "0 0 2vw 0 rgba(0, 155, 0, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";
      _DEFAULT_MARK_BOX_SHADOW_REMOVE = "0 0 2vw 0 rgba(255, 0, 0, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";
      _DEFAULT_MARK_BOX_SHADOW_MOVE_UP = "0 0 2vw 0 rgba(0, 50, 255, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";
      _DEFAULT_MARK_BOX_SHADOW_MOVE_DOWN = "0 0 2vw 0 rgba(0, 50, 255, .5), inset 0 0 2vw 0 rgba(0, 0, 0, .5)";
    }
  }

  web.HTMLElement? _addDomElement;
  web.HTMLElement? _removeDomElement;
  web.HTMLElement? _moveUpDomElement;
  web.HTMLElement? _moveDownDomElement;

  static const _REPEAT_BUTTON_SIZE = 20;

  String _originalBoxShadow = '';

  static const _DEFAULT_BUTTON_COLOR_ADD = "#0a0";
  static const _DEFAULT_BUTTON_COLOR_REMOVE = "#f00";
  static const _DEFAULT_BUTTON_COLOR_MOVE_UP = "#06f";
  static const _DEFAULT_BUTTON_COLOR_MOVE_DOWN = "#00f";
  static const _DEFAULT_BUTTON_OPACITY = ".3";

  String _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  String _DEFAULT_MARK_BOX_SHADOW_ADD = "0 0 2vw 0 rgba(0, 155, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  String _DEFAULT_MARK_BOX_SHADOW_REMOVE = "0 0 2vw 0 rgba(255, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  String _DEFAULT_MARK_BOX_SHADOW_MOVE_UP = "0 0 2vw 0 rgba(0, 50, 255, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  String _DEFAULT_MARK_BOX_SHADOW_MOVE_DOWN = "0 0 2vw 0 rgba(0, 50, 255, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";

  static const SVG_UP =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M16 1l-15 15h9v16h12v-16h9z"></path></svg>';
  static const SVG_DOWN =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"> <path d="M16 31l15-15h-9v-16h-12v16h-9z"></path></svg>';
  static const SVG_PLUS =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M31 12h-11v-11c0-0.552-0.448-1-1-1h-6c-0.552 0-1 0.448-1 1v11h-11c-0.552 0-1 0.448-1 1v6c0 0.552 0.448 1 1 1h11v11c0 0.552 0.448 1 1 1h6c0.552 0 1-0.448 1-1v-11h11c0.552 0 1-0.448 1-1v-6c0-0.552-0.448-1-1-1z"></path></svg>';
  static const SVG_REMOVE =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M31.708 25.708c-0-0-0-0-0-0l-9.708-9.708 9.708-9.708c0-0 0-0 0-0 0.105-0.105 0.18-0.227 0.229-0.357 0.133-0.356 0.057-0.771-0.229-1.057l-4.586-4.586c-0.286-0.286-0.702-0.361-1.057-0.229-0.13 0.048-0.252 0.124-0.357 0.228 0 0-0 0-0 0l-9.708 9.708-9.708-9.708c-0-0-0-0-0-0-0.105-0.104-0.227-0.18-0.357-0.228-0.356-0.133-0.771-0.057-1.057 0.229l-4.586 4.586c-0.286 0.286-0.361 0.702-0.229 1.057 0.049 0.13 0.124 0.252 0.229 0.357 0 0 0 0 0 0l9.708 9.708-9.708 9.708c-0 0-0 0-0 0-0.104 0.105-0.18 0.227-0.229 0.357-0.133 0.355-0.057 0.771 0.229 1.057l4.586 4.586c0.286 0.286 0.702 0.361 1.057 0.229 0.13-0.049 0.252-0.124 0.357-0.229 0-0 0-0 0-0l9.708-9.708 9.708 9.708c0 0 0 0 0 0 0.105 0.105 0.227 0.18 0.357 0.229 0.356 0.133 0.771 0.057 1.057-0.229l4.586-4.586c0.286-0.286 0.362-0.702 0.229-1.057-0.049-0.13-0.124-0.252-0.229-0.357z"></path></svg>';

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
    _domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW_MOVE_DOWN;
  }

  void _clearMark() {
    _domElement.style.boxShadow = _isShown ? _DEFAULT_MARK_BOX_SHADOW : _originalBoxShadow;
  }

  void _bindControls() {
    _originalBoxShadow = _domElement.style.boxShadow;

    _addDomElement = web.document.createElement('div') as web.HTMLElement;
    _addDomElement!.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _DEFAULT_BUTTON_COLOR_ADD
      ..width = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..height = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..borderRadius = (_REPEAT_BUTTON_SIZE).toString() + "px"
      ..opacity = _DEFAULT_BUTTON_OPACITY
      ..cursor = "pointer";
    _addDomElement!.innerHTML = SVG_PLUS.toJS;
    _addDomElement!.addEventListener('mouseover', ((web.Event m) => _markForAdd()).toJS);
    _addDomElement!.addEventListener('mouseleave', ((web.Event m) => _clearMark()).toJS);
    _addDomElement!.addEventListener('click', _addCopy.toJS);
    _addDomElement!.addEventListener('contextmenu', _addCopy.toJS);
    web.document.body?.appendChild(_addDomElement!);

    if (_canBeDeleted) {
      _removeDomElement = web.document.createElement('div') as web.HTMLElement;
      _removeDomElement!.style
        ..display = "none"
        ..position = "absolute"
        ..backgroundColor = _DEFAULT_BUTTON_COLOR_REMOVE
        ..width = _REPEAT_BUTTON_SIZE.toString() + "px"
        ..height = _REPEAT_BUTTON_SIZE.toString() + "px"
        ..borderRadius = (_REPEAT_BUTTON_SIZE).toString() + "px"
        ..opacity = _DEFAULT_BUTTON_OPACITY
        ..cursor = "pointer";
      _removeDomElement!.innerHTML = SVG_REMOVE.toJS;
      _removeDomElement!.addEventListener('mouseover', ((web.Event m) => _markForRemove()).toJS);
      _removeDomElement!.addEventListener('mouseleave', ((web.Event m) => _clearMark()).toJS);
      _removeDomElement!.addEventListener('click', _removeCopy.toJS);
      _removeDomElement!.addEventListener('contextmenu', _removeCopy.toJS);
      web.document.body?.appendChild(_removeDomElement!);
    }

    _moveUpDomElement = web.document.createElement('div') as web.HTMLElement;
    _moveUpDomElement!.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _DEFAULT_BUTTON_COLOR_MOVE_UP
      ..width = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..height = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..borderRadius = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..opacity = _DEFAULT_BUTTON_OPACITY
      ..cursor = "pointer";
    _moveUpDomElement!.innerHTML = SVG_UP.toJS;
    _moveUpDomElement!.addEventListener('mouseover', ((web.Event m) => _markForMoveUp()).toJS);
    _moveUpDomElement!.addEventListener('mouseleave', ((web.Event m) => _clearMark()).toJS);
    _moveUpDomElement!.addEventListener('click', _moveUp.toJS);
    _moveUpDomElement!.addEventListener('contextmenu', _moveUp.toJS);
    web.document.body?.appendChild(_moveUpDomElement!);

    _moveDownDomElement = web.document.createElement('div') as web.HTMLElement;
    _moveDownDomElement!.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _DEFAULT_BUTTON_COLOR_MOVE_DOWN
      ..width = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..height = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..borderRadius = _REPEAT_BUTTON_SIZE.toString() + "px"
      ..opacity = _DEFAULT_BUTTON_OPACITY
      ..cursor = "pointer";
    _moveDownDomElement!.innerHTML = SVG_DOWN.toJS;
    _moveDownDomElement!.addEventListener('mouseover', ((web.Event m) => _markForMoveDown()).toJS);
    _moveDownDomElement!.addEventListener('mouseleave', ((web.Event m) => _clearMark()).toJS);
    _moveDownDomElement!.addEventListener('click', _moveDown.toJS);
    _moveDownDomElement!.addEventListener('contextmenu', _moveDown.toJS);
    web.document.body?.appendChild(_moveDownDomElement!);
  }

  int _getOffsetTop(web.HTMLElement? el) {
    int res = 0;

    while (el != null) {
      res += el.offsetTop.round();
      el = el.offsetParent as web.HTMLElement?;
    }

    return res;
  }

  int _getOffsetLeft(web.HTMLElement? el) {
    int res = 0;

    while (el != null) {
      res += el.offsetLeft.round();
      el = el.offsetParent as web.HTMLElement?;
    }

    return res;
  }

  void show() {
    _isShown = true;

    _mark();

    final addStyle = _addDomElement?.style;
    if (addStyle != null) {
      addStyle.left = (_getOffsetLeft(_domElement) + _domElement.offsetWidth.round() - _REPEAT_BUTTON_SIZE * 4).toString() + "px";
      addStyle.top = (_getOffsetTop(_domElement) - _REPEAT_BUTTON_SIZE / 2).toString() + "px";
      addStyle.display = "block";
    }

    if (_canBeDeleted) {
      final removeStyle = _removeDomElement?.style;
      if (removeStyle != null) {
        removeStyle.left = (_getOffsetLeft(_domElement) + _domElement.offsetWidth.round() - _REPEAT_BUTTON_SIZE * 2.5).toString() + "px";
        removeStyle.top = (_getOffsetTop(_domElement) - _REPEAT_BUTTON_SIZE / 2).toString() + "px";
        removeStyle.display = "block";
      }
    }

    final moveUpStyle = _moveUpDomElement?.style;
    if (moveUpStyle != null) {
      moveUpStyle.left = (_getOffsetLeft(_domElement) + _domElement.offsetWidth.round() - _REPEAT_BUTTON_SIZE / 2).toString() + "px";
      moveUpStyle.top = (_getOffsetTop(_domElement) - _REPEAT_BUTTON_SIZE / 2).toString() + "px";
      moveUpStyle.display = "block";
    }

    final moveDownStyle = _moveDownDomElement?.style;
    if (moveDownStyle != null) {
      moveDownStyle.left = (_getOffsetLeft(_domElement) + _domElement.offsetWidth.round() - _REPEAT_BUTTON_SIZE / 2).toString() + "px";
      moveDownStyle.top = (_getOffsetTop(_domElement) + _REPEAT_BUTTON_SIZE * 0.6).toString() + "px";
      moveDownStyle.display = "block";
    }
  }

  void hide() {
    _isShown = false;

    _clearMark();

    _addDomElement?.style.display = "none";

    if (_canBeDeleted) {
      _removeDomElement?.style.display = "none";
    }

    _moveUpDomElement?.style.display = "none";

    _moveDownDomElement?.style.display = "none";
  }

  void _addCopy(web.MouseEvent e) {
    hide();
    _repeat.addCopy(_key, _domElement);
    show();

    // Ensure that no links will be triggered
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();

    _repeat._page.save(() {}, () {});
  }

  void _removeCopy(web.MouseEvent e) {
    _repeat.removeCopy(_key, _domElement);

    _addDomElement?.remove();

    _moveUpDomElement?.remove();
    _moveDownDomElement?.remove();

    if (_canBeDeleted) {
      _removeDomElement?.remove();
    }

    // Ensure that no links will be triggered
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();

    _repeat._page.save(() {}, () {});
  }

  void _moveUp(web.MouseEvent e) {
    _repeat.moveCopyUp(_key);

    // Ensure that no links will be triggered
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();

    _repeat._page.save(() {}, () {});
  }

  void _moveDown(web.MouseEvent e) {
    _repeat.moveCopyDown(_key);

    // Ensure that no links will be triggered
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();

    _repeat._page.save(() {}, () {});
  }

  void prepareDomForHtmlSave() {
    _addDomElement?.remove();
    _removeDomElement?.remove();
    _moveUpDomElement?.remove();
    _moveDownDomElement?.remove();
  }

  void restoreDomAfterHtmlSave() {
    if (_addDomElement != null) web.document.body?.appendChild(_addDomElement!);
    if (_removeDomElement != null) web.document.body?.appendChild(_removeDomElement!);
    if (_moveUpDomElement != null) web.document.body?.appendChild(_moveUpDomElement!);
    if (_moveDownDomElement != null) web.document.body?.appendChild(_moveDownDomElement!);
  }
}

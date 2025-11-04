// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Image {
  static const SUPPORTED_IMAGE_TAGS = const ["img"];

  static const _DEFAULT_BUTTON_ADD_COLOR = "#0a0";
  static const _DEFAULT_BUTTON_REMOVE_COLOR = "#a00";
  static const _DEFAULT_BUTTON_OPACITY = ".3";
  static const _DEFAULT_BUTTON_SIZE = 20;

  static const _DEFAULT_MARK_BOX_SHADOW = "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  static const _DEFAULT_MARK_BOX_SHADOW_ADD = "0 0 2vw 0 rgba(0, 155, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  static const _DEFAULT_MARK_BOX_SHADOW_REMOVE = "0 0 2vw 0 rgba(255, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";

  static const SVG_UPLOAD =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M29.996 4c0.001 0.001 0.003 0.002 0.004 0.004v23.993c-0.001 0.001-0.002 0.003-0.004 0.004h-27.993c-0.001-0.001-0.003-0.002-0.004-0.004v-23.993c0.001-0.001 0.002-0.003 0.004-0.004h27.993zM30 2h-28c-1.1 0-2 0.9-2 2v24c0 1.1 0.9 2 2 2h28c1.1 0 2-0.9 2-2v-24c0-1.1-0.9-2-2-2v0z"></path><path d="M26 9c0 1.657-1.343 3-3 3s-3-1.343-3-3 1.343-3 3-3 3 1.343 3 3z"></path><path d="M28 26h-24v-4l7-12 8 10h2l7-6z"></path></svg>';
  static const SVG_REVERT =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:2px 1px 2px 1px;position:absolute" width="16" height="14" viewBox="0 0 34 32"><path d="M20 2c7.732 0 14 6.268 14 14s-6.268 14-14 14v-3c2.938 0 5.701-1.144 7.778-3.222s3.222-4.84 3.222-7.778c0-2.938-1.144-5.701-3.222-7.778s-4.84-3.222-7.778-3.222c-2.938 0-5.701 1.144-7.778 3.222-1.598 1.598-2.643 3.601-3.041 5.778h5.819l-7 8-7-8h5.143c0.971-6.784 6.804-12 13.857-12zM26 14v4h-8v-10h4v6z"></path></svg>';

  Page _page;

  String _key;
  String _type = '';
  List<int>? _width;
  List<double>? _pixelDensity;

  bool _hasContent = false;

  bool _isHighlighted = false;

  final web.HTMLElement _domElement;
  String _originalBoxShadow = '';

  web.HTMLInputElement? _fileUploadDomElement;
  web.HTMLElement? _addImageDomElement;
  web.HTMLElement? _removeImageDomElement;

  String _templateSrc = '';
  String _templateSrcset = '';

  Image.fromMap(this._page, this._key, this._domElement, Map<dynamic, dynamic>? map) {
    _isHighlighted = false;

    if (map != null) {
      _hasContent = true;
      _width = (map[IMAGE_WIDTH] as List?)?.cast<int>();
      _pixelDensity = (map[IMAGE_PIXEL_DENSITY] as List?)?.cast<double>();
      _type = (map[IMAGE_TYPE] as String?) ?? '';
      _hasContent = (_type != "" && _type.isNotEmpty);
    } else {
      _hasContent = false;
    }

    _syncElement();

    _bindControls();
  }

  Map<String, dynamic> toMap() {
    final map = <String, dynamic>{};

    map[IMAGE_KEY] = _key;
    map[IMAGE_TYPE] = _type;
    map[IMAGE_WIDTH] = _width;
    map[IMAGE_PIXEL_DENSITY] = _pixelDensity;

    return map;
  }

  void _syncElement() {
    _originalBoxShadow = _domElement.style.boxShadow;

    final image = _domElement as web.HTMLImageElement;
    _templateSrc = image.src;
    _templateSrcset = image.srcset;
  }

  void _bindControls() {
    _fileUploadDomElement = web.document.createElement('input') as web.HTMLInputElement;
    _fileUploadDomElement!.type = 'file';
    _fileUploadDomElement!.style
      ..position = "absolute"
      ..width = "1px"
      ..height = "1px"
      ..overflow = "hidden"
      ..opacity = "0";
    _fileUploadDomElement!.addEventListener('change', _uploadFile.toJS);
    web.document.body?.appendChild(_fileUploadDomElement!);

    _addImageDomElement = web.document.createElement('div') as web.HTMLElement;
    _addImageDomElement!.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _DEFAULT_BUTTON_ADD_COLOR
      ..width = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..height = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..borderRadius = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..opacity = _DEFAULT_BUTTON_OPACITY
      ..cursor = "pointer";
    _addImageDomElement!.innerHTML = SVG_UPLOAD.toJS;
    _addImageDomElement!.addEventListener('mouseover', ((web.Event e) => _markForAdd()).toJS);
    _addImageDomElement!.addEventListener('mouseleave', ((web.Event e) => _clearMark()).toJS);
    _addImageDomElement!.addEventListener('mousedown', _addImage.toJS);
    _addImageDomElement!.addEventListener('contextmenu', _stopEvent.toJS);
    web.document.body?.appendChild(_addImageDomElement!);

    _removeImageDomElement = web.document.createElement('div') as web.HTMLElement;
    _removeImageDomElement!.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _DEFAULT_BUTTON_REMOVE_COLOR
      ..width = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..height = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..borderRadius = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..opacity = _DEFAULT_BUTTON_OPACITY
      ..cursor = "pointer";
    _removeImageDomElement!.innerHTML = SVG_REVERT.toJS;
    _removeImageDomElement!.addEventListener('mouseover', ((web.Event e) => _markForRemove()).toJS);
    _removeImageDomElement!.addEventListener('mouseleave', ((web.Event e) => _clearMark()).toJS);
    _removeImageDomElement!.addEventListener('click', _removeImage.toJS);
    _removeImageDomElement!.addEventListener('contextmenu', _removeImage.toJS);
    web.document.body?.appendChild(_removeImageDomElement!);
  }

  void _mark() {
    _domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW;
  }

  void _markForAdd() {
    _domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW_ADD;
  }

  void _markForRemove() {
    _domElement.style.boxShadow = _DEFAULT_MARK_BOX_SHADOW_REMOVE;
  }

  void _clearMark() {
    _domElement.style.boxShadow = _isHighlighted ? _DEFAULT_MARK_BOX_SHADOW : _originalBoxShadow;
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

  void highlight() {
    _isHighlighted = true;

    _mark();

    final addStyle = _addImageDomElement?.style;
    if (addStyle != null) {
      addStyle.left = (_getOffsetLeft(_domElement) + _domElement.offsetWidth.round() - _DEFAULT_BUTTON_SIZE * 2).toString() + "px";
      addStyle.top = (_getOffsetTop(_domElement) - _DEFAULT_BUTTON_SIZE / 2).toString() + "px";
      addStyle.display = "block";
    }

    final removeStyle = _removeImageDomElement?.style;
    if (removeStyle != null) {
      removeStyle.left = (_getOffsetLeft(_domElement) + _domElement.offsetWidth.round() - _DEFAULT_BUTTON_SIZE / 2).toString() + "px";
      removeStyle.top = (_getOffsetTop(_domElement) - _DEFAULT_BUTTON_SIZE / 2).toString() + "px";
      removeStyle.display = "block";
    }
  }

  void normalise() {
    _isHighlighted = false;

    _clearMark();

    _addImageDomElement?.style.display = "none";

    _removeImageDomElement?.style.display = "none";
  }

  void render() {
    final image = _domElement as web.HTMLImageElement;

    if (!_hasContent) {
      image.src = _templateSrc;
      image.srcset = _templateSrcset;
      return;
    }

    final suffix = "?" + DateTime.now().millisecondsSinceEpoch.toString();

    image.src = "./" + _key + "." + _type + suffix;
    final sb = StringBuffer();
    _width?.forEach((i) => sb.write("./" + _key + "-" + i.toString() + "w." + _type + suffix + " " + i.toString() + "w,"));
    _pixelDensity?.forEach(
      (f) => sb.write("./" + _key + "-" + f.toStringAsFixed(1) + "x." + _type + suffix + " " + f.toStringAsFixed(1) + "x,"),
    );
    image.srcset = sb.toString();
  }

  void _stopEvent(web.MouseEvent e) {
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();
  }

  void _addImage(web.MouseEvent e) {
    _stopEvent(e);

    final addElement = _addImageDomElement;
    final uploadElement = _fileUploadDomElement;
    if (addElement != null && uploadElement != null) {
      uploadElement.style.left = addElement.offsetLeft.toString() + "px";
      uploadElement.style.top = addElement.offsetTop.toString() + "px";
      uploadElement.focus();
      uploadElement.click();
    }
  }

  void _removeImage(web.MouseEvent e) {
    _type = "";
    _hasContent = false;

    render();

    _page.save(() {}, () {});

    _stopEvent(e);
  }

  void _uploadFile(web.Event e) {
    final files = _fileUploadDomElement?.files;
    if (files == null || files.length == 0) {
      return;
    }

    final file = files.item(0);
    if (file == null) {
      return;
    }

    final reader = web.FileReader();
    reader.addEventListener(
      'load',
      ((web.Event e) {
        _sendDataToServer(file.name, reader.result);
      }).toJS,
    );
    reader.readAsArrayBuffer(file);
  }

  void _sendDataToServer(String name, JSAny? data) async {
    final url =
        web.window.location.protocol +
        "//" +
        web.window.location.host +
        "/~?k=" +
        _key +
        "&n=" +
        name +
        "&p=" +
        web.window.location.pathname;

    try {
      final response = await web.window.fetch(url.toJS, web.RequestInit(method: 'POST', body: data)).toDart;

      if (response.ok) {
        final text = await response.text().toDart;
        _processImageUploadSuccess(text as String);
      } else {
        _processImageUploadError();
      }
    } catch (e) {
      _processImageUploadError();
    }
  }

  void _processImageUploadSuccess(String responseText) {
    final Map<String, dynamic> map = convert.json.decode(responseText) as Map<String, dynamic>;

    _type = (map[IMAGE_TYPE] as String?) ?? '';
    _width = (map[IMAGE_WIDTH] as List?)?.cast<int>();
    _pixelDensity = (map[IMAGE_PIXEL_DENSITY] as List?)?.cast<double>();
    _hasContent = true;

    render();

    print("upload complete");

    _page.save(() {}, () {});
  }

  void _processImageUploadError() {
    print("fail");
  }

  void prepareDomForHtmlSave() {
    _fileUploadDomElement?.remove();
    _addImageDomElement?.remove();
    _removeImageDomElement?.remove();
  }

  void restoreDomAfterHtmlSave() {
    if (_fileUploadDomElement != null) web.document.body?.appendChild(_fileUploadDomElement!);
    if (_addImageDomElement != null) web.document.body?.appendChild(_addImageDomElement!);
    if (_removeImageDomElement != null) web.document.body?.appendChild(_removeImageDomElement!);
  }
}

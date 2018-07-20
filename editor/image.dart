// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
part of wedit;

class Image {
  static const SUPPORTED_IMAGE_TAGS = const ["img"];

  static const _DEFAULT_BUTTON_ADD_COLOR = "#0a0";
  static const _DEFAULT_BUTTON_REMOVE_COLOR = "#a00";
  static const _DEFAULT_BUTTON_OPACITY = ".3";
  static const _DEFAULT_BUTTON_SIZE = 20;

  static const _DEFAULT_MARK_BOX_SHADOW =
      "0 0 2vw 0 rgba(0, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  static const _DEFAULT_MARK_BOX_SHADOW_ADD =
      "0 0 2vw 0 rgba(0, 155, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";
  static const _DEFAULT_MARK_BOX_SHADOW_REMOVE =
      "0 0 2vw 0 rgba(255, 0, 0, .5), inset 0 0 2vw 0 rgba(255, 255, 255, .5)";

  static const SVG_UPLOAD =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:3px;position:absolute" width="14" height="14" viewBox="0 0 32 32"><path d="M29.996 4c0.001 0.001 0.003 0.002 0.004 0.004v23.993c-0.001 0.001-0.002 0.003-0.004 0.004h-27.993c-0.001-0.001-0.003-0.002-0.004-0.004v-23.993c0.001-0.001 0.002-0.003 0.004-0.004h27.993zM30 2h-28c-1.1 0-2 0.9-2 2v24c0 1.1 0.9 2 2 2h28c1.1 0 2-0.9 2-2v-24c0-1.1-0.9-2-2-2v0z"></path><path d="M26 9c0 1.657-1.343 3-3 3s-3-1.343-3-3 1.343-3 3-3 3 1.343 3 3z"></path><path d="M28 26h-24v-4l7-12 8 10h2l7-6z"></path></svg>';
  static const SVG_REVERT =
      '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#fff" style="margin:2px 1px 2px 1px;position:absolute" width="16" height="14" viewBox="0 0 34 32"><path d="M20 2c7.732 0 14 6.268 14 14s-6.268 14-14 14v-3c2.938 0 5.701-1.144 7.778-3.222s3.222-4.84 3.222-7.778c0-2.938-1.144-5.701-3.222-7.778s-4.84-3.222-7.778-3.222c-2.938 0-5.701 1.144-7.778 3.222-1.598 1.598-2.643 3.601-3.041 5.778h5.819l-7 8-7-8h5.143c0.971-6.784 6.804-12 13.857-12zM26 14v4h-8v-10h4v6z"></path></svg>';

  Page _page;

  String _key;
  String _type;
  List<int> _width;
  List<double> _pixelDensity;

  bool _hasContent;

  bool _isHighlighted;

  html.Element _domElement;
  String _originalBoxShadow;

  html.InputElement _fileUploadDomElement;
  html.Element _addImageDomElement;
  html.Element _removeImageDomElement;

  String _templateSrc;
  String _templateSrcset;

  Image.fromMap(this._page, this._key, this._domElement, Map map) {
    _isHighlighted = false;
    
    if (map != null) {
      _hasContent= true;
      _width = map[IMAGE_WIDTH];
      _pixelDensity = map[IMAGE_PIXEL_DENSITY];
      _type = map[IMAGE_TYPE];
      _hasContent = (_type != "" && _type.isNotEmpty);
    } else {
      _hasContent = false;
    }

    _syncElement();

    _bindControls();
  }

  Map toMap() {
    var map = new Map();

    map[IMAGE_KEY] = _key;
    map[IMAGE_TYPE] = _type;
    map[IMAGE_WIDTH] = _width;
    map[IMAGE_PIXEL_DENSITY] = _pixelDensity;

    return map;
  }

  void _syncElement() {
    _originalBoxShadow = _domElement.style.boxShadow;
    
     var image = _domElement as html.ImageElement;
    _templateSrc = image.src;
    _templateSrcset = image.srcset;
  }

  void _bindControls() {
    _fileUploadDomElement = new html.InputElement(type: "file");
    _fileUploadDomElement.style
      ..position = "absolute"
      ..width = "1px"
      ..height = "1px"
      ..overflow = "hidden"
      ..opacity = "0";
    _fileUploadDomElement.onChange.listen(_uploadFile);
    html.document.body.children.add(_fileUploadDomElement);

    _addImageDomElement = new html.Element.div();
    _addImageDomElement.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _DEFAULT_BUTTON_ADD_COLOR
      ..width = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..height = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..borderRadius = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..opacity = _DEFAULT_BUTTON_OPACITY
      ..cursor = "pointer";
    _addImageDomElement
      ..children.add(new svg.SvgElement.svg(SVG_UPLOAD))
      ..onMouseOver.listen((m) => _markForAdd())
      ..onMouseLeave.listen((m) => _clearMark())
      ..onMouseDown.listen(_addImage)
      ..onContextMenu.listen(_stopEvent);
    html.document.body.children.add(_addImageDomElement);

    _removeImageDomElement = new html.Element.div();
    _removeImageDomElement.style
      ..display = "none"
      ..position = "absolute"
      ..backgroundColor = _DEFAULT_BUTTON_REMOVE_COLOR
      ..width = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..height = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..borderRadius = _DEFAULT_BUTTON_SIZE.toString() + "px"
      ..opacity = _DEFAULT_BUTTON_OPACITY
      ..cursor = "pointer";
    _removeImageDomElement
      ..children.add(new svg.SvgElement.svg(SVG_REVERT))
      ..onMouseOver.listen((m) => _markForRemove())
      ..onMouseLeave.listen((m) => _clearMark())
      ..onClick.listen(_removeImage)
      ..onContextMenu.listen(_removeImage);
    html.document.body.children.add(_removeImageDomElement);
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
    _domElement.style.boxShadow =
        _isHighlighted ? _DEFAULT_MARK_BOX_SHADOW : _originalBoxShadow;
  }

  int _getOffsetTop(html.Element el) {
    int res = 0;

    while(el != null) {
      res += el.offsetTop;
      el = el.offsetParent;
    }

    return res;
  }

  int _getOffsetLeft(html.Element el) {
    int res = 0;

    while(el != null) {
      res += el.offsetLeft;
      el = el.offsetParent;
    }

    return res;
  }


  void highlight() {
    _isHighlighted = true;

    _mark();

    _addImageDomElement.style
      ..left = (_getOffsetLeft(_domElement) +
                  _domElement.offsetWidth -
                  _DEFAULT_BUTTON_SIZE * 2)
              .toString() +
          "px"
      ..top =
          (_getOffsetTop(_domElement) - _DEFAULT_BUTTON_SIZE / 2).toString() + "px"
      ..display = "block";

    _removeImageDomElement.style
      ..left = (_getOffsetLeft(_domElement) +
                  _domElement.offsetWidth -
                  _DEFAULT_BUTTON_SIZE / 2)
              .toString() +
          "px"
      ..top =
          (_getOffsetTop(_domElement) - _DEFAULT_BUTTON_SIZE / 2).toString() + "px"
      ..display = "block";
  }

  void normalise() {
    _isHighlighted = false;

    _clearMark();

    _addImageDomElement.style..display = "none";

    _removeImageDomElement.style..display = "none";
  }

  void render() {
    
    var image = _domElement as html.ImageElement;

    if(!_hasContent) {
      image.src = _templateSrc;
      image.srcset = _templateSrcset;
      return;
    }

    var suffix = "?" + new DateTime.now().millisecondsSinceEpoch.toString();

    image.src = "./" + _key + "." + _type + suffix;
    StringBuffer sb = new StringBuffer();
    if(_width != null ) {
      _width.forEach((i)=>sb.write("./" + _key + "-" + i.toString() + "w." + _type + suffix +" " + i.toString() + "w,"));
    }
    if(_pixelDensity != null ) {
      _pixelDensity.forEach((f)=>sb.write("./" + _key + "-" + f.toStringAsFixed(1) + "x." + _type + suffix + " " + f.toStringAsFixed(1) + "x,"));
    }
    image.srcset = sb.toString();
  }

  void _stopEvent(html.MouseEvent e) {
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.preventDefault();
  }

  void _addImage(html.MouseEvent e) {
    _stopEvent(e);

    _fileUploadDomElement.style
      ..left = _addImageDomElement.offsetLeft.toString() + "px"
      ..top = _addImageDomElement.offsetTop.toString() + "px";
    _fileUploadDomElement.focus();
    _fileUploadDomElement.click();
  }

  void _removeImage(html.MouseEvent e) {

    _type = "";
    _hasContent = false;

    render();

  _page.save(null, null);

    _stopEvent(e);
  }

  void _uploadFile(html.Event e) {
    final html.File file = _fileUploadDomElement.files.first;

    if (file == null) {
      return;
    }

    final reader = new html.FileReader();
    reader.onLoad.listen((e) {
      _sendDataToServer(file.name, reader.result);
    });
    reader.readAsArrayBuffer(file);
  }

  void _sendDataToServer(String name, dynamic data) {
    final request = new html.HttpRequest();
    request.onReadyStateChange.listen((html.Event e) {
      if (request.readyState != html.HttpRequest.DONE) {
        return;
      }

      if (request.status == 200 || request.status == 0) {
        _processImageUploadSuccess(request);
      } else {
        _processImageUploadError(request);
      }
    });

    var url = html.window.location.href.replaceAll("/!/", "/!image/") + name + "?key=" + _key;
    request.open("POST", url);
    request.send(data);
  }

  void _processImageUploadSuccess(html.HttpRequest request) {
    Map map = convert.JSON.decode(request.responseText);
    
    _type = map[IMAGE_TYPE];
    _width = map[IMAGE_WIDTH];
    _pixelDensity = map[IMAGE_PIXEL_DENSITY];
    _hasContent = true;
    
    render();

    print("upload complete");

    _page.save(null, null);
  }

  void _processImageUploadError(html.HttpRequest request) {
    print("fail");
  }

  void prepareDomForHtmlSave() {
    _fileUploadDomElement.remove();
    _addImageDomElement.remove();
    _removeImageDomElement.remove();
  }

  void restoreDomAfterHtmlSave() {
    html.document.body.children.add(_fileUploadDomElement);
    html.document.body.children.add(_addImageDomElement);
    html.document.body.children.add(_removeImageDomElement);
  }
}

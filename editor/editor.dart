// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
library wedit;

import 'dart:html' as html;
import 'dart:convert' as convert;
import 'dart:svg' as svg;
import 'package:markdown/markdown.dart' as md;
import 'package:html_unescape/html_unescape.dart' as uesc;


part 'dataMap.dart';
part 'element.dart';
part 'image.dart';
part 'repeat.dart';
part 'repeatShadow.dart';
part 'page.dart';
part 'pageMenu.dart';

void main() {
  var url = html.window.location.protocol + "//" + html.window.location.host + "/~?p=" + html.window.location.pathname;
  html.HttpRequest.getString(url).then(onDataLoaded);
}

void onDataLoaded(String responseText) {
  Map map = convert.json.decode(responseText);
  new Page.fromMap(map);
}

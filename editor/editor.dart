// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
library wedit;

import 'dart:convert' as convert;
import 'dart:js_interop';
import 'package:web/web.dart' as web;
import 'package:markdown/markdown.dart' as md;
import 'package:html_unescape/html_unescape.dart' as uesc;

part 'dataMap.dart';
part 'element.dart';
part 'image.dart';
part 'repeat.dart';
part 'repeatShadow.dart';
part 'class.dart';
part 'page.dart';
part 'pageMenu.dart';

bool isDemo = false;

void main() {
  isDemo = web.document.body?.getAttribute("wedit-demo") == "true";
  if (isDemo) {
    final url = "./wedit.json";
    _fetchData(url).then(onDataLoaded);
  } else {
    final url = "${web.window.location.protocol}//${web.window.location.host}/~?p=${web.window.location.pathname}";
    _fetchData(url).then(onDataLoaded);
  }
}

Future<String> _fetchData(String url) async {
  final response = await web.window.fetch(url.toJS).toDart;
  final text = await response.text().toDart;
  return text as String;
}

void onDataLoaded(String responseText) {
  final Map<String, dynamic> map = convert.json.decode(responseText) as Map<String, dynamic>;
  Page.fromMap(map);
}

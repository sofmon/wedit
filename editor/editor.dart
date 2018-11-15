// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
library wedit;

import 'dart:html' as html;
import 'dart:convert' as convert;
import 'dart:async' as async;
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

const PAGE_DATA_TEMPLATE =
    '{"h":"","s":"","p":"","t":"","e":[],"r":[],"s":{"e":"","r":""}}';

void main() {
  String url = getLoadUrl(html.window.location.href);
  html.HttpRequest.getString(url).then(onDataLoaded);
}

String getLoadUrl(String url) {
  var sb = new StringBuffer();

  if (url.startsWith("http://", 0)) {
    sb.write("http://");
    url = url.substring(7);
  } else if (url.startsWith("https://", 0)) {
    sb.write("https://");
    url = url.substring(8);
  }

  List<String> split = url.split("/");
  if(split.length <= 2) {
    print("unable to parse current browser URL. Trying demo mode.");
    return "/!load.json";
  }

  sb.write(split[0]);
  split.removeAt(0); // {doamin}
  split.removeAt(0); // /!/

  sb.write("/!load/");
  sb.write(split.join("/"));
  
  return sb.toString();
}

void onDataLoaded(String responseText) {
  Map map = convert.json.decode(responseText);
  new Page.fromMap(map);
}

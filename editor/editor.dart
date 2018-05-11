// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
library wedit;

import 'dart:html' as html;
import 'dart:convert' as convert;
import 'dart:async' as async;
import 'dart:svg' as svg;
import 'package:markdown/markdown.dart' as md;

part 'dataMap.dart';
part 'element.dart';
part 'image.dart';
part 'repeat.dart';
part 'repeatShadow.dart';
part 'page.dart';
part 'message.dart';

const PAGE_DATA_TEMPLATE =
    '{"h":"","s":"","p":"","t":"","e":[],"r":[],"s":{"e":"","r":""}}';

void main() {
  Map map = convert.JSON.decode(PAGE_DATA_TEMPLATE);
  new Page.fromMap(map);
}

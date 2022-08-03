import 'package:localstorage/localstorage.dart';

class AppLocalStorage {
  static final LocalStorage _singleton = LocalStorage("simple_poker_game");

  static Future<void> setItem(String key, dynamic value) async {
    if (await _singleton.ready) {
      await _singleton.setItem(key, value);
    }
  }

  static dynamic getItem(String key) {
    return _singleton.getItem(key);
  }
}

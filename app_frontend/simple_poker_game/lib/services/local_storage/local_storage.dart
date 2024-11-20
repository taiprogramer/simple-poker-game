import 'package:localstorage/localstorage.dart';

class AppLocalStorage {
  static Future<void> setItem(String key, dynamic value) async {
    localStorage.setItem(key, value);
  }

  static dynamic getItem(String key) {
    return localStorage.getItem(key);
  }
}

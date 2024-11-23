import 'package:localstorage/localstorage.dart';

class AppLocalStorage {
  static void setItem(String key, dynamic value) {
    localStorage.setItem(key, value.toString());
  }

  static dynamic getItem(String key) {
    return localStorage.getItem(key);
  }
}

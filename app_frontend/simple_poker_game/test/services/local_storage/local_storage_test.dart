import 'package:flutter_test/flutter_test.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';

void main() {
  test('LocalStorage can be set & get value', () {
    const key = 'key';
    const value = 'value';
    AppLocalStorage.setItem(key, value);
    expect(AppLocalStorage.getItem(key), value);
  });
}

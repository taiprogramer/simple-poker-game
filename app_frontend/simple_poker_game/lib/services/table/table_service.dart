import 'dart:convert';

import 'package:simple_poker_game/models/table.dart';
import 'package:simple_poker_game/services/compact_http_client.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';

class TableService {
  static const String _tableEndpoint = '/table';
  static Future<PokerTable> getTable(int tableID, int userID) async {
    String accessToken = AppLocalStorage.getItem('access_token');
    final res = await CompactHttpClient.get(
        '/$tableID?userID=$userID', _tableEndpoint, accessToken);
    String stringData = await res.transform(utf8.decoder).join();
    dynamic body = jsonDecode(stringData);
    PokerTable table = PokerTable.fromMap(body);
    return table;
  }

  static Future<PokerTable> performAction(
      int tableID, int userID, String action, int amount) async {
    String accessToken = AppLocalStorage.getItem('access_token');
    final payload =
        '{"user_id": $userID, "action": "$action", "amount": $amount}';
    final res = await CompactHttpClient.post(
        payload, '$_tableEndpoint/$tableID', accessToken);
    String stringData = await res.transform(utf8.decoder).join();
    dynamic body = jsonDecode(stringData);
    PokerTable table = PokerTable.fromMap(body);
    return table;
  }
}

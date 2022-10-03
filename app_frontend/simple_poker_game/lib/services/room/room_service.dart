import 'dart:convert';

import 'package:simple_poker_game/services/local_storage/local_storage.dart';

import '../../models/room.dart';
import '../compact_http_client.dart';

class RoomService {
  static const String _roomEndPoint = '/room';

  static Future<List<Room>> listRoom(int offset, int limit) async {
    final res = await CompactHttpClient.get(
        '?offset=$offset&limit=$limit', _roomEndPoint);
    String stringData = await res.transform(utf8.decoder).join();
    List<dynamic> body = jsonDecode(stringData);
    List<Room> rooms = body.map((item) => Room.fromMap(item)).toList();
    return rooms;
  }

  static Future<Room> newRoom({int userID = 0, String password = ''}) async {
    String accessToken = AppLocalStorage.getItem('access_token');
    final res = await CompactHttpClient.post(
        '{"user_id": $userID, "password": "$password"}',
        _roomEndPoint,
        accessToken);
    String stringData = await res.transform(utf8.decoder).join();
    dynamic body = jsonDecode(stringData);
    Room room = Room.fromMap(body);
    return room;
  }

  static Future<Room> getRoom({int roomID = 0}) async {
    String accessToken = AppLocalStorage.getItem('access_token');
    final res =
        await CompactHttpClient.get('/$roomID', _roomEndPoint, accessToken);
    String stringData = await res.transform(utf8.decoder).join();
    dynamic body = jsonDecode(stringData);
    Room room = Room.fromMap(body);
    return room;
  }

  static Future<Room> joinRoom(
      {int roomID = 0, int userID = 0, int money = 0}) async {
    String accessToken = AppLocalStorage.getItem('access_token');
    final res = await CompactHttpClient.post(
        '{"user_id": $userID, "money": $money}',
        '$_roomEndPoint/$roomID',
        accessToken);
    String stringData = await res.transform(utf8.decoder).join();
    dynamic body = jsonDecode(stringData);
    Room room = Room.fromMap(body);
    return room;
  }
}

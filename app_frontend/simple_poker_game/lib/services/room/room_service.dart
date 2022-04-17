import 'dart:convert';

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
}

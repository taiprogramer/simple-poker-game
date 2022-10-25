class _User {
  final int id;
  final bool ready;
  final String username;

  _User({this.id = 0, this.ready = false, this.username = ''});

  Map<String, dynamic> toJson() =>
      {'id': id, 'ready': ready, 'username': username};

  factory _User.fromMap(Map data) {
    return _User(
        id: data['id'], ready: data['ready'], username: data['username']);
  }
}

class Room {
  final int id;
  final String code;
  final bool playing;
  final bool private;
  final int owner;
  final int table;
  final List<_User> users;

  Room(
      {this.id = 0,
      this.code = '',
      this.playing = false,
      this.private = false,
      this.owner = 0,
      this.table = 0,
      this.users = const []});
  Map<String, dynamic> toJson() => {
        'id': id,
        'code': code,
        'playing': playing,
        'private': private,
        'owner': owner,
        'table': table,
        'users': users
      };

  factory Room.fromMap(Map data) {
    return Room(
        id: data['id'],
        code: data['code'],
        playing: data['playing'],
        private: data['private'],
        owner: data['owner'],
        table: data['table'],
        users: convertListDynamicToUsers(data['users']));
  }

  static List<_User> convertListDynamicToUsers(List<dynamic> usersData) {
    List<_User> users = List.empty(growable: true);
    for (var userData in usersData) {
      users.add(_User.fromMap(userData));
    }
    return users;
  }
}

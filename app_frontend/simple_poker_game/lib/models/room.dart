class Room {
  final int id;
  final String code;
  final bool playing;
  final bool private;
  final int owner;
  final int table;

  Room(
      {this.id = 0,
      this.code = '',
      this.playing = false,
      this.private = false,
      this.owner = 0,
      this.table = 0});
  Map<String, dynamic> toJson() => {
        'id': id,
        'code': code,
        'playing': playing,
        'private': private,
        'owner': owner,
        'table': table
      };

  factory Room.fromMap(Map data) {
    return Room(
        id: data['id'],
        code: data['code'],
        playing: data['playing'],
        private: data['private'],
        owner: data['owner'],
        table: data['table']);
  }
}

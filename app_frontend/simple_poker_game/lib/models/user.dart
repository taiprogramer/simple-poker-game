class User {
  final int id;
  final String username;
  final int money;

  User({this.id = 0, this.username = '', this.money = 0});
  Map<String, dynamic> toJson() =>
      {'id': id, 'username': username, 'money': money};

  factory User.fromMap(Map data) {
    return User(
        id: data['id'], username: data['username'], money: data['money']);
  }
}

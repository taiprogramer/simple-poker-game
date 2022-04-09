import 'dart:convert';
import 'dart:io';

import '../../models/user.dart';
import '../config.dart';

class UserCredential {
  final String username;
  final String password;

  UserCredential({this.username = '', this.password = ''});
}

class AuthService {
  static final String host = ServiceConfig.getHost();
  static final int port = ServiceConfig.getPort();
  static const String userEndPoint = '/user';
  static final HttpClient http = HttpClient();

  static Future<User> signUp(UserCredential credential) async {
    HttpClientRequest request = await http.post(host, port, userEndPoint);

    request.headers.add(HttpHeaders.contentTypeHeader, 'application/json');
    request.write(
        '{"username": "${credential.username}", "password": "${credential.password}"}');
    final response = await request.close();
    List<String> l = await response.transform(utf8.decoder).toList();

    if (response.statusCode == HttpStatus.badRequest) {
      throw Exception(l.elementAt(0));
    }

    return User.fromMap(json.decode(l.elementAt(0)));
  }
}

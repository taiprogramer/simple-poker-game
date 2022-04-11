import 'dart:convert';
import 'dart:io';

import '../../models/user.dart';
import '../compact_http_client.dart';

class UserCredential {
  final String username;
  final String password;

  UserCredential({this.username = '', this.password = ''});
}

class AuthService {
  static const String _userEndPoint = '/user';
  static const String _authEndPoint = '/auth';

  static Future<User> signUp(UserCredential credential) async {
    final response = await CompactHttpClient.post(
        '{"username": "${credential.username}", "password": "${credential.password}"}',
        _userEndPoint);
    List<String> l = await response.transform(utf8.decoder).toList();

    if (response.statusCode == HttpStatus.badRequest) {
      List<dynamic> errors = json.decode(l.elementAt(0))['error_messages'];
      throw Exception(errors.elementAt(0));
    }

    return User.fromMap(json.decode(l.elementAt(0)));
  }

  // return: access_token (JWT format) when succeed
  // throw: Exception when fail
  static Future<String> signIn(UserCredential credential) async {
    final response = await CompactHttpClient.post(
        '{"username": "${credential.username}", "password": "${credential.password}"}',
        _authEndPoint);
    List<String> l = await response.transform(utf8.decoder).toList();

    if (response.statusCode == HttpStatus.badRequest) {
      List<dynamic> errors = json.decode(l.elementAt(0))['error_messages'];
      throw Exception(errors.elementAt(0));
    }
    return json.decode(l.elementAt(0))['access_token'];
  }
}

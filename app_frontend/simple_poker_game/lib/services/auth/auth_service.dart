import 'dart:convert';
import 'dart:io';

import '../../models/user.dart';
import '../compact_http_client.dart';

class UserCredential {
  final String username;
  final String password;

  UserCredential({this.username = '', this.password = ''});
}

class Authentication {
  final String accessToken;
  final int userID;

  Authentication({this.accessToken = '', this.userID = 0});
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

  // return: Authentication object when succeed
  // throw: Exception when fail
  static Future<Authentication> signIn(UserCredential credential) async {
    final response = await CompactHttpClient.post(
        '{"username": "${credential.username}", "password": "${credential.password}"}',
        _authEndPoint);
    List<String> l = await response.transform(utf8.decoder).toList();

    if (response.statusCode == HttpStatus.badRequest) {
      List<dynamic> errors = json.decode(l.elementAt(0))['error_messages'];
      throw Exception(errors.elementAt(0));
    }

    String accessToken = json.decode(l.elementAt(0))['access_token'];
    int userID = json.decode(l.elementAt(0))['user_id'];

    Authentication authentication =
        Authentication(userID: userID, accessToken: accessToken);
    return authentication;
  }
}

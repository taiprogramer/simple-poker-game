import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:simple_poker_game/pages/sign_in_page.dart';
import 'package:simple_poker_game/pages/texas_holdem_page.dart';
import 'package:simple_poker_game/services/auth/auth_service.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';
import 'package:simple_poker_game/shared/string_utils/string_utils.dart';

class SignUpForm extends StatefulWidget {
  const SignUpForm({Key? key}) : super(key: key);

  @override
  SignUpFormState createState() {
    return SignUpFormState();
  }
}

class SignUpFormState extends State<SignUpForm> {
  final _formKey = GlobalKey<FormState>();

  String username = '';
  String password = '';
  bool isLoading = false;

  @override
  Widget build(BuildContext context) {
    return isLoading
        ? const Center(
            child: SizedBox(
                width: 55,
                height: 55,
                child: CircularProgressIndicator(
                  color: Colors.pink,
                  strokeWidth: 6,
                )))
        : Form(
            key: _formKey,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text(
                  'Sign up',
                  style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 36.0,
                      color: Colors.pink),
                ),
                Container(
                    margin: const EdgeInsets.only(top: 10.0),
                    child: TextFormField(
                      onSaved: (String? value) {
                        username = value ?? "";
                      },
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter your username';
                        }
                        return null;
                      },
                      decoration: const InputDecoration(
                          labelText: 'Username',
                          icon: Icon(
                            Icons.person,
                            color: Colors.pinkAccent,
                          )),
                    )),
                Container(
                    margin: const EdgeInsets.only(top: 10.0, bottom: 10.0),
                    child: TextFormField(
                      obscureText: true,
                      onSaved: (String? value) {
                        password = value ?? "";
                      },
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter your password';
                        }
                        return null;
                      },
                      decoration: const InputDecoration(
                          labelText: 'Password',
                          icon: Icon(
                            Icons.lock,
                            color: Colors.pinkAccent,
                          )),
                    )),
                SizedBox(
                    width: 200,
                    child: ElevatedButton(
                        onPressed: () async {
                          if (_formKey.currentState!.validate()) {
                            _formKey.currentState!.save();
                            try {
                              final user = await AuthService.signUp(
                                  UserCredential(
                                      username: username, password: password));
                              try {
                                setState(() {
                                  isLoading = true;
                                });
                                // call login api for auto login
                                final authentication = await AuthService.signIn(
                                    UserCredential(
                                        username: user.username,
                                        password: password));
                                AppLocalStorage.setItem(
                                    'access_token', authentication.accessToken);
                                AppLocalStorage.setItem(
                                    'user_id', authentication.userID);
                                await Future.delayed(
                                    const Duration(microseconds: 500));
                                Navigator.pushReplacementNamed(
                                    context, TexasHoldemPage.routeName);
                              } catch (_) {
                                Navigator.pushReplacementNamed(
                                    context, SignInPage.routeName);
                              }
                            } catch (e) {
                              // remove 'Exception: ' in message
                              final message = StringUtils.cleanExceptionMessage(
                                  e.toString());
                              Fluttertoast.showToast(
                                  msg: message, gravity: ToastGravity.CENTER);
                            } finally {
                              setState(() {
                                isLoading = false;
                              });
                            }
                          }
                        },
                        style: ButtonStyle(
                            shape: MaterialStateProperty.all<
                                    RoundedRectangleBorder>(
                                RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20)))),
                        child: const Text('Sign up')))
              ],
            ),
          );
  }
}

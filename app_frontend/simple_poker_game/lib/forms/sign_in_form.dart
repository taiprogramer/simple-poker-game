import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:simple_poker_game/pages/texas_holdem_page.dart';
import 'package:simple_poker_game/services/auth/auth_service.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';
import 'package:simple_poker_game/shared/string_utils/string_utils.dart';

class SignInForm extends StatefulWidget {
  const SignInForm({Key? key}) : super(key: key);

  @override
  SignInFormState createState() {
    return SignInFormState();
  }
}

class SignInFormState extends State<SignInForm> {
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
                Container(
                    margin: const EdgeInsets.only(top: 10.0),
                    child: TextFormField(
                      onSaved: (String? value) {
                        username = value ?? '';
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
                    margin: const EdgeInsets.only(top: 10.0),
                    child: TextFormField(
                      onSaved: (String? value) {
                        password = value ?? '';
                      },
                      obscureText: true,
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
                              setState(() {
                                isLoading = true;
                              });
                              // must be saved in somewhere (on local machine)
                              final authentication = await AuthService.signIn(
                                  UserCredential(
                                      username: username, password: password));
                              await Future.delayed(
                                  const Duration(microseconds: 500));
                              AppLocalStorage.setItem(
                                  'access_token', authentication.accessToken);
                              AppLocalStorage.setItem(
                                  'user_id', authentication.userID);
                              Navigator.pushNamed(
                                  context, TexasHoldemPage.routeName);
                            } catch (e) {
                              Fluttertoast.showToast(
                                  msg: StringUtils.cleanExceptionMessage(
                                      e.toString()));
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
                        child: const Text('Sign in')))
              ],
            ),
          );
  }
}

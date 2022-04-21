import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:simple_poker_game/pages/texas_holdem_page.dart';
import 'package:simple_poker_game/services/auth/auth_service.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';

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

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text(
            'Sign in',
            style: TextStyle(fontWeight: FontWeight.bold, fontSize: 36.0),
          ),
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
                    border: OutlineInputBorder(),
                    labelText: 'Enter your username'),
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
                    border: OutlineInputBorder(),
                    labelText: 'Enter your password'),
              )),
          ElevatedButton(
              onPressed: () async {
                if (_formKey.currentState!.validate()) {
                  _formKey.currentState!.save();
                  try {
                    // must be saved in somewhere (on local machine)
                    final accessToken = await AuthService.signIn(
                        UserCredential(username: username, password: password));
                    AppLocalStorage.setItem('access_token', accessToken);
                    Navigator.pushReplacementNamed(
                        context, TexasHoldemPage.routeName);
                  } catch (e) {
                    Fluttertoast.showToast(
                        msg: e.toString(), gravity: ToastGravity.CENTER);
                  }
                }
              },
              child: const Text('Sign in'))
        ],
      ),
    );
  }
}

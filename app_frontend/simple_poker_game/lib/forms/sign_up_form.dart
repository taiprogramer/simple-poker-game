import 'package:flutter/material.dart';
import 'package:simple_poker_game/pages/sign_in_page.dart';
import 'package:simple_poker_game/services/auth/auth_service.dart';

class SignUpForm extends StatefulWidget {
  const SignUpForm({Key? key}) : super(key: key);

  @override
  SignUpFormState createState() {
    return SignUpFormState();
  }
}

class SignUpFormState extends State<SignUpForm> {
  final _formKey = GlobalKey<FormState>();

  String username = "";
  String password = "";

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text(
            'Sign up',
            style: TextStyle(fontWeight: FontWeight.bold, fontSize: 36.0),
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
                    border: OutlineInputBorder(),
                    labelText: 'Enter your username'),
              )),
          Container(
              margin: const EdgeInsets.only(top: 10.0),
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
                    border: OutlineInputBorder(),
                    labelText: 'Enter your password'),
              )),
          ElevatedButton(
              onPressed: () async {
                if (_formKey.currentState!.validate()) {
                  _formKey.currentState!.save();
                  try {
                    await AuthService.signUp(
                        UserCredential(username: username, password: password));
                    Navigator.pushNamed(context, SignInPage.routeName);
                  } catch (e) {
                    // Toast here
                  }
                }
              },
              child: const Text('Sign up'))
        ],
      ),
    );
  }
}

import 'package:flutter/material.dart';

import '../forms/sign_in_form.dart';

class SignInPage extends StatelessWidget {
  static const String routeName = '/signIn';

  const SignInPage({Key? key, this.title = 'Simple Poker Game'})
      : super(key: key);

  final String title;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Container(
            padding: const EdgeInsets.only(left: 10.0, right: 10.0),
            child: const SignInForm()));
  }
}

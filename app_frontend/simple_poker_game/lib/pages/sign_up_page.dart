import 'package:flutter/material.dart';

import '../forms/sign_up_form.dart';

class SignUpPage extends StatelessWidget {
  static const String routeName = '/signUp';
  const SignUpPage({Key? key, this.title = 'Simple Poker Game'})
      : super(key: key);

  final String title;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: Text(title)),
        body: Container(
            padding: const EdgeInsets.all(20.0), child: const SignUpForm()));
  }
}

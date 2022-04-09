import 'package:flutter/material.dart';

class SignInPage extends StatelessWidget {
  static const String routeName = '/signIn';

  const SignInPage({Key? key, this.title = 'Simple Poker Game'})
      : super(key: key);

  final String title;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(title),
      ),
      body: const Center(child: Text('Login Page works!')),
    );
  }
}

import 'package:flutter/material.dart';

import 'sign_in_page.dart';
import 'sign_up_page.dart';

class HomePage extends StatelessWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: const Text('Simple Poker Game')),
        body: Center(
            child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ElevatedButton(
                onPressed: () {
                  Navigator.pushNamed(context, SignUpPage.routeName);
                },
                child: const Text('Sign up')),
            ElevatedButton(
                onPressed: () {
                  Navigator.pushNamed(context, SignInPage.routeName);
                },
                child: const Text('Sign in')),
          ],
        )));
  }
}

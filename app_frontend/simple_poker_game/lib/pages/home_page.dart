import 'package:flutter/material.dart';
import 'package:simple_poker_game/pages/setting_page.dart';

import 'sign_in_page.dart';
import 'sign_up_page.dart';

class HomePage extends StatelessWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar:
            AppBar(title: const Text('Simple Poker Game'), actions: <Widget>[
          IconButton(
            onPressed: () {
              Navigator.pushNamed(context, SettingPage.routeName);
            },
            icon: const Icon(Icons.settings, color: Colors.white),
          )
        ]),
        body: Center(
            child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SizedBox(
              width: 200,
              height: 45,
              child: ElevatedButton(
                  style: ButtonStyle(
                      shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                          RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(20)))),
                  onPressed: () {
                    Navigator.pushNamed(context, SignInPage.routeName);
                  },
                  child: const Text('Sign in')),
            ),
            Container(
              padding: const EdgeInsets.all(5),
            ),
            SizedBox(
              width: 200,
              height: 45,
              child: OutlinedButton(
                  style: ButtonStyle(
                      shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                          RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(20)))),
                  onPressed: () {
                    Navigator.pushNamed(context, SignUpPage.routeName);
                  },
                  child: const Text('Create new account')),
            ),
          ],
        )));
  }
}

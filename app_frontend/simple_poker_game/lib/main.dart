import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:simple_poker_game/pages/texas_holdem_page.dart';

import 'pages/home_page.dart';
import 'pages/sign_in_page.dart';
import 'pages/sign_up_page.dart';

Future main() async {
  await dotenv.load(fileName: '.env');
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    var routes = <String, WidgetBuilder>{
      SignInPage.routeName: (BuildContext context) => const SignInPage(),
      SignUpPage.routeName: (BuildContext context) => const SignUpPage(),
      TexasHoldemPage.routeName: (BuildContext context) =>
          const TexasHoldemPage(),
    };
    return MaterialApp(
      title: 'Simple Poker Game',
      theme: ThemeData(
        primarySwatch: Colors.pink,
      ),
      home: const HomePage(),
      routes: routes,
    );
  }
}

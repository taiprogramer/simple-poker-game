import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    var routes = <String, WidgetBuilder>{
      SignInPage.routeName: (BuildContext context) => const SignInPage(),
      SignUpPage.routeName: (BuildContext context) => const SignUpPage(),
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

class SignUpPage extends StatelessWidget {
  static const String routeName = '/signUp';
  const SignUpPage({Key? key, this.title = 'Simple Poker Game'})
      : super(key: key);

  final String title;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: Text(title)),
        body: const Center(child: Text('Sign up page works!')));
  }
}

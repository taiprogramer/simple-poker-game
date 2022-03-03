import 'package:flutter/material.dart';

void main() {
  runApp(const Center(child: Text("Hello, Flutter")));
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Simple Poker Game'),
      ),
      body: const Center(child: Text("Hello, Flutter")),
    );
  }
}

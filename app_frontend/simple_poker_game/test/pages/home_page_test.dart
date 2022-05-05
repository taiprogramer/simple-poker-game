import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:simple_poker_game/pages/home_page.dart';

void main() {
  testWidgets('Homepage should display sign up & sign in button',
      (WidgetTester tester) async {
    await tester.pumpWidget(const MaterialApp(home: HomePage()));
    expect(find.text('Sign in'), findsOneWidget);
    expect(find.text('Create new account'), findsOneWidget);
  });
}

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simple_poker_game/pages/room_texas_holdem_page.dart';
import 'package:simple_poker_game/pages/texas_holdem_page.dart';
import 'package:simple_poker_game/services/config.dart';

import 'pages/setting_page.dart';
import 'pages/home_page.dart';
import 'pages/sign_in_page.dart';
import 'pages/sign_up_page.dart';

Future main() async {
  await dotenv.load(fileName: '.env');
  await init();
  WidgetsFlutterBinding.ensureInitialized();
  SystemChrome.setPreferredOrientations([DeviceOrientation.landscapeLeft])
      .then((_) {
    runApp(const MyApp());
  });
}

Future init() async {
  // get host config
  // search precedence: shared preferences -> .env -> dummy value
  final serviceConfig = ServiceConfig.getInstance();
  final prefs = await SharedPreferences.getInstance();
  var host =
      prefs.getString('host') ?? dotenv.env['API_SERVER_HOST'] ?? '10.10.10.10';
  var port = prefs.getInt('port') ??
      int.parse(dotenv.env['API_SERVER_PORT'] ?? '1975');
  serviceConfig.setHost(host);
  serviceConfig.setPort(port);
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
      RoomTexasHoldemPage.routeName: (BuildContext context) =>
          const RoomTexasHoldemPage(),
      SettingPage.routeName: (BuildContext context) => const SettingPage(),
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

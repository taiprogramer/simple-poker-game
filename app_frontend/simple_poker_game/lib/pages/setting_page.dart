import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simple_poker_game/services/config.dart';

class SettingPage extends StatefulWidget {
  static const String routeName = '/setting';
  const SettingPage({Key? key}) : super(key: key);

  @override
  State<SettingPage> createState() => _SettingPageState();
}

class _SettingPageState extends State<SettingPage> {
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  TextEditingController hostController = TextEditingController();
  TextEditingController portController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    hostController.text = ServiceConfig.getInstance().getHost();
    portController.text = ServiceConfig.getInstance().getPort().toString();
    return Scaffold(
        body: Container(
            padding: const EdgeInsets.only(top: 50.0, right: 25.0, left: 25.0),
            child: Column(children: [
              const Text('Settings',
                  style: TextStyle(color: Colors.pink, fontSize: 36.0)),
              Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: <Widget>[
                    TextFormField(
                      controller: hostController,
                      onSaved: (String? host) async {
                        // save info to shared_preferences
                        host = host ?? '';
                        final prefs = await SharedPreferences.getInstance();
                        prefs.setString('host', host);
                        // update ServiceConfig
                        ServiceConfig.getInstance().setHost(host);
                      },
                      decoration: const InputDecoration(
                        hintText: 'Enter the host: 10.10.10.10',
                      ),
                      validator: (String? value) {
                        if (value!.isEmpty) {
                          return 'The host can not be empty';
                        }
                        return null;
                      },
                    ),
                    TextFormField(
                      controller: portController,
                      onSaved: (String? port) async {
                        // save info to shared_preferences
                        port = port ?? '0';
                        final portInt = int.parse(port);
                        final prefs = await SharedPreferences.getInstance();
                        prefs.setInt('port', portInt);
                        // update ServiceConfig
                        ServiceConfig.getInstance().setPort(portInt);
                      },
                      decoration: const InputDecoration(
                        hintText: 'Enter the port: 1975',
                      ),
                      validator: (String? value) {
                        if (value!.isEmpty) {
                          return 'The host can not be empty';
                        }
                        return null;
                      },
                    ),
                    Center(
                        child: ElevatedButton(
                      onPressed: () {
                        if (_formKey.currentState!.validate()) {
                          _formKey.currentState!.save();
                        }
                      },
                      child: const Text('Save'),
                    )),
                  ],
                ),
              )
            ])));
  }
}

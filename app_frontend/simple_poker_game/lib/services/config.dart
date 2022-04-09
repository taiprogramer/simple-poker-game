import 'package:flutter_dotenv/flutter_dotenv.dart';

class ServiceConfig {
  static String getHost() {
    return dotenv.env['API_SERVER_HOST'] ?? '192.168.1.10';
  }

  static int getPort() {
    return int.parse(dotenv.env['API_SERVER_PORT'] ?? '3000');
  }
}

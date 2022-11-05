class ServiceConfig {
  late String host;
  late int port;
  static final ServiceConfig _singleton = ServiceConfig();

  static ServiceConfig getInstance() {
    return _singleton;
  }

  void setHost(String host) {
    this.host = host;
  }

  void setPort(int port) {
    this.port = port;
  }

  String getHost() {
    return host;
  }

  int getPort() {
    return port;
  }
}

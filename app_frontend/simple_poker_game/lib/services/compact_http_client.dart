// Wrapper of 'dart:io' HttpClient
// provide: simple static methods with default value for common parameters
import 'dart:io';

import 'config.dart';

class CompactHttpClient {
  static final HttpClient _http = HttpClient();
  static final String _host = ServiceConfig.getHost();
  static final int _port = ServiceConfig.getPort();

  static Future<HttpClientResponse> post(String body, String endPoint,
      [String accessToken = '']) async {
    HttpClientRequest request = await _http.post(_host, _port, endPoint);
    request.headers.add(HttpHeaders.contentTypeHeader, 'application/json');
    request.headers.add(HttpHeaders.authorizationHeader, 'Bearer $accessToken');
    request.write(body);
    final response = await request.close();
    return response;
  }

  static Future<HttpClientResponse> put(String body, String endPoint,
      [String accessToken = '']) async {
    HttpClientRequest request = await _http.put(_host, _port, endPoint);
    request.headers.add(HttpHeaders.contentTypeHeader, 'application/json');
    request.headers.add(HttpHeaders.authorizationHeader, 'Bearer $accessToken');
    request.write(body);
    final response = await request.close();
    return response;
  }

  static Future<HttpClientResponse> get(String query, String endPoint,
      [String accessToken = '']) async {
    String url = endPoint + query;
    HttpClientRequest request = await _http.get(_host, _port, url);
    return await request.close();
  }
}

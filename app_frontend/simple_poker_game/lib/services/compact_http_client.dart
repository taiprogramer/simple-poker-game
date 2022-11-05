// Wrapper of 'dart:io' HttpClient
// provide: simple static methods with default value for common parameters
import 'dart:io';

import 'config.dart';

class CompactHttpClient {
  static final HttpClient _http = HttpClient();
  static final String _host = ServiceConfig.getInstance().getHost();
  static final int _port = ServiceConfig.getInstance().getPort();

  static Future<HttpClientResponse> post(String body, String endPoint,
      [String accessToken = '']) async {
    HttpClientRequest request = await _http.post(
        ServiceConfig.getInstance().getHost(),
        ServiceConfig.getInstance().getPort(),
        endPoint);
    request.headers.add(HttpHeaders.contentTypeHeader, 'application/json');
    request.headers.add(HttpHeaders.authorizationHeader, 'Bearer $accessToken');
    request.write(body);
    final response = await request.close();
    return response;
  }

  static Future<HttpClientResponse> put(String body, String endPoint,
      [String accessToken = '']) async {
    HttpClientRequest request = await _http.put(
        ServiceConfig.getInstance().getHost(),
        ServiceConfig.getInstance().getPort(),
        endPoint);
    request.headers.add(HttpHeaders.contentTypeHeader, 'application/json');
    request.headers.add(HttpHeaders.authorizationHeader, 'Bearer $accessToken');
    request.write(body);
    final response = await request.close();
    return response;
  }

  static Future<HttpClientResponse> get(String query, String endPoint,
      [String accessToken = '']) async {
    String url = endPoint + query;
    HttpClientRequest request = await _http.get(
        ServiceConfig.getInstance().getHost(),
        ServiceConfig.getInstance().getPort(),
        url);
    request.headers.add(HttpHeaders.contentTypeHeader, 'application/json');
    request.headers.add(HttpHeaders.authorizationHeader, 'Bearer $accessToken');
    return await request.close();
  }
}

// contains useful functions for manipulate string
class StringUtils {
  // Exception message format: 'Exception: %s'
  // After cleaning process: '%s'
  static String cleanExceptionMessage(String exceptionMessage) {
    return exceptionMessage.replaceFirst('Exception: ', '');
  }
}

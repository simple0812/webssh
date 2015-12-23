function setTop() {
  var top = (get_text_rows('#messages') - 1) * 13;
  $('form').css('top', Math.min(top, $('#messages').get(0).clientHeight - 13));
}

function get_text_rows(ta) {
  return $(ta).val().split("\n").length;
}

function get_last_cols(ta) {
  var p = $(ta).val().split("\n");
  return p[p.length - 1];
}

function getQueryString(key) {
  var value = "";
  var sURL = window.document.URL;

  if (sURL.indexOf("?") > 0) {
    var arrayParams = sURL.split("?");
    var arrayURLParams = arrayParams[1].split("&");

    for (var i = 0; i < arrayURLParams.length; i++) {
      var sParam = arrayURLParams[i].split("=");

      if ((sParam[0] == key) && (sParam[1] != "")) {
        value = sParam[1];
        break;
      }
    }
  }
  return value;
}
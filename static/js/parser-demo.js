var beforePredict = function () {
  $('#result').show();
  $('.result-container').hide();
  $('.waiting-message').show();
  $('.error-message').hide();
  $('#predict').blur();
}

var request = function (url, success) {
  $.ajax(url).done(function (data) {
    $('.waiting-message').hide();
    success(data);
  }).fail(function () {
    $('.waiting-message').hide();
    $('.error-message').show();
  });
};

var doParse = function () {
  beforePredict();

  $('#predict').text($(this).text()).unbind().click(doParse);
  var sentenceText = $('#sentence').val();
  sentenceText = sentenceText.replace(/[><]/, '');
  if (sentenceText == '') return ;

  var treeUrl = 'tree2svg?q=' + encodeURIComponent(sentenceText);
  request(treeUrl, function (data) {
    $('.tree-container').show();
    $('.tree').html(data);
  });
};

var doPredict = function (textMapFunc) {
  beforePredict();

  var sentenceText = $('#sentence').val();
  sentenceText = sentenceText.replace(/[><]/, '');
  if (sentenceText == '') return ;

  var predictUrl = 'predict?q=' + encodeURIComponent(sentenceText);
  request(predictUrl, function (data) {
    var jsonObj = JSON.parse(data);
    var text = _.map(jsonObj, textMapFunc);
    $('.text-container').show().html(text.join('&nbsp; '));
  });
};

var doSeg = function() {
  $('#predict').text($(this).text()).unbind().click(doSeg);
  doPredict(function (item) {
    return item['word'];
  });
};

var doPosTag = function() {
  $('#predict').text($(this).text()).unbind().click(doPosTag);
  doPredict(function (item) {
    return item['word'] + '/' + item['postag'];
  });
};

$(function () {
  $('.btn').removeClass('disabled');
  $('#predict').click(doParse);
  $('#dr-seg').click(doSeg);
  $('#dr-postag').click(doPosTag);
  $('#dr-depparse').click(doParse);
  $('#sentence').keypress(function (event) {
    if (event.which == 13) $('#predict').click();
  });
});

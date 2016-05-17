$(function() {

  function formatTimestamp(unix_timestamp) {
    var date = new Date(unix_timestamp*1000)

    var year = date.getFullYear()
    var month = "0" + date.getMonth()
    var hours = "0" + date.getHours()
    var day = "0" + date.getDate()
    var minutes = "0" + date.getMinutes()
    var seconds = "0" + date.getSeconds()

    function dd(overlong) { return overlong.substr(-2) }

    return year + "-" + dd(month) + "-" + dd(day) + " " + dd(hours) + ':' + dd(minutes) + ':' + dd(seconds)

  }

  $.getJSON($('#baseUrl').val() + 'api/logs', function(entries) {
    $('#loading-logs').hide()

    if (!entries) {
      $('#no-logs').removeClass('hidden')
    } else {
      $('#available-logs').removeClass('hidden')
    }
    entries.reverse()
    entries.forEach(function(entry) {
      var $li = $('<li></li>').appendTo($('#available-logs'))
      $li.append(
        '<span class="date">' + formatTimestamp(entry.Timestamp) + '</span>'
                 + ' <a href="' + $('#baseUrl').val() + 'log/' + entry.PushId + '">' + entry.Title + '</a>'
                 + ' <span class="cmd"><pre>' + entry.Cmd + '</pre></span>'
      )
    })
  })
})

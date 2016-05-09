$(function() {
  function makeButtonPushable($btn) {
    $btn.click(function(event) {
      var $button = $(event.target)
      $button.addClass("running")
      $.post('push/' + $button.data('buttonid'), {}, null, "json"
      ).done(function(data) {
        console.log(data)
        alert(data.buttonId + " was pressed! -> " + data.pushId)
      }).fail(function(error) {
        console.log(error)
        alert(error)
        alert(error.responseText)
      }).always(function() {
        $button.removeClass("running")
      })
    })
  }

  $.getJSON('buttons', function(buttons) {
    $('#loading-buttons').hide()

    if (!buttons) {
      $('#no-buttons').removeClass('hidden')
    } else {
      $('#available-buttons').removeClass('hidden')
    }
    buttons.forEach(function(button) {
      var $button = $('<li></li>')
                      .text(button.Title)
                      .data('buttonid', button.Id)
                      .appendTo($('#available-buttons'))
      makeButtonPushable($button)
    })
  })
})

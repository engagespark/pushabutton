$(function() {
  function closeParameterModal() {
    $('#parameterModal').hide()
  }

  function showParameterModal($button, buttonDef, sendFunc) {

    var $modal = $('#parameterModal');

    $('#parameterModal ol').empty()
    buttonDef.Parameters.forEach(function(parameterDef) {
      $('<li>' + parameterDef.Name + '</li>').appendTo($('#parameterModal ol'))
    })

    $('#btnPushWithParameters').unbind('click').click(function() {sendFunc('text1')})
    $('#btnCancelParameters').unbind('click').click(closeParameterModal)

    $modal.removeClass('hidden')
    $modal.show()
  }

  function makeButtonPushable($button, buttonDef) {
    function pushFunc(pushArguments) {
      $.post('push/' + $button.data('buttonid'), pushArguments || {}, null, "json"
      ).done(function(data) {
        alert(data.buttonId + " was pressed! -> " + data.pushId)
      }).fail(function(error) {
        alert(error)
        alert(error.responseText)
      }).always(function() {
        $button.removeClass("running")
        closeParameterModal()
      })
    }

    $button.click(function(event) {
      $button.addClass("running")
      if (buttonDef.Parameters && buttonDef.Parameters.length > 0) {
        showParameterModal($button, buttonDef, pushFunc)
      } else {
        pushFunc()
      }
    })
  }

  $.getJSON('buttons', function(buttons) {
    $('#loading-buttons').hide()

    if (!buttons) {
      $('#no-buttons').removeClass('hidden')
    } else {
      $('#available-buttons').removeClass('hidden')
    }
    buttons.forEach(function(buttonDef) {
      var $button = $('<li></li>')
                      .text(buttonDef.Title)
                      .data('buttonid', buttonDef.Id)
                      .appendTo($('#available-buttons'))
      makeButtonPushable($button, buttonDef)
    })
  })
})

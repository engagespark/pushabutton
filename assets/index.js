$(function() {
  function htmlForParamDef(parameterDef) {
    var html = (
      "<p>" + parameterDef.Title + "</p>"
                + '<p class="description">' + parameterDef.Description + "</p>"
    )

    if (parameterDef.Type == "string") {
      html += '<p><input type="text" maxlength="200"></input> <span class="description">(max. 200 characters)</span></p>'
    }

    if (parameterDef.Type == "choice") {
      html += '<select>' + parameterDef.Details.choices.map(function(choice) {
        return '<option value="' + choice + '">' + choice + '</option>'
      }) + '</select>'

    }

    return html
  }

  function closeParameterModal() {
    $('#parameterModal').hide()
  }

  function redirectToLog(response) {
    window.location = $('#baseUrl').val() + 'log/' + response.pushId + '?autorefresh=10'
  }

  function showParameterModal($button, buttonDef, sendFunc) {

    var $modal = $('#parameterModal');

    $('#parameterModal ol').empty()

    $modal.find('#modal-description').text('Please provide input to push the button ' + buttonDef.Title)

    buttonDef.Parameters.forEach(function(parameterDef, idx) {
      var elementId = 'param-' + idx ;
      parameterDef.elementId = elementId
      $('<li id="' + elementId + '">' + htmlForParamDef(parameterDef) + '</li>').appendTo($('#parameterModal ol'))
    })

    function getArguments() {
      return buttonDef.Parameters.map(function(parameterDef, idx) {
        if (parameterDef.Type == 'string') {
          return $('#' + parameterDef.elementId + ' input').val()
        } else if (parameterDef.Type == 'choice') {
          return $('#' + parameterDef.elementId + ' select').val()
        } else {
          return ""
        }
      })
    }

    $('#btnPushWithParameters').unbind('click').click(
      function() {sendFunc({"pushArguments": getArguments()}).then(redirectToLog)}
    ).text(buttonDef.Title + ' — now for real')
    $('#btnCancelParameters').unbind('click').click(closeParameterModal)

    $modal.removeClass('hidden')
    $modal.show()

    $('.modal-background').height($(document).height())
  }

  function makeButtonPushable($button, buttonDef) {
    function pushFunc(pushArguments) {
      return $.ajax(
        $('#baseUrl').val() + 'api/push/' + $button.data('buttonid'), {
          method: 'POST',
          contentType: 'application/json; charset=UTF-8',
          data: JSON.stringify(pushArguments || {pushArguments: []}),
          dataType: "json"
        }).done(function(data) {
      }).fail(function(error) {
      }).always(function() {
        $button.removeClass("running")
        closeParameterModal()
      })
    }

    $button.click(function(event) {
      function chooseFunc() {
        if (buttonDef.Parameters.length > 0) {
          return function() {
            return showParameterModal($button, buttonDef, pushFunc)
          }
        }
        return pushFunc().then(redirectToLog)
      }

      $button.addClass("running")
      chooseFunc()()
    })
  }

  function toButtonText(buttonDef) {
    var title = buttonDef.Title
    if (buttonDef.Parameters.length > 0) {
      return title.substring(0, title.length - 1) + " … " + title.substr(-1)
    }
    return title
  }

  $('.modal-background').click(closeParameterModal)

  $(document).keyup(function(e) {
     if (e.keyCode == 27) { // escape key maps to keycode `27`
       closeParameterModal()
    }
  });

  $.getJSON($('#baseUrl').val() + 'api/buttons', function(buttons) {
    $('#loading-buttons').hide()

    if (!buttons) {
      $('#no-buttons').removeClass('hidden')
    } else {
      $('#available-buttons').removeClass('hidden')
    }
    buttons.forEach(function(buttonDef) {
      var $button = $('<li></li>')
                      .text(toButtonText(buttonDef))
                      .data('buttonid', buttonDef.Id)
                      .appendTo($('#available-buttons'))
      makeButtonPushable($button, buttonDef)
    })
  })
})

$(function() {
  function htmlForParamDef(parameterDef) {
    var html = (
      "<p>" + parameterDef.Name + "</p>"
                + '<p class="description">' + parameterDef.Description + "</p>"
    )

    if (parameterDef.Type == "string") {
      html += '<input type="text"></input>'
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

  function showParameterModal($button, buttonDef, sendFunc) {

    var $modal = $('#parameterModal');

    $('#parameterModal ol').empty()
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
      function() {sendFunc({"pushArguments": getArguments()})})
    $('#btnCancelParameters').unbind('click').click(closeParameterModal)

    $modal.removeClass('hidden')
    $modal.show()
  }

  function makeButtonPushable($button, buttonDef) {
    function pushFunc(pushArguments) {
      $.ajax(
        'push/' + $button.data('buttonid'), {
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

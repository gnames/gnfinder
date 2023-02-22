$(function () {
  $('ul.tabs')
    .tabs('div.panes > div');
  $('#advanced_options').click(function (e) {
    e.preventDefault();
    $('#advanced_selections').toggle();
  });

  $('#with_verification').click(function () {
    if ($(this).is(':checked')) {
      $.each($('[id^=preferred_data_sources]'), function () {
        $(this).attr('checked', false).attr('disabled', true);
      });
    } else {
      $.each($('[id^=preferred_data_sources]'), function () {
        $(this).removeAttr('disabled');
      });
    }
  });
});

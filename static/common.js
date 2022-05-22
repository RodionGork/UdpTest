function updater() {
    const argType = $('#type').val();
    const argLevel = $('#level').val();
    $.get('/api/latest?type=' + argType + '&level=' + argLevel, data => {
        const tbody = $('table tbody');
        const toRemove = tbody.find('tr');
        for (let i = 0; i < data.length; i++) {
            const event = data[i];
            const tr = $('<tr/>').attr('id', event.UUID);
            $('<td/>').text(event.Type).appendTo(tr);
            $('<td/>').text(event.Severity).addClass(event.Severity.toLowerCase()).appendTo(tr);
            $('<td/>').text(event.EventID).appendTo(tr);
            $('<td/>').text(new Date(event.UnixTimestamp * 1000).toLocaleTimeString()).appendTo(tr);
            tbody.append(tr);
        }
        toRemove.remove();
    });
    setTimeout(updater, 3000);
}

function counter() {
    $.get('/api/count', data => {
        $('#total').text(data['total']);
    });
    setTimeout(counter, 230);
}

$(function() {
    updater();
    counter();
    $('tbody').on('click', 'tr', function () {
        window.open('/api/get/' + $(this).attr('id'));
    });
});

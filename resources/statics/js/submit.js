(function (){
    document.addEventListener('DOMContentLoaded', init, false);

    function init(){
        var button = document.getElementById('submit-button');
        button.onclick = sendForm;
    }

    function sendForm() {
        button = this
        button.disabled = 'disabled';
  
        var form = document.getElementById('submit-form');
        var formData = new FormData(form);
        $.ajax({
            url: form.action + '',
            type: form.method + '',
            dataType: 'json',
            data: formData,
            processData: false,
            contentType: false,
            xhrFields: {
                withCredentials: true
            }
        }).done(function(res, status, xhr){
            console.log('ok');
            $(form).hide();
            var success = $('<div/>');
            success.text('success!');
            success.addClass('bg-success col md-8 offset-md-2');
            $('#container').append(success);
        }).fail(function(jqXHR, textStatus, errorThrown){
            console.log('error');
            console.log(textStatus)
            console.log(jqXHR);
            button.removeAttribute('disabled')
            var success = $('<div/>');
            success.text('fail!');
            success.addClass('bg-error col md-8 offset-md-2');
            $('#container').append(success);
        });

    }
})();

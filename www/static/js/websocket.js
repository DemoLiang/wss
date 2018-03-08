var socket;

$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://127.0.0.1:7777/ws');
    // Message received on the socket
    socket.onmessage = function (event) {
        var data = JSON.parse(event.data);
        console.log(data);
        $("#chatbox li").first().before("<li>"+ data.errcode + "</li>");
    };

    // Send messages.
    var postConecnt = function () {
        var content = $('#sendbox').val();
        socket.send(content);
        $('#sendbox').val("");
    }

    $('#sendbtn').click(function () {
        postConecnt();
    });
});
var socket;

$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://127.0.0.1:7777/ws');
    // Message received on the socket
    socket.onmessage = function (event) {
        console.log(event.data);
        $("#chatbox li").first().before("<li>"+ event.data.toLocaleString() + "</li>");
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
$(document).ready(function () {
    var userName;
    var finalConexion;

    $("#chat").hide();

    $('#frmRegistro').on("submit", function (e) {
        e.preventDefault();
        userName = $("#user_name").val();

        $.ajax({
            type: 'POST',
            url: 'http://localhost:3030/validate',
            data: {
                "user_name": userName
            },
            success: function (response) {
                result(response);
            },
            error: function (result) {
                console.log(result);
            }
        });
    });

    function result(response) {
        obj = JSON.parse(response);
        if (obj.isvalid === true) {
            createConexion();
        } else {
            console.log('Intentalo de nuevo');
        }
    }

    function createConexion() {
        $("#registro").hide();
        $("#chat").show();
        var conexion = new WebSocket("ws://localhost:3030/chat/" + userName);
        finalConexion = conexion;

        conexion.onopen = function (response) {
            conexion.onmessage = function (response) {
                console.log(response.data);
                val = $("#chat_area").val();
                $("#chat_area").val(val + "\n" + response.data);
            }
        }
    }

    $("#frmChat").on("submit", function (e) {
        e.preventDefault();
        mensaje = $("#msg").val();
        finalConexion.send(mensaje);
        $("#msg").val("");
    })
});
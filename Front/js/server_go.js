$(document).ready(()=>{
    var user_name;
    var final_conexion;
    $("#form_registro").on("submit", (e)=>{
        e.preventDefault();
        user_name = $('#username').val();
        $.ajax({
            type: "POST",
            url: "http://localhost:8000/validate",
            data:{
                "user_name": user_name
            },
            success: function (data) {
                result(data);
            }
        })
    })

    function result(data) {
       let obj = JSON.parse(data)
        if (obj.isvalid === true){
            create_conexion();
        }else{
         console.log("Fallo de conexion")
        }
    }

    function create_conexion() {
        $("#registro").hide();
        $("#container_chat").show();
        let conexion = new WebSocket("ws://localhost:8000/chat/" + user_name )
        final_conexion = conexion
        conexion.onopen = function (response) {
            conexion.onmessage = function (response) {
                console.log(response.data)
                let val = $('#chat_area').val()
                $('#chat_area').val(val + "\n" + response.data)

            }
        }
    }
    
    $('#form_message').on('submit', (e)=>{
        e.preventDefault();
        let message = $('#msg').val()
        final_conexion.send(message)
        $('#msg').val('');
    })
    
    
});
const NewClient=function () {
    var ws = new WebSocket("ws://localhost:8080/echo");
    ws.onopen = function(){
        console.log("open");
    }
    ws.onclose = function(e){
        console.log("close");
    }
    ws.onerror = function(e){

        console.log(e);
    }
    return ws
}

window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var name = document.getElementById("name");
    var log = document.getElementById("log");
    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }
    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        if (!name.value) {
            return false;
        }
        var data = '{"name":"'+name.value+'","msg":"'+msg.value+'"}'
        conn.send(data);
        msg.value = "";
        name.value = "";
        return false;
    };
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://localhost:8080/socket");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>我们无法连接至OnlineChat，我们正在尝试重连...</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                var obj=JSON.parse(messages);
                item.innerText = obj.name + "：" + obj.msg;
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>你的浏览器不支持 WebSockets.</b>";
        appendLog(item);
    }
}
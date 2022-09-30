window.onload = function () {

    var request = new XMLHttpRequest();
      request.open("get", "/api/past_msg");
      request.send(null);
      request.onload = function () {
        if (request.status == 200) {
          var json = JSON.parse(request.responseText);
          var ol = document.getElementById('ol');
          var frag = document.createDocumentFragment();
          json.person.map(person => {
            var li = document.createElement("li");
            li.innerHTML = `名字是 ${person.name} 图片是 ${person.msg}`;
            frag.append(li);
          })
          ol.append(frag)
        }
      }

      disabledButton("send", true)

    var conn;
    var msg = document.getElementById("msg");
    var name = document.getElementById("user");
    var log = document.getElementById("log");
    var notice = document.getElementById("notice");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function getTime() {
        var mydate = new Date();
        var year = mydate.getFullYear();
        var month = mydate.getMonth() + 1;
        var day = mydate.getDate();
        var hour = mydate.getHours();
        var minutes = mydate.getMinutes();
        var now = "";
        if (minutes < 10) {
            minutes = '0' + minutes;
        }

        if (hour <= 12) {
            now = "上午";
        } else {
            now = "下午";
        }
        var seconds = mydate.getSeconds();
        if (seconds < 10) {
            seconds = '0' + seconds;
        }
        var weekday = mydate.getDay();
        var arr = new Array('星期天', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六');
        var time = year + '/' + month + '/' + day + "  " + now + "  " + hour + ':' + minutes + ':'
            + seconds + "  " + arr[weekday];
        return time
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        var data = '{"name":"' + name.value + '","msg":"' + msg.value + '","time":"' + getTime() + '"}'
        conn.send(data);
        // conn.send(msg.value);
        msg.value = "";
        return false;
    };

    function disabledButton(id, data) {
        var button = document.getElementById(id);
        button.disabled = data; //使用true或false，控制是否让按钮禁用
    }

    function checkUserName() {
        if (name.value == "") {
            disabledButton("send", true)
            notice.innerHTML = `
            <div class="alert alert-warning" role="alert">
            <svg style="width:24px;height:24px;margin-right:5px;" viewBox="0 0 24 24">
            <path fill="currentColor" d="M10 4A4 4 0 0 1 14 8A4 4 0 0 1 10 12A4 4 0 0 1 6 8A4 4 0 0 1 10 4M10 14C14.42 14 18 15.79 18 18V20H2V18C2 15.79 5.58 14 10 14M20 12V7H22V13H20M20 17V15H22V17H20Z" />
        </svg>您还没有登录，没有发言的权限。
            </div>`;
        }else{
            disabledButton("send", false)
        }
    }

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://192.168.1.2:8080/socket");
        conn.onclose = function () {
            notice.innerHTML = `
                <div class="alert alert-danger" role="alert">
                <svg style="width:24px;height:24px;margin-right:5px;" viewBox="0 0 24 24">
                <path fill="currentColor" d="M4,1C2.89,1 2,1.89 2,3V7C2,8.11 2.89,9 4,9H1V11H13V9H10C11.11,9 12,8.11 12,7V3C12,1.89 11.11,1 10,1H4M4,3H10V7H4V3M14,13C12.89,13 12,13.89 12,15V19C12,20.11 12.89,21 14,21H11V23H23V21H20C21.11,21 22,20.11 22,19V15C22,13.89 21.11,13 20,13H14M3.88,13.46L2.46,14.88L4.59,17L2.46,19.12L3.88,20.54L6,18.41L8.12,20.54L9.54,19.12L7.41,17L9.54,14.88L8.12,13.46L6,15.59L3.88,13.46M14,15H20V19H14V15Z" />
            </svg> 连接意外的断开，无法与 WebSocket 服务器连接，刷新页面重试。
                </div>`;
            disabledButton("send", true)
        };
        conn.onopen = function () {
            checkUserName()
            // disabledButton("send", false)
        }
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                var obj = JSON.parse(messages);
                // item.innerHTML = `obj.name obj.msg`;
                item.innerHTML = `
                    <div class="log-item">
                        <p class="user-name"><a href="#"><strong>` + obj.name + `</strong></a></p>
                        <div class="message">` + obj.msg + `</div>
                        <p class="time"><small>` + obj.time + `</small></p>
                    </div>
                `;
                // item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        notice.innerHTML = `
        <div class="alert alert-danger" role="alert">
            貌似出现了一些问题
        </div>`;
    }
};
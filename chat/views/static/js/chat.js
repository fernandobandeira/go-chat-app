(function() {
    window.qs = (function(a) {
        if (a == "") return {};
        var b = {};
        for (var i = 0; i < a.length; ++i)
        {
            var p=a[i].split('=', 2);
            if (p.length == 1)
                b[p[0]] = "";
            else
                b[p[0]] = decodeURIComponent(p[1].replace(/\+/g, " "));
        }
        return b;
    })(window.location.search.substr(1).split('&'));

    setUserName();

    window.chat = createChatController();
    window.chat.startListenWS();

    document.getElementById("btnSend").onclick = function() {
        window.chat.sendMessage();
    };
    document.getElementById("txtMessage").onkeyup = function(event) {
        if (event.key == "Enter") {
            event.preventDefault();
            window.chat.sendMessage();
        }
    };
    document.getElementById("txtMessage").focus();
})();

function basePath() {
    return document.getElementsByTagName("base")[0].attributes["href"].value;
}

function setUserName() {
    if(window.qs["name"] == undefined || window.qs["name"] == null) {
        window.location = window.location.origin + basePath();
    }
    document.getElementById("author").value = window.qs["name"];
    document.getElementById("greet").innerText = `Chat - Hello, ${window.qs["name"]}`;
}

function createChatController() {
    var state = {};

    return {
        webSocket: undefined,
        state: state,
        startListenWS: function() {
            // this.webSocket = new WebSocket("ws://localhost:80/ws");
            this.webSocket = new WebSocket(`ws://${window.location.host}${basePath()}ws`);

            this.webSocket.onmessage = function(event) {
                msg = JSON.parse(event.data);
                addMessage(msg);
            }

            this.webSocket.onopen = function() {
                this.webSocket.send("Hello!\0");
            }
        },
        sendMessage: function() {
            let text = document.getElementById("txtMessage").value;
            let author = document.getElementById("author").value;

            msg = {
                author: author,
                text: text,
            };

            this.webSocket.send(JSON.stringify(msg));

            document.getElementById("txtMessage").value = "";
        }
    }
}

function addMessage(msg) {
    let author = msg.author;
    let text = msg.text;

    let chatMessage = document.createTextNode(`${author}: ${text}`);
    let liElement = document.createElement("li");
    liElement.appendChild(chatMessage);

    document.getElementById("msgList").appendChild(liElement);
}

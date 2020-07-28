var term,
    id,
    nodeIP,
    socket,
    charWidth,
    charHeight;

var terminalContainer = document.getElementById('terminal-container');

createTerminal();

function getUrlParam(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
    var r = window.location.search.substr(1).match(reg);  //匹配目标参数
    if (r != null) return unescape(r[2]); return null; //返回参数值
}

function setTerminalSize() {
    var initialGeometry = term.proposeGeometry(),
        cols = initialGeometry.cols,
        rows = initialGeometry.rows;
    charWidth = Math.ceil(term.element.offsetWidth / cols);
    charHeight = Math.ceil(term.element.offsetHeight / rows);
    var width = (cols * charWidth).toString() + 'px';
    var height = (rows * charHeight).toString() + 'px';

    terminalContainer.style.width = width;
    terminalContainer.style.height = height;

    var url = '/nodes/containers/shell/resize';
    var returnCode;
    $.ajaxSettings.async = false;
    $.getJSON(url, {id: id, nodeIP: nodeIP, cols: cols, rows: rows})
        .done(function (data) {
            returnCode = data.return_code;
            if (returnCode != 200) {
                console.log(data.return_message);
            }
        })
        .fail(function (jqxhr, textStatus, error) {
            var err = textStatus + ", " + error;
            console.log( "Request Failed: " + err );
        });
    if (returnCode != 200) {
        return false
    } else {
        return true
    }
}

function createTerminal() {
    // Clean terminal
    while (terminalContainer.children.length) {
        terminalContainer.removeChild(terminalContainer.children[0]);
    }
    // get id
    nodeIP = getUrlParam("nodeIP");
    containerID = getUrlParam("containerID");
    command = "/bin/bash";
    user = "root";
    $.ajaxSettings.async = false;
    $.getJSON("/nodes/containers/shell/create", {nodeIP: nodeIP, containerID: containerID, command: command, user: user})
        .done(function (data) {
            // console.log(data);
            if (data.return_code != 200) {
                console.log(data.return_message);
                return
            }
            id = data.data.id;
            return
        })
        .fail(function (jqxhr, textStatus, error) {
            var err = textStatus + ", " + error;
            console.log( "Request Failed: " + err );
        });
    if (id == "" || id == undefined) {
        $(".terminal-container").text("无法获取ID");
        return
    }
    // create terminal object
    term = new Terminal(
    );
    term.open(terminalContainer);
    term.fit();
    // term.refresh();
    // build websocket
    var protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';
    var socketURL = protocol + location.hostname + ((location.port) ? (':' + location.port) : '') + '/nodes/containers/shell/ws';
    socketURL += '?id=' + id + '&nodeIP=' + nodeIP;
    socket = new WebSocket(socketURL);
    // attach xterm to websocket
    socket.onopen = runRealTerminal;
    socket.onclose = runFakeTerminal;
    socket.onerror = runFakeTerminal;
    // console.log(term.cols);
    // console.log(term.rows);
    // send resize request to docker rest api
    setTerminalSize();
}

function runRealTerminal() {
    term.attach(socket);
    term._initialized = true;
}

function runFakeTerminal() {
    if (term._initialized) {
        return;
    }

    term._initialized = true;

    var shellprompt = '$ ';

    term.prompt = function () {
        term.write('\r\n' + shellprompt);
    };

    term.writeln('Welcome to xterm.js');
    term.writeln('This is a local terminal emulation, without a real terminal in the back-end.');
    term.writeln('Type some keys and commands to play around.');
    term.writeln('');
    term.prompt();

    term.on('key', function (key, ev) {
        var printable = (
            !ev.altKey && !ev.altGraphKey && !ev.ctrlKey && !ev.metaKey
        );

        if (ev.keyCode == 13) {
            term.prompt();
        } else if (ev.keyCode == 8) {
            // Do not delete the prompt
            if (term.x > 2) {
                term.write('\b \b');
            }
        } else if (printable) {
            term.write(key);
        }
    });

    term.on('paste', function (data, ev) {
        term.write(data);
    });
}

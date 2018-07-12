(function () {

    var openWsConnection = function () {
        var socket = new WebSocket("ws://localhost:10000/io");

        socket.onopen = function () {
            console.log("connected");
            socket.send('ping');
        };

        socket.onclose = function (e) {
            console.log("connection closed (" + e.code + ")");
        };

        socket.onmessage = function (e) {
            console.log("message received: " + e.data);
        };

        return socket;
    };

    window.App = Vue.component('app', {
        template: '#template_app',
        data: function () {
            return {
                links: [1, 2, 3]
            }
        },
        created() {
            openWsConnection();
        },
    });
})();

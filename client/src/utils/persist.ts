export function wrapper(socket: WebSocket, callback: Function) {
    setTimeout(function () {
        if (socket.readyState === 1) {
            if (callback != null) {
                callback();
            }
        } else {
            wrapper(socket, callback);
        }
    }, 5);
}

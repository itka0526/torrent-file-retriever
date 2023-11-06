<script lang="ts">
    import List from "./main/list.svelte";

    let socket = new WebSocket("ws://" + document.location.host + "/api/ws");
    let list: string[] = [];

    function wrapper(socket: WebSocket, callback: Function) {
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

    wrapper(socket, () => {
        socket.send("get_files");
        socket.onmessage = (e) => {
            console.log(JSON.parse(e.data));
            list = JSON.parse(e.data);
        };
    });
</script>

{#if window["WebSocket"] && socket.readyState == WebSocket.OPEN}
    <div class="h-screen w-screen flex">
        <div class="basis-1/4" />
        <div class="basis-3/4">
            <List bind:list />
        </div>
    </div>
{/if}

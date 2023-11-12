import { readable, writable } from "svelte/store";
import type { WSMessage, MyFileInfo, ParsedData, UploadingFiles, DownloadingFiles } from "../types";

export let file_list = writable<MyFileInfo[]>([]);

export let socket = readable<WebSocket | null>(null, (set) => {
    let ws = new WebSocket("ws://" + document.location.host + "/api/ws");

    ws.onmessage = function (event: MessageEvent<string>) {
        const message = JSON.parse(event.data) as WSMessage;
        switch (message.response_type) {
            case "get_files_res":
                file_list.update(() => JSON.parse(message.data) as ParsedData as MyFileInfo[]);
        }
    };

    set(ws);

    return () => {
        ws.close();
    };
});

export let uploadingFiles = writable<UploadingFiles>({});
export let downloadingFiles = writable<DownloadingFiles>({});

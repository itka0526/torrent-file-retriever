<script lang="ts">
    import { toast } from "@zerodevx/svelte-toast";
    import { DownloadIcon } from "svelte-feather-icons";
    import type { Message, MyFileInfo } from "../../types";
    import axios, { HttpStatusCode } from "axios";
    import { downloadingFiles } from "../../utils/store";

    export let file: MyFileInfo;

    const handleDownload = async () => {
        const fileName = file.is_directory ? file.name + ".zip" : file.name;

        const pendingResponse = await axios.post("/api/download", file, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            responseType: "blob",
            onDownloadProgress: (event) => {
                event.total = file.size;
                $downloadingFiles[fileName] = event;
            },
        });

        if (pendingResponse.status == HttpStatusCode.BadRequest) {
            const reader = new FileReader();
            reader.onload = function () {
                const response: Message = JSON.parse(reader.result as string);
                toast.push(response.message);
            };
            return;
        }

        const blob = pendingResponse.data;
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.download = fileName;
        link.href = url;
        link.click();
        URL.revokeObjectURL(url);

        downloadingFiles.update((df) => {
            console.log(df);
            delete df[fileName];
            return df;
        });
    };
</script>

<button on:click={handleDownload}>
    <DownloadIcon class="max-md:w-5" />
</button>

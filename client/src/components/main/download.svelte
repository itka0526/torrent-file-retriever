<script lang="ts">
    import { toast } from "@zerodevx/svelte-toast";
    import { DownloadIcon } from "svelte-feather-icons";
    import type { Message, MyFileInfo } from "../../types";
    import { HttpStatusCode } from "axios";

    export let file: MyFileInfo;

    const handleDownload = async () => {
        const pendingResponse = await fetch("/api/download", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(file),
        });

        if (pendingResponse.status == HttpStatusCode.BadRequest) {
            const response: Message = await pendingResponse.json();
            toast.push(response.message);
            return;
        }

        const blob = await pendingResponse.blob();
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.download = file.is_directory ? file.name + ".zip" : file.name;
        link.href = url;
        link.click();
        URL.revokeObjectURL(url);
    };
</script>

<button on:click={handleDownload}>
    <DownloadIcon class="max-md:w-5" />
</button>

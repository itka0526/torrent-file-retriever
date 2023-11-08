<script lang="ts">
    import { TrashIcon, XIcon } from "svelte-feather-icons";
    import type { Message, MyFileInfo } from "../../types";
    import { toast } from "@zerodevx/svelte-toast";

    export let file: MyFileInfo;
    let prompt = false;

    const handleDelete = async () => {
        if (!prompt) {
            prompt = true;
            return;
        }

        const pending_promise = await fetch("/api/delete", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(file),
        });

        const response: Message = await pending_promise.json();
        toast.push(response.message);
        prompt = false;
    };
</script>

<button on:click={handleDelete}>
    {#if prompt && file}
        <XIcon />
    {:else}
        <TrashIcon />
    {/if}
</button>

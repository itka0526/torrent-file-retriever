<script lang="ts">
    import { toast } from "@zerodevx/svelte-toast";
    import type { Message } from "../../types";
    import { ClipboardIcon, LoaderIcon } from "svelte-feather-icons";

    let url = "";
    let loading = false;

    const handleChange = (
        e: Event & {
            currentTarget: EventTarget & HTMLInputElement;
        }
    ) => {
        if (e.target) url = (e.target as HTMLInputElement).value;
    };

    const downloadTorrent = async (u: string) => {
        if (!u || !u.match(/magnet:\?xt=urn:[a-z0-9]+:[a-z0-9]{32}/i)) {
            u.length && toast.push("Invalid magnet URL.");
            return;
        }
        loading = true;
        const pendingResponse = await fetch("/api/magnet", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                url: u,
            }),
        });
        loading = false;
        const response: Message = await pendingResponse.json();
        toast.push(response.message);
        url = "";
    };

    $: downloadTorrent(url);
</script>

<div class="flex bg-gray-200 px-4 py-2 gap-2 shadow-md">
    {#if loading}
        <LoaderIcon class="animate-spin" strokeWidth={1} />
    {:else}
        <ClipboardIcon strokeWidth={1} />
    {/if}
    <input placeholder="Paste link URL Here" class="outline-none bg-transparent p-0 m-0" on:change={handleChange} type="text" bind:value={url} />
</div>

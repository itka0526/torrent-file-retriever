<script lang="ts">
    import { store } from "./store";
    import type { ItemInfo } from "./types";

    let url: string = "";
    let timer: number;
    let state: boolean = true;

    $: url, debouncedSearch();

    function debouncedSearch() {
        clearTimeout(timer);
        timer = setTimeout(search, 750);
    }

    type Result = { data: ItemInfo | null; message: string };

    const search = async () => {
        if (!url) return;

        const response = await fetch("/api/metadata", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ URL: url }),
        });

        const result: Result = await response.json();

        if (!result.data) {
            state = false;
            url = "";
        } else {
            state = true;
            url = "";

            store.update((prev) => {
                prev.push(result.data as ItemInfo);
                localStorage.setItem("store", JSON.stringify(prev));
                return prev;
            });
        }
    };
</script>

<div class={`overflow-hidden rounded-full w-full h-12 border ${state ? "border-black" : "border-red-400"} my-4`}>
    <input
        bind:value={url}
        placeholder="Enter a magnet URL..."
        class={`outline-none w-full h-full px-8 text-lg ${state ? "" : "placeholder:text-red-400"}`}
    />
</div>

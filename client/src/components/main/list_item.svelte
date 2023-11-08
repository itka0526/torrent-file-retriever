<script lang="ts">
    import { ChevronRightIcon } from "svelte-feather-icons";
    import { humanFileSize } from "../../utils/humanFileSize";
    import Download from "./download.svelte";
    import Delete from "./delete.svelte";
    import type { MyFileInfo } from "../../types";

    export let indent = 0;
    export let children: any = [];

    export let path = "";
    export let name = "";
    export let modified_date = "";
    export let is_directory = false;
    export let size = 0;

    let open = false;

    function toggleOpen() {
        open = !open;
    }

    let file: MyFileInfo = { path, name, modified_date, is_directory, size };
</script>

<li class="px-4 py-2 hover:bg-gray-200 flex items-center">
    <div class="flex justify-stretch basis-3/4 flex-grow-0 overflow-hidden">
        <div class="basis-1/2 flex" style="padding-left: {indent}px">
            {#if is_directory}
                <button on:click={toggleOpen}>
                    <ChevronRightIcon class={`transition-transform ${open ? "rotate-90" : ""}`} />
                </button>
            {/if}
            <span>
                {name}
            </span>
        </div>
        <span class="basis-1/4">
            {new Date(modified_date).toLocaleString()}
        </span>
        <span class="basis-1/4">
            {humanFileSize(size)}
        </span>
    </div>
    <div class={"basis-1/4 flex gap-2"}>
        <Download bind:file />
        <Delete bind:file />
    </div>
</li>

{#if open}
    {#each children as child}
        <svelte:self {...child} indent={indent + 16} />
    {/each}
{/if}

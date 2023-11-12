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

    $: file = { path, name, modified_date, is_directory, size } as MyFileInfo;
</script>

<li class="px-4 py-2 hover:bg-gray-200 flex items-center max-md:text-sm max-md:px-2 max-md:py-1">
    <div class="grid grid-cols-[2fr,1fr,1fr] gap-2 basis-3/4 max-md:grid-cols-[3fr,0,1fr]">
        <div class="flex w-full overflow-hidden items-center" style="padding-left: {indent}px">
            {#if is_directory}
                <button on:click={toggleOpen}>
                    <ChevronRightIcon class={`transition-transform ${open ? "rotate-90" : ""} max-md:w-5`} />
                </button>
            {/if}
            <div class="text-ellipsis overflow-hidden whitespace-nowrap w-full">
                <span title={name}>
                    {name}
                </span>
            </div>
        </div>
        <span class="max-md:hidden whitespace-nowrap">
            {new Date(modified_date).toLocaleString()}
        </span>
        <span class="whitespace-nowrap">
            {humanFileSize(size)}
        </span>
    </div>
    <div class={"basis-1/4 flex gap-2 items-center justify-center"}>
        <Download {file} />
        <Delete {file} />
    </div>
</li>

{#if open}
    {#each children as child}
        <svelte:self {...child} indent={indent + 16} />
    {/each}
{/if}

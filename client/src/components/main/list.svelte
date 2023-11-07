<script lang="ts">
    import { humanFileSize } from "../../utils/humanFileSize";
    import { file_list } from "../../utils/store";
</script>

<ul>
    {#each $file_list as item}
        <li class="px-4 py-2 hover:bg-gray-200 flex items-center">
            <div class="flex justify-stretch basis-3/4 flex-grow-0 overflow-hidden">
                <span class="basis-1/2">
                    {item.name}
                </span>
                <span class="basis-1/4">
                    {new Date(item.modified_date).toLocaleString()}
                </span>
                <span class="basis-1/4">
                    {humanFileSize(item.size)}
                </span>
            </div>
            <div class="basis-1/4" />
        </li>
    {/each}
</ul>

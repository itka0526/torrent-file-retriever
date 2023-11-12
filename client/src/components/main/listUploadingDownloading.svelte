<script lang="ts">
    import { humanFileSize } from "../../utils/humanFileSize";
    import { uploadingFiles, downloadingFiles } from "../../utils/store";

    let colA = "w-1/2 text-ellipsis whitespace-nowrap max-w-0 overflow-hidden",
        colB = "w-1/2 text-right whitespace-nowrap max-w-0 overflow-hidden";
</script>

<ul class="px-2 py-4 w-full flex flex-col gap-4">
    {#each Object.entries($downloadingFiles).reverse() as [name, info]}
        <li class="bg-gray-200 shadow-md px-4 py-2">
            <table class="w-full text-sm border-collapse">
                <tr>
                    <th colspan="2" title={name} class={" text-ellipsis whitespace-nowrap max-w-0 overflow-hidden w-full "}>
                        {name}
                    </th>
                </tr>
                <tr>
                    <td class={colA}> Progress: </td>
                    <td class={colB}>
                        {info.progress ? Math.floor(info.progress * 100) : "0"}%
                    </td>
                </tr>
                <tr>
                    <td class={colA}> Loaded: </td>
                    <td class={colB}>
                        {humanFileSize(info.loaded)}/{humanFileSize(info.total ?? 0)}
                    </td>
                </tr>
                <tr>
                    <td class={colA}> Rate: </td>
                    <td class={colB}>
                        {humanFileSize(info.rate ?? 0)}
                    </td>
                </tr>
                <tr>
                    <td class={colA}> Estimated: </td>
                    <td class={colB}>
                        {Math.floor(info.estimated ?? 0)} secs
                    </td>
                </tr>
            </table>
        </li>
    {/each}
    {#each Object.entries($uploadingFiles).reverse() as [name, info]}
        <li class="bg-gray-200 shadow-md px-4 py-2">
            <table class="w-full text-sm border-collapse">
                <tr>
                    <td colspan="2" title={name} class={" text-ellipsis whitespace-nowrap max-w-0 overflow-hidden w-full "}>
                        {name}
                    </td>
                </tr>
                <tr>
                    <td class={colA}> Progress: </td>
                    <td class={colB}>
                        {info.progress ? Math.floor(info.progress * 100) : "0"}%
                    </td>
                </tr>
                <tr>
                    <td class={colA}> Loaded: </td>
                    <td class={colB}>
                        {humanFileSize(info.loaded)}/{humanFileSize(info.total ?? 0)}
                    </td>
                </tr>
                <tr>
                    <td class={colA}> Rate: </td>
                    <td class={colB}>
                        {humanFileSize(info.rate ?? 0)}
                    </td>
                </tr>
                <tr>
                    <td class={colA}> Estimated: </td>
                    <td class={colB}>
                        {Math.floor(info.estimated ?? 0)} secs
                    </td>
                </tr>
            </table>
        </li>
    {/each}
</ul>

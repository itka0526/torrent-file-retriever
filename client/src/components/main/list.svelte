<script lang="ts">
    export let list: string[];
</script>

<ul>
    {#each list as item}
        <li class="px-4 py-2 hover:bg-gray-200">{item}</li>
    {/each}
</ul>

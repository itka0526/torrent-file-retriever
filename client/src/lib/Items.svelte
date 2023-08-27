<script lang="ts">
    import { store } from "./store";
    import type { ItemInfo } from "./types";

    let localItems: ItemInfo[] = [];

    store.subscribe((val) => (localItems = val));
</script>

<div class="">
    <ul>
        {#each localItems as item}
            <li>{item.Name}</li>
        {/each}
    </ul>
</div>

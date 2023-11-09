<script lang="ts">
    import { file_list } from "../../utils/store";
    import { transform } from "../../utils/transform";
    import ListItem from "./listItem.svelte";

    $: files = transform($file_list);
</script>

<ul>
    <ListItem {...files} />
</ul>

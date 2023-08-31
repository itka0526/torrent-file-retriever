<script lang="ts">
    import Items from "./lib/Items.svelte";
    import Search from "./lib/Search.svelte";
    import TailwindCss from "./lib/TailwindCSS.svelte";
    import { store } from "./lib/store";

    store.subscribe((v) => console.log(v));
</script>

<TailwindCss />

<main class="flex flex-col min-h-screen w-full px-[15%]">
    <Search />
    <Items />
</main>

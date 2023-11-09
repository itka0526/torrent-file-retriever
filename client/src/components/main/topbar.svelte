<script lang="ts">
    import { MenuIcon } from "svelte-feather-icons";
    import MagnetUrl from "./magnetURL.svelte";

    export let open: Boolean;

    const handleClick = () => {
        open = !open;
    };
</script>

<div class="w-full px-6 py-6 flex justify-end items-center max-md:justify-start">
    <button class="hidden max-md:flex max-md:pr-6" on:click={handleClick}>
        <MenuIcon />
    </button>
    <MagnetUrl />
</div>

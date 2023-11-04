<script lang="ts">
    // import Items from "./lib/Items.svelte";
    // import Search from "./lib/Search.svelte";
    import Login from "./components/login.svelte";
    import TailwindCss from "./lib/TailwindCSS.svelte";
    import type { Message } from "./types";
    // import { store } from "./lib/store";
    let authorized = false;

    const init = async () => {
        const pendingResponse = await fetch("/api/auth");
        const response: Message = await pendingResponse.json();

        if (response.status) {
            authorized = true;
        } else {
            authorized = false;
        }
    };
    init();
</script>

<TailwindCss />

<main class="min-h-screen w-full">
    {#if authorized}
        <p>You are in!</p>
    {:else}
        <Login />
    {/if}
</main>

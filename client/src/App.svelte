<script lang="ts">
    import { SvelteToast } from "@zerodevx/svelte-toast";
    import Login from "./components/login.svelte";
    import TailwindCss from "./lib/TailwindCSS.svelte";
    import type { Message } from "./types";
    import { loginReqOpt } from "./lib/types";
    import Main from "./components/main.svelte";
    let loaded = false,
        authorized = false;
    const init = async () => {
        const pendingResponse = await fetch("/api/auth", loginReqOpt({ username: "", password: "" }));
        const response: Message = await pendingResponse.json();
        loaded = true;
        authorized = response.status;
    };
    init();
</script>

<TailwindCss />
<SvelteToast />
{#if loaded}
    <main class="min-h-screen w-full">
        {#if authorized}
            <Main />
        {:else}
            <Login bind:authorized />
        {/if}
    </main>
{/if}

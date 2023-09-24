<script lang="ts">
    // import Items from "./lib/Items.svelte";
    // import Search from "./lib/Search.svelte";
    import TailwindCss from "./lib/TailwindCSS.svelte";
    // import { store } from "./lib/store";

    // store.subscribe((v) => console.log(v));
    let username = "admin",
        password = "123";

    async function login() {
        // let socket = new WebSocket("/api/ws");
        //You should do the authentication through web,
        // return a cookie and then connect to the websocket server again, carrying the cookie.
        // The WS server can validate the cookie
        const options: RequestInit = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username,
                password,
            }),
        };
        await fetch("/api/auth", options);
    }
</script>

<TailwindCss />

<main class="flex flex-col min-h-screen w-full px-[15%]">
    <div class="h-screen w-screen flex items-center">
        <div class="flex flex-col">
            Username: <input value={username} />
            <br />
            Password: <input value={password} />
            <br />
            <button on:click={login}>Login</button>
        </div>
    </div>
</main>

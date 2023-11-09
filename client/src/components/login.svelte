<script lang="ts">
    import { UsersIcon } from "svelte-feather-icons";
    import { loginReqOpt } from "../lib/types";
    export let authorized: boolean;
    import type { Message } from "../types";
    import { toast } from "@zerodevx/svelte-toast";

    let username = "",
        password = "";

    async function login() {
        const pendingResponse = await fetch("/api/auth", loginReqOpt({ username, password }));
        const response: Message = await pendingResponse.json();

        authorized = response.status;
        toast.push(response.message);
    }
</script>

<div class="h-screen w-screen flex items-center justify-center">
    <div class="flex flex-col px-16 py-8 min-w-[30%]">
        <div class="px-4 py-6 gap-2 flex justify-center items-center">
            <h1 class="text-3xl">Welcome back!</h1>
        </div>
        <div class="flex flex-col gap-3">
            <input class="bg-gray-200 px-4 py-2 shadow-md" bind:value={username} placeholder="username" />
            <input class="bg-gray-200 px-4 py-2 shadow-md" bind:value={password} placeholder="password" />
            <button class="bg-green-400 px-4 py-2 shadow-md" on:click={login}>Login</button>
        </div>
    </div>
</div>

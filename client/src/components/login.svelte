<script lang="ts">
    import { toasts, ToastContainer, FlatToast, BootstrapToast } from "svelte-toasts";
    import type { Message } from "../types";
    const showToast = () => {
        const toast = toasts.add({
            title: "Message title",
            description: "Message body",
            duration: 10000, // 0 or negative to avoid auto-remove
            placement: "bottom-right",
            type: "info",
            theme: "dark",
            onClick: () => {},
            onRemove: () => {},
        });
    };
    let username = "admin",
        password = "123";

    async function login() {
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
        const pendingResponse = await fetch("/api/auth", options);
        const response: Message = await pendingResponse.json();

        if (response.status) {
        }
    }
</script>

<div class="h-screen w-screen flex items-center justify-center">
    <div class="flex flex-col px-16 py-8 rounded-sm min-w-[30%]">
        <div class="px-4 py-6 flex justify-center items-center">
            <h1 class="text-3xl">Welcome back!</h1>
        </div>

        <div class="flex flex-col gap-3">
            <input class="bg-gray-200 px-4 py-2 shadow-md" value={username} placeholder="username" />

            <input class="bg-gray-200 px-4 py-2 shadow-md" value={password} placeholder="password" />

            <button class="bg-green-400 px-4 py-2 shadow-md" on:click={login}>Login</button>
        </div>
    </div>
</div>

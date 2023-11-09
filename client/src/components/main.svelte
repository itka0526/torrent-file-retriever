<script lang="ts">
    import { wrapper } from "../utils/persist";
    import { socket } from "../utils/store";
    import List from "./main/list.svelte";
    import Topbar from "./main/topbar.svelte";
    import Upload from "./main/upload.svelte";

    let open = false;

    $socket &&
        wrapper($socket, () => {
            if ($socket) {
                $socket.send("get_files");
            }
        });
</script>

<Topbar bind:open />
<div class="h-screen w-full grid grid-cols-[1fr,3fr] max-md:grid-cols-[100%,100%] max-md:relative">
    <div
        class={`w-full max-md:absolute max-md:bg-white max-md:h-full max-md:z-10 max-md:transition-transform
        ${open ? "max-md:translate-x-0" : "max-md:-translate-x-full"}`}
    >
        <Upload />
    </div>
    <div
        class={`w-full max-md:absolute max-md:transition-transform
        ${open ? "max-md:translate-x-full" : "max-md:-translate-x-0"}
        `}
    >
        <List />
    </div>
</div>

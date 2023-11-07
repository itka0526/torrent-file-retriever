<script lang="ts">
    import { wrapper } from "../utils/persist";
    import { socket } from "../utils/store";
    import List from "./main/list.svelte";
    import Upload from "./main/upload.svelte";

    $socket &&
        wrapper($socket, () => {
            if ($socket) {
                $socket.send("get_files");
            }
        });
</script>

<div class="w-full py-12" />
<div class="h-screen w-screen flex">
    <div class="basis-1/4">
        <Upload />
    </div>
    <div class="basis-3/4">
        <List />
    </div>
</div>

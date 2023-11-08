<script lang="ts">
    import { UploadCloudIcon } from "svelte-feather-icons";
    import axios from "axios";
    import { toast } from "@zerodevx/svelte-toast";
    import { type Message } from "../../types";

    const handleUpload = async (
        e: Event & {
            currentTarget: EventTarget & HTMLInputElement;
        }
    ) => {
        let files = (e.target as HTMLInputElement)?.files;

        if (files) {
            let formData = new FormData();
            const file_names = [];

            for (const f of files) {
                file_names.push(f.name);
                formData.append(f.name, f);
            }

            formData.append("file_names", JSON.stringify(file_names));

            axios.defaults.headers.post["Access-Control-Allow-Origin"] = "*";
            const pending_promise = await axios.post<Message>("/api/upload", formData, {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
            });
            files = null;

            toast.push(pending_promise.data.message);
        }
    };
</script>

<form class="cursor-pointer">
    <input type="file" class="hidden" id="input-file" on:change={handleUpload} />
    <label for="input-file" class="flex justify-center items-center cursor-pointer">
        <div class="flex justify-center items-center gap-2 px-4 py-2 bg-gray-200 shadow-md">
            <UploadCloudIcon strokeWidth={1} />
            <span>Upload</span>
        </div>
    </label>
</form>

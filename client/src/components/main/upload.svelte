<script lang="ts">
    import { UploadCloudIcon } from "svelte-feather-icons";
    import axios from "axios";
    import { toast } from "@zerodevx/svelte-toast";
    import { uploadingFiles } from "../../utils/store";
    import type { Message } from "../../types";

    const handleUpload = async (
        e: Event & {
            currentTarget: EventTarget & HTMLInputElement;
        }
    ) => {
        let files = (e.target as HTMLInputElement)?.files;

        if (files && files?.length >= 1) {
            let formData = new FormData();
            const fileNames: string[] = [];

            for (const f of files) {
                fileNames.push(f.name);
                formData.append(f.name, f);
            }

            formData.append("fileNames", JSON.stringify(fileNames));

            const pendingResponse = await axios.post<Message>("/api/upload", formData, {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
                onUploadProgress: (event) => {
                    $uploadingFiles[fileNames.toString()] = event;
                },
            });

            files = null;
            uploadingFiles.update((uf) => {
                delete uf[fileNames.toString()];
                return uf;
            });
            toast.push(pendingResponse.data.message);
        }
    };
</script>

<form class="cursor-pointer">
    <input type="file" class="hidden" id="input-file" on:change={handleUpload} multiple />
    <label for="input-file" class="flex justify-center items-center cursor-pointer">
        <div class="flex justify-center items-center gap-2 px-4 py-2 bg-gray-200 shadow-md">
            <UploadCloudIcon strokeWidth={1} />
            <span>Upload</span>
        </div>
    </label>
</form>

import type { AxiosProgressEvent } from "axios";

export type Message = {
    status: boolean;
    message: string;
};

export type MyFileInfo = {
    path: string;
    name: string;
    size: number;
    modified_date: string;
    is_directory: boolean;
};

export type WSMessage = {
    response_type: "get_files_res";
    data: string;
};

export type UploadingFiles = Record<string, AxiosProgressEvent>;
export type DownloadingFiles = Record<string, AxiosProgressEvent>;

export type ParsedData = MyFileInfo[];

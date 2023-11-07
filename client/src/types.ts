export type Message = {
    status: boolean;
    message: string;
};

export type MyFileInfo = {
    name: string;
    size: number;
    modified_date: string;
};

export type WSMessage = {
    response_type: "get_files_res";
    data: string;
};

export type ParsedData = MyFileInfo[];

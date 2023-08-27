interface FileInfo {
    Length: number;
    Path: string[];
    PathUtf8: string[];
}

export interface ItemInfo {
    PieceLength: number;
    Pieces: string;
    Name: string;
    NameUtf8: string;
    Length: number;
    Private?: boolean;
    Source: string;
    Files: FileInfo[];
}

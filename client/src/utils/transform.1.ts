import type { MyFileInfo } from "../types";

export const transform = (s: MyFileInfo[]) => {
    const result = {};

    s.forEach((item: any) => {
        const parts: any = item.path.split("/").filter((part: any) => part !== "." && part !== "");
        let currentLevel = result;
        parts.forEach((part: string | number, index: number) => {
            if (!currentLevel[part]) {
                if (index === parts.length - 1 && item.is_directory) {
                    currentLevel[part] = { children: {} };
                } else if (index === parts.length - 1 && !item.is_directory) {
                    currentLevel[part] = { name: item.name };
                } else {
                    currentLevel[part] = { children: {} };
                }
            }
            currentLevel = currentLevel[part].children;
        });
    });
    console.log(result);
    return result;
};

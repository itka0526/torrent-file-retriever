import type { MyFileInfo } from "../types";

export const transform = (mfis: MyFileInfo[]) => {
    const result = {};

    for (const item of mfis) {
        const paths: string[] = item.path.split("/").filter((part: any) => part !== "." && part !== "");
        let curr: any = result;
        for (let index = 0; index < paths.length; index++) {
            const part = paths[index];
            if (!curr[part]) {
                if (index === paths.length - 1 && item.is_directory) {
                    curr[part] = { dummy: {}, ...item };
                } else if (index === paths.length - 1 && !item.is_directory) {
                    curr[part] = item;
                }
            }
            curr = curr[part].dummy;
        }
    }

    let otoa = (obj: any) => {
        if (typeof obj !== "object") return;

        for (let k1 in obj) {
            if (typeof obj[k1] == "object" && k1 == "dummy") {
                obj["children"] = [];
                for (let k2 in obj[k1]) {
                    obj["children"].push(obj[k1][k2]);
                }
            }
            otoa(obj[k1]);
            delete obj["dummy"];
        }
    };

    otoa(result);
    return Object.values(result)[0] as MyFileInfo & { children: MyFileInfo[] };
};

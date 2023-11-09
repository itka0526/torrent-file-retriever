import { toast } from "@zerodevx/svelte-toast";

export const newToast = (m: string) => {
    if (m.length > 30) m = m.substring(0, 30) + "...";
    toast.push(m);
};

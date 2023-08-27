import { writable } from "svelte/store";
import type { ItemInfo } from "./types";

export const store = writable((localStorage.getItem("store") ? JSON.parse(localStorage.getItem("store") as string) : []) as ItemInfo[]);

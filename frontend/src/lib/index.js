import { readable } from "svelte/store";

// place files you want to import through the `$lib` alias in this folder.
export const api_url = readable("http://localhost:8000/api")

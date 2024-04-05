<script>
    import { NewsSource } from "$lib/rss.js";
    import { NewsSourcefromURL } from "$lib/rss";
    /**
     * @type {NewsSource}
     */
    export let filter;
</script>

<div class="navbar bg-neutral p-5 font-serif">
    <div class="flex-1 items-center justify-evenly">
        <a class="btn btn-accent text-xl" href="/">Home Page</a>
        <a class="btn btn-accent text-xl" href="/shorts">Short News</a>
        <select class="select bg-gray-300 text-black text-xl" bind:value={filter}>
            <option selected value="">Filter</option>
            {#await NewsSourcefromURL()}
                <p>Waiting...</p>
            {:then options} 
                {#each options as option}
                    <option value={option.pubID}>{option.pub}</option>
                {/each}
            {/await}
        </select>
        <h1 class="text-4xl w-1/4 font-extrabold">Home Page</h1>
        <a class="btn btn-accent w-1/12 text-xl" href="/rss">Your RSS</a>
    </div>
</div>
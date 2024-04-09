<script>
    import { NewsSource } from "$lib/rss.js";
    /** @type {string}*/
    export let src_filter;
    /** @type {NewsSource[]}*/
    export let src_list;
</script>
<div class="navbar bg-neutral p-5 font-serif">
    <div class="flex-1 items-center justify-evenly gap-5">
        <a class="btn w-[10%] btn-accent text-xl" href="/">Home Page</a>
        <a class="btn w-[10%] btn-accent text-xl" href="/shorts">Short News</a>
        <select class="select w-[25%] bg-gray-300 text-black text-xl" bind:value={src_filter}>
            <option selected value="">Filter</option>
            {#each src_list as news_src}
            <option value={news_src.pubID}>{news_src.link}</option>
            {/each}
        </select>
        <h1 class="text-4xl w-[25%] font-extrabold">Home Page</h1>
        <a class="btn btn-accent w-[10%] text-xl" href="/rss">Your RSS</a>
    </div>
</div>
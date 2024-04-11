<script>
    import { Article, articleFromApi } from "$lib/article";
    import { api_url } from "$lib";
    import { newsSourcefromURL } from "$lib/rss";
    /**
     * @type {Article | undefined}
     */
    export let clickedArticle;
    /** @type {string}*/
    export let pub;
    let cur_page = 0;
</script>

{#await newsSourcefromURL()}
<span class="loading loading-spinner loading-lg"></span>
{:then newsSrc} 
    <div class="mx-auto h-[10%] flex justify-center items-center">
        <div class="join">
            <button class="join-item btn" on:click={()=>{cur_page = cur_page-1<0? cur_page: cur_page-1}}>«</button>
            <button class="join-item btn text-xl">Page {cur_page+1}</button>
            <button class="join-item btn" on:click={()=>{cur_page++}}>»</button>
        </div>
    </div>
    <div class="grid h-[90%] overflow-y-scroll gap-5">
    {#await articleFromApi(10, cur_page, pub, true)}
        wait
    {:then articles}
        {#each articles as article}
        <div aria-hidden="true" class="card card-side bg-neutral sm:w-full lg:w-3/4 justify-center mx-auto btn-ghost" 
            on:click={() => {clickedArticle=article}}>
            <div class="p-4 w-1/3 card-title">
                <img class="w-full" src={$api_url+`/articles/thumbnail/${article.id}`} alt={article.title}/>
            </div>
            <div class="p-4 w-2/3">
                <h2 class="card-title">{article.title}</h2>
                <p class="badge badge-secondary">
                {#if newsSrc.filter(data => data.pubID === article.pubid)[0] !== undefined}
                    {newsSrc.filter(data => data.pubID === article.pubid)[0].link}                        
                {/if}
                </p>
                <p class="pt-5 text-justify">{article.summary}</p>
            </div>
        </div>
        {/each}
    {/await}
    </div>
{/await}

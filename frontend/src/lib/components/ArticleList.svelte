<script>
    import { Article, articleFromApi } from "$lib/article";
    import { api_url } from "$lib";
    /**
     * @type {Article | undefined}
     */
    export let clickedArticle;
    let cur_page = 0;
</script>
<div class="mx-auto w-1/2 h-[10%] flex justify-center items-center">
    <div class="join">
        <button class="join-item btn" on:click={()=>{cur_page = cur_page-1<0? cur_page: cur_page-1}}>«</button>
        <button class="join-item btn text-xl">Page {cur_page+1}</button>
        <button class="join-item btn" on:click={()=>{cur_page++}}>»</button>
    </div>
</div>
<div class="h-[90%] overflow-y-scroll">
{#await articleFromApi(cur_page)}
    wait
{:then articles} 
{#each articles as article}
<div aria-hidden="true" class="card bg-neutral mx-5 mb-10 shadow-xl btn-ghost" on:click={() => {clickedArticle=article; console.log(clickedArticle)}}>
    <div class="p-4 card-title flex flex-row justify-left">
        <figure class="w-fit h-fit">
            <img src={$api_url+`/articles/thumbnail/${article.id}`} alt=""/>
        </figure>
        <h2 class="text-2xl w-2/3">{article.title}</h2>
    </div>
    <div class="p-4">
        <p>{article.summary}</p>
    </div>
</div>
{/each}
{/await}
</div>

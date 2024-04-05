<script>
    import ArticleList from "$lib/components/ArticleList.svelte";
    import HomeBar from "$lib/components/HomeBar.svelte";
    import { NewsSource } from "$lib/rss.js";
    import { Article } from "$lib/article";
    import { api_url } from "$lib";

    /**
     * @type {NewsSource}
     */
    let sourceFilter;
    /**
     * @type {Article | undefined}
     */
    let article;
</script>

<div class="flex flex-col w-full h-screen">
    <div class="grid card h-[10%]">
        <HomeBar filter={sourceFilter}/>
    </div>
    <div class="flex card flex-row gap-10 p-5 h-[90%]">
        <div class="border w-1/2">
            <ArticleList bind:clickedArticle={article}/>
        </div>
        <div class="border w-1/2 h-full">
            {#if article !== undefined}
                <div class="h-[10%] flex flex-row items-center justify-evenly">
                    <h1 class="text-4xl">🔊</h1>
                    <audio controls autoplay={false} src={$api_url+`/articles/audio/${article.id}`}>
                    </audio>
                    <a class=" btn btn-primary" href={article.link}>Go to site</a>
                    <button on:click={() => {article=undefined}} class="text-4xl font-extrabold">X</button>
                </div>
                <embed src={$api_url+`/articles/full/${article.id}`} height="fit" class="w-full h-[90%]">
            {/if}
        </div>
    </div>
</div>
<script>
    import HomeBar from "$lib/components/HomeBar.svelte";
    import ArticleList from "$lib/components/ArticleList.svelte";
    import ArticleFull from "$lib/components/ArticleFull.svelte";
    import { newsSourcefromURL } from "$lib/rss.js";
    import { Article } from "$lib/article";
    /** @type {string}*/
    let publisher = "";
    /** @type {Article | undefined}*/
    let article;
</script>

{#await newsSourcefromURL()}
<span class="loading loading-spinner loading-lg"></span>
{:then newsSrc} 
<div class="flex flex-col w-full h-screen">
    <div class="grid card h-[10%]">
        <HomeBar src_list={newsSrc} bind:src_filter={publisher}/>
    </div>
    <div class="flex card flex-row gap-10 h-[90%]">
        <div class="w-full">
            <ArticleList pub={publisher} bind:clickedArticle={article}/>
        </div>
        <ArticleFull article={article}/>
    </div>
</div>
{/await}

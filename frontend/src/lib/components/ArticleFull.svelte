<script>
    import { api_url } from "$lib";
    import { Article } from "$lib/article";
    /** @type {Article | undefined}} */
    export let article;
</script>

{#if article !== undefined}
<div class="absolute w-full mx-auto h-full">
    <div class="w-5/6 mx-auto h-full jusify-center">
        <div class="h-[10%] flex flex-row items-center justify-between">
            <div class="flex">
                <h1 class="text-4xl">🔊</h1>
                <audio controls autoplay={false} src={$api_url+`/articles/audio/${article.id}`}>
                </audio>
            </div>
            <a class=" btn btn-primary" href={article.link}>Go to site</a>
            <button on:click={() => {article=undefined}} class="text-4xl font-extrabold">✕</button>
        </div>
        <embed src={$api_url+`/articles/full/${article.id}`} height="fit" class="w-full h-[90%]">
    </div>
</div>
{/if}
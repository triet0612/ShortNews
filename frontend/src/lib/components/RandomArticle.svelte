<script>
    import { randomArticle } from "$lib/article";
    import { api_url } from "$lib";
    import { newsSourcefromURL } from "$lib/rss";

    /**
     * @type {HTMLParagraphElement}
     */
    let title;
    /**
     * @param {WheelEvent} event
     */
    function scrollLeftRight(event) {
        title.scrollBy(event.deltaY, 0)
    }
    /**
     * @type {any}
     */
     export let load;
</script>

{#key load}
{#await newsSourcefromURL() }
    wait
{:then newsSrc} 
    {#await randomArticle()}
        wait
    {:then article}
        <div class="bg-neutral h-[100%] p-5 rounded-3xl mx-auto justify-center">
            <h1 class="h-[5%] text-xl font-bold text-nowrap overflow-hidden hover:overflow-scroll"
                on:wheel={scrollLeftRight} bind:this={title}>
                {article.title}
            </h1>
            <p class="h-[3%] badge badge-secondary my-3">
            {#if newsSrc.filter(data => data.pubID === article.pubid)[0] !== undefined}
                {newsSrc.filter(data => data.pubID === article.pubid)[0].pub}                        
            {/if}
            </p>
            <a class="btn bg-neutral h-[48%] w-full rounded-3xl" href={article.link}>
                <img class="h-[100%] w-[100%] object-fill" src={$api_url+`/articles/thumbnail/${article.id}`} alt={`image: ${article.title}`}/>
            </a>
            <div class="relative p-4 h-[35%] overflow-y-scroll">
                <p class="text-justify text-lg">{article.summary}</p>
            </div>
            <audio class="h-[5%] mt-2 w-full rounded-3xl" controls autoplay={true} src={$api_url+`/articles/audio/${article.id}`} on:ended={()=>{location.reload()}}>
        </div>
    {/await}
{/await}
{/key}
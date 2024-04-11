<script>
    import { articleFromApi } from "$lib/article";
    import { api_url } from "$lib";
    import { newsSourcefromURL } from "$lib/rss";

    /*** @type {HTMLParagraphElement}*/
    let titleElem;
    /*** @param {WheelEvent} event*/
    function scrollLeftRight(event) {
        titleElem.scrollBy(event.deltaY, 0)
    }
    export let cur_page = 0;
</script>

{#await newsSourcefromURL()}
    wait
{:then newsSrc} 
    {#await articleFromApi(1, cur_page, "", true, true)}
        wait
    {:then article}
    {#each article as a}
    <div class="bg-neutral h-[100%] p-5 rounded-3xl mx-auto justify-center">
        <h1 class="h-[5%] text-xl font-bold text-nowrap overflow-hidden hover:overflow-scroll"
            on:wheel={scrollLeftRight} bind:this={titleElem}>
            {a.title}
        </h1>
        <p class="h-[3%] badge badge-secondary my-3">
        {#if newsSrc.filter(data => data.pubID === a.pubid)[0] !== undefined}
            {newsSrc.filter(data => data.pubID === a.pubid)[0].link}                        
        {/if}
        </p>
        <button class="btn bg-neutral h-[48%] w-full rounded-3xl" on:click={()=> {window.open(a.link, '_blank')?.focus();}}>
            <img class="h-[100%] w-[100%] object-fill" src={$api_url+`/articles/thumbnail/${a.id}`} alt={`image: ${a.title}`}/>
        </button>
        <div class="relative p-4 h-[35%] overflow-y-scroll">
            <p class="text-justify text-lg">{a.summary}</p>
        </div>
        <audio class="h-[5%] mt-2 w-full rounded-3xl" controls autoplay={true} src={$api_url+`/articles/audio/${a.id}`} on:ended={()=>{cur_page++}}>
    </div>
    {/each}
    {/await}
{/await}


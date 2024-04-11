<script>
    import RandomArticle from "$lib/components/RandomArticle.svelte";
    import NavBar from "$lib/components/NavBar.svelte";

    let page = 0;
</script>
<div class="flex flex-col w-full h-screen">
    <div class="grid card h-[10%]">
        <NavBar page_name="Shorts"/>
    </div>
    <div class="h-[90%] w-2/5 p-3 mx-auto justify-center">
        <RandomArticle bind:cur_page={page}/>
        <button class="absolute bg-neutral bottom-5 right-5 hover:opacity-50 text-3xl rounded-full w-20 h-20" 
            on:click={() => {page++}}>↓</button>
    </div>
</div>
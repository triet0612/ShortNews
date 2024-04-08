<script>
    import CreateForm from "$lib/components/CreateForm.svelte";
    import NavBar from "$lib/components/NavBar.svelte";
    import { deleteSource, newsSourcefromURL } from "$lib/rss";
</script>

<div class="flex flex-col w-full h-screen">
    <div class="grid card h-[10%]">
        <NavBar page_name="Shorts"/>
    </div>
    <div class="h-[90%] w-full flex">
        <div class="w-1/2 fill-black">
            <h1 class="p-5 text-center text-3xl font-extrabold">Your RSS sources</h1>
        {#await newsSourcefromURL()}
            wait
        {:then rssSource}
            {#each rssSource as rss}
            <div class="flex flex-row justify-center p-5">
                <div class="w-10/12 flex h-full gap-4">
                    <div class="w-2/3 flex flex-row bg-neutral border-2 rounded-xl items-center font-serif text-lg">
                        <p class="pl-2 py-2 w-1/4 font-bold text-base text-accent">RSS Link: </p>
                        <p class="py-2 w-3/4 hover:text-white"><a href={rss.link}>{rss.link}</a></p>
                    </div>
                    <div class="w-1/4 flex flex-row bg-neutral border-2 rounded-xl text-center items-center justify-evenly font-serif text-lg">
                        <p class="py-4 text-base text-accent font-bold">Language: </p>                    
                        <p class="py-4">{rss.lang}</p>
                    </div>
                </div>
                <div class="w-2/12 flex gap-2">
                    <button on:click={async () => {await deleteSource(rss.pubID); location.reload();}} class="btn h-full hover:bg-red-700">
                        <svg class="fill-white" xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 -960 960 960" width="24"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg>
                    </button>
                </div>
            </div>
            {/each}
        {/await}
        </div>
        <div class="w-1/2">
            <CreateForm />
        </div>
    </div>
</div>
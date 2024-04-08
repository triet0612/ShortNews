<script>
    import { NewsSource, createSource } from "$lib/rss";
    /**
     * @type {NewsSource}
     */
    let rss = new NewsSource("", "", "");
</script>
<form>
    <h1 class="p-5 text-center text-3xl font-extrabold">Create a new RSS:</h1>
    <div class="flex flex-col gap-10 text-center py-5 px-10">
        <label class="w-full text-2xl">
            <input bind:value={rss.link} type="text" class="w-full text-center p-4 bg-neutral rounded-xl" placeholder="RSS Link:"/>
        </label>
        <select bind:value={rss.lang} class="p-5 text-2xl w-full text-center bg-neutral rounded-xl">
            <option disabled selected class="">Language:</option>
            <option value="vi">Vietnamese</option>
            <option value="en">English</option>
        </select>
        <button class="btn btn-accent w-1/2 mx-auto text-2xl" on:click={async () => {await createSource(rss); location.reload();}}>
            Submit
        </button>
    </div>
</form>
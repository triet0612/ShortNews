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
        <select bind:value={rss.voiceType} class="p-5 text-2xl w-full text-center bg-neutral rounded-xl">
            <option disabled selected class="">Voice Type:</option>
            <option value="model_1">model_1</option>
            <option value="model_2">model_2</option>
            <option value="model_2">model_2</option>
        </select>
        <button class="btn btn-accent w-1/2 mx-auto text-2xl" on:click={async () => {await createSource(rss); location.reload();}}>
            Submit
        </button>
    </div>
</form>
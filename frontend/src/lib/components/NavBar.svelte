<script>
    export let page_name = ""
</script>

<div class="navbar bg-neutral p-5 font-serif">
    <div class="flex-1 items-center justify-evenly gap-5">
        <a class="btn w-[10%] btn-accent text-xl" href="/">Home Page</a>
        <a class="btn w-[10%] btn-accent text-xl" href="/shorts">Short News</a>
        <div class="w-[25%] bg-gray-300 text-black text-xl">
        </div>
        <h1 class="text-4xl w-[25%] font-extrabold">{page_name}</h1>
        <a class="btn btn-accent w-[10%] text-xl" href="/rss">Your RSS</a>
    </div>
</div>
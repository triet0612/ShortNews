<script>
    export let page_name = ""
</script>

<div class="navbar bg-neutral p-5 font-serif">
    <div class="flex-1 items-center justify-evenly">
        <a class="btn btn-accent text-xl" href="/">Home Page</a>
        <a class="btn btn-accent text-xl" href="/shorts">Short News</a>
        <h1 class="text-4xl w-1/4 font-extrabold">{page_name}</h1>
        <a class="btn btn-accent w-1/12 text-xl" href="/rss">Your RSS</a>
    </div>
</div>
<script>
    import Header from "./Header.svelte";
    // import './styles.css';
</script>

<div class="h-screen w-screen flex flex-col">
    <Header />

    <main class="flex flex-1 flex-col justify-center items-center">
        <slot />
    </main>
</div>

<style>
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
</style>

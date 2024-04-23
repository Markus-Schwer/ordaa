<script>
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import { token } from "$lib/auth";
    import {
        Button,
        Navbar,
        NavBrand,
        NavLi,
        NavUl,
        NavHamburger,
    } from "flowbite-svelte";
    $: activeUrl = $page.url.pathname;

    function logout() {
        token.set(null);
        goto("/");
    }
</script>

<Navbar>
    <NavBrand href="/">
        <span
            class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
            >Dotinder</span
        >
    </NavBrand>
    <NavHamburger />
    <NavUl {activeUrl}>
        <NavLi href="/">Home</NavLi>
        <NavLi href="/menus">Menus</NavLi>
        <NavLi href="/orders">Orders</NavLi>
    </NavUl>
    {#if $token != null}
        <Button size="sm" on:click={logout}>Logout</Button>
    {:else}
        <Button size="sm"><a href="/login">Login</a></Button>
    {/if}
</Navbar>

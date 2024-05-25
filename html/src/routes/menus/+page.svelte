<script>
    import { get } from "svelte/store";
    import { token } from "$lib/auth";
    import { Alert, Spinner } from "flowbite-svelte";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import axios from "$lib/api";

    onMount(() => {
        if (get(token) == null) {
            goto("/login");
        }
    });

    let promise = getMenus();

    async function getMenus() {
        return await axios.get("/api/menus").then((res) => {
            return res.data;
        }).catch((err) => {
            throw new Error(err);
        });
    }
</script>

{#await promise}
    <Spinner />
{:then menus}
    <ul>
        {#each menus as menu}
            <a href="/menus/{menu.uuid}">{menu.name}</a>
        {/each}
    </ul>
{:catch error}
    <Alert>
        <span class="font-medium">{error.message.status}</span>
    </Alert>
{/await}

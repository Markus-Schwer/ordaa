<script>
    import { get } from "svelte/store";
    import { token } from "$lib/auth";
    import { Alert, Spinner } from "flowbite-svelte";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    onMount(() => {
        if (get(token) == null) {
            goto("/login");
        }
    });

    let promise = getMenus();

    async function getMenus() {
        const res = await fetch("http://localhost:8080/api/menus", {
            headers: {
                Authorization: `Bearer ${get(token)}`,
            },
        });
        const resObj = await res.json();

        if (res.ok) {
            return resObj;
        } else {
            throw new Error(res);
        }
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

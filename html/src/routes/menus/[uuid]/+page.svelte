<script>
    import { onMount } from "svelte";
    import { page } from "$app/stores";

    async function getMenus(uuid) {
        const res = await fetch(`/api/menus/${uuid}`);
        let menuObj = await res.json();

        if (res.ok) {
            menuObj.items = menuObj.items.sort((a, b) =>
                a.short_name.localeCompare(b.short_name, undefined, {
                    numeric: true,
                    sensitivity: "base",
                })
            );
            return menuObj;
        } else {
            throw new Error(menuObj);
        }
    }

    onMount(async () => {
        console.log(await getMenus($page.params.uuid));
    });
</script>

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Price</th>
            <th>Id</th>
        </tr>
    </thead>

    {#await getMenus($page.params.uuid)}
        <p>...waiting</p>
    {:then menu}
        {#each menu.items as item}
            <tr>
                <td>{item.name}</td>
                <td>{item.price}</td>
                <td>{item.short_name}</td>
            </tr>
        {/each}
        <!-- <ul> -->
        <!--     {#each menus as menu} -->
        <!--         <a href="/menus/{menu.uuid}">{menu.name}</a> -->
        <!--     {/each} -->
        <!-- </ul> -->
    {:catch error}
        <p style="color: red">{error.message}</p>
    {/await}
</table>

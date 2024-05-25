<script>
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import axios from "$lib/api";

    async function getMenu(uuid) {
        return await axios.get(`/api/menus/${uuid}`).then((res) => {
            let menuObj = res.data;

            menuObj.items = menuObj.items.sort((a, b) =>
                a.short_name.localeCompare(b.short_name, undefined, {
                    numeric: true,
                    sensitivity: "base",
                })
            );
            return menuObj;
        }).catch((err) => {
            throw new Error(menuObj);
        });
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

    {#await getMenu($page.params.uuid)}
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

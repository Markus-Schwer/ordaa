<script>
    import { Button, Card, Label, Input } from "flowbite-svelte";
    import { token } from "$lib/auth";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { get } from "svelte/store";
    import axios from "$lib/api";

    let username, password;

    onMount(() => {
        if (get(token) != null) {
            goto("/");
        }
    });

    async function handleLogin() {
        return await axios.post("/api/login", {
            username,
            password,
        }).then((res) => {
            token.set(res.data.jwt);
            goto("/");
        }).catch((err) => {
            alert(err);
        });
    }
</script>

<Card class="flex flex-col space-y-6">
    <h3 class="text-xl font-medium text-gray-900 dark:text-white">Login</h3>
    <Label for="username" class="space-y-2">
        <span>Username</span>
        <Input
            name="username"
            placeholder="username"
            bind:value={username}
            required
        />
    </Label>
    <Label for="password" class="space-y-2">
        <span>Password</span>
        <Input
            name="password"
            placeholder="password"
            type="password"
            bind:value={password}
            required
        />
    </Label>
    <Button on:click={handleLogin} class="w-full">Login</Button>
    <div class="text-sm">
        Not registered? <a href="/register" class="hover:underline text-primary-700">
            Create account
        </a>
    </div>
</Card>

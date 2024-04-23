<script>
    import { Button, Card, Label, Input } from "flowbite-svelte";
    import { token } from "$lib/auth";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { get } from "svelte/store";

    let username, password;

    onMount(() => {
        if (get(token) != null) {
            goto("/");
        }
    });

    async function handleLogin() {
        const res = await fetch(`http://localhost:8080/api/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username,
                password,
            }),
        });
        if (res.ok) {
            const body = await res.json();
            token.set(body.jwt);
            goto("/");
        } else {
            const text = await res.text();
            alert(text);
        }
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

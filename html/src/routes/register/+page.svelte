<script>
    import { Button, Card, Helper, Label, Input } from "flowbite-svelte";
    import { token } from "$lib/auth";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { get } from "svelte/store";
    import axios from "$lib/api";

    let username, password, password2, passwordsMatch;
    $: passwordsMatch = password === password2;

    onMount(() => {
        if (get(token) != null) {
            goto("/");
        }
    });

    async function handleRegister() {
        if (!passwordsMatch) {
            return;
        }
        return await axios.post(`/api/users`, {
            username,
            password,
        }).then((res) => {
            goto("/login");
        }).catch((err) => {
            alert(err);
        });
    }
</script>

<Card class="flex flex-col space-y-6">
    <h3 class="text-xl font-medium text-gray-900 dark:text-white">Register</h3>
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
    <Label for="password2" class="space-y-2">
        <span>Repeat Password</span>
        <Input
            name="password2"
            placeholder="password"
            type="password"
            bind:value={password2}
            required
            valid={passwordsMatch}
            color={passwordsMatch ? "" : "red"}
        />
    </Label>
    {#if !passwordsMatch}
        <Helper class="mt-2" color="red">
            <span class="font-medium">Oh, snapp!</span>
            Passwords don't match.
        </Helper>
    {/if}
    <Button on:click={handleRegister} class="w-full">Register</Button>
    <div class="text-sm">
        Already registered? <a
            href="/login"
            class="hover:underline text-primary-700">Login</a
        >
    </div>
</Card>

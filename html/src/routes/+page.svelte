<script>
    import { token } from "$lib/auth";

    import { Alert } from 'flowbite-svelte';

    let username, password;
    let register = false;

    async function handleLogin() {
        const res = await fetch(
            `/api/${register ? "users" : "login"}`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username,
                    password,
                }),
            }
        );
        if (res.ok) {
            const body = await res.json();
            token.set(body.jwt);
        } else {
            const text = await res.text();
            alert(text);
        }
    }
</script>
<div class="p-8">
  <Alert>
    <span class="font-medium">Info alert!</span>
    Change a few things up and try submitting again.
  </Alert>
</div>

<div>
    <label for="username">Username</label>
    <input type="text" bind:value={username} name="username" />
    <label for="password">Password</label>
    <input type="password" bind:value={password} name="password" />
    <label for="register">Register</label>
    <input type="checkbox" bind:checked={register} name="register" />
    <button on:click={handleLogin}>Login</button>
</div>

<p>token: {$token}</p>

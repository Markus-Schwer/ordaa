<script>
    import { token } from "$lib/auth";
    import axios from "$lib/api";

    import { Alert } from 'flowbite-svelte';

    let username, password;
    let register = false;

    async function handleLogin() {
        return await axios.post(`/api/${register ? "users" : "login"}`, {
            username,
            password,
        }).then((res) => {
            // TODO: this won't work when registering a user...
            token.set(res.data.jwt);
        }).catch((err) => {
            alert(err);
        });
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

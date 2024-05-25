import axios from "axios";
import { get } from "svelte/store";
import { goto } from "$app/navigation";
import { token } from "$lib/auth";

const instance = axios.create();

instance.defaults.headers.post['Content-Type'] = 'application/json';
instance.interceptors.request.use(
    (config) => {
        let storedToken = get(token);

        if (storedToken) {
            config.headers["Authorization"] = `Bearer ${storedToken}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

instance.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        if (error.response?.status === 401) {
            // logout
            token.set(null);
            goto("/login");
            return Promise.reject(error);
        }

        return Promise.reject(error);
    }
);

export default instance;

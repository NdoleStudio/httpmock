import axios from "axios";

const client = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:8000",
  headers: {
    "X-Client-Version": process.env.NEXT_PUBLIC_GITHUB_SHA || "dev",
  },
});

export function setAuthHeader(token: string | null) {
  client.defaults.headers.common.Authorization = "Bearer " + token;
}

export default client;

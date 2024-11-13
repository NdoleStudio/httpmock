"use client";

import axios from "axios";

const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:8000",
  headers: {
    "X-Client-Version": process.env.NEXT_PUBLIC_GITHUB_SHA || "dev",
  },
});

axiosInstance.interceptors.request.use(async (request) => {
  const token = await (window as any).Clerk.session.getToken();
  if (token) {
    request.headers.Authorization = `Bearer ${token}`;
  }
  return request;
});

export default axiosInstance;

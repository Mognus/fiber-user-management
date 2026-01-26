// Server-side Auth API functions
// Uses the shared serverFetch utility

import { serverFetch } from "@/lib/api/api-server";
import type { User } from "../types";

const API_URL = process.env.API_URL || "http://backend:8080/api";

export async function fetchMe(): Promise<User> {
    return serverFetch<User>("/auth/me", { withAuth: true });
}

export async function fetchMeWithToken(token?: string): Promise<User | null> {
    if (!token) return null;
    const res = await fetch(`${API_URL}/auth/me`, {
        headers: { Authorization: token },
    });
    if (!res.ok) return null;
    return res.json();
}

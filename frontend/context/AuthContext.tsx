"use client";

import { createContext, useContext, useMemo, useState } from "react";
import type { User } from "../types";

interface AuthContextValue {
    user: User | null;
    setUser: (user: User | null) => void;
    isAdmin: boolean;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

interface AuthProviderProps {
    initialUser?: User | null;
    children: React.ReactNode;
}

export function AuthProvider({ initialUser = null, children }: AuthProviderProps) {
    const [user, setUser] = useState<User | null>(initialUser);
    const value = useMemo<AuthContextValue>(
        () => ({
            user,
            setUser,
            isAdmin: user?.role?.name === "admin",
        }),
        [user]
    );

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
    const ctx = useContext(AuthContext);
    if (!ctx) {
        throw new Error("useAuth must be used within AuthProvider");
    }
    return ctx;
}

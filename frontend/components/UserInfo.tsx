"use client";

import { useAuth } from "../context/AuthContext";

export function UserInfo() {
    const { user } = useAuth();

    if (!user) {
        return (
            <div className="text-sm text-muted-foreground">
                Not signed in
            </div>
        );
    }

    const name = [user.first_name, user.last_name].filter(Boolean).join(" ");
    const displayName = name || user.email;

    return (
        <div className="text-sm">
            <span className="font-medium">{displayName}</span>
            <span className="text-muted-foreground"> · {user.role}</span>
        </div>
    );
}

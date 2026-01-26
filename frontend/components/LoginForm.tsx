"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { authAPI } from "../lib/api";
import { useAuth } from "../context/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";

const loginSchema = z.object({
    email: z.string().email("Bitte eine gültige E-Mail eingeben."),
    password: z.string().min(1, "Passwort ist erforderlich."),
});

type LoginValues = z.infer<typeof loginSchema>;

export function LoginForm() {
    const router = useRouter();
    const { setUser } = useAuth();
    const [error, setError] = useState<string | null>(null);
    const form = useForm<LoginValues>({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    });

    const handleSubmit = async (values: LoginValues) => {
        setError(null);
        try {
            const response = await authAPI.login(values);
            setUser(response.user);
            router.replace("/");
        } catch (err: any) {
            setError(err?.message || "Login failed");
        }
    };

    return (
        <div className="min-h-screen flex flex-col">
            <div className="flex-1 min-h-0 flex items-center justify-center p-6">
                <Card className="w-full max-w-md">
                    <CardHeader>
                        <CardTitle className="text-2xl">Login</CardTitle>
                        <CardDescription>Sign in with your account.</CardDescription>
                    </CardHeader>
                    <CardContent>

                        {error && (
                            <div className="mt-4 rounded border border-red-200 bg-red-50 p-3 text-sm text-red-700">
                                {error}
                            </div>
                        )}

                        <Form {...form}>
                            <form onSubmit={form.handleSubmit(handleSubmit)} className="mt-6 space-y-4">
                                <FormField
                                    control={form.control}
                                    name="email"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Email</FormLabel>
                                            <FormControl>
                                                <Input type="email" {...field} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <FormField
                                    control={form.control}
                                    name="password"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Password</FormLabel>
                                            <FormControl>
                                                <Input type="password" {...field} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <Button
                                    type="submit"
                                    className="w-full"
                                    disabled={form.formState.isSubmitting}
                                >
                                    {form.formState.isSubmitting ? "Signing in..." : "Sign in"}
                                </Button>
                            </form>
                        </Form>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}

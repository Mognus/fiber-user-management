"use client";

import { useState } from "react";
import { useRouter } from "@/i18n/navigation";
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
import { useTranslations } from "next-intl";

export function LoginForm() {
    const t = useTranslations("Auth.Login");
    const router = useRouter();
    const { setUser } = useAuth();
    const [error, setError] = useState<string | null>(null);
    const loginSchema = z.object({
        email: z.string().email(t("validation.email")),
        password: z.string().min(1, t("validation.password")),
    });

    type LoginValues = z.infer<typeof loginSchema>;
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
            setError(err?.message || t("error.generic"));
        }
    };

    return (
        <div className="min-h-screen flex flex-col">
            <div className="flex-1 min-h-0 flex items-center justify-center p-6">
                <Card className="w-full max-w-md">
                    <CardHeader>
                        <CardTitle className="text-2xl">{t("title")}</CardTitle>
                        <CardDescription>{t("description")}</CardDescription>
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
                                            <FormLabel>{t("fields.email")}</FormLabel>
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
                                            <FormLabel>{t("fields.password")}</FormLabel>
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
                                    {form.formState.isSubmitting
                                        ? t("actions.submitting")
                                        : t("actions.submit")}
                                </Button>
                            </form>
                        </Form>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}

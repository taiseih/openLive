"use client";

import { Button } from "@/components/common/Button";
import { Input } from "@/components/common/Input";
import { Loading } from "@/components/common/Loading";
import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function LoginPage() {
  const router = useRouter();
  const {
    isLoading: authLoading,
    isAuthenticated,
    signIn,
    loginWithEmail,
    registerWithEmail,
  } = useAuth();

  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loading text="読み込み中..." />
      </div>
    );
  }

  if (isAuthenticated) {
    router.push("/");
    return null;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      if (mode === "login") {
        await loginWithEmail(email, password);
      } else {
        await registerWithEmail(email, password, displayName);
      }
      router.push("/");
    } catch (err) {
      setError("認証に失敗しました。入力内容を確認してください。");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md bg-white rounded-xl shadow-lg p-8 space-y-6">
        <h1 className="text-2xl font-bold text-center text-primary-600">
          openLive ログイン
        </h1>

        <div className="flex space-x-2 bg-gray-100 rounded-lg p-1">
          <button
            className={`flex-1 py-2 text-sm font-medium rounded-md ${
              mode === "login"
                ? "bg-white text-primary-600 shadow"
                : "text-gray-500"
            }`}
            onClick={() => setMode("login")}
          >
            ログイン
          </button>
          <button
            className={`flex-1 py-2 text-sm font-medium rounded-md ${
              mode === "register"
                ? "bg-white text-primary-600 shadow"
                : "text-gray-500"
            }`}
            onClick={() => setMode("register")}
          >
            新規登録
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {mode === "register" && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                表示名
              </label>
              <Input
                type="text"
                value={displayName}
                onChange={setDisplayName}
                placeholder="例: 山田 太郎"
                required
              />
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              メールアドレス
            </label>
            <Input
              type="email"
              value={email}
              onChange={setEmail}
              placeholder="you@example.com"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              パスワード
            </label>
            <Input
              type="password"
              value={password}
              onChange={setPassword}
              placeholder="********"
              required
            />
          </div>

          {error && (
            <p className="text-sm text-red-600 bg-red-50 border border-red-100 rounded-md p-2">
              {error}
            </p>
          )}

          <Button
            type="submit"
            disabled={submitting}
            className="w-full justify-center"
          >
            {submitting
              ? "処理中..."
              : mode === "login"
              ? "メールアドレスでログイン"
              : "メールアドレスで登録"}
          </Button>
        </form>

        <div className="flex items-center">
          <div className="flex-1 h-px bg-gray-200" />
          <span className="px-3 text-xs text-gray-400">または</span>
          <div className="flex-1 h-px bg-gray-200" />
        </div>

        <Button
          variant="secondary"
          className="w-full justify-center"
          onClick={signIn}
        >
          Google でログイン
        </Button>

        <p className="text-xs text-gray-500 text-center">
          ログインまたは登録することで、利用規約およびプライバシーポリシーに同意したものとみなされます。
        </p>
      </div>
    </div>
  );
}

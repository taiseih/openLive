"use client";

// Single Responsibility Principle (SRP): 認証フロー（OAuthリダイレクト）のみを担当
import { IAuthService } from "@/types/interfaces";

// Next.js のビルド時にインライン展開される環境変数用の型定義（Node型定義に依存しない）
declare const process: { env: { NEXT_PUBLIC_API_URL?: string } };

const API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

// API_BASE が `http://localhost:8080/api/v1` のような形式である前提で、ベースURLを組み立てる
const resolveAuthUrl = (path: string) => {
  try {
    const url = new URL(API_BASE);
    return `${url.origin}${url.pathname.replace(/\/$/, "")}${path}`;
  } catch {
    // NEXT_PUBLIC_API_URL が不正な場合のフォールバック
    return `http://localhost:8080/api/v1${path}`;
  }
};

export class AuthService implements IAuthService {
  async signInWithGoogle(): Promise<void> {
    if (typeof window === "undefined") return;
    const loginUrl = resolveAuthUrl("/auth/login");
    window.location.href = loginUrl;
  }

  async signInWithAccessToken(token: string): Promise<void> {
    if (typeof window === "undefined") return;
    // アクセストークンを直接指定してログインする方式
    localStorage.setItem("auth_token", token);
  }

  async signOut(): Promise<void> {
    const logoutUrl = resolveAuthUrl("/auth/logout");
    try {
      await fetch(logoutUrl, {
        method: "POST",
        credentials: "include",
      });

      if (typeof window !== "undefined") {
        // トークンをローカルストレージから削除（後続のuseAuthで状態リセット）
        localStorage.removeItem("auth_token");
      }
    } catch (error) {
      console.error("Sign out error:", error);
      throw error;
    }
  }

  async getCurrentToken(): Promise<string | null> {
    if (typeof window === "undefined") return null;
    // OAuth コールバックで保存されたトークンを参照する想定
    return localStorage.getItem("auth_token");
  }

  // 旧Firebase用のAPIとの互換用ダミー実装（現状は未使用）
  // 認証状態は useAuth 内で getCurrentToken を通じて判定する
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onAuthStateChanged(_callback: (user: any) => void): () => void {
    return () => {};
  }
}

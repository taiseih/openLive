// カスタムフック - 認証ロジックの分離
import { authService, httpClient, userService } from "@/services";
import { useAuthStore } from "@/store/authStore";
import { useEffect } from "react";

export const useAuth = () => {
  const {
    user,
    isAuthenticated,
    isLoading,
    setUser,
    setToken,
    setLoading,
    reset,
  } = useAuthStore();

  useEffect(() => {
    // 初期ロード時にトークンとユーザー情報を取得
    const initializeAuth = async () => {
      try {
        const token = await authService.getCurrentToken();
        if (token) {
          setToken(token);
          httpClient.setAuthToken(token);

          const userData = await userService.getCurrentUser();
          setUser(userData);
        } else {
          reset();
        }
      } catch (error) {
        console.error("Failed to get user data:", error);
        reset();
      } finally {
        setLoading(false);
      }
    };

    void initializeAuth();
  }, [setUser, setToken, setLoading, reset]);

  const signIn = async () => {
    try {
      await authService.signInWithGoogle();
    } catch (error) {
      console.error("Sign in error:", error);
      throw error;
    }
  };

  const signInWithAccessToken = async (token: string) => {
    try {
      await authService.signInWithAccessToken(token);

      // トークン反映後にユーザー情報を即座に取得
      setToken(token);
      httpClient.setAuthToken(token);
      const userData = await userService.getCurrentUser();
      setUser(userData);
    } catch (error) {
      console.error("Sign in with access token error:", error);
      reset();
      throw error;
    }
  };

  const loginWithEmail = async (email: string, password: string) => {
    try {
      const { accessToken } = await httpClient.post<{ accessToken: string }>(
        "/auth/login",
        {
          email,
          password,
        }
      );

      await signInWithAccessToken(accessToken);
    } catch (error) {
      console.error("Login with email error:", error);
      reset();
      throw error;
    }
  };

  const registerWithEmail = async (
    email: string,
    password: string,
    displayName: string
  ) => {
    try {
      const { accessToken } = await httpClient.post<{ accessToken: string }>(
        "/auth/register",
        {
          email,
          password,
          displayName,
        }
      );

      await signInWithAccessToken(accessToken);
    } catch (error) {
      console.error("Register with email error:", error);
      reset();
      throw error;
    }
  };

  const signOut = async () => {
    try {
      await authService.signOut();
      reset();
    } catch (error) {
      console.error("Sign out error:", error);
      throw error;
    }
  };

  return {
    user,
    isAuthenticated,
    isLoading,
    signIn,
    signInWithAccessToken,
    loginWithEmail,
    registerWithEmail,
    signOut,
  };
};

// カスタムフック - 認証ロジックの分離
import { useEffect } from 'react';
import { useAuthStore } from '@/store/authStore';
import { authService, userService, httpClient } from '@/services';

export const useAuth = () => {
  const { user, isAuthenticated, isLoading, setUser, setToken, setLoading, reset } = useAuthStore();

  useEffect(() => {
    // 認証状態の監視
    const unsubscribe = authService.onAuthStateChanged(async (firebaseUser) => {
      if (firebaseUser) {
        try {
          // トークン取得
          const token = await authService.getCurrentToken();
          if (token) {
            setToken(token);
            httpClient.setAuthToken(token);

            // ユーザー情報取得
            const userData = await userService.getCurrentUser();
            setUser(userData);
          }
        } catch (error) {
          console.error('Failed to get user data:', error);
          reset();
        }
      } else {
        reset();
      }
      setLoading(false);
    });

    return () => unsubscribe();
  }, [setUser, setToken, setLoading, reset]);

  const signIn = async () => {
    try {
      await authService.signInWithGoogle();
    } catch (error) {
      console.error('Sign in error:', error);
      throw error;
    }
  };

  const signOut = async () => {
    try {
      await authService.signOut();
      reset();
    } catch (error) {
      console.error('Sign out error:', error);
      throw error;
    }
  };

  return {
    user,
    isAuthenticated,
    isLoading,
    signIn,
    signOut,
  };
};


'use client';

// Single Responsibility Principle (SRP): 認証のみに責任を持つ
import { initializeApp, FirebaseApp } from 'firebase/app';
import {
  getAuth,
  Auth,
  GoogleAuthProvider,
  signInWithPopup,
  signOut as firebaseSignOut,
  onAuthStateChanged,
  User,
  connectAuthEmulator,
} from 'firebase/auth';
import { IAuthService } from '@/types/interfaces';

export class AuthService implements IAuthService {
  private auth: Auth;
  private app: FirebaseApp;

  constructor() {
    // Firebase初期化
    this.app = initializeApp({
      apiKey: 'demo-api-key', // Emulator使用時はダミーでOK
      projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID || 'openlive-dev',
    });

    this.auth = getAuth(this.app);

    // Emulator接続
    const emulatorHost = process.env.NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST;
    if (emulatorHost && typeof window !== 'undefined') {
      connectAuthEmulator(this.auth, `http://${emulatorHost}`, { disableWarnings: true });
    }
  }

  async signInWithGoogle(): Promise<void> {
    const provider = new GoogleAuthProvider();
    try {
      await signInWithPopup(this.auth, provider);
    } catch (error) {
      console.error('Sign in error:', error);
      throw error;
    }
  }

  async signOut(): Promise<void> {
    try {
      await firebaseSignOut(this.auth);
    } catch (error) {
      console.error('Sign out error:', error);
      throw error;
    }
  }

  async getCurrentToken(): Promise<string | null> {
    const user = this.auth.currentUser;
    if (!user) return null;

    try {
      const token = await user.getIdToken();
      return token;
    } catch (error) {
      console.error('Get token error:', error);
      return null;
    }
  }

  onAuthStateChanged(callback: (user: User | null) => void): () => void {
    return onAuthStateChanged(this.auth, callback);
  }

  getCurrentUser(): User | null {
    return this.auth.currentUser;
  }
}


'use client';

// Dependency Injection Container (DIP)
// サービスのシングルトンインスタンスを管理

import { HttpClient } from './HttpClient';
import { AuthService } from './AuthService';
import { UserService } from './UserService';
import { LiveAreaService } from './LiveAreaService';
import { LiveStreamService } from './LiveStreamService';
import { WebSocketService } from './WebSocketService';

// HTTPクライアントの初期化
const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const httpClient = new HttpClient(apiUrl);

// サービスの初期化
export const authService = new AuthService();
export const userService = new UserService(httpClient);
export const liveAreaService = new LiveAreaService(httpClient);
export const liveStreamService = new LiveStreamService(httpClient);
export const webSocketService = new WebSocketService();

// HTTPクライアントを外部からアクセス可能にする（認証トークンの設定用）
export { httpClient };

// サービスの型をエクスポート
export type {
  IUserService,
  ILiveAreaService,
  ILiveStreamService,
  IAuthService,
  IWebSocketService,
} from '@/types/interfaces';


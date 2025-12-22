// SOLID原則に基づくインターフェース定義
// Interface Segregation Principle (ISP): 細かいインターフェースに分割

import {
  User,
  LiveArea,
  LiveAreaMember,
  LiveStream,
  LiveStreamViewer,
  ViewerCount,
} from './models';
import {
  UserUpdateRequest,
  LiveAreaCreateRequest,
  LiveAreaUpdateRequest,
  InviteMemberRequest,
  LiveStreamCreateRequest,
  PaginationParams,
} from './api';

// ユーザーサービスインターフェース
export interface IUserService {
  getCurrentUser(): Promise<User>;
  updateUser(data: UserUpdateRequest): Promise<User>;
}

// ライブエリアサービスインターフェース
export interface ILiveAreaService {
  createLiveArea(data: LiveAreaCreateRequest): Promise<LiveArea>;
  getLiveAreas(params?: PaginationParams): Promise<LiveArea[]>;
  getLiveAreaById(id: string): Promise<LiveArea>;
  updateLiveArea(id: string, data: LiveAreaUpdateRequest): Promise<LiveArea>;
  deleteLiveArea(id: string): Promise<void>;
  inviteMember(liveAreaId: string, data: InviteMemberRequest): Promise<void>;
  removeMember(liveAreaId: string, userId: string): Promise<void>;
  getMembers(liveAreaId: string): Promise<LiveAreaMember[]>;
}

// ライブ配信サービスインターフェース
export interface ILiveStreamService {
  createStream(liveAreaId: string, data: LiveStreamCreateRequest): Promise<LiveStream>;
  getStreamsByLiveArea(liveAreaId: string): Promise<LiveStream[]>;
  getLiveStreams(params?: PaginationParams): Promise<LiveStream[]>;
  getStreamById(id: string): Promise<LiveStream>;
  startStream(id: string): Promise<LiveStream>;
  endStream(id: string): Promise<LiveStream>;
  getViewers(streamId: string): Promise<LiveStreamViewer[]>;
  getViewerCount(streamId: string): Promise<ViewerCount>;
}

// 認証サービスインターフェース
export interface IAuthService {
  signInWithGoogle(): Promise<void>;
  signInWithAccessToken(token: string): Promise<void>;
  signOut(): Promise<void>;
  getCurrentToken(): Promise<string | null>;
  onAuthStateChanged(callback: (user: any) => void): () => void;
}

// WebSocketサービスインターフェース
export interface IWebSocketService {
  connect(streamId: string, token: string): void;
  disconnect(): void;
  sendMessage(message: any): void;
  onMessage(callback: (message: any) => void): void;
  onConnect(callback: () => void): void;
  onDisconnect(callback: () => void): void;
}

// HTTPクライアントインターフェース（Dependency Inversion Principle）
export interface IHttpClient {
  get<T>(url: string, config?: any): Promise<T>;
  post<T>(url: string, data?: any, config?: any): Promise<T>;
  put<T>(url: string, data?: any, config?: any): Promise<T>;
  delete<T>(url: string, config?: any): Promise<T>;
  setAuthToken(token: string): void;
}

// ストレージインターフェース
export interface IStorageService {
  getItem(key: string): string | null;
  setItem(key: string, value: string): void;
  removeItem(key: string): void;
  clear(): void;
}


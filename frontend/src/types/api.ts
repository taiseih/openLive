// API リクエスト・レスポンスの型定義

// ユーザー関連
export interface UserUpdateRequest {
  displayName: string;
  photoUrl: string;
}

// ライブエリア関連
export interface LiveAreaCreateRequest {
  name: string;
  description: string;
  isPublic: boolean;
}

export interface LiveAreaUpdateRequest {
  name: string;
  description: string;
  isPublic: boolean;
}

export interface InviteMemberRequest {
  userId: string;
  role: 'member' | 'moderator';
}

// ライブ配信関連
export interface LiveStreamCreateRequest {
  title: string;
  description: string;
}

// エラーレスポンス
export interface ApiError {
  message: string;
  statusCode: number;
}

// ページネーション
export interface PaginationParams {
  limit?: number;
  offset?: number;
}


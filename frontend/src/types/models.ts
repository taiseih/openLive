// ドメインモデルの型定義

export interface User {
  id: string;
  email: string;
  displayName: string;
  photoUrl: string;
  createdAt: string;
  updatedAt: string;
}

export interface LiveArea {
  id: string;
  ownerId: string;
  name: string;
  description: string;
  isPublic: boolean;
  thumbnailUrl: string;
  createdAt: string;
  updatedAt: string;
}

export interface LiveAreaMember {
  id: string;
  liveAreaId: string;
  userId: string;
  role: 'owner' | 'moderator' | 'member';
  joinedAt: string;
}

export interface LiveStream {
  id: string;
  liveAreaId: string;
  streamerId: string;
  title: string;
  description: string;
  status: 'scheduled' | 'live' | 'ended';
  thumbnailUrl: string;
  startedAt?: string;
  endedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface LiveStreamViewer {
  id: string;
  liveStreamId: string;
  userId?: string;
  joinedAt: string;
  isActive: boolean;
}

export interface ViewerCount {
  count: number;
}

// WebSocketメッセージの型定義
export type WebSocketMessageType = 'chat' | 'offer' | 'answer' | 'ice-candidate';

export interface WebSocketMessage {
  type: WebSocketMessageType;
  to?: string;
  data: any;
}

export interface ChatMessage {
  userId: string;
  userName: string;
  message: string;
  timestamp: string;
}

export interface RTCSignalData {
  sdp?: string;
  type?: 'offer' | 'answer';
  candidate?: string;
  sdpMid?: string;
  sdpMLineIndex?: number;
}


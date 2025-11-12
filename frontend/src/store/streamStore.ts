// ライブ配信の状態管理
import { create } from 'zustand';
import { LiveStream, ChatMessage } from '@/types/models';

interface StreamState {
  currentStream: LiveStream | null;
  chatMessages: ChatMessage[];
  viewerCount: number;
  isConnected: boolean;
  setCurrentStream: (stream: LiveStream | null) => void;
  addChatMessage: (message: ChatMessage) => void;
  setViewerCount: (count: number) => void;
  setConnected: (connected: boolean) => void;
  reset: () => void;
}

export const useStreamStore = create<StreamState>((set) => ({
  currentStream: null,
  chatMessages: [],
  viewerCount: 0,
  isConnected: false,
  setCurrentStream: (stream) => set({ currentStream: stream }),
  addChatMessage: (message) =>
    set((state) => ({
      chatMessages: [...state.chatMessages, message],
    })),
  setViewerCount: (count) => set({ viewerCount: count }),
  setConnected: (connected) => set({ isConnected: connected }),
  reset: () =>
    set({
      currentStream: null,
      chatMessages: [],
      viewerCount: 0,
      isConnected: false,
    }),
}));


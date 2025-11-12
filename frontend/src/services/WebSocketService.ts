'use client';

// Single Responsibility Principle (SRP): WebSocket通信のみに責任を持つ
import { IWebSocketService } from '@/types/interfaces';

export class WebSocketService implements IWebSocketService {
  private ws: WebSocket | null = null;
  private messageCallbacks: ((message: any) => void)[] = [];
  private connectCallbacks: (() => void)[] = [];
  private disconnectCallbacks: (() => void)[] = [];

  connect(streamId: string, token: string): void {
    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/api/v1/ws';
    const url = `${wsUrl}/streams/${streamId}?token=${token}`;

    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.connectCallbacks.forEach((callback) => callback());
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        this.messageCallbacks.forEach((callback) => callback(message));
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.disconnectCallbacks.forEach((callback) => callback());
    };
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  sendMessage(message: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  onMessage(callback: (message: any) => void): void {
    this.messageCallbacks.push(callback);
  }

  onConnect(callback: () => void): void {
    this.connectCallbacks.push(callback);
  }

  onDisconnect(callback: () => void): void {
    this.disconnectCallbacks.push(callback);
  }

  removeAllListeners(): void {
    this.messageCallbacks = [];
    this.connectCallbacks = [];
    this.disconnectCallbacks = [];
  }
}


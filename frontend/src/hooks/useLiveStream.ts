// カスタムフック - ライブ配信ロジックの分離
import { useState, useEffect } from 'react';
import { liveStreamService, webSocketService } from '@/services';
import { useStreamStore } from '@/store/streamStore';
import { LiveStream } from '@/types/models';
import { LiveStreamCreateRequest } from '@/types/api';

export const useLiveStreams = () => {
  const [streams, setStreams] = useState<LiveStream[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchStreams = async () => {
    try {
      setIsLoading(true);
      const data = await liveStreamService.getLiveStreams();
      setStreams(data);
      setError(null);
    } catch (err) {
      setError(err as Error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchStreams();
  }, []);

  return {
    streams,
    isLoading,
    error,
    refetch: fetchStreams,
  };
};

export const useLiveStream = (streamId: string, token?: string) => {
  const {
    currentStream,
    chatMessages,
    viewerCount,
    isConnected,
    setCurrentStream,
    addChatMessage,
    setViewerCount,
    setConnected,
    reset,
  } = useStreamStore();

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchStream = async () => {
      try {
        setIsLoading(true);
        const stream = await liveStreamService.getStreamById(streamId);
        setCurrentStream(stream);
        setError(null);

        // 視聴者数を取得
        const count = await liveStreamService.getViewerCount(streamId);
        setViewerCount(count.count);
      } catch (err) {
        setError(err as Error);
      } finally {
        setIsLoading(false);
      }
    };

    if (streamId) {
      fetchStream();
    }

    return () => {
      reset();
    };
  }, [streamId]);

  useEffect(() => {
    if (streamId && token) {
      // WebSocket接続
      webSocketService.connect(streamId, token);

      webSocketService.onConnect(() => {
        setConnected(true);
      });

      webSocketService.onDisconnect(() => {
        setConnected(false);
      });

      webSocketService.onMessage((message) => {
        if (message.type === 'chat') {
          addChatMessage(message.data);
        } else if (message.type === 'viewer-count') {
          setViewerCount(message.data.count);
        }
      });

      return () => {
        webSocketService.disconnect();
      };
    }
  }, [streamId, token]);

  const sendChatMessage = (message: string) => {
    webSocketService.sendMessage({
      type: 'chat',
      data: { message },
    });
  };

  const startStream = async () => {
    try {
      const stream = await liveStreamService.startStream(streamId);
      setCurrentStream(stream);
    } catch (err) {
      setError(err as Error);
      throw err;
    }
  };

  const endStream = async () => {
    try {
      const stream = await liveStreamService.endStream(streamId);
      setCurrentStream(stream);
    } catch (err) {
      setError(err as Error);
      throw err;
    }
  };

  return {
    stream: currentStream,
    chatMessages,
    viewerCount,
    isConnected,
    isLoading,
    error,
    sendChatMessage,
    startStream,
    endStream,
  };
};


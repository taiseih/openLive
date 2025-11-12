// カスタムフック - ライブエリアロジックの分離
import { useState, useEffect } from 'react';
import { liveAreaService } from '@/services';
import { LiveArea } from '@/types/models';
import { LiveAreaCreateRequest, PaginationParams } from '@/types/api';

export const useLiveAreas = (params?: PaginationParams) => {
  const [liveAreas, setLiveAreas] = useState<LiveArea[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchLiveAreas = async () => {
    try {
      setIsLoading(true);
      const data = await liveAreaService.getLiveAreas(params);
      setLiveAreas(data);
      setError(null);
    } catch (err) {
      setError(err as Error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchLiveAreas();
  }, [params?.limit, params?.offset]);

  const createLiveArea = async (data: LiveAreaCreateRequest) => {
    try {
      const newArea = await liveAreaService.createLiveArea(data);
      setLiveAreas((prev) => [newArea, ...prev]);
      return newArea;
    } catch (err) {
      setError(err as Error);
      throw err;
    }
  };

  return {
    liveAreas,
    isLoading,
    error,
    refetch: fetchLiveAreas,
    createLiveArea,
  };
};

export const useLiveArea = (id: string) => {
  const [liveArea, setLiveArea] = useState<LiveArea | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchLiveArea = async () => {
    try {
      setIsLoading(true);
      const data = await liveAreaService.getLiveAreaById(id);
      setLiveArea(data);
      setError(null);
    } catch (err) {
      setError(err as Error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (id) {
      fetchLiveArea();
    }
  }, [id]);

  return {
    liveArea,
    isLoading,
    error,
    refetch: fetchLiveArea,
  };
};


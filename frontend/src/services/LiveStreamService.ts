'use client';

// Single Responsibility Principle (SRP): ライブ配信管理のみに責任を持つ
import { ILiveStreamService, IHttpClient } from '@/types/interfaces';
import { LiveStream, LiveStreamViewer, ViewerCount } from '@/types/models';
import { LiveStreamCreateRequest, PaginationParams } from '@/types/api';

export class LiveStreamService implements ILiveStreamService {
  constructor(private httpClient: IHttpClient) {}

  async createStream(liveAreaId: string, data: LiveStreamCreateRequest): Promise<LiveStream> {
    return this.httpClient.post<LiveStream>(`/liveareas/${liveAreaId}/streams`, data);
  }

  async getStreamsByLiveArea(liveAreaId: string): Promise<LiveStream[]> {
    return this.httpClient.get<LiveStream[]>(`/liveareas/${liveAreaId}/streams`);
  }

  async getLiveStreams(params?: PaginationParams): Promise<LiveStream[]> {
    const queryParams = new URLSearchParams();
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const query = queryParams.toString();
    return this.httpClient.get<LiveStream[]>(`/streams/live${query ? `?${query}` : ''}`);
  }

  async getStreamById(id: string): Promise<LiveStream> {
    return this.httpClient.get<LiveStream>(`/streams/${id}`);
  }

  async startStream(id: string): Promise<LiveStream> {
    return this.httpClient.post<LiveStream>(`/streams/${id}/start`);
  }

  async endStream(id: string): Promise<LiveStream> {
    return this.httpClient.delete<LiveStream>(`/streams/${id}`);
  }

  async getViewers(streamId: string): Promise<LiveStreamViewer[]> {
    return this.httpClient.get<LiveStreamViewer[]>(`/streams/${streamId}/viewers`);
  }

  async getViewerCount(streamId: string): Promise<ViewerCount> {
    return this.httpClient.get<ViewerCount>(`/streams/${streamId}/viewers/count`);
  }
}


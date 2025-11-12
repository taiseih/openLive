'use client';

// Single Responsibility Principle (SRP): ライブエリア管理のみに責任を持つ
import { ILiveAreaService, IHttpClient } from '@/types/interfaces';
import { LiveArea, LiveAreaMember } from '@/types/models';
import {
  LiveAreaCreateRequest,
  LiveAreaUpdateRequest,
  InviteMemberRequest,
  PaginationParams,
} from '@/types/api';

export class LiveAreaService implements ILiveAreaService {
  constructor(private httpClient: IHttpClient) {}

  async createLiveArea(data: LiveAreaCreateRequest): Promise<LiveArea> {
    return this.httpClient.post<LiveArea>('/liveareas', data);
  }

  async getLiveAreas(params?: PaginationParams): Promise<LiveArea[]> {
    const queryParams = new URLSearchParams();
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const query = queryParams.toString();
    return this.httpClient.get<LiveArea[]>(`/liveareas${query ? `?${query}` : ''}`);
  }

  async getLiveAreaById(id: string): Promise<LiveArea> {
    return this.httpClient.get<LiveArea>(`/liveareas/${id}`);
  }

  async updateLiveArea(id: string, data: LiveAreaUpdateRequest): Promise<LiveArea> {
    return this.httpClient.put<LiveArea>(`/liveareas/${id}`, data);
  }

  async deleteLiveArea(id: string): Promise<void> {
    return this.httpClient.delete<void>(`/liveareas/${id}`);
  }

  async inviteMember(liveAreaId: string, data: InviteMemberRequest): Promise<void> {
    return this.httpClient.post<void>(`/liveareas/${liveAreaId}/invite`, data);
  }

  async removeMember(liveAreaId: string, userId: string): Promise<void> {
    return this.httpClient.delete<void>(`/liveareas/${liveAreaId}/members/${userId}`);
  }

  async getMembers(liveAreaId: string): Promise<LiveAreaMember[]> {
    return this.httpClient.get<LiveAreaMember[]>(`/liveareas/${liveAreaId}/members`);
  }
}


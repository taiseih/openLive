'use client';

// Single Responsibility Principle (SRP): ユーザー管理のみに責任を持つ
import { IUserService, IHttpClient } from '@/types/interfaces';
import { User } from '@/types/models';
import { UserUpdateRequest } from '@/types/api';

export class UserService implements IUserService {
  constructor(private httpClient: IHttpClient) {}

  async getCurrentUser(): Promise<User> {
    return this.httpClient.get<User>('/users/me');
  }

  async updateUser(data: UserUpdateRequest): Promise<User> {
    return this.httpClient.put<User>('/users/me', data);
  }
}


import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom } from 'rxjs';
import { Server, LeaseRequest, ReleaseRequest, LoadSnapshot } from './servers.interface';

@Injectable()
export class ServersService {
  private readonly goApiUrl = 'http://localhost:8080';

  constructor(private readonly httpService: HttpService) {}

  async getServers(): Promise<Server[]> {
    const response = await firstValueFrom(
      this.httpService.get(`${this.goApiUrl}/servers`)
    );
    return response.data;
  }

  async getServer(id: string): Promise<Server> {
    const response = await firstValueFrom(
      this.httpService.get(`${this.goApiUrl}/server?id=${id}`)
    );
    return response.data;
  }

  async getLoadHistory(id: string): Promise<LoadSnapshot[]> {
    const response = await firstValueFrom(
      this.httpService.get(`${this.goApiUrl}/load-history?id=${id}`)
    );
    return response.data;
  }

  async leaseResources(lease: LeaseRequest): Promise<{ message: string }> {
    try {
      const response = await firstValueFrom(
        this.httpService.post(`${this.goApiUrl}/lease`, lease)
      );
      return response.data;
    } catch (error) {
      throw new HttpException(
        error.response?.data || 'Failed to lease resources',
        error.response?.status || HttpStatus.BAD_REQUEST
      );
    }
  }

  async releaseResources(release: ReleaseRequest): Promise<{ message: string }> {
    try {
      const response = await firstValueFrom(
        this.httpService.post(`${this.goApiUrl}/release`, release)
      );
      return response.data;
    } catch (error) {
      throw new HttpException(
        error.response?.data || 'Failed to release resources',
        error.response?.status || HttpStatus.BAD_REQUEST
      );
    }
  }
}
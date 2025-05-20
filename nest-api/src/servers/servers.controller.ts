import { Controller, Get, Param, Post, Body, HttpException, HttpStatus } from '@nestjs/common';
import { ServersService } from './servers.service';
import { Server, LeaseRequest, ReleaseRequest, LoadSnapshot } from './servers.interface';
import { IsNotEmpty, IsString, IsNumber, Min } from 'class-validator';

class LeaseDto {
  @IsNotEmpty()
  @IsString()
  user_id: string;

  @IsNotEmpty()
  @IsString()
  server_id: string;

  @IsNumber()
  @Min(0)
  cpu: number;

  @IsNumber()
  @Min(0)
  ram: number;

  @IsNumber()
  @Min(0)
  storage: number;
}

class ReleaseDto {
  @IsNotEmpty()
  @IsString()
  user_id: string;

  @IsNotEmpty()
  @IsString()
  server_id: string;
}

@Controller('servers')
export class ServersController {
  constructor(private readonly serversService: ServersService) {}

  @Get()
  async getServers(): Promise<Server[]> {
    return this.serversService.getServers();
  }

  @Get(':id')
  async getServer(@Param('id') id: string): Promise<Server> {
    if (!id) {
      throw new HttpException('Server ID is required', HttpStatus.BAD_REQUEST);
    }
    return this.serversService.getServer(id);
  }

  @Get(':id/load-history')
  async getLoadHistory(@Param('id') id: string): Promise<LoadSnapshot[]> {
    if (!id) {
      throw new HttpException('Server ID is required', HttpStatus.BAD_REQUEST);
    }
    return this.serversService.getLoadHistory(id);
  }

  @Post('lease')
  async leaseResources(@Body() leaseDto: LeaseDto): Promise<{ message: string }> {
    return this.serversService.leaseResources(leaseDto);
  }

  @Post('release')
  async releaseResources(@Body() releaseDto: ReleaseDto): Promise<{ message: string }> {
    return this.serversService.releaseResources(releaseDto);
  }
}
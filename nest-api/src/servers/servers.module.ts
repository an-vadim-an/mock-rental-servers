import { Module } from '@nestjs/common';
import { HttpModule } from '@nestjs/axios';
import { ServersService } from './servers.service';
import { ServersController } from './servers.controller';

@Module({
  imports: [HttpModule],
  controllers: [ServersController],
  providers: [ServersService],
})
export class ServersModule {}
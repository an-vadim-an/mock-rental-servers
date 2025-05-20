import { Module } from '@nestjs/common';
import { HttpModule } from '@nestjs/axios';
import { ServersModule } from './servers/servers.module';

@Module({
  imports: [HttpModule, ServersModule],
})
export class AppModule {}
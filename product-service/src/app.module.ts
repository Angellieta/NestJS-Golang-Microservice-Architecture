// src/app.module.ts

import { Module } from '@nestjs/common';
import { RabbitMQModule } from '@golevelup/nestjs-rabbitmq';
import { TypeOrmModule, TypeOrmModuleOptions } from '@nestjs/typeorm';
import { CacheModule } from '@nestjs/cache-manager';
import { redisStore } from 'cache-manager-redis-store';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ProductsModule } from './products/products.module';

// Mendefinisikan konfigurasi DB di sini dengan tipe yang jelas
const dbConfig: TypeOrmModuleOptions = {
  type: 'postgres',
  host: 'postgres',
  port: 5432,
  username: 'user',
  password: 'password',
  database: 'microservice_db',
  autoLoadEntities: true,
  synchronize: true,
};

@Module({
  imports: [
    TypeOrmModule.forRoot(dbConfig),
    // Konfigurasi RabbitMQ Module
    RabbitMQModule.forRoot({
      exchanges: [
        {
          name: 'orders_exchange', // Pastikan nama sama dengan di Go
          type: 'topic',
        },
      ],
      uri: 'amqp://guest:guest@rabbitmq:5672', // URL dari docker-compose
      connectionInitOptions: { wait: false },
    }),
    CacheModule.registerAsync({
      isGlobal: true,
      useFactory: async () => {
        const store = await redisStore({
          socket: {
            host: 'redis',
            port: 6379,
          },
          ttl: 300, // 5 menit
        });
        return {
          store: () => store,
        };
      },
    }),
    ProductsModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}

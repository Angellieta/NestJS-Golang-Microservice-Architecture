// src/products/products.service.ts

import { RabbitSubscribe } from '@golevelup/nestjs-rabbitmq';
import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { CreateProductDto } from './dto/create-product.dto';
import { Product } from './entities/product.entity';

// DTO sederhana untuk payload event dari order-service
class OrderEvent {
  id: string;
  productId: string;
  qty: number;
  totalPrice: number;
  status: string;
  correlationId: string;
  createdAt: string;
}

@Injectable()
export class ProductsService {
  constructor(
    @InjectRepository(Product)
    private productsRepository: Repository<Product>,
  ) {}

  create(createProductDto: CreateProductDto): Promise<Product> {
    const newProduct = this.productsRepository.create(createProductDto);
    return this.productsRepository.save(newProduct);
  }

  async findOne(id: string, correlationId: string): Promise<Product> {
    console.log(
      `[CorrelationID: ${correlationId}] CACHE MISS: Fetching product ${id} from database...`,
    );

    const product = await this.productsRepository.findOne({ where: { id } });
    if (!product) {
      throw new NotFoundException(`Product with ID "${id}" not found`);
    }
    return product;
  }

  // --- METHOD BARU UNTUK MENDENGARKAN EVENT ---
  @RabbitSubscribe({
    exchange: 'orders_exchange',
    routingKey: 'order.created',
    queue: 'products_queue', // Nama "kotak surat" untuk service ini
  })
  public async handleOrderCreated(msg: OrderEvent) {
    const correlationId = msg.correlationId; // Ambil ID dari pesan
    console.log(
      `[CorrelationID: ${correlationId}] Received order.created event: ${JSON.stringify(msg)}`,
    );

    // Logika untuk mengurangi stok produk
    const product = await this.productsRepository.findOne({
      where: { id: msg.productId },
    });

    if (!product) {
      console.error(
        `[CorrelationID: ${correlationId}] Product with ID ${msg.productId} not found.`,
      );
      return;
    }

    // Misalkan setiap order hanya berisi 1 item dari produk tersebut
    // Mengurangi qty produk sesuai dengan qty di order
    if (product.qty >= msg.qty) {
      product.qty -= msg.qty;
      await this.productsRepository.save(product);
      console.log(
        `[CorrelationID: ${correlationId}] Product ${product.id} quantity updated to ${product.qty}`,
      );
    } else {
      console.warn(
        `[CorrelationID: ${correlationId}] Product ${product.id} has insufficient stock.`,
      );
    }
  }
}

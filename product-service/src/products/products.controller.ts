// src/products/products.controller.ts

import {
  Controller,
  Get,
  Post,
  Body,
  Param,
  UseInterceptors,
} from '@nestjs/common';
import { CacheInterceptor } from '@nestjs/cache-manager';
import { ProductsService } from './products.service';
import { CreateProductDto } from './dto/create-product.dto';

@Controller('products')
export class ProductsController {
  constructor(private readonly productsService: ProductsService) {}

  @Post()
  create(@Body() createProductDto: CreateProductDto) {
    return this.productsService.create(createProductDto);
  }

  // Menerapkan decorator caching di sini
  @UseInterceptors(CacheInterceptor)
  @Get(':id')
  findOne(@Param('id') id: string) {
    // Menambahkan log ini untuk membuktikan cache bekerja
    console.log(
      `[product-service] CACHE MISS: Fetching product ${id} from database...`,
    );
    return this.productsService.findOne(id);
  }
}

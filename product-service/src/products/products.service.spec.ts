// product-service/src/products/products.service.spec.ts

import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { NotFoundException } from '@nestjs/common';
import { ProductsService } from './products.service';
import { Product } from './entities/product.entity';

const mockProductRepository = {
  findOne: jest.fn(),
  create: jest.fn(),
  save: jest.fn(),
};

describe('ProductsService', () => {
  let service: ProductsService;

  const mockProduct = {
    id: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
    name: 'Test Product',
    price: 100,
    qty: 10,
    createdAt: new Date(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ProductsService,
        {
          provide: getRepositoryToken(Product),
          useValue: mockProductRepository, // Menggunakan mock object yang sudah didefinisikan
        },
      ],
    }).compile();

    service = module.get<ProductsService>(ProductsService);
    // Casting ke any lalu ke Repository<Product> untuk menyesuaikan dengan Jest's mock
  });

  // Tes 1: Kasus jika produk berhasil ditemukan
  describe('findOne', () => {
    it('should return a product if it exists', async () => {
      mockProductRepository.findOne.mockResolvedValue(mockProduct);

      // Tambahkan argumen kedua (misal: 'test-correlation-id')
      const result = await service.findOne(
        mockProduct.id,
        'test-correlation-id',
      );

      expect(result).toEqual(mockProduct);
      // ...
    });

    // Tes 2: Kasus jika produk tidak ditemukan
    it('should throw a NotFoundException if the product does not exist', async () => {
      mockProductRepository.findOne.mockResolvedValue(null);

      // Tambahkan argumen kedua di sini juga
      await expect(
        service.findOne('non-existent-id', 'test-correlation-id'),
      ).rejects.toThrow(NotFoundException);
    });
  });
});

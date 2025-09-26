import { Test, TestingModule } from '@nestjs/testing';
import { ProductsController } from './products.controller';
import { ProductsService } from './products.service';
import { CacheModule } from '@nestjs/cache-manager';

describe('ProductsController', () => {
  let controller: ProductsController;

  // Membuat mock sederhana untuk ProductsService karena dibutuhkan controller
  const mockProductsService = {
    findOne: jest.fn(),
    create: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [ProductsController],
      providers: [
        {
          provide: ProductsService,
          useValue: mockProductsService,
        },
      ],
      // Menambahkan CacheModule ke dalam imports
      imports: [CacheModule.register({ isGlobal: true })],
    }).compile();

    controller = module.get<ProductsController>(ProductsController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});

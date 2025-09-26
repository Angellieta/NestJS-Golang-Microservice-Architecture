import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
} from 'typeorm';

@Entity('products') // Nama tabel di database
export class Product {
  @PrimaryGeneratedColumn('uuid') // ID unik otomatis (UUID)
  id: string;

  @Column()
  name: string;

  @Column('decimal')
  price: number;

  @Column('int')
  qty: number;

  @CreateDateColumn() // Tanggal dibuat otomatis
  createdAt: Date;
}

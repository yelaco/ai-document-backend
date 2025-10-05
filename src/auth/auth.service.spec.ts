import { Test, TestingModule } from '@nestjs/testing';
import { AuthService } from './auth.service';
import { getRepositoryToken } from '@nestjs/typeorm';
import { User } from '../users/entities/user.entity';
import { JwtService } from '@nestjs/jwt';
import { ARGON2_OPTIONS } from './auth.constants';
import { argon2id } from 'argon2';

const mockUserRepository = {
  save: jest.fn(),
  findOne: jest.fn(),
  find: jest.fn(),
};

describe('AuthService', () => {
  let service: AuthService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AuthService,
        JwtService,
        {
          provide: getRepositoryToken(User),
          useValue: mockUserRepository,
        },
        {
          provide: ARGON2_OPTIONS,
          useValue: {
            type: argon2id,
            memoryCost: 65536, // 64 MiB
            timeCost: 4, // 4 iterations
            parallelism: 1,
          },
        },
      ],
    }).compile();

    service = module.get<AuthService>(AuthService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});

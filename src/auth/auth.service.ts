import {
  ConflictException,
  Inject,
  Injectable,
  InternalServerErrorException,
} from '@nestjs/common';
import { Repository } from 'typeorm';
import { User } from '../users/entities/user.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { CreateUserDto } from '../users/dto/create-user.dto';
import { ARGON2_OPTIONS } from './auth.constants';
import * as argon2 from 'argon2';
import { JwtService } from '@nestjs/jwt';

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,

    @Inject(ARGON2_OPTIONS)
    private options: argon2.Options,

    private jwtService: JwtService,
  ) {}

  async registerUser(createUserDto: CreateUserDto): Promise<User> {
    const existingUser = await this.userRepository.findOne({
      where: { email: createUserDto.email },
    });
    if (existingUser) {
      throw new ConflictException('User already exists.');
    }
    const newUser = this.userRepository.create();
    newUser.email = createUserDto.email;
    newUser.fullName = createUserDto.fullName;
    newUser.passwordHash = await this.hashPassword(createUserDto.password);
    return await this.userRepository.save(newUser);
  }

  async login(user: any) {
    const payload = { username: user.email, sub: user.id };
    return {
      access_token: this.jwtService.sign(payload),
    };
  }

  async validateUser(email: string, password: string): Promise<User | null> {
    const user = await this.userRepository.findOne({ where: { email } });
    if (user && (await this.verifyPassword(user.passwordHash, password))) {
      return user;
    }
    return null;
  }

  private async hashPassword(password: string): Promise<string> {
    try {
      // The salt is generated automatically by the library
      const hash = await argon2.hash(password, this.options);
      return hash;
    } catch {
      throw new InternalServerErrorException('Password hashing failed.');
    }
  }

  private async verifyPassword(
    storedHash: string,
    passwordAttempt: string,
  ): Promise<boolean> {
    try {
      const isMatch = await argon2.verify(storedHash, passwordAttempt);
      return isMatch;
    } catch {
      return false;
    }
  }
}

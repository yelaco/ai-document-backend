import {
  Body,
  Controller,
  HttpCode,
  HttpStatus,
  Post,
  Request,
  UseGuards,
} from '@nestjs/common';
import { ApiOperation, ApiResponse, ApiTags, ApiBody } from '@nestjs/swagger';
import { AuthService } from './auth.service';
import { CreateUserDto } from '../users/dto/create-user.dto';
import { toUserDto, UserDto } from '../users/dto/user.dto';
import { LocalAuthGuard } from './guard/local-auth.guard';

@ApiTags('auth')
@Controller('auth')
export class AuthController {
  constructor(private authService: AuthService) {}

  @ApiOperation({ summary: 'User login', description: 'Authenticate and obtain a JWT token.' })
  @ApiResponse({ status: 200, description: 'Login successful', schema: { example: { access_token: 'ey123...', user: { id: 1, email: 'user@example.com', fullName: 'John Doe' } } } })
  @UseGuards(LocalAuthGuard)
  @Post('login')
  @HttpCode(HttpStatus.OK)
  login(@Request() req) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    return this.authService.login(req.user);
  }

  @ApiOperation({ summary: 'User logout', description: 'Log out the current user.' })
  @ApiResponse({ status: 200, description: 'Logout successful', schema: { example: { message: 'Logged out' } } })
  @Post('logout')
  logout(@Request() req) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-return, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
    return req.logout();
  }

  @ApiOperation({ summary: 'Register user', description: 'Create a new user account.' })
  @ApiBody({ type: CreateUserDto, examples: { default: { value: { email: 'user@example.com', password: 'password123', fullName: 'John Doe' } } } })
  @ApiResponse({ status: 201, description: 'User created', type: UserDto, schema: { example: { id: 1, email: 'user@example.com', fullName: 'John Doe', documents: [] } } })
  @Post('register')
  async register(@Body() createUserDto: CreateUserDto): Promise<UserDto> {
    const user = await this.authService.registerUser(createUserDto);
    return toUserDto(user);
  }
}

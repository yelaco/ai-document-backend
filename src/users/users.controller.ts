import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
} from '@nestjs/common';
import { UsersService } from './users.service';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { ApiTags, ApiOperation, ApiBody, ApiResponse } from '@nestjs/swagger';

@ApiTags('users')
@Controller('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @ApiOperation({ summary: 'Create user', description: 'Create a new user.' })
  @ApiBody({
    type: CreateUserDto,
    examples: {
      default: {
        value: {
          email: 'user@example.com',
          password: 'password123',
          fullName: 'John Doe',
        },
      },
    },
  })
  @ApiResponse({
    status: 201,
    description: 'User created',
    schema: {
      example: {
        id: 1,
        email: 'user@example.com',
        fullName: 'John Doe',
        documents: [],
      },
    },
  })
  @Post()
  create(@Body() createUserDto: CreateUserDto) {
    return this.usersService.create(createUserDto);
  }

  @ApiOperation({ summary: 'List users', description: 'Get all users.' })
  @ApiResponse({
    status: 200,
    description: 'List of users.',
    schema: {
      example: [
        {
          id: 1,
          email: 'user@example.com',
          fullName: 'John Doe',
          documents: [],
        },
      ],
    },
  })
  @Get()
  findAll() {
    return this.usersService.findAll();
  }

  @ApiOperation({
    summary: 'Get user by id',
    description: 'Get a single user by id.',
  })
  @ApiResponse({
    status: 200,
    description: 'User details.',
    schema: {
      example: {
        id: 1,
        email: 'user@example.com',
        fullName: 'John Doe',
        documents: [],
      },
    },
  })
  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.usersService.findOne(+id);
  }

  @ApiOperation({ summary: 'Update user', description: 'Update a user by id.' })
  @ApiBody({
    type: UpdateUserDto,
    examples: { default: { value: { fullName: 'Jane Doe' } } },
  })
  @ApiResponse({
    status: 200,
    description: 'Updated user.',
    schema: {
      example: {
        id: 1,
        email: 'user@example.com',
        fullName: 'Jane Doe',
        documents: [],
      },
    },
  })
  @Patch(':id')
  update(@Param('id') id: string, @Body() updateUserDto: UpdateUserDto) {
    return this.usersService.update(+id, updateUserDto);
  }

  @ApiOperation({ summary: 'Delete user', description: 'Delete a user by id.' })
  @ApiResponse({
    status: 200,
    description: 'User deleted.',
    schema: { example: { success: true } },
  })
  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.usersService.remove(+id);
  }
}

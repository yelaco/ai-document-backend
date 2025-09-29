import { BadRequestException, ValidationPipe } from '@nestjs/common';
import { ValidationError } from 'class-validator';

export default class CustomValidationPipe extends ValidationPipe {
  constructor() {
    super({
      transform: true,
      exceptionFactory: (errors: ValidationError[]) => {
        const formattedErrors = errors.flatMap((error) => {
          if (!error.constraints) {
            return [];
          }
          return Object.keys(error.constraints).map((constraint) => ({
            code: this.getErrorCode(constraint),
            target: error.property,
            message: error.constraints?.[constraint] ?? '',
          }));
        });

        return new BadRequestException({
          message: 'Invalid request data',
          details: formattedErrors,
        });
      },
    });
  }

  private getErrorCode(constraint: string): string {
    switch (constraint) {
      case 'isNotEmpty':
        return 'MISSING_REQUIRED_FIELD';
      case 'maxLength':
      case 'minLength':
        return 'INVALID_FIELD_LENGTH';
      case 'isEnum':
        return 'INVALID_ENUM_VALUE';
      case 'isDateString':
        return 'INVALID_DATE_FORMAT';
      case 'isNumber':
        return 'INVALID_NUMBER_FORMAT';
      default:
        return 'UNKNOWN_VALIDATION_ERROR';
    }
  }
}

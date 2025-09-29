import {
  CallHandler,
  ExecutionContext,
  Injectable,
  NestInterceptor,
  Logger, // Use built-in Logger
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Request, Response } from 'express'; // or @nestjs/platform-express

@Injectable()
export class LoggingInterceptor implements NestInterceptor {
  private readonly logger = new Logger(LoggingInterceptor.name);

  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const ctx = context.switchToHttp();
    const request = ctx.getRequest<Request>();
    const response = ctx.getResponse<Response>();
    const now = Date.now();

    const { method, url } = request;

    // 1. Log BEFORE the request hits the route handler
    this.logger.log(`[REQUEST] ${method} ${url} started...`);

    return next.handle().pipe(
      tap(() => {
        const responseTime = Date.now() - now;
        const statusCode = response.statusCode;

        this.logger.log(
          `[RESPONSE] ${method} ${url} | Status: ${statusCode} | Time: ${responseTime}ms`,
          'HttpAccess',
        );
      }),
    );
  }
}

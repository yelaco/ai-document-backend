interface DatabaseConfig {
  host: string | undefined;
  port: number;
  user: string | undefined;
  pass: string | undefined;
  name: string | undefined;
}

interface VectorDatabaseConfig {
  host: string | undefined;
  port: number;
}

interface RedisConfig {
  host: string | undefined;
  port: number;
  pass: string | undefined;
}

interface AiConfig {
  serviceType: string | undefined;
  geminiApiKey: string | undefined;
  anthropicApiKey: string | undefined;
}

export function databaseConfig(): DatabaseConfig {
  return {
    host: process.env.DB_HOST,
    port: process.env.DB_PORT ? parseInt(process.env.DB_PORT) : 5432,
    user: process.env.DB_USER,
    pass: process.env.DB_PASS,
    name: process.env.DB_NAME,
  };
}

export function vectorDatabaseConfig(): VectorDatabaseConfig {
  return {
    host: process.env.VDB_HOST || 'vector_db',
    port: process.env.VDB_PORT ? parseInt(process.env.VDB_PORT) : 8000,
  };
}

export function redisConfig(): RedisConfig {
  return {
    host: process.env.REDIS_HOST || 'redis',
    port: process.env.REDIS_PORT ? parseInt(process.env.REDIS_PORT) : 6379,
    pass: process.env.REDIS_PASS,
  };
}

export function aiConfig(): AiConfig {
  return {
    serviceType: process.env.AI_SERVICE_TYPE,
    geminiApiKey: process.env.GEMINI_API_KEY,
    anthropicApiKey: process.env.ANTHROPIC_API_KEY,
  };
}

export default () => ({
  appEnv: process.env.APP_ENV || 'development',
  port: process.env.APP_PORT ? parseInt(process.env.APP_PORT) : 3000,
  database: databaseConfig(),
  vectorDatabase: vectorDatabaseConfig(),
  redis: redisConfig(),
  ai: aiConfig(),
});

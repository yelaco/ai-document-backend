interface DatabaseConfig {
  host: string | undefined;
  port: number;
  user: string | undefined;
  pass: string | undefined;
  name: string | undefined;
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

export default () => ({
  port: process.env.APP_PORT ? parseInt(process.env.APP_PORT) : 3000,
  database: databaseConfig(),
});

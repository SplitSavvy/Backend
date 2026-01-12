import pg from "pg";
import "dotenv/config";

//database connection details stored in .env file and processed and used here
const db = new pg.Pool({
  user: process.env.PGUSER,
  host: process.env.PGHOST,
  database: process.env.PGDATABASE,
  password: process.env.PGPASSWORD,
  port: process.env.PGPORT,
});

export default db;

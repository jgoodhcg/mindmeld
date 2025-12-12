const fs = require('node:fs');
const path = require('node:path');
const process = require('node:process');
const { Client } = require('pg');

const migrationsDir = path.join(process.cwd(), 'db', 'migrations');
const envFiles = [path.join(process.cwd(), '.env.local'), path.join(process.cwd(), '.env')];

function loadEnv() {
  for (const envPath of envFiles) {
    if (!fs.existsSync(envPath)) continue;

    const lines = fs.readFileSync(envPath, 'utf8').split(/\r?\n/);
    for (const line of lines) {
      const trimmed = line.trim();
      if (!trimmed || trimmed.startsWith('#')) continue;
      const eqIndex = trimmed.indexOf('=');
      if (eqIndex === -1) continue;

      const key = trimmed.slice(0, eqIndex).trim();
      let value = trimmed.slice(eqIndex + 1).trim();

      // Strip surrounding quotes
      if (
        (value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.endsWith("'"))
      ) {
        value = value.slice(1, -1);
      }

      if (!process.env[key]) {
        process.env[key] = value;
      }
    }
  }
}

function getMigrationFiles() {
  if (!fs.existsSync(migrationsDir)) {
    return [];
  }

  return fs
    .readdirSync(migrationsDir)
    .filter((file) => file.endsWith('.sql'))
    .sort();
}

async function ensureMigrationsTable(client) {
  await client.query(`
    CREATE TABLE IF NOT EXISTS migrations (
      name text PRIMARY KEY,
      applied_at timestamptz NOT NULL DEFAULT now()
    )
  `);
}

async function main() {
  loadEnv();
  const databaseUrl = process.env.DATABASE_URL;

  if (!databaseUrl) {
    console.error('DATABASE_URL is not set; cannot run migrations.');
    process.exit(1);
  }

  const client = new Client({ connectionString: databaseUrl });
  await client.connect();

  try {
    await ensureMigrationsTable(client);

    const applied = new Set(
      (await client.query('SELECT name FROM migrations')).rows.map((row) => row.name)
    );

    const files = getMigrationFiles();

    for (const file of files) {
      if (applied.has(file)) {
        continue;
      }

      const fullPath = path.join(migrationsDir, file);
      const sql = fs.readFileSync(fullPath, 'utf8');

      console.log(`Applying migration ${file}...`);
      try {
        await client.query(sql);
      } catch (err) {
        console.error(`Migration ${file} failed.`);
        throw err;
      }

      await client.query('INSERT INTO migrations (name) VALUES ($1)', [file]);
      console.log(`Applied migration ${file}`);
    }

    console.log('Migrations complete.');
  } finally {
    await client.end();
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});

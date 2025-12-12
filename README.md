This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Local database (Postgres)

- Preferred: run Postgres via Docker so dev matches prod.
- Copy `.env.example` to `.env.local` and adjust if needed.
- Start DB: `docker compose up -d db`
- Stop DB: `docker compose down` (add `--volumes` to reset data)
- If you prefer native Postgres, skip Docker and point `DATABASE_URL` to your local instance.

## Migrations

- Install deps if you haven’t since pulling changes: `npm install`
- Make sure `DATABASE_URL` is set (via `.env.local` or env).
- Apply migrations: `npm run db:migrate` (uses `db/migrations/*.sql`, tracks state in `migrations` table).
- On DigitalOcean App Platform, run `npm run db:migrate` as a deploy/release command with the managed Postgres `DATABASE_URL` and a strong `LOBBY_TOKEN_SECRET`.

## Deployment (DigitalOcean App Platform)

- Attach a managed Postgres and use its connection string as `DATABASE_URL`.
- Set `LOBBY_TOKEN_SECRET` to a strong random string.
- Build command: `npm run build`
- Run command: `npm run start` (will switch to `node server.js` once the custom WS server is added).
- App Platform sets `PORT` for you; Next.js reads it automatically.

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

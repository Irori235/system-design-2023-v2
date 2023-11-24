FROM node:20 AS base

WORKDIR /app
COPY ./package.json ./yarn.lock ./
RUN yarn install --frozen-lockfile --production
COPY . .

FROM node:20 AS dev
WORKDIR /app
COPY . .
CMD ["yarn", "dev"]

FROM base AS builder

WORKDIR /app
RUN yarn build 

FROM gcr.io/distroless/nodejs20

WORKDIR /app

# Copy artifacts
COPY ./next.config.js .

# Copy node_modules
COPY --from=base /app/node_modules ./node_modules

# Copy build
COPY --from=builder /app/.next ./.next
EXPOSE 3000
CMD ["./node_modules/next/dist/bin/next", "start"]

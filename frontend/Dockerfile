FROM node:24-bookworm-slim AS deps

WORKDIR /src

COPY package*.json ./
RUN npm ci

FROM deps AS builder

COPY . .
RUN npm run build

FROM node:24-bookworm-slim

WORKDIR /app
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=80

COPY --from=builder /src/.output /app/.output

EXPOSE 80

CMD ["node", ".output/server/index.mjs"]

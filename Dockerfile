FROM docker.io/directus/directus:10.13

USER root
RUN corepack enable
USER node

RUN pnpm install directus-extension-sync

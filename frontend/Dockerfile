FROM node:20-alpine AS build

# Bust Docker layer cache when sources change (pass --build-arg CACHEBUST=$(date +%s))
ARG CACHEBUST=1

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
COPY . .
RUN echo "build cache bust: ${CACHEBUST}" && npm run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]

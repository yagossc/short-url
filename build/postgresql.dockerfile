FROM postgres:12.1-alpine

# Create initialization scripts dir
# for more info see: https://hub.docker.com/_/postgres
RUN mkdir -p /docker-entrypoint-initdb.d

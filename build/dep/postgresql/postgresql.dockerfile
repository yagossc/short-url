FROM postgres:12.1-alpine

# Create initialization scripts directory.
# For more info see: https://hub.docker.com/_/postgres
RUN mkdir -p /docker-entrypoint-initdb.d

# Create postgreSQL configuration directory.
# This is because of:
#   Important note: you must set listen_addresses = '*'
#   so that other containers will be able to access postgres.
# For more info see: https://hub.docker.com/_/postgres
RUN mkdir -p /etc/postgresql/

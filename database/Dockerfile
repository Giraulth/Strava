FROM postgres:15

# Table creation
COPY dump/ /docker-entrypoint-initdb.d/

# Passwordless access from host
COPY enable_passwordless.sh /docker-entrypoint-initdb.d/
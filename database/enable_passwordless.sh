sed -i 's#host    all             all             127.0.0.1/32            trust#host    all             all             all            trust#' /data/postgres/pg_hba.conf

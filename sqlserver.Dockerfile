FROM mcr.microsoft.com/mssql/server:2019-CU16-ubuntu-18.04

EXPOSE 1433

ARG ACCEPT_EULA \
    MSSQL_SA_PASSWORD \
    MSSQL_PID \
    SECS_AWAIT_SVR_BOOTUP

ENV ACCEPT_EULA=$ACCEPT_EULA \
    MSSQL_SA_PASSWORD=$MSSQL_SA_PASSWORD \
    MSSQL_PID=$MSSQL_PID

#WORKDIR /var/opt/mssql/backup
WORKDIR /var/opt/mssql/create_scripts

COPY ./db/sqlserver/db_create_scripts .

RUN /opt/mssql/bin/sqlservr --accept-eula & \
    ( \
    echo "awaiting server bootup..." && sleep $SECS_AWAIT_SVR_BOOTUP && \
    echo "creating database" && /opt/mssql-tools/bin/sqlcmd -S tcp:127.0.0.1,1433 -U SA -P "$MSSQL_SA_PASSWORD" -d master -i 01_create_db.sql && \
    echo "inserting into HttpMethods" && /opt/mssql-tools/bin/sqlcmd -S tcp:127.0.0.1,1433 -U SA -P "$MSSQL_SA_PASSWORD" -d access_control -i 02_insert_HttpMethods.sql && \
    echo "inserting into Admins" && /opt/mssql-tools/bin/sqlcmd -S tcp:127.0.0.1,1433 -U SA -P "$MSSQL_SA_PASSWORD" -d access_control -i 03_insert_Admins.sql \
    )

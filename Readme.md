TODO: 

- Setup frontend next js
- Setup golang http server
- integrate wapi.go client
- Next.js should be statically served via the same golang server, we are opting for this to reduce the complexity, we can add a flag in the configs if the user want to host the frontend separately.
- optional caching server if enabled, than docker compose should deploy a caching server as well



Flow for the database is: 

1. Prepare the models using GORM
2. use atlas CLI to generate and apply migrations to the database
3. use go-jet lib to build SQL queries
4. use go-jet to execute the query on database using a db connection using native database/sql package.
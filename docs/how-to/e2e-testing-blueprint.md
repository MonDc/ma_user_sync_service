# E2E [DATABASE]
# E2E Testing — MA User Sync [DATABASE]

**Reset (optional, destructive):** drop table, tracking, and database (We will get truly clean slate).

```bash
mysql -h 192...x -u user...x -p'pass...x' -e "USE ma_db; DROP TABLE IF EXISTS ma_users; DROP TABLE IF EXISTS schema_migrations;"
mysql -h 192...x -u user...x -p'pass...x' -e "DROP DATABASE IF EXISTS ma_db;"



**Provision the database manually (00001 migration. This must be done manually). 
**"Provision ma_db manually so migrations always track in ma_db itself — otherwise you must connect through any db e.g. mysql, which leaves a stray, disconnected schema_migrations ledger there. (leftover)"**
**(Note: after this, 000001_create_database will still run on the first migrate up — but it's CREATE DATABASE IF NOT EXISTS ma_db, so it's a no-op. This manual command is kept as a reference.)
mysql -h 192...x -u user...x -p'pass...x' -e "CREATE DATABASE IF NOT EXISTS ma_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"



**Apply migrations (up 1 = one version per run):
**Each migrate up applies the next pending migration, increasing the version by one. If no update(pending), the version stays at the highest applied file's number — running it again does nothing.
migrate -database 'mysql://user...x:pass...x@tcp(192...x:3306)/ma_db' -path database/migrations up 1



**Verify indexes and constraints in phpMyAdmin — open ma_db → ma_users → Structure/Indexes tab.
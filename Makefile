## Goose
# create:
# 	goose mysql "mysql://root:password@tcp(localhost:3306)/gamedb" -dir repository/mysql/migrations create 00001_init sql


## sql-migrate
create_migrations:
	sql-migrate new -env="production" -config=repository/mysql/dbconfig.yml
	
migrate_up:
	sql-migrate up -env="production" -config=repository/mysql/dbconfig.yml

migrate_down:
	sql-migrate down -env="production" -config=repository/mysql/dbconfig.yml
#sql-migrate down -env="production" -config=./repository/mysql/dbconfig.yml -limit=1
	
migrate_status:
	sql-migrate status -env="production" -config=./repository/mysql/dbconfig.yml
	
	
## sqlx-cli
# create_migrations:
# 	sqlx migrate add -r --source repository/mysql/migrations init
	
# migrate_up:
# 	sqlx migrate run --source repository/mysql/migrations --database-url "mysql://root:password@localhost:3306/gamedb"

# migrate_down:
# 		sqlx migrate revert --source repository/mysql/migrations --database-url "mysql://root:password@localhost:3306/gamedb"

# migrate_info:
# 	sqlx migrate info --source repository/mysql/migrations --database-url "mysql://root:password@localhost:3306/gamedb"

##
# sqlx migrate run -env="production" -config=dbconfig.yml

## migrate
# create_migrate:
# 	migrate create -ext sql -dir repository/mysql/migrations -seq init_schema


# migrate_up:
# 	 migrate -path repository/mysql/migrations -database "mysql://root:password@tcp(localhost:3306)/gamedb" -verbose up


# migrate_down:
# 	 migrate -path repository/mysql/migrations -database "mysql://root:password@tcp(localhost:3306)/gamedb" -verbose down
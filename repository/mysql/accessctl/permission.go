package accessctl

import (
	"game-app/entity"
	"game-app/repository/mysql"
	"time"
)

func scanPermission(scanner mysql.Scanner) (entity.Permission, error) {
	var createdAt time.Time
	var p entity.Permission

	err := scanner.Scan(&p.ID, &p.Title, &createdAt)

	return p, err
}

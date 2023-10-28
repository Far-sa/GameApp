package mysql

import (
	"game-app/entity"
	"time"
)

func scanPermission(scanner Scanner) (entity.Permission, error) {
	var createdAt time.Time
	var p entity.Permission

	err := scanner.Scan(&p.ID, &p.Title, &createdAt)

	return p, err
}

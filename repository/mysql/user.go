package mysql

import (
	"database/sql"
	"fmt"
	"game-app/entity"
)

func (r MySQLDB) RegisterUser(user entity.User) (entity.User, error) {

	query := "insert into users(name,phone_number) values(?,?)"
	res, err := r.Db.Exec(query, user.Name, user.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("can not execute commnad %w", err)
	}

	id, _ := res.LastInsertId()
	user.ID = uint(id)

	return user, nil

}

func (r MySQLDB) UniquenePhonenumber(phoneNumer string) (bool, error) {
	user := entity.User{}
	var createdAt []uint8

	query := "select id,name,phone_number,created_at from users where phone_number=?"
	row := r.Db.QueryRow(query, phoneNumer)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, fmt.Errorf("can not scan query result :%w", err)
	}

	return false, nil
}

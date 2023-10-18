package mysql

import (
	"database/sql"
	"fmt"
	"game-app/entity"
)

func (r MySQLDB) RegisterUser(user entity.User) (entity.User, error) {

	query := "insert into users(name,phone_number,password) values(?,?,?)"
	res, err := r.db.Exec(query, user.Name, user.PhoneNumber, user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can not execute commnad %w", err)
	}

	id, _ := res.LastInsertId()
	user.ID = uint(id)

	return user, nil

}

func (r MySQLDB) UniquenePhonenumber(phoneNumer string) (bool, error) {

	query := "select id,name,password,phone_number,created_at from users where phone_number=?"
	row := r.db.QueryRow(query, phoneNumer)

	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, fmt.Errorf("can not scan query result :%w", err)
	}

	return false, nil
}

func (r MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {

	row := r.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		return entity.User{}, false, fmt.Errorf("can not scan query result :%w", err)
	}

	return user, true, nil
}

func (r MySQLDB) GetUserById(userID uint64) (entity.User, error) {

	row := r.db.QueryRow(`select * from users where id=?`, userID)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("record not found")
		}

		return entity.User{}, fmt.Errorf("can not scan row for user %w", err)
	}
	return user, nil
}

func scanUser(row *sql.Row) (entity.User, error) {
	var createdAt []uint8
	var user entity.User

	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.PhoneNumber, &createdAt)

	return user, err
}

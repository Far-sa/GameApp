package mysql

import (
	"database/sql"
	"fmt"
	"game-app/entity"
	"game-app/pkg/errs"
	"game-app/pkg/richerror"
)

// TODO implement sqlx

func (r MySQLDB) RegisterUser(user entity.User) (entity.User, error) {

	query := "insert into users(phone_number,name,password) values(?,?,?)"
	res, err := r.db.Exec(query, user.PhoneNumber, user.Name, user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can not execute commnad %w", err)
	}

	id, _ := res.LastInsertId()
	user.ID = uint(id)

	return user, nil

}

func (r MySQLDB) UniquenePhonenumber(phoneNumer string) (bool, error) {
	const op = "mysql.UniquenePhonenumber"

	query := "select * from users where phone_number=?"
	row := r.db.QueryRow(query, phoneNumer)

	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, richerror.New(op).WithMessage(errs.ErrorMsgCantQuery).
			WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (r MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := r.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		//fmt.Errorf("can not scan query result :%w", err)
		return entity.User{}, false, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgCantQuery).WithKind(richerror.KindUnexpected)
	}

	return user, true, nil
}

func (r MySQLDB) GetUserById(userID uint64) (entity.User, error) {
	const op = "mysql.GetUserById"

	row := r.db.QueryRow(`select * from users where id=?`, userID)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errs.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgCantQuery).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(row *sql.Row) (entity.User, error) {
	var createdAt []uint8
	var user entity.User

	err := row.Scan(&user.ID, &user.PhoneNumber, &user.Name, &user.Password, &createdAt)

	return user, err
}

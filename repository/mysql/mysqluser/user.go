package mysqluser

import (
	"database/sql"
	"fmt"
	"game-app/entity"
	"game-app/pkg/errs"
	"game-app/pkg/richerror"
	"game-app/repository/mysql"
)

// TODO implement sqlx

func (r *DB) RegisterUser(user entity.User) (entity.User, error) {

	query := "insert into users(phone_number,name,password,role) values(?,?,?,?)"
	res, err := r.conn.Conn().Exec(query, user.PhoneNumber, user.Name, user.Password, user.Role.String())
	if err != nil {
		return entity.User{}, fmt.Errorf("can not execute commnad %w", err)
	}

	id, _ := res.LastInsertId()
	user.ID = uint(id)

	return user, nil

}

func (r *DB) UniquenePhonenumber(phoneNumer string) (bool, error) {
	const op = "mysql.UniquenePhonenumber"

	query := "select * from users where phone_number=?"
	row := r.conn.Conn().QueryRow(query, phoneNumer)

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

func (r *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := r.conn.Conn().QueryRow(`select * from users where phone_number=?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errs.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}

		// TODO --> add loger for unexpected errors for better observability
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgCantQuery).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (r *DB) GetUserById(userID uint) (entity.User, error) {
	const op = "mysql.GetUserById"

	row := r.conn.Conn().QueryRow(`select * from users where id=?`, userID)

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

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var createdAt []uint8
	var user entity.User

	var roleStr string

	err := scanner.Scan(&user.ID, &user.PhoneNumber, &user.Name, &user.Password, &createdAt, &roleStr)

	user.Role = entity.MapToRoleEntity(roleStr)

	return user, err
}

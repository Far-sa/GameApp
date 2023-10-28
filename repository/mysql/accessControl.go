package mysql

import (
	"game-app/entity"
	"game-app/pkg/errs"
	"game-app/pkg/richerror"
	"game-app/pkg/slice"
	"strings"
	"time"
)

func (r *MySQLDB) GetUserPermissionTitles(userID uint) ([]entity.PermissionTitle, error) {
	const op = "mysql.GetUserPermissionTitles"
	// get user
	user, err := r.GetUserById(userID)
	if err != nil {
		return nil, richerror.New(op).WithErr(err)
	}

	roleACL := make([]entity.AccessControl, 0)

	rows, err := r.db.Query("select * from access_controls where actor_type = ? and actor_id=?", entity.RoleActorType, user.Role)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := scanAccessControl(rows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
		}

		roleACL = append(roleACL, acl)
	}

	// check error of above loop
	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
	}

	userACL := make([]entity.AccessControl, 0)

	userRows, err := r.db.Query("select * from access_controls where actor_type = ? and actor_id=?",
		entity.UserActorType, user.ID)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
	}

	defer userRows.Close()

	for userRows.Next() {
		acl, err := scanAccessControl(userRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
		}

		userACL = append(userACL, acl)
	}
	// check error of above loop
	if err := userRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
	}

	//* merge ACLs by Permissions
	permissionIDs := make([]uint, 0)

	for _, r := range roleACL {
		if !slice.DoesExist(permissionIDs, r.PermissionID) {
			permissionIDs = append(permissionIDs, r.PermissionID)
		}
	}

	if len(permissionIDs) == 0 {
		return nil, nil
	}

	//* cause of dynamic args of permissionsID
	args := make([]any, len(permissionIDs))

	for i, id := range permissionIDs {
		args[i] = id
	}

	//* this query works if we have one or more permissions id
	query := "SELECT * FROM permissions WHERE id in (?" + strings.Repeat(",?", len(permissionIDs)-1) + ")"

	pRows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
	}

	defer pRows.Close()

	permissionTtiles := make([]entity.PermissionTitle, 0)
	for pRows.Next() {
		permission, err := scanPermission(pRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
		}
		permissionTtiles = append(permissionTtiles, permission.Title)
	}

	if err := pRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errs.ErrorMsgSomethingWrong).WithKind(richerror.KindUnexpected)
	}

	return permissionTtiles, nil

}

func scanAccessControl(scanner Scanner) (entity.AccessControl, error) {
	var createdAt time.Time
	var acl entity.AccessControl

	err := scanner.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, &acl.PermissionID, &createdAt)

	return acl, err

}

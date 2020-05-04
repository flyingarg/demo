package data

import (
	"demo/internal/logger"
	"errors"
	"fmt"
)

type Plan string

const (
	GOLD   Plan = "GOLD"
	SILVER      = "SILVER"
	BRONZE      = "BRONZE"
)

func (p Plan) IsValid(plan string) error {
	switch plan {
	case "GOLD":
		p = GOLD
	case "SILVER":
		p = SILVER
	case "BRONZE":
		p = BRONZE
	default:
		return errors.New(fmt.Sprintf("plan type not recognized %s", plan))
	}
	return nil
}

type User struct {
	ID   int64
	Name string
	Plan Plan
}

//GetUser returns a user, if the user's name and password hashes are correct, else returns empty user struct.
//Note that hashes and user details are in separate tables.
func GetUser(userName string, passwordHash string) (User, error) {
	query := `select u.id, u.username, u.plan 
				from tokens as t, user_details as u 
				where 
					t.user_id = (select id from user_details where username = ? ) and 
					t.user_secret = ?;`
	var id int64
	var name, plan string

	db := GetConnection()
	err := db.QueryRow(query, userName, passwordHash).Scan(&id, &name, &plan)
	if err != nil {
		logger.Log.Sugar().Errorf("failed to execute query %s, with error %e", query, err)
		return User{}, nil
	} else {
		p := Plan(plan)
		err = p.IsValid(plan)
		if err == nil {
			return User{
				ID:   id,
				Name: name,
				Plan: p}, nil
		} else {
			return User{}, err
		}
	}
}

//WithinRateLimits checks the plan and the permissible per hour query count to determine if the user has exceeded
//the rate limit and hence return a bool and error
func WithinRateLimits(username string, plan string) (bool, error) {
	db := GetConnection()
	checkRateCountQuery := `select count(*) from requests where request_time > sysdate() - 3600 and user_id = (select id from user_details where username = ?);`
	var count int64
	err := db.QueryRow(checkRateCountQuery, username).Scan(&count)
	if err == nil {
		if (plan == "GOLD" && count > 100) ||
			(plan == "SILVER" && count > 50) ||
			(plan == "BRONZE" && count > 25) {
			return false, errors.New("quota exceeded")
		} else {
			return true, nil
		}
	}else {
		return false, err
	}
}

//THonorRateLimitPostProcessing does 2 things.
//1. insert the request struct and request into the database.
//2. check if the entry into the database, violated the rate limit.
//3. if rate limits are correct, the transaction is committed.
//4. if rate limit has been violated, the transaction is rolled back.
func HonorRateLimitPostProessing(username string, plan string, request Request, dataType string, data string) (bool, error) {
	db := GetConnection()
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	insertQuery := `
		insert into requests(id, user_id, data_format, data, status, request_time) values(?, (select id from user_details where username = ?), ?, ?, true, sysdate());
	`
	r, err := tx.Exec(insertQuery, request.ID, username, dataType, data)
	logger.Log.Sugar().Infof("%v", r)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			logger.Log.Sugar().Error("failed to roll back transaction with error %e", err2)
		}
		return false, err
	}

	checkRateCountQuery := `select count(*) from requests where request_time > sysdate() - 3600 and user_id = (select id from user_details where username = ?);`
	var count int64
	err = tx.QueryRow(checkRateCountQuery, username).Scan(&count)
	if err != nil {
		return false, errors.New("failed to get fetch quota count for the user")
	}
	if ( plan == "GOLD" && count > 100 ) ||
		( plan == "SILVER" && count > 50 ) ||
		( plan == "BRONZE" && count > 25 ) {
		err2 := tx.Rollback()
		if err2 != nil {
			logger.Log.Sugar().Error("failed to roll back transaction with error %e", err2)
			return false, err
		}else {
			return false, errors.New("quota exceeded")
		}
	}else{
		tx.Commit()
		return true, nil
	}
}


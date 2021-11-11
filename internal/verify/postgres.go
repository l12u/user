package verify

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/l12u/userm/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type PostgresVerifier struct {
	db *pg.DB
}

func NewPostgresVerifier(opt *pg.Options) *PostgresVerifier {
	db := pg.Connect(opt)
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	// setup tables if not exists
	err = db.Model(&model.RegisteredUser{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		panic(err)
	}

	return &PostgresVerifier{
		db: db,
	}
}

func (p *PostgresVerifier) Verify(user string, pw string) (bool, error) {
	u := &model.RegisteredUser{}
	err := p.db.Model(u).Where("? = ?", pg.Ident("username"), user).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return false, ErrNotFound
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pw))
	if err != nil {
		return false, ErrHashDoesntMatch
	}

	return true, nil
}

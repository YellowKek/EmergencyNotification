package repository

import (
	"app/internal/model"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	CreateTable() error
	Create(user model.User) error
	//GetById(id int32) (*model.User, error)
	//GetEmergencyGroup(id int32) (map[string]string, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() ([]*model.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateTable() error {
	_, err := r.db.Exec(`create table if not exists users (
    id serial primary key,
    email varchar(100) not null,
    name varchar(100) not null,
    surname varchar(100) not null,
    password varchar(100) not null)`)
	_, err = r.db.Exec(`create table if not exists emergency_groups (
    user_id integer,
    type varchar(100) not null,
    value varchar(100) not null
)`)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) Create(user model.User) error {
	var userID int
	err := r.db.QueryRow(`insert into users (name, surname, email, password) values ($1, $2, $3, $4) returning id;`,
		user.Name, user.Surname, user.Email, user.Password).Scan(&userID) // Сохраняем ID нового пользователя
	if err != nil {
		return err
	}
	user.Id = int32(userID)

	logrus.Printf("Created user: %+v with ID: %d", user, userID)
	return nil
}

func (r *UserRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	row := r.db.QueryRow(`select * from users where email = $1`, email)
	user := &model.User{}
	err := row.Scan(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetAll() ([]*model.User, error) {
	rows, err := r.db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Email, &user.Password)
		users = append(users, user)
	}
	return users, nil
}

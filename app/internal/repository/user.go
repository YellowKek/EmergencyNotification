package repository

import (
	"app/internal/model"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	CreateTable() error
	Create(user model.User) error
	GetById(id int32) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() ([]*model.User, error)
	GetEmergencyGroups(userID int32) (map[string]string, error)
	AddEmergencyGroup(userID int32, groupType, value string) error
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
	err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetAll() ([]*model.User, error) {
	rows, err := r.db.Query(`select * from users`)
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

func (r *UserRepositoryImpl) GetEmergencyGroups(userID int32) (map[string]string, error) {
	rows, err := r.db.Query(`select * from emergency_groups where user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	groups := make(map[string]string)
	for rows.Next() {
		var group, value string
		err = rows.Scan(&group, &value)
		if err != nil {
			return nil, err
		}
		groups[group] = value
	}
	return groups, nil
}

func (r *UserRepositoryImpl) GetById(id int32) (*model.User, error) {
	row := r.db.QueryRow(`select * from users where id = $1`, id)
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) AddEmergencyGroup(userID int32, groupType, value string) error {
	_, err := r.db.Exec(`insert into emergency_groups (user_id, type, value) values ($1, $2, $3)`, userID, groupType, value)
	if err != nil {
		return err
	}
	logrus.Print("inserted emergency group")
	return nil
}

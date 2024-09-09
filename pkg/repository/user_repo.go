package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/harshgupta9473/Social/pkg/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUserAccount(user models.TempUser) error {
	query := `insert into users
	(email,encrypted_password,created_at,verified)
	values($1,$2,$3,$4)`
	// encpw, err := bcrypt.GenerateFromPassword([]byte(user.Encrypted_Password), bcrypt.DefaultCost)
	
	_, err := repo.DB.Exec(query, user.Email, user.Encrypted_Password, time.Now().UTC(), true)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) CreateUserProfile(userProfile models.ProfileReq) error {
	query := `insert into userprofiles
	(username,email,firstname,lastname,bio,interests,location,created_at,updated_at)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	interestsJSON, err := json.Marshal(userProfile.Interests)
	if err != nil {
		return err
	}
	_, err = repo.DB.Exec(query, userProfile.UserName, userProfile.Email, userProfile.FirsName, userProfile.LastName, userProfile.Bio,string(interestsJSON),userProfile.Location,  time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `select id, email,encrypted_password,created_at,verified from users where email=$1`
	row, err := repo.DB.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var user = new(models.User)
	if row.Next() {
		err := row.Scan(&user.ID, &user.Email, &user.Encrypted_Password, &user.CreatedAt, &user.Verified)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("user not found")
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) GetUserProfileByEmail(email string) (*models.UserProfile, error) {
	query := `select id,username,email,firstname,lastname,bio,interests,location,created_at,updated_at from userprofiles where email=$1`
	row, err := repo.DB.Query(query, email)
	if err != nil {
		return nil, err
	}
	var userProfile = new(models.UserProfile)
	var interestsJSON string
	if row.Next() {
		err = row.Scan(&userProfile.ID,
			&userProfile.UserName,
			&userProfile.Email,
			&userProfile.FirsName,
			&userProfile.LastName,
			&userProfile.Bio,
			&interestsJSON,
			&userProfile.Location,
			&userProfile.CreatedAt,
			&userProfile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("user not found")
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(interestsJSON), &userProfile.Interests)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (repo *UserRepository) InsertIntoTempUser(email, password, token string) error {
	query := `insert into tempusers
	(email,encrypted_password,token,expires_at)
	values($1,$2,$3,$4)`
	encryptPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var expires_at time.Time = time.Now().UTC().Add(time.Minute * 30)
	_, err = repo.DB.Exec(query, email, string(encryptPass), token, expires_at)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) UpdateTempUser(email, password, token string) error {
	query := `update tempusers set  encrypted_password=$1, token=$2,expires_at=$3 where email=$4`
	encryptPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var expires_at time.Time = time.Now().UTC().Add(time.Minute * 30)
	_, err = repo.DB.Exec(query, string(encryptPass), token, expires_at,email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetTempUserByEmail(email string) (*models.TempUser, error) {
	query := `select id, email,encrypted_password,token,expires_at from tempusers where email=$1`
	rows, err := repo.DB.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tmpUser := new(models.TempUser)
	if rows.Next() {
		err = rows.Scan(&tmpUser.ID, &tmpUser.Email, &tmpUser.Encrypted_Password, &tmpUser.Token, &tmpUser.ExpiresAt)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("not found")
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tmpUser, nil
}

func (repo *UserRepository) GetUserProfileByUserName(username string) (*models.UserProfile, error) {
	query := `select id,username,email,firstname,lastname,bio,interests,location,created_at,updated_at from userprofiles where username=$1`
	row, err := repo.DB.Query(query, username)
	if err != nil {
		return nil, err
	}
	var userProfile = new(models.UserProfile)
	var interestsJSON string
	if row.Next() {
		err = row.Scan(userProfile.ID,
			&userProfile.UserName,
			&userProfile.Email,
			&userProfile.FirsName,
			&userProfile.LastName,
			&userProfile.Bio,
			&interestsJSON,
			&userProfile.Location,
			&userProfile.CreatedAt,
			&userProfile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("user not found")
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(interestsJSON), &userProfile.Interests)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (repo *UserRepository) UpdateUserProfile(user models.ProfileReq) error {
	query := `update userprofiles set username=$1,firstname=$2,lastname=$3,bio=$4,interests=$5,location=$6,updated_at=$7 where email=$8`
	interestsJSON, err := json.Marshal(user.Interests)
	if err != nil {
		interestsJSON = []byte("")
	}
	_, err = repo.DB.Exec(query, user.UserName, user.FirsName, user.LastName, user.Bio, string(interestsJSON), user.Location, time.Now().UTC(),user.Email)
	if err != nil {
		return err
	}
	return nil
}

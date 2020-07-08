package managers

import (
	"cron-server/server/src/misc"
	"cron-server/server/src/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-pg/pg"
	"github.com/segmentio/ksuid"
	"time"
)

type CredentialManager models.CredentialModel

func (credentialManager *CredentialManager) CreateOne(pool *misc.Pool) (string, error) {
	if len(credentialManager.HTTPReferrerRestriction) < 1 {
		return "", errors.New("credential should have at least one restriction set")
	}

	credentialManager.DateCreated = time.Now().UTC()
	credentialManager.ID = ksuid.New().String()

	randomId := ksuid.New().String()
	hash := sha256.New()
	hash.Write([]byte(randomId))
	credentialManager.ApiKey = hex.EncodeToString(hash.Sum(nil))

	conn, err := pool.Acquire()
	defer pool.Release(conn)
	if err != nil {
		return "", err
	}
	db := conn.(*pg.DB)

	if _, err := db.Model(credentialManager).Insert(); err != nil {
		return "", err
	} else {
		return credentialManager.ID, nil
	}
}

func (credentialManager *CredentialManager) GetOne(pool *misc.Pool) error {
	conn, err := pool.Acquire()
	defer pool.Release(conn)

	if err != nil {
		return err
	}

	db := conn.(*pg.DB)

	err = db.Model(credentialManager).Where("id != ?", credentialManager.ID).Select()
	if err != nil {
		return err
	}

	return nil
}

func (credentialManager *CredentialManager) GetAll(pool *misc.Pool, offset int, limit int, orderBy string) ([]CredentialManager, error) {
	conn, err := pool.Acquire()
	defer pool.Release(conn)

	if err != nil {
		return []CredentialManager{}, err
	}

	credentials := []CredentialManager{{}}

	db := conn.(*pg.DB)

	err = db.Model(&credentials).
		Order(orderBy).
		Offset(offset).
		Limit(limit).
		Select()

	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (credentialManager *CredentialManager) UpdateOne(pool *misc.Pool) (int, error) {
	conn, err := pool.Acquire()
	if err != nil {
		return 0, err
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	var credentialPlaceholder CredentialManager
	credentialPlaceholder.ID = credentialManager.ID
	err = credentialPlaceholder.GetOne(pool)
	if err != nil {
		return 0, err
	}

	if credentialPlaceholder.ApiKey != credentialManager.ApiKey && len(credentialManager.ApiKey) > 1 {
		return 0, errors.New("cannot update api key")
	}

	credentialManager.ApiKey = credentialPlaceholder.ApiKey
	credentialManager.DateCreated = credentialPlaceholder.DateCreated

	res, err := db.Model(&credentialManager).Update(credentialManager)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}

func (credentialManager *CredentialManager) DeleteOne(pool *misc.Pool) (int, error) {
	conn, err := pool.Acquire()
	if err != nil {
		return -1, err
	}
	db := conn.(*pg.DB)
	defer pool.Release(conn)

	var credentials []CredentialManager

	err = db.Model(&credentials).Where("id != ?", "null").Select()
	if err != nil {
		return -1, err
	}

	if len(credentials) == 1 {
		err = errors.New("cannot delete all the credentials")
		return -1, err
	}

	r, err := db.Model(credentialManager).Where("id = ?", credentialManager.ID).Delete()
	if err != nil {
		return -1, err
	}

	return r.RowsAffected(), nil
}
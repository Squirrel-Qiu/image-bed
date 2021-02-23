package dbb

import (
	"database/sql"

	"golang.org/x/xerrors"

	"github.com/Squirrel-Qiu/image-bed/handle"
)

type Impl struct {
	DB *sql.DB
}

func (db *Impl) GetIdValue(idType string) (idValue int64, err error) {
	row, err := db.DB.Exec("UPDATE generate_id SET id_value=last_insert_id(id_value+10) WHERE id_type=?", idType)
	if err != nil {
		return -1, xerrors.Errorf("update id value failed: %w", err)
	}

	idValue, err = row.LastInsertId()
	if err != nil {
		return -1, xerrors.Errorf("get last insert id value failed: %w", err)
	}

	return idValue, nil
}

func (db *Impl) Store(resource *handle.Resource) (err error) {
	_, err = db.DB.Exec("INSERT INTO resource (id, bucket, hash, create_time, size) VALUES (?,?,?,?,?)",
		resource.Id, resource.Bucket, resource.Hash, resource.CreateTime, resource.Size)
	if err != nil {
		return xerrors.Errorf("insert resource failed: %w", err)
	}

	return nil
}

func (db *Impl) FileIsExistByHash(hash string) (ok bool, resourceId string, err error) {
	err = db.DB.QueryRow("SELECT id FROM resource WHERE hash=?", hash).Scan(&resourceId)
	switch {
	case xerrors.Is(err, sql.ErrNoRows):
		return true, "", nil

	default:
		return false, "", xerrors.Errorf("scan failed: %w", err)

	case err == nil:
	}

	return false, resourceId, nil
}

func (db *Impl) FileIsExistById(resourceId string) (ok bool, bucket string, err error) {
	err = db.DB.QueryRow("SELECT bucket FROM resource WHERE id=?", resourceId).Scan(&bucket)
	switch {
	case xerrors.Is(err, sql.ErrNoRows):
		return false, "", nil

	default:
		return false, "", xerrors.Errorf("scan failed: %w", err)

	case err == nil:
	}

	return true, bucket, nil
}

func (db *Impl) Close() error {
	return db.DB.Close()
}

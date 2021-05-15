package filestore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/PSauerborn/gamma-project/internal/pkg/filestore"
	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
)

type PostgresPersistence struct {
	*utils.BasePostgresPersistence

	// define settings for file storage
	BaseFilePath string
}

// db function used to retrieve metadata for all files
func (db *PostgresPersistence) ListFiles() ([]filestore.FileMetadata, error) {
	log.Debug("fetching files from postgres storage...")
	files := []filestore.FileMetadata{}

	query := `SELECT file_id,file_name,created,size,metadata FROM file_metadata
	WHERE archived=false`
	rows, err := db.Session.Query(context.Background(), query)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return files, nil
		default:
			return files, err
		}
	}

	for rows.Next() {
		var (
			meta     filestore.FileMetadata
			jsonMeta []byte
		)

		if err := rows.Scan(&meta.FileId, &meta.FileName, &meta.Created,
			&meta.Size, &jsonMeta); err != nil {
			log.Error(fmt.Errorf("unable to read data into local variables: %+v", err))
			continue
		}

		if err := json.Unmarshal(jsonMeta, &meta.Meta); err != nil {
			log.Error(fmt.Errorf("unable to parse JSON metadata: %+v", err))
			continue
		}
		files = append(files, meta)
	}
	return files, nil
}

// db function used to retrieve metadata for a single file
// with a given file ID
func (db *PostgresPersistence) GetFileMetadata(fileId uuid.UUID) (filestore.FileMetadata, error) {
	log.Debug(fmt.Sprintf("fetching file %s metadata from postgres storage...",
		fileId))
	var meta filestore.FileMetadata

	query := `SELECT file_id,file_name,created,size,metadata FROM file_metadata
	WHERE file_id = $1 AND archived=false`
	row := db.Session.QueryRow(context.Background(), query, fileId)
	if err := row.Scan(&meta.FileId, &meta.FileName, &meta.Created,
		&meta.Size, &meta.Meta); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return meta, filestore.ErrFileNotFound
		default:
			return meta, err
		}
	}
	return meta, nil
}

// db function used to retrieve file contents from local disk storage
func (db *PostgresPersistence) GetFileContents(meta filestore.FileMetadata) ([]byte, error) {
	log.Debug(fmt.Sprintf("fetching file contents for %+v", meta))
	// open file with given file path
	contents, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", db.BaseFilePath, meta.FileId))
	if err != nil {
		log.Error(fmt.Errorf("unable to open file %s: %+v", meta.FileName, err))
		return []byte{}, err
	}
	return contents, nil
}

// db function used to create a new file
func (db *PostgresPersistence) CreateFile(content []byte, fileName string,
	meta map[string]interface{}) (uuid.UUID, error) {
	log.Debug("inserting new file into postgres storage...")
	fileId := uuid.New()

	jsonBody, err := json.Marshal(meta)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert file metadata to JSON: %+v", err))
		return fileId, errors.New("invalid file metadata")
	}

	query := `INSERT INTO file_metadata(file_id,file_name,size,metadata)
	VALUES($1,$2,$3,$4)`
	_, err = db.Session.Exec(context.Background(), query, fileId, fileName,
		len(content), jsonBody)
	if err != nil {
		return fileId, err
	}
	// construct file path and write to files
	fpath := fmt.Sprintf("%s/%s", db.BaseFilePath, fileId)
	return fileId, ioutil.WriteFile(fpath, content, 0644)
}

// db function used to modify an existing file metadata
func (db *PostgresPersistence) ModifyFile(meta filestore.FileMetadata, contents []byte) error {
	log.Debug(fmt.Sprintf("modifying file %s postgres storage...", meta.FileId))
	return filestore.ErrFeatureNotSupported
}

// db function used to delete a particular file with given file ID
func (db *PostgresPersistence) DeleteFile(meta filestore.FileMetadata) error {
	log.Debug(fmt.Sprintf("deleting file %s from postgres storage...", meta.FileId))

	fpath := fmt.Sprintf("%s/%s", db.BaseFilePath, meta.FileId)
	if err := os.Remove(fpath); err != nil {
		log.Error(fmt.Errorf("unable to delete file: %+v", err))
		return filestore.ErrCannotDeleteFile
	}

	query := `DELETE FROM file_metadata WHERE file_id = $1`
	_, err := db.Session.Exec(context.Background(), query, meta.FileId)
	if err != nil {
		return err
	}
	return nil
}

// db function used to archive file
func (db *PostgresPersistence) ArchiveFile(meta filestore.FileMetadata) error {
	log.Debug("archieving file %+v...", meta)
	// construct current and target directories
	current := fmt.Sprintf("%s/%s", db.BaseFilePath, meta.FileId)
	target := fmt.Sprintf("%s/archive/%s", db.BaseFilePath, meta.FileId)
	if err := os.Rename(current, target); err != nil {
		log.Error(fmt.Errorf("cannot move files: %+v", err))
		return err
	}

	query := `UPDATE file_metadata SET archived=true WHERE file_id=$1`
	_, err := db.Session.Exec(context.Background(), query, meta.FileId)
	return err
}

// db function to search meta
func (db *PostgresPersistence) SearchFilesByMetadata(terms map[string]interface{}) (
	[]filestore.FileMetadata, error) {
	log.Debug(fmt.Sprintf("searching files with terms %+v", terms))
	matches := []filestore.FileMetadata{}
	// retrieve file list from database
	files, err := db.ListFiles()
	if err != nil {
		return matches, err
	}

	// iterate over files and append to results
	// if provided metadata fields all match
	for _, f := range files {
		if filestore.MapMatchesTerms(f.Meta, terms, filestore.CompleteMatch) {
			matches = append(matches, f)
		}
	}
	return matches, nil
}

package repository

import (
	"context"
	"database/sql"

	db "go-template/db/sqlc"
	"go-template/internal/entity"
)

type FileRepository interface {
	Create(ctx context.Context, fileName, originalName, filePath string, fileSize int64, mimeType, description, category string, uploadedBy int) (*entity.File, error)
	GetByID(ctx context.Context, id int) (*entity.File, error)
	GetByUserID(ctx context.Context, userID int) ([]entity.File, error)
	Update(ctx context.Context, id int, description, category string) (*entity.File, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.File, error)
}

type fileRepository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewFileRepository(dbConn *sql.DB) FileRepository {
	return &fileRepository{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (r *fileRepository) Create(ctx context.Context, fileName, originalName, filePath string, fileSize int64, mimeType, description, category string, uploadedBy int) (*entity.File, error) {
	createdFile, err := r.queries.CreateFile(ctx, db.CreateFileParams{
		FileName:     fileName,
		OriginalName: originalName,
		FilePath:     filePath,
		FileSize:     fileSize,
		MimeType:     mimeType,
		Description:  sql.NullString{String: description, Valid: description != ""},
		Category:     sql.NullString{String: category, Valid: category != ""},
		UploadedBy:   int32(uploadedBy),
	})
	if err != nil {
		return nil, err
	}

	return r.mapDBFileToEntity(&createdFile), nil
}

func (r *fileRepository) GetByID(ctx context.Context, id int) (*entity.File, error) {
	file, err := r.queries.GetFile(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return r.mapDBFileToEntity(&file), nil
}

func (r *fileRepository) GetByUserID(ctx context.Context, userID int) ([]entity.File, error) {
	dbFiles, err := r.queries.GetFilesByUser(ctx, int32(userID))
	if err != nil {
		return nil, err
	}

	files := make([]entity.File, len(dbFiles))
	for i, dbFile := range dbFiles {
		files[i] = *r.mapDBFileToEntity(&dbFile)
	}

	return files, nil
}

func (r *fileRepository) Update(ctx context.Context, id int, description, category string) (*entity.File, error) {
	updatedFile, err := r.queries.UpdateFile(ctx, db.UpdateFileParams{
		ID:          int32(id),
		Description: sql.NullString{String: description, Valid: description != ""},
		Category:    sql.NullString{String: category, Valid: category != ""},
	})
	if err != nil {
		return nil, err
	}

	return r.mapDBFileToEntity(&updatedFile), nil
}

func (r *fileRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteFile(ctx, int32(id))
}

func (r *fileRepository) GetAll(ctx context.Context) ([]entity.File, error) {
	dbFiles, err := r.queries.GetAllFiles(ctx)
	if err != nil {
		return nil, err
	}

	files := make([]entity.File, len(dbFiles))
	for i, dbFile := range dbFiles {
		files[i] = *r.mapDBFileToEntity(&dbFile)
	}

	return files, nil
}

func (r *fileRepository) mapDBFileToEntity(dbFile *db.Files) *entity.File {
	return &entity.File{
		ID:           int(dbFile.ID),
		FileName:     dbFile.FileName,
		OriginalName: dbFile.OriginalName,
		FilePath:     dbFile.FilePath,
		FileSize:     dbFile.FileSize,
		MimeType:     dbFile.MimeType,
		Description:  dbFile.Description.String,
		Category:     dbFile.Category.String,
		UploadedBy:   int(dbFile.UploadedBy),
		CreatedAt:    dbFile.CreatedAt.Time,
		UpdatedAt:    dbFile.UpdatedAt.Time,
	}
}
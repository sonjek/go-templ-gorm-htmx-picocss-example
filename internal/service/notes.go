package service

import (
	"errors"
	"fmt"

	database "github.com/sonjek/go-full-stack-example/internal/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const MsgTitleAlreadyExists = "Title already used by another note"

var (
	ErrRecordNotFound = fmt.Errorf("record not found: %w", gorm.ErrRecordNotFound)
	ErrDuplicateTitle = errors.New("duplicate title: title already used by another note")
)

type NoteService struct {
	db       *gorm.DB
	pageSize int
}

func NewNoteService(db *gorm.DB, pageSize int) *NoteService {
	return &NoteService{
		db:       db,
		pageSize: pageSize,
	}
}

func (s *NoteService) LoadMore(cursorID int) ([]database.Note, error) {
	var notes []database.Note

	var result *gorm.DB
	if cursorID < 1 {
		result = s.db.Limit(s.pageSize).Order("id DESC").Find(&notes)
	} else {
		result = s.db.Where("id < ?", cursorID).Order("id DESC").Limit(s.pageSize).Find(&notes)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return notes, nil
}

func (s *NoteService) Create(title, body string) (database.Note, error) {
	note := database.Note{
		Title: title,
		Body:  body,
	}
	if err := s.db.Create(&note).Error; err != nil {
		if isUniqueConstraintViolation(err) {
			return database.Note{}, ErrDuplicateTitle
		}
		return database.Note{}, err
	}
	return note, nil
}

func (s *NoteService) Get(noteID int) (database.Note, error) {
	var note database.Note
	result := s.db.First(&note, noteID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return database.Note{}, ErrRecordNotFound
		}
		return database.Note{}, result.Error
	}
	return note, nil
}

func (s *NoteService) FindAndUpdate(noteID int, title, body string) (database.Note, error) {
	var note database.Note

	updateData := map[string]any{
		"title": title,
		"body":  body,
	}

	// Updates the record and populates 'note' via RETURNING clause
	result := s.db.Model(&note).
		Clauses(clause.Returning{}).
		Where("id = ?", noteID).
		Updates(updateData)

	if result.Error != nil {
		if isUniqueConstraintViolation(result.Error) {
			return database.Note{}, ErrDuplicateTitle
		}
		return database.Note{}, result.Error
	}

	if result.RowsAffected == 0 {
		return database.Note{}, ErrRecordNotFound
	}

	return note, nil
}

func (s *NoteService) Delete(noteID int) error {
	result := s.db.Delete(&database.Note{}, noteID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}

func IsDuplicateTitle(err error) bool {
	return errors.Is(err, ErrDuplicateTitle)
}

func isUniqueConstraintViolation(err error) bool {
	if sqliteErr, ok := errors.AsType[*sqlite.Error](err); ok {
		return sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
	}
	return false
}

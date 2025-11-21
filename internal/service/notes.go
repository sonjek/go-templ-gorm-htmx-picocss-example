package service

import (
	database "github.com/sonjek/go-templ-htmx-picocss-example/internal/storage"
	"gorm.io/gorm"
)

type NoteService struct {
	db *gorm.DB
}

func NewNoteService(db *gorm.DB) *NoteService {
	return &NoteService{
		db: db,
	}
}

func (s *NoteService) LoadMore(cursorID, pageSize int) ([]database.Note, error) {
	var notes []database.Note

	var result *gorm.DB
	if cursorID < 1 {
		result = s.db.Limit(pageSize).Order("id DESC").Find(&notes)
	} else {
		result = s.db.Where("id < ?", cursorID).Order("id DESC").Limit(pageSize).Find(&notes)
	}

	if result.Error != nil {
		return notes, result.Error
	}

	return notes, nil
}

func (s *NoteService) Create(title, body string) database.Note {
	note := database.Note{
		Title: title,
		Body:  body,
	}
	s.db.Create(&note)

	return note
}

func (s *NoteService) Get(noteID int) database.Note {
	note := database.Note{ID: noteID}
	s.db.First(&note)
	return note
}

func (s *NoteService) Update(note database.Note, title, body string) {
	note.Title = title
	note.Body = body
	s.db.Save(&note)
}

func (s *NoteService) FindAndUpdate(noteID int, title, body string) database.Note {
	note := s.Get(noteID)
	s.Update(note, title, body)
	return s.Get(noteID)
}

func (s *NoteService) Delete(noteID string) {
	s.db.Delete(&database.Note{}, noteID)
}

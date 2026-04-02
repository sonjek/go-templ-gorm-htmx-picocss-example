package storage

import (
	"time"

	"gorm.io/gorm"
)

var NotesSeed = []Note{
	{
		Title:     "gorm",
		Body:      "The fantastic ORM library for Golang, aims to be developer friendly.",
		CreatedAt: time.Date(2013, 10, 25, 8, 24, 0, 0, time.UTC),
	},
	{
		Title:     "Ionic",
		Body:      "Premium hand-crafted icons built by Ionic, for Ionic apps and web apps everywhere.",
		CreatedAt: time.Date(2013, 10, 30, 12, 34, 0, 0, time.UTC),
	},
	{
		Title:     "cznic/sqlite",
		Body:      "Package sqlite is a CGo-free port of SQLite/SQLite3.",
		CreatedAt: time.Date(2017, 4, 20, 23, 17, 29, 0, time.UTC),
	},
	{
		Title:     "Tailwind CSS",
		Body:      "A utility-first CSS framework for rapid UI development.",
		CreatedAt: time.Date(2017, 11, 1, 21, 11, 0, 0, time.UTC),
	},
	{
		Title:     "GoLang",
		Body:      "The Go programming language.",
		CreatedAt: time.Date(2018, 9, 25, 22, 20, 0, 0, time.UTC),
	},
	{
		Title:     "Pico CSS",
		Body:      "Minimal CSS Framework for semantic HTML.",
		CreatedAt: time.Date(2019, 12, 11, 10, 8, 0, 0, time.UTC),
	},
	{
		Title:     "templ",
		Body:      "A language for writing HTML user interfaces in Go.",
		CreatedAt: time.Date(2021, 5, 16, 22, 33, 0, 0, time.UTC),
	},
	{
		Title:     "htmx.js",
		Body:      "</> htmx - high power tools for HTML.",
		CreatedAt: time.Date(2022, 8, 10, 21, 21, 0, 0, time.UTC),
	},
	{
		Title:     "daisyUI",
		Body:      "The most popular, free and open-source Tailwind CSS component library.",
		CreatedAt: time.Date(2023, 11, 13, 12, 49, 0, 0, time.UTC),
	},
}

type Note struct {
	ID        int            `form:"id" json:"id" gorm:"primaryKey"`
	Title     string         `form:"title" json:"title" binding:"required"`
	Body      string         `form:"body" json:"body" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

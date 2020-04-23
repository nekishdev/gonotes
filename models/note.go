package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Note struct {
	gorm.Model
	UserId uint `json:"user_id"`
	Text string `json:"text"`
}

func CreateNote(userId uint, text string) (*Note, error) {
	note := &Note{
		UserId: userId,
		Text: text,
	}
	GetDB().Create(note)
	if note.ID <= 0 {
		return nil, errors.New("Failed to create note")
	}
	return note, nil
}

func GetNote(userId uint, id int) *Note {
	note := &Note{}
	GetDB().First(note, "user_id = ? and id = ?", userId, id)
	if note.ID <= 0 {
		return nil
	}
	return note
}

func GetNotesByUser(userId uint) []Note {
	notes := make([]Note, 0)
	GetDB().Table("notes").Where("user_id = ?", userId).Find(&notes)
	return notes
}

func UpdateNote(note *Note, text string) error {
	return GetDB().Model(note).Update("Text", text).Error
}

func DeleteNote(note *Note) error {
	return GetDB().Delete(note).Error
}

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testProject/models"
	u "testProject/utils"
)

func NoteDetailHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	id, _ := strconv.Atoi(u.GetUrlVar(r, "id"))
	note := models.GetNote(userId, id)
	if note == nil {
		u.Abort(w, 404, "Not found")
		return
	}
	u.Respond(w, note)
}

func NoteCreateHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	note := &models.Note{}
	err := json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		u.Abort(w, 400, fmt.Sprintf("Invalid request. Field to parse request body. %s", err))
		return
	}
	if len(note.Text) <= 0 {
		u.Abort(w, 400, "Text required")
		return
	}

	_note, err := models.CreateNote(userId, note.Text)
	if err != nil {
		u.Abort(w, 400, err)
		return
	}
	u.Respond(w, _note)
}

func NoteListHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	notes := models.GetNotesByUser(userId)

	if len(notes) <= 0 {
		u.Abort(w, 404, "Not found")
		return
	}
	u.Respond(w, notes)
}

func NoteUpdateHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	request := &struct {
		Text string
	}{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		u.Abort(w, 400, fmt.Sprintf("Invalid request. Field to parse request body. %s", err))
		return
	}
	if len(request.Text) <= 0 {
		u.Abort(w, 400, "Text required")
		return
	}
	id, _ := strconv.Atoi(u.GetUrlVar(r, "id"))

	note := models.GetNote(userId, id)
	if note == nil {
		u.Abort(w, 404, "Not found")
		return
	}

	err = models.UpdateNote(note, request.Text)
	if err != nil {
		u.Abort(w, 400, "Failed to update note")
		return
	}
	u.Respond(w, "Updated")
}

func NoteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	id, _ := strconv.Atoi(u.GetUrlVar(r, "id"))

	note := models.GetNote(userId, id)
	if note == nil {
		u.Abort(w, 404, "Not found")
		return
	}

	err := models.DeleteNote(note)
	if err != nil {
		u.Abort(w, 400, "Failed to delete note")
		return
	}
	u.Respond(w, "Deleted")
}

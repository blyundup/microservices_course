package note

import (
    "github.com/blyundup/microservices_course/week_3/internal/repository"
    def "github.com/blyundup/microservices_course/week_3/internal/service"
)

var _ def.NoteService = (*serv)(nil)

type serv struct {
	noteRepository repository.NoteRepository
}

func NewService(noteRepository repository.NoteRepository) *serv {
	return &serv{
		noteRepository: noteRepository,
	}
}

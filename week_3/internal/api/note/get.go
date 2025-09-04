package note

import (
    "context"
    "log"

    "github.com/blyundup/microservices_course/week_3/internal/converter"
    desc "github.com/blyundup/microservices_course/week_3/pkg/note_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	noteObj, err := i.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %v, title: %v, body: %v, created_at: %v, updated_at: %v", noteObj.ID, noteObj.Info.Title, noteObj.Info.Content, noteObj.CreatedAt, noteObj.UpdatedAt)

	return &desc.GetResponse{
		Note: converter.ToNoteFromService(noteObj),
	}, nil
}

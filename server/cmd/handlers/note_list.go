package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListNote(ctx *mhttp.Context) {
	req := &api.ListNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if req.Page.Size <= 0 || req.Page.Num <= 0 {
		e := utils.ErrInvalidParams().WithParam("pagination", req.Page)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	writer := ""
	if req.OnlyOperator && len(ctx.UserName) > 0 {
		writer = ctx.UserName
	}

	count, notes, e := dal.ListNotes(req.Page.Size, req.Page.Num, writer)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ListNoteRes{
		Count: count,
		Notes: notesDBToHTTP(notes),
	}
}

func notesDBToHTTP(notes []*model.Note) []*api.Note {
	res := make([]*api.Note, len(notes))
	for i, v := range notes {
		res[i] = &api.Note{
			NoteID:      v.NoteID,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Writer:      v.WriterName,
			IsAnonymous: v.IsAnonymous,
			Title:       v.Title,
			Content:     v.Content,
		}
	}

	return res
}

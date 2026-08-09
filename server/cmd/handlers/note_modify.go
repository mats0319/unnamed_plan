package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyNote(ctx *mhttp.Context) {
	req := &api.ModifyNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(ctx.UserName) < 1 || len(req.NoteID) < 1 || len(req.Content) < 1 {
		e := utils.ErrInvalidParams().WithParam("operator", ctx.UserName).
			WithParam("note id", req.NoteID).WithParam("content", req.Content)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	note, e := dal.GetNote(req.NoteID)
	if e != nil {
		ctx.ResData = e
		return
	}

	if ctx.UserName != note.Writer {
		e := utils.ErrPermissionDenied().WithParam("need writer but get", ctx.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	if req.IsAnonymous == note.IsAnonymous && req.Title == note.Title && req.Content == note.Content {
		e := utils.ErrNoChanges().WithParam("operator", ctx.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	note.IsAnonymous = req.IsAnonymous
	note.Title = req.Title
	note.Content = req.Content

	e = dal.UpdateNote(note)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ModifyNoteRes{}
}

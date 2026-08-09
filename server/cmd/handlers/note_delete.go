package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func DeleteNote(ctx *mhttp.Context) {
	req := &api.DeleteNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(ctx.UserName) < 1 || len(req.NoteID) < 1 {
		e := utils.ErrInvalidParams().WithParam("operator", ctx.UserName).WithParam("note id", req.NoteID)
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
		e := utils.ErrPermissionDenied().WithParam("need owner but get", ctx.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	e = dal.DeleteNote(req.NoteID)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.DeleteNoteRes{}
}

package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func CreateNote(ctx *mhttp.Context) {
	req := &api.CreateNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(ctx.UserName) < 1 || len(req.Content) < 1 {
		e := utils.ErrInvalidParams().WithParam("operator", ctx.UserName).WithParam("content", req.Content)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	operator, e := dal.GetUser(ctx.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	note := &model.Note{
		Writer:      operator.UserName,
		WriterName:  operator.Nickname,
		IsAnonymous: req.IsAnonymous,
		Title:       req.Title,
		Content:     req.Content,
	}

	e = dal.CreateNote(note)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.CreateNoteRes{}
}

package structure

// user

func MakeLoginRes(userID string, nickname string, permission uint8) interface{} {
	return &struct {
		*UserID
		*Nickname
		*Permission
	}{
		&UserID{UserID: userID},
		&Nickname{Nickname: nickname},
		&Permission{Permission: permission},
	}
}

func MakeListUserRes(total int, users []*UserListRes) interface{} {
	return &struct {
		*Total
		*Users
	}{
		&Total{Total: total},
		&Users{Users: users},
	}
}

func MakeCreateUserRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeLockUserRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeUnlockUserRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeModifyUserInfoRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeModifyUserPermissionRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

// cloud file

func MakeListCloudFileByUploaderRes(total int, files []*FileListRes) interface{} {
	return &struct {
		*Total
		*Files
	}{
		&Total{Total: total},
		&Files{Files: files},
	}
}

func MakeListPublicCloudFileRes(total int, files []*FileListRes) interface{} {
	return &struct {
		*Total
		*Files
	}{
		&Total{Total: total},
		&Files{Files: files},
	}
}

func MakeUploadCloudFileRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeModifyCloudFileRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeDeleteCloudFileRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

// note

func MakeListThinkingNoteByWriterRes(total int, notes []*NoteListRes) interface{} {
	return &struct {
		*Total
		*Notes
	}{
		&Total{Total: total},
		&Notes{Notes: notes},
	}
}

func MakeListPublicThinkingNoteRes(total int, notes []*NoteListRes) interface{} {
	return &struct {
		*Total
		*Notes
	}{
		&Total{Total: total},
		&Notes{Notes: notes},
	}
}

func MakeCreateThinkingNoteRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeModifyThinkingNoteRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

func MakeDeleteThinkingNoteRes(isSuccess bool) interface{} {
	return &struct {
		*IsSuccess
	}{
		&IsSuccess{IsSuccess: isSuccess},
	}
}

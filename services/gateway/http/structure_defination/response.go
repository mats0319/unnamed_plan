package structure

// user

func MakeLoginRes(userID string, nickname string, permission uint32) interface{} {
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

func MakeListUserRes(total uint32, users []*UserRes) interface{} {
	return &struct {
		*Total
		*Users
	}{
		&Total{Total: total},
		&Users{Users: users},
	}
}

// cloud file

func MakeListCloudFileByUploaderRes(total uint32, files []*FileRes) interface{} {
	return &struct {
		*Total
		*Files
	}{
		&Total{Total: total},
		&Files{Files: files},
	}
}

func MakeListPublicCloudFileRes(total uint32, files []*FileRes) interface{} {
	return &struct {
		*Total
		*Files
	}{
		&Total{Total: total},
		&Files{Files: files},
	}
}

// note

func MakeListNoteByWriterRes(total uint32, notes []*NoteRes) interface{} {
	return &struct {
		*Total
		*Notes
	}{
		&Total{Total: total},
		&Notes{Notes: notes},
	}
}

func MakeListPublicNoteRes(total uint32, notes []*NoteRes) interface{} {
	return &struct {
		*Total
		*Notes
	}{
		&Total{Total: total},
		&Notes{Notes: notes},
	}
}

package structure

func MakeLoginRes(userID string, userName string) interface{} {
	return &struct {
		*UserID
		*UserName
	}{
		&UserID{UserID: userID},
		&UserName{UserName: userName},
	}
}

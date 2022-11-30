package main

func whiteSpaces(times int) []byte {
	res := make([]byte, 0, times)
	for i := 0; i < times; i++ {
		res = append(res, ' ')
	}

	return res
}

// isEnum return if target payload type is enum(by payload name)
func isEnum(payloadList []*payload, payloadName string) bool {
	isEnumFlag := false
ALL:
	for i := range payloadList {
		item := payloadList[i]
		if item.name == payloadName {
			isEnumFlag = item.typ == PayloadType_Enum
			break
		}

		for j := range item.children {
			if item.children[j].name == payloadName {
				isEnumFlag = item.children[j].typ == PayloadType_Enum
				break ALL
			}
		}
	}

	return isEnumFlag
}

package responses

type OK struct {
	OK bool `json:"ok"`
}

func NewOK(ok bool) OK {
	return OK{ok}
}

type IDList struct {
	IDList []uint `json:"id_list"`
}

func NewIDList(idList []uint) IDList {
	return IDList{idList}
}

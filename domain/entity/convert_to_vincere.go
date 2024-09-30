package entity

import (
	"gopkg.in/guregu/null.v4"
)

// 有or無のチェック
func ConvertGenderToVincere(inputValue null.Int) string {
	// 空文字の場合は無しとする
	if inputValue == null.NewInt(0, true) {
		return "MALE"
	} else if inputValue == null.NewInt(1, true) {
		return "FEMALE"
	} else {
		return ""
	}
}

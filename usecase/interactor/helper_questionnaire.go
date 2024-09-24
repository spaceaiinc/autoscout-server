package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

// 志望度
func getStrMyRanking(inputNumber null.Int) string {
	if inputNumber.Valid {
		for index, value := range entity.MyRanking {
			if inputNumber == null.NewInt(int64(index), true) {
				return value
			}
		}
	}

	return ""
}

// 他社の志望度
func getStrOtherMyRanking(inputNumber null.Int) string {
	if inputNumber.Valid {
		for index, value := range entity.OtherMyRanking {
			if inputNumber == null.NewInt(int64(index), true) {
				return value
			}
		}
	}

	return ""
}

// 内定が出た場合の意向
func getStrIntentionToJobOffer(inputNumber null.Int) string {
	if inputNumber.Valid {
		for index, value := range entity.IntentionToJobOffer {
			if inputNumber == null.NewInt(int64(index), true) {
				return value
			}
		}
	}

	return ""
}

// 選考継続希望有無
func getStrContinueSelectionType(inputNumber null.Int) string {
	if inputNumber.Valid {
		for index, value := range entity.ContinueSelectionType {
			if inputNumber == null.NewInt(int64(index), true) {
				return value
			}
		}
	}

	return ""
}

// 選考後アンケートの志望企業のフェーズ
func getStrSelectionPhase(inputNumber null.Int) string {
	if inputNumber.Valid {
		for index, value := range entity.SelectionPhase {
			if inputNumber == null.NewInt(int64(index), true) {
				return value
			}
		}
	}

	return ""
}

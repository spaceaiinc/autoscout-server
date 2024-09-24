package interactor

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 求人の最大ページ数を返す（本番実装までは1ページあたり50件）
func getSaleListMaxPage(saleList []*entity.Sale) uint {
	var maxPage = len(saleList) / 50

	if 0 < (len(saleList) % 50) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの求人一覧を返す（本番実装までは1ページあたり50件）
func getSaleListWithPage(saleList []*entity.Sale, page uint) []*entity.Sale {
	var (
		perPage uint = 50
		listLen uint = uint(len(saleList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return saleList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.Sale{}
	}

	if (listLen - first) <= perPage {
		return saleList[first:]
	}
	return saleList[first:last]
}

// 検索の期間(period)の値に応じて日付（形式: 2023-03）の文字列を返す
func changePeriodValue(fiscalYear string, period uint) ([]string, error) {
	var p []string

	layout := "2006-01" // 入力文字列のフォーマット
	t, err := time.Parse(layout, fiscalYear)
	if err != nil {
		fmt.Println("Error:", err)
		return p, err
	}

	switch period {
	case 0:
		// 決算月の11ヶ月前
		elevenMonthsAgo := t.AddDate(0, -11, 0)
		elevenMonthsAgoStr := elevenMonthsAgo.Format("2006-01")
		p = append(p, elevenMonthsAgoStr)
		break
	case 1:
		// 決算月の10ヶ月前
		tenMonthsAgo := t.AddDate(0, -10, 0)
		tenMonthsAgoStr := tenMonthsAgo.Format("2006-01")
		p = append(p, tenMonthsAgoStr)
		break
	case 2:
		// 決算月の9ヶ月前
		nineMonthsAgo := t.AddDate(0, -9, 0)
		nineMonthsAgoStr := nineMonthsAgo.Format("2006-01")
		p = append(p, nineMonthsAgoStr)
		break
	case 3:
		// 決算月の8ヶ月前
		eightMonthsAgo := t.AddDate(0, -8, 0)
		eightMonthsAgoStr := eightMonthsAgo.Format("2006-01")
		p = append(p, eightMonthsAgoStr)
		break
	case 4:
		// 決算月の7ヶ月前
		sevenMonthsAgo := t.AddDate(0, -7, 0)
		sevenMonthsAgoStr := sevenMonthsAgo.Format("2006-01")
		p = append(p, sevenMonthsAgoStr)
		break
	case 5:
		// 決算月の6ヶ月前
		sixMonthsAgo := t.AddDate(0, -6, 0)
		sixMonthsAgoStr := sixMonthsAgo.Format("2006-01")
		p = append(p, sixMonthsAgoStr)
		break
	case 6:
		// 決算月の5ヶ月前
		fiveMonthsAgo := t.AddDate(0, -5, 0)
		fiveMonthsAgoStr := fiveMonthsAgo.Format("2006-01")
		p = append(p, fiveMonthsAgoStr)
		break
	case 7:
		// 決算月の4ヶ月前
		fourMonthsAgo := t.AddDate(0, -4, 0)
		fourMonthsAgoStr := fourMonthsAgo.Format("2006-01")
		p = append(p, fourMonthsAgoStr)
		break
	case 8:
		// 決算月の3ヶ月前
		threeMonthsAgo := t.AddDate(0, -3, 0)
		threeMonthsAgoStr := threeMonthsAgo.Format("2006-01")
		p = append(p, threeMonthsAgoStr)
		break
	case 9:
		// 決算月の2ヶ月前
		twoMonthsAgo := t.AddDate(0, -2, 0)
		twoMonthsAgoStr := twoMonthsAgo.Format("2006-01")
		p = append(p, twoMonthsAgoStr)
		break
	case 10:
		// 決算月の1ヶ月前
		oneMonthAgo := t.AddDate(0, -1, 0)
		oneMonthAgoStr := oneMonthAgo.Format("2006-01")
		p = append(p, oneMonthAgoStr)
		break
	case 11:
		// 決算月
		p = append(p, fiscalYear)
	case 12:
		// 決算月の11ヶ月前 〜 決算月の9ヶ月前
		elevenMonthsAgo := t.AddDate(0, -11, 0)
		elevenMonthsAgoStr := elevenMonthsAgo.Format("2006-01")

		nineMonthsAgo := t.AddDate(0, -9, 0)
		nineMonthsAgoStr := nineMonthsAgo.Format("2006-01")

		p = append(p, elevenMonthsAgoStr, nineMonthsAgoStr)
		break
	case 13:
		// 決算月の8ヶ月前 〜 決算月の6ヶ月前
		eightMonthsAgo := t.AddDate(0, -8, 0)
		eightMonthsAgoStr := eightMonthsAgo.Format("2006-01")

		sixMonthsAgo := t.AddDate(0, -6, 0)
		sixMonthsAgoStr := sixMonthsAgo.Format("2006-01")

		p = append(p, eightMonthsAgoStr, sixMonthsAgoStr)
		break
	case 14:
		// 決算月の5ヶ月前 〜 決算月の3ヶ月前"
		fiveMonthsAgo := t.AddDate(0, -5, 0)
		fiveMonthsAgoStr := fiveMonthsAgo.Format("2006-01")

		threeMonthsAgo := t.AddDate(0, -3, 0)
		threeMonthsAgoStr := threeMonthsAgo.Format("2006-01")

		p = append(p, fiveMonthsAgoStr, threeMonthsAgoStr)
		break
	case 15:
		// 決算月の2ヶ月前 〜 決算月"
		twoMonthsAgo := t.AddDate(0, -2, 0)
		twoMonthsAgoStr := twoMonthsAgo.Format("2006-01")

		p = append(p, twoMonthsAgoStr, fiscalYear)
		break
	case 16:
		// 決算月の11ヶ月前 〜 決算月"
		elevenMonthsAgo := t.AddDate(0, -11, 0)
		elevenMonthsAgoStr := elevenMonthsAgo.Format("2006-01")

		p = append(p, elevenMonthsAgoStr, fiscalYear)
		break
	}

	return p, nil
}

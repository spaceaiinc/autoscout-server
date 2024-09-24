package utility

import "regexp"

// regexpのコンパイル処理は時間がかかるため、サーバー起動時に一度だけの実行で済むように関数外で宣言しておく
// ref:https://budougumi0617.github.io/2020/08/20/regexponce/
var (
	RegexpForLineBreak            = regexp.MustCompile(`\r\n|\n`) // 改行コードの正規表現
	RegexpForAmbiUserID           = regexp.MustCompile(`会員No\.(\d+)`)
	RegexpForMynaviScoutingUserID = regexp.MustCompile(`会員No\.[\s　]*：[\s　]*(\d+)`) // 全角スペースを含む正規表現（直接全角スペースを入力）
	RegexpForMynaviAgentScoutUserID = regexp.MustCompile(`求職者ID]\s*(\d+)`)
)

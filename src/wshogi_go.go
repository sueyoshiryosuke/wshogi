package main

/*
	Cのdllとしてコンパイルするときのコマンドは以下。
	go build -buildmode=c-shared -o wshogi_go.dll wshogi_go.go

	以下、注意事項。
	import "C"
	と関数定義の直前の行に、次の文字列が必要。
	「//export sfen2csa」
*/

// #include <stdlib.h>
import "C"

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

/*
	盤面の配置は空っぽ。
	盤面の番地は0から120。壁は-1。
	成り駒の種類の値は先手駒+100、後手駒-100。
	番地121から136は、持ち駒の枚数。
	番地137は、手番。先手番1、後手番-1。
	先手の持ち駒。
		121: 歩,
		122: 香,
		123: 桂,
		124: 銀,
		125: 金,
		126: 角,
		127: 飛,
	後手の持ち駒
		129: 歩,
		130: 香,
		（中略）
		135: 飛,
	将棋盤とコマンドの変換表
	 壁|  9,  8,  7,  6,  5,  4,  3, 2 ,  1| 壁|
	---|-----------------------------------|---|
	  0|  1,  2,  3,  4,  5,  6,  7,  8,  9| 10|壁
	---|-----------------------------------|---|
	 11| 12, 13, 14, 15, 16, 17, 18, 19, 20| 21|一
	   | 9a, 8a, 7a, 6a, 5a, 4a, 3a, 2a, 1a|   |a↑
	 22| 23, 24, 25, 26, 27, 28, 29, 30, 31| 32|二
	   | 9b, 8b, 7b, 6b, 5b, 4b, 3b, 2b, 1b|   |b↑
	 33| 34, 35, 36, 37, 38, 39, 40, 41, 42| 43|三
	   | 9c, 8c, 7c, 6c, 5c, 4c, 3c, 2c, 1c|   |c↑
	 44| 45, 46, 47, 48, 49, 50, 51, 52, 53| 54|四
	   | 9d, 8d, 7d, 6d, 5d, 4d, 3d, 2d, 1d|   |d↑
	 55| 56, 57, 58, 59, 60, 61, 62, 63, 64| 65|五
	   | 9e, 8e, 7e, 6e, 5e, 4e, 3e, 2e, 1e|   |e↑
	 66| 67, 68, 69, 70, 71, 72, 73, 74, 75| 76|六
	   | 9f, 8f, 7f, 6f, 5f, 4f, 3f, 2f, 1f|   |f↑
	 77| 78, 79, 80, 81, 82, 83, 84, 85, 86| 87|七
	   | 9g, 8g, 7g, 6g, 5g, 4g, 3g, 2g, 1g|   |g↑
	 88| 89, 90, 91, 92, 93, 94, 95, 96, 97| 98|八
	   | 9h, 8h, 7h, 6h, 5h, 4h, 3h, 2h, 1h|   |h↑
	 99|100,101,102,103,104,105,106,107,108|109|九
	   | 9i, 8i, 7i, 6i, 5i, 4i, 3i, 2i, 1i|   |i↑
	---|-----------------------------------|---|
	110|111,112,113,114,115,116,117,118,119|120|壁

	壁には-1、空っぽは0です。
*/

// dllにするためmain関数は使わない。
func main() {}

//export sfen2csa
func sfen2csa(sfen_input_str *C.char) *C.char {
	/*
		文字列sfenから表示用のcsaを返す関数。
		args:
			string sfen: 以下、その例
				"+l+n+sgk1snl/1r4g2/+p1pppp1+Rp/6p2/1p7/2P6/P+b1PPPP1P/2G6/LNS1KGSNL b 2Pbp 1"
				"sfen +l+n+sgk1snl/1r4g2/+p1pppp1+Rp/6p2/1p7/2P6/P+b1PPPP1P/2G6/LNS1KGSNL b 2Pbp 1"
		return string: csa形式での盤面の文字列。以下、その例。
				"P1-GI *  * -KE-NY-NK-NG-KI
				P2-OU * -GI-KE-KY * -HI *  *
				P3 *  * -KI *  * -TO * -FU-FU
				P4-FU-FU * +RY-FU *  *  *  *
				P5 *  * -FU *  *  * -FU *  *
				P6 *  *  *  *  *  *  * +FU *
				P7 *  *  *  *  * +FU-UM * +FU
				P8+FU+FU+FU * +FU *  * +KI *
				P9 *  *  *  *  *  *  *  *  *
				P+
				P-00FU"
	*/
	// ここで、引数が無いことは想定しない。
	/*
	   sfen_strの例：
	       "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	   moves_strの例：
	       "7g7f 1a1b"
	*/
	sfen_str, moves_str := formatSfen(C.GoString(sfen_input_str))

	// 局面を示すposition_arr（スライス）この数行下に例がある。
	var position_arr []int

	// sfen_strをスラッシュや半角スペースで区切って配列sfen_arrに格納する。
	sfen_arr := strings.FieldsFunc(sfen_str, func(r rune) bool {
		return r == '/' || r == ' '
	})

	// 平手局面を配列に設定する。
	if sfen_str == "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1" {
		position_arr = append(position_arr,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -102, -103, -104, -105, -108, -105, -104, -103, -102, -1,
			-1, 0, -107, 0, 0, 0, 0, 0, -106, 0, -1,
			-1, -101, -101, -101, -101, -101, -101, -101, -101, -101, -1,
			-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1,
			-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1,
			-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1,
			-1, 101, 101, 101, 101, 101, 101, 101, 101, 101, -1,
			-1, 0, 106, 0, 0, 0, 0, 0, 107, 0, -1,
			-1, 102, 103, 104, 105, 108, 105, 104, 103, 102, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			1)
	} else {
		// 持ち駒を除いた盤面の配列変換
		position_arr = append(position_arr, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1)
		for i := 0; i < 9; i++ {
			//position_arr = append(position_arr, sfen_char2arr_convert(sfen_arr[i])...)
			// 以下の行はsfen_char2arr_convertを高速化したもの。
			position_arr = append(position_arr, sfen_char2arr_convert2(sfen_arr[i])...)
		}
		position_arr = append(position_arr, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1)
		// 持ち駒の盤面の配列変換
		// 2Pbp :先手は歩2枚、後手は角、歩
		position_arr = append(position_arr, sfen_hand2arr_convert(sfen_arr[10])...)
		// 手番の盤面の配列変換
		// b:先手1、w:後手-1
		if sfen_arr[9] == "b" {
			position_arr = append(position_arr, 1)
		} else {
			position_arr = append(position_arr, -1)
		}
	}

	// moves_strが"2g2f 3c3d 7g7f 4c5c"の場合、
	// moves_arrは[]string{"2g2f", "3c3d", "7g7f", "4c5c"}になるようにする。
	moves_arr := strings.Fields(moves_str)

	before_pos := ""      // 移動前の駒の位置
	after_pos := ""       // 移動前の駒の位置
	bePromoted := ""      // 成る場合は「+」
	tmp_before_piece := 0 // 一時的に記録しておく移動前の駒の種類
	tmp_after_piece := 0  // 一時的に記録しておく移動先にある駒の種類

	// 手番の確認。
	csfen_str := C.CString(sfen_str)
	defer C.free(unsafe.Pointer(csfen_str))
	turn_now := C.GoString(turn(csfen_str)) // "BLACK"（先手）か"WHITE"（後手）

	for _, move := range moves_arr {
		if len(move) == 4 {
			before_pos = move[:2]
			after_pos = move[2:]
		} else {
			before_pos = move[:2]
			after_pos = move[2:4]
			bePromoted = move[len(move)-1:]
		}

		// 駒の位置の記録の処理。
		posIndexMap := map[string]int{
			"1a": 20,
			"2a": 19,
			"3a": 18,
			"4a": 17,
			"5a": 16,
			"6a": 15,
			"7a": 14,
			"8a": 13,
			"9a": 12,
			"1b": 31,
			"2b": 30,
			"3b": 29,
			"4b": 28,
			"5b": 27,
			"6b": 26,
			"7b": 25,
			"8b": 24,
			"9b": 23,
			"1c": 42,
			"2c": 41,
			"3c": 40,
			"4c": 39,
			"5c": 38,
			"6c": 37,
			"7c": 36,
			"8c": 35,
			"9c": 34,
			"1d": 53,
			"2d": 52,
			"3d": 51,
			"4d": 50,
			"5d": 49,
			"6d": 48,
			"7d": 47,
			"8d": 46,
			"9d": 45,
			"1e": 64,
			"2e": 63,
			"3e": 62,
			"4e": 61,
			"5e": 60,
			"6e": 59,
			"7e": 58,
			"8e": 57,
			"9e": 56,
			"1f": 75,
			"2f": 74,
			"3f": 73,
			"4f": 72,
			"5f": 71,
			"6f": 70,
			"7f": 69,
			"8f": 68,
			"9f": 67,
			"1g": 86,
			"2g": 85,
			"3g": 84,
			"4g": 83,
			"5g": 82,
			"6g": 81,
			"7g": 80,
			"8g": 79,
			"9g": 78,
			"1h": 97,
			"2h": 96,
			"3h": 95,
			"4h": 94,
			"5h": 93,
			"6h": 92,
			"7h": 91,
			"8h": 90,
			"9h": 89,
			"1i": 108,
			"2i": 107,
			"3i": 106,
			"4i": 105,
			"5i": 104,
			"6i": 103,
			"7i": 102,
			"8i": 101,
			"9i": 100,
		}
		// 移動元の駒について
		// before_posが"1a"なら、posIndexに20が入る。
		// posIndexはposition_arrのindexと同じ値。
		beforePosIndex, ok := posIndexMap[before_pos]
		if !ok {
			// moves_strがmapに存在しない場合の処理
			// 例）"P*"
			// 先手、後手の区別がないので、区別するために先手を小文字にする処理を入れる。
			if turn_now == "BLACK" {
				// 先手の場合。小文字にする処理。
				before_pos = strings.ToLower(before_pos)
			}
			// 駒を駒台から打つ場合。
			beforeHandIndexMap := map[string]int{
				"p*": 121,
				"l*": 122,
				"n*": 123,
				"s*": 124,
				"g*": 125,
				"b*": 126,
				"r*": 127,
				"P*": 129,
				"L*": 130,
				"N*": 131,
				"S*": 132,
				"G*": 133,
				"B*": 134,
				"R*": 135,
			}
			// before_posが"p*"なら、beforeHandIndexに121が入る。
			// beforeHandIndexはposition_arrのindexと同じ値。
			beforeHandIndex, ok := beforeHandIndexMap[before_pos]
			if !ok {
				// 打つコマの種類がmapに存在しないのは想定外。
				fmt.Println("error. beforeHandIndex")
			}
			// 駒台から駒を減らす処理。
			position_arr[beforeHandIndex] += -1
			// 移動元の駒の種類を記録する処理。
			tmpBeforePieceMap := map[string]int{
				"p*": 101,
				"l*": 102,
				"n*": 103,
				"s*": 104,
				"g*": 105,
				"b*": 106,
				"r*": 107,
				"P*": -101,
				"L*": -102,
				"N*": -103,
				"S*": -104,
				"G*": -105,
				"B*": -106,
				"R*": -107,
			}
			// 移動元が持ち駒ではない場合。
			// before_posが"p*"なら、tmp_before_pieceに101が入る。
			tmp_before_piece, ok = tmpBeforePieceMap[before_pos]
			if !ok {
				// 打つコマの種類がmapに存在しないのは想定外。
				fmt.Println("error. tmpBeforePieceMap")
			}
		}
		// 移動元の自分の駒を記録する。
		if 0 < beforePosIndex && beforePosIndex < 121 {
			// 移動元が持ち駒ではない場合。
			tmp_before_piece = position_arr[beforePosIndex]
		}
		// 駒を動かすと、そのマスには何もないので、0にする。
		// beforePosIndex=0は壁なので変更しない。
		if beforePosIndex != 0 {
			position_arr[beforePosIndex] = 0
		}
		// 移動先の駒について
		afterPosIndex, ok := posIndexMap[after_pos]
		if !ok {
			// 移動先がmapに存在しないのは想定外
			fmt.Println("error. afterPosIndex")
		}
		// 移動先にある相手の駒を記録する。
		tmp_after_piece = position_arr[afterPosIndex]
		// 移動元の駒を移動先にコピーする。
		// 移動先に相手の駒がない場合。
		if tmp_after_piece == 0 {
			if bePromoted != "+" {
				// 成らない場合。
				position_arr[afterPosIndex] = tmp_before_piece
			} else {
				// 成る場合。
				if turn_now == "BLACK" {
					// 先手の場合。
					position_arr[afterPosIndex] = tmp_before_piece + 100
				} else {
					// 後手の場合。
					position_arr[afterPosIndex] = tmp_before_piece - 100
				}
				// 成るフラグを戻す。
				bePromoted = ""
			}
		} else {
			// 移動先に相手の駒がある場合、取る処理をする。
			// 先手の場合。
			if turn_now == "BLACK" {
				// 移動先に後手の駒があれば、先手の駒の種類に変える。
				// 成り駒じゃない場合。例）後手の歩：-101 -> 先手の歩：101
				if tmp_after_piece > -200 {
					tmp_after_piece *= -1
				} else {
					// 成り駒の場合。例）後手のと金：-201 -> 先手の歩：101
					tmp_after_piece *= -1
					tmp_after_piece += -100
				}
			} else {
				// 後手の場合。
				// 移動先に先手の駒があれば、後手の駒の種類に変える。
				// 成り駒じゃない場合。例）先手の歩：101 -> 後手の歩：-101
				if tmp_after_piece < 200 {
					tmp_after_piece *= -1
				} else {
					// 成り駒の場合。例）先手のと金：201 -> 後手の歩：-101
					tmp_after_piece *= -1
					tmp_after_piece += 100
				}
			}

			// 移動先にあった相手の駒の処理。
			afterHandIndexMap := map[int]int{
				101:  121,
				102:  122,
				103:  123,
				104:  124,
				105:  125,
				106:  126,
				107:  127,
				-101: 129,
				-102: 130,
				-103: 131,
				-104: 132,
				-105: 133,
				-106: 134,
				-107: 135,
			}
			// tmp_before_pieceが101なら、afterHandIndexに121が入る。
			// afterHandIndexはposition_arrのindexと同じ値。
			afterHandIndex, ok := afterHandIndexMap[tmp_after_piece]
			if !ok {
				// tmp_before_pieceがmapに存在しない場合は想定外。
				fmt.Println("error. afterHandIndex")
			}
			// 持ち駒を増やす処理。
			position_arr[afterHandIndex] += 1
			// 移動先の駒を変更する処理。
			if bePromoted != "+" {
				// 成らない場合。
				position_arr[afterPosIndex] = tmp_before_piece
			} else {
				// 成る場合。
				if turn_now == "BLACK" {
					// 先手の場合。
					position_arr[afterPosIndex] = tmp_before_piece + 100
				} else {
					// 後手の場合。
					position_arr[afterPosIndex] = tmp_before_piece - 100
				}
				// 成るフラグを戻す。
				bePromoted = ""
			}
		}
		// 次のmoveになるので、先手と後手を変える。
		if turn_now == "BLACK" {
			turn_now = "WHITE"
		} else {
			turn_now = "BLACK"
		}
	}

	// 盤面を表示する文字列
	var position_str string
	/*
		for i, num := range position_arr {
			//num_str := convertNumToStr(num)
			num_str := convertNumToStr2(num)
			if i <= 119 {
				if (i+1)%11 == 0 && i <= 108 {
					position_str += num_str + "\nP" + strconv.Itoa((i+1)/11)
				} else {
					position_str += num_str
				}
			} else {
				if (i+1)%11 == 0 && i <= 121 {
					position_str += "\nP+"
				} else if i == 121 && num > 0 {
					position_str += strings.Repeat("00FU", num)
				} else if i == 122 && num > 0 {
					position_str += strings.Repeat("00KY", num)
				} else if i == 123 && num > 0 {
					position_str += strings.Repeat("00KE", num)
				} else if i == 124 && num > 0 {
					position_str += strings.Repeat("00GI", num)
				} else if i == 125 && num > 0 {
					position_str += strings.Repeat("00KI", num)
				} else if i == 126 && num > 0 {
					position_str += strings.Repeat("00KA", num)
				} else if i == 127 && num > 0 {
					position_str += strings.Repeat("00HI", num)
				} else if i == 128 {
					position_str += "\nP-"
				} else if i == 129 && num > 0 {
					position_str += strings.Repeat("00FU", num)
				} else if i == 130 && num > 0 {
					position_str += strings.Repeat("00KY", num)
				} else if i == 131 && num > 0 {
					position_str += strings.Repeat("00KE", num)
				} else if i == 132 && num > 0 {
					position_str += strings.Repeat("00GI", num)
				} else if i == 133 && num > 0 {
					position_str += strings.Repeat("00KI", num)
				} else if i == 134 && num > 0 {
					position_str += strings.Repeat("00KA", num)
				} else if i == 135 && num > 0 {
					position_str += strings.Repeat("00HI", num)
				}
			}
		}
	*/

	// 上記のコメントアウトの部分をmapとswicth等を使って高速化した。
	var positionBuilder strings.Builder
	pieceMap := map[int]string{
		121: "00FU",
		122: "00KY",
		123: "00KE",
		124: "00GI",
		125: "00KI",
		126: "00KA",
		127: "00HI",
		129: "00FU",
		130: "00KY",
		131: "00KE",
		132: "00GI",
		133: "00KI",
		134: "00KA",
		135: "00HI",
	}
	for i, num := range position_arr {
		num_str := convertNumToStr2(num)
		switch {
		case i <= 119:
			if (i+1)%11 == 0 && i == 10 {
				positionBuilder.WriteString(num_str + "P" + strconv.Itoa((i+1)/11))
			} else if (i+1)%11 == 0 && i >= 12 && i <= 108 {
				positionBuilder.WriteString(num_str + "\nP" + strconv.Itoa((i+1)/11))
			} else {
				positionBuilder.WriteString(num_str)
			}
		case (i+1)%11 == 0 && i <= 121:
			positionBuilder.WriteString("\nP+")
		case i >= 121 && i <= 127 && num > 0:
			positionBuilder.WriteString(strings.Repeat(pieceMap[i], num))
		case i == 128:
			positionBuilder.WriteString("\nP-")
		case i >= 129 && i <= 135 && num > 0:
			positionBuilder.WriteString(strings.Repeat(pieceMap[i], num))
		default:
		}
	}
	position_str = positionBuilder.String()

	return C.CString(position_str)
}

/*
// 独自の駒の数値からCSA形式の文字列に変換する
func convertNumToStr(num int) string {
	var num_str string
	if num == -1 {
		num_str = ""
	} else if num == 0 {
		num_str = " * "
	} else {
		if num > 0 {
			num_str = "+"
		} else {
			num_str = "-"
			num = -num
		}
		if num == 101 {
			num_str += "FU"
		} else if num == 102 {
			num_str += "KY"
		} else if num == 103 {
			num_str += "KE"
		} else if num == 104 {
			num_str += "GI"
		} else if num == 105 {
			num_str += "KI"
		} else if num == 106 {
			num_str += "KA"
		} else if num == 107 {
			num_str += "HI"
		} else if num == 108 {
			num_str += "OU"
		} else if num == 201 {
			num_str += "TO"
		} else if num == 202 {
			num_str += "NY"
		} else if num == 203 {
			num_str += "NK"
		} else if num == 204 {
			num_str += "NG"
		} else if num == 206 {
			num_str += "UM"
		} else if num == 207 {
			num_str += "RY"
		} else {
			num_str += strconv.Itoa(num)
		}
	}
	return num_str
}
*/

// 独自の駒の数値からCSA形式の文字列に変換するconvertNumToStr関数を高速化したもの。
func convertNumToStr2(num int) string {
	numToStrMap := map[int]string{
		0:   " * ",
		101: "FU",
		102: "KY",
		103: "KE",
		104: "GI",
		105: "KI",
		106: "KA",
		107: "HI",
		108: "OU",
		201: "TO",
		202: "NY",
		203: "NK",
		204: "NG",
		206: "UM",
		207: "RY",
	}
	var num_str string
	if num > 0 {
		num_str = "+"
	} else if num < 0 {
		if num == -1 {
			return ""
		}
		num_str = "-"
		num = -num
	}
	if val, ok := numToStrMap[num]; ok {
		num_str += val
	} else if num != -1 {
		num_str += strconv.Itoa(num)
	}
	return num_str
}

/*
// 持ち駒を除いた盤面の配列変換
func sfen_char2arr_convert(input string) []int {
	output := []int{} // returnする配列

	isPromoted := false // 成り駒かの判定。

	// 入力文字列を "/" で分割
	splitInput := strings.Split(input, "/")
	for _, str := range splitInput {
		output = append(output, -1)
		for _, char := range str {
			if char == '+' {
				isPromoted = true
				continue
			}
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				var pieceValue int
				if char == 'p' || char == 'P' {
					pieceValue = 101
				} else if char == 'l' || char == 'L' {
					pieceValue = 102
				} else if char == 'n' || char == 'N' {
					pieceValue = 103
				} else if char == 's' || char == 'S' {
					pieceValue = 104
				} else if char == 'g' || char == 'G' {
					pieceValue = 105
				} else if char == 'b' || char == 'B' {
					pieceValue = 106
				} else if char == 'r' || char == 'R' {
					pieceValue = 107
				} else if char == 'k' || char == 'K' {
					pieceValue = 108
				} else {
					pieceValue = 0
				}
				if isPromoted {
					pieceValue += 100
					isPromoted = false
				}
				if char >= 'a' && char <= 'z' {
					pieceValue *= -1
				}
				output = append(output, pieceValue)
			} else if char >= '0' && char <= '9' {
				numZeroes, _ := strconv.Atoi(string(char))
				for i := 0; i < numZeroes; i++ {
					output = append(output, 0)
				}
			}
		}
		output = append(output, -1)
	}

	return output
}
*/

// 持ち駒を除いた盤面の配列変換
// mapを使って高速化したもの。
func sfen_char2arr_convert2(input string) []int {
	output := []int{} // returnする配列

	// mapを作成
	pieceMap := map[rune]int{
		'p': 101,
		'l': 102,
		'n': 103,
		's': 104,
		'g': 105,
		'b': 106,
		'r': 107,
		'k': 108,
	}

	isPromoted := false

	splitInput := strings.Split(input, "/")
	for _, str := range splitInput {
		output = append(output, -1)
		for _, char := range str {
			if char == '+' {
				isPromoted = true
				continue
			}
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				pieceValue, ok := pieceMap[unicode.ToLower(char)]
				if !ok {
					pieceValue = 0
				}
				if isPromoted {
					pieceValue += 100
					isPromoted = false
				}
				if char >= 'a' && char <= 'z' {
					pieceValue *= -1
				}
				output = append(output, pieceValue)
			} else if char >= '0' && char <= '9' {
				numZeroes, _ := strconv.Atoi(string(char))
				for i := 0; i < numZeroes; i++ {
					output = append(output, 0)
				}
			}
		}
		output = append(output, -1)
	}

	return output
}

// 持ち駒の盤面の配列変換
// 2Pb11p :先手は歩2枚、後手は角1枚、歩11枚
// 先手の持ち駒、インデックス121は歩、122 香
// 後手の持ち駒、インデックス129は歩、130 香
func sfen_hand2arr_convert(input string) []int {
	if input == "-" {
		if_hyphen_output := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		return if_hyphen_output
	}

	output := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} // returnするスライス

	hand_piece_kind := "PLNSGBRKplnsgbrk"

	cnt_str := ""

	// charはルーン型。シングルクォートで定義され、UTF-8エンコーディングのUnicodeコードポイントを返す。
	for _, char := range input {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			index := strings.Index(hand_piece_kind, string(char))
			if cnt_str == "" {
				output[index] = 1
			} else {
				num, err := strconv.Atoi(cnt_str)
				if err != nil {
					fmt.Println("error. hand convert：", err)
				} else {
					output[index] = num
				}
			}
			cnt_str = ""
		} else {
			cnt_str += string(char)
		}
	}

	return output
}

func formatSfen(sfen_str string) (string, string) {
	/*
	   局面の文字列をsfen形式でstartposもsfenにする。
	   args:
	        string: 局面の文字列
	            例）
	            "position startpos"
	            "position startpos moves 7g7f"
	            "position sfen lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	            "sfen lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1 moves 7g7f 1a1b"
	            "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	   return:
	       sfen_str(string): 「sfen」がないsfenの文字列。movesを含まない。
	            例）
	                "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	       moves_str(string): sfenのmovesより後ろの文字列。
	            例）
	                "7g7f 1a1b"
	*/
	// sfen_strに「startpos」があれば
	// 「sfen lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1」に置き換える。
	sfen_str = strings.Replace(sfen_str, "startpos", "sfen lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1", 1)

	// sfen_strに「sfen」があれば
	// 局面表示に「sfen」は不要なので、そこまで消す。
	if strings.Contains(sfen_str, "sfen") {
		index := strings.Index(sfen_str, "sfen")
		sfen_str = sfen_str[index+len("sfen"):]
	}

	/*
		sfen_strに「moves」が含まれている場合、
		「moves」より後を変数moves_strに入れ、
		sfen_strから「moves」以降を削除する。
		例）
		  "sfen +l+n+sgk1snl/1r4g2/+p1pppp1+Rp/6p2/1p7/2P6/P+b1PPPP1P/2G6/LNS1KGSNL b 2Pbp 1 moves 2g2f 3c3d 7g7f 4c5c"の場合
		  moves_strは"2g2f 3c3d 7g7f 4c5c"
		  sfen_strは"sfen +l+n+sgk1snl/1r4g2/+p1pppp1+Rp/6p2/1p7/2P6/P+b1PPPP1P/2G6/LNS1KGSNL b 2Pbp 1 "
	*/
	moves_str := ""
	if strings.Contains(sfen_str, "moves") {
		index := strings.Index(sfen_str, "moves")
		moves_str = strings.TrimSpace(sfen_str[index+len("moves"):])
		sfen_str = sfen_str[:index]
	}
	return sfen_str, moves_str
}

//export turn
func turn(sfen_input_str *C.char) *C.char {
	// sfenの文字列から手番を返す。
	/*
	   args:
	       string: 局面の文字列
	           例）
	           "position startpos"
	           "position startpos moves 7g7f"
	           "position sfen lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	           "sfen lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1 moves 7g7f 1a1b"
	           "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	   return:
	       string: "BLACK"（先手）か、"WHITE"（後手）
	*/

	/*
	   sfen_strの例：
	       "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	   moves_strの例：
	       "7g7f 1a1b"
	*/
	sfen_str, moves_str := formatSfen(C.GoString(sfen_input_str))

	// それぞれ、スラッシュやスペースで区切って配列に格納する。
	sfen_arr := strings.FieldsFunc(sfen_str, func(c rune) bool {
		return c == '/' || c == ' '
	})

	moves_arr := strings.FieldsFunc(moves_str, func(c rune) bool {
		return c == '/' || c == ' '
	})
	if sfen_arr[9] == "b" && len(moves_arr)%2 == 0 {
		return C.CString("BLACK")
	} else if sfen_arr[9] == "w" && len(moves_arr)%2 == 1 {
		return C.CString("BLACK")
	} else {
		return C.CString("WHITE")
	}

}

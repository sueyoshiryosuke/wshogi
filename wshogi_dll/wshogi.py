from ctypes import *
import re


wshogi_cpp_dll = cdll.LoadLibrary('./wshogi_dll/wshogi_cpp.dll')
wshogi_cpp_dll.legal_moves.restype = c_char_p
wshogi_cpp_dll.legal_moves.argtypes = [c_char_p]

wshogi_go_dll = cdll.LoadLibrary('./wshogi_dll/wshogi_go.dll')
wshogi_go_dll.sfen2csa.restype = c_char_p
wshogi_go_dll.sfen2csa.argtypes = [c_char_p]


def legal_moves(input_str="position startpos"):
	# 文字列sfenから合法手を生成する関数。usi形式のものをリストで返す。
    result = wshogi_cpp_dll.legal_moves(input_str.encode()).decode().split()

    return result


def wshogi_push(move_str, sfen_str="position startpos"):
    # 手を指す関数。sfenに1手追加する。
    if sfen_str=="":
        sfen_str="position startpos"
    rtn_str = sfen_str + " " + move_str

    return rtn_str


def wshogi_pop(sfen_str):
    # sfenから1手を戻す関数。1手分を消す。
    if "moves " in sfen_str:
        rtn_str = sfen_str.rsplit(' ',1)[0]
    else:
        rtn_str = sfen_str
        
    return rtn_str


def sfen2csa(sfen_str="position startpos"):
    # 文字列sfenから表示用のcsaを返す関数。盤面表示などにどうぞ。
    rtn_str = wshogi_go_dll.sfen2csa(sfen_str.encode()).decode()

    return rtn_str


def wshogi_turn(sfen_str="position startpos"):
    # sfenの文字列から手番を返す。"BLACK"（先手）か、"WHITE"（後手）
    cmd_lst = list(sfen_str.split(" "))  # スペース区切りのリスト。
    # 例1）position startpos
    if len(cmd_lst) == 2:
        rtn_str = "BLACK"
    # 例2）position startpos moves 7g7f 8d8e
    elif cmd_lst[1] == "startpos":
        if (len(cmd_lst)-2) % 2 == 1:
            rtn_str = "BLACK"
        else:
            rtn_str = "WHITE"

    # sfenで局面が送られてくるとき
    elif cmd_lst[1] == "sfen":  # 指定局面
        # 例1）position sfen lnsgkgsnl/1r5b1/p1ppppppp/1p7/9/7P1/PPPPPPP1P/1B5R1/LNSGKGSNL b - 1
        if len(cmd_lst) == 6:
            if cmd_lst[3] == "b":
                rtn_str = "BLACK"
            else:
                rtn_str = "WHITE"
        # 例2）position sfen lnsgkgsnl/1r5b1/p1ppppppp/1p7/9/7P1/PPPPPPP1P/1B5R1/LNSGKGSNL b - 1 moves 7g7f 8d8e
        else:
            if (len(cmd_lst)-6) % 2 == 1 and cmd_lst[3] == "b":
                rtn_str = "BLACK"
            elif (len(cmd_lst)-6) % 2 == 0 and cmd_lst[3] == "w":
                rtn_str = "BLACK"
            elif (len(cmd_lst)-6) % 2 == 1 and cmd_lst[3] == "w":
                rtn_str = "WHITE"
            elif (len(cmd_lst)-6) % 2 == 0 and cmd_lst[3] == "b":
                rtn_str = "WHITE"
            else:
                rtn_str = "BLACK"  # ここは通らない想定

    return rtn_str

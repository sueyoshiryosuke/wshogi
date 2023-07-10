from wshogi_dll.wshogi import *
import cshogi
import shogi
import time

# 以下、チェックしたい局面の「sfen_move_str」のみコメントを外してください。
# 平手局面チェック
#sfen_move_str = "position startpos"
# 平手局面チェック2
sfen_move_str = "position startpos moves 7g7f 3c3d 2g2f 4c4d"
# 打ち歩詰めチェック
#sfen_move_str = "position startpos moves 7g7f 3c3d 2g2f 4c4d 2f2e 2b3c 3i4h 8b4b 5i6h 5a6b 6h7h 3a3b 5g5f 7a7b 4i5h 3b4c 8h7g 6b7a 6g6f 4d4e 7h8h 4c5d 4h5g 6c6d 6i7h 7a8b 5h6g 4a5b 9i9h 9c9d 8h9i 9d9e 7i8h 1c1d 7h7i 1d1e 3g3f 4b2b 2i3g 2b4b 6g6h 4b4a 2h2f 3c4d 2f2g 4d3c 2e2d 2c2d 6f6e 3d3e 7g3c+ 2a3c B*2b 3e3f 2b3c+ 5d6e N*4d 5b4b 3c2b 9e9f 9g9f P*9g 9h9g B*4i 2g2d 3f3g+ P*6b 6a6b 4d3b+ 4b3b 2b3b 4a6a 3b4c 4i7f+ 1i1h P*9h 9i9h N*8f 9h9i"
# 連続王手の千日手チェック：「4c3c」で連続王手の千日手（反則）も返される。
#sfen_move_str = "position startpos moves 7g7f 3c3d 8h2b+ 1a1b 2b3c 5a5b 3c4c 5b5a 4c3c 5a5b 3c4c 5b5a 4c3c 5a5b 3c4c 5b5a 4c3c 5a5b"
# 指定局面チェック
#sfen_move_str = "position sfen l6nl/6gkp/pr1+P5/3b1gspP/3ppN1PL/PpP1S1p2/2KG5/1S2G4/LN5R1 w SN2Pb4p 1"
# 指定局面チェック2
#sfen_move_str = "position sfen l6nl/6gkp/pr1+P5/3b1gspP/3ppN1PL/PpP1S1p2/2KG5/1S2G4/LN5R1 w SN2Pb4p 1 moves 6d4b 6c5c"
# 詰み状態チェック
#sfen_move_str = "position sfen l6nl/6gkp/pr1+P5/3b1gspP/3ppN1PL/PpP1S1p2/2KG5/1S2G4/LN5R1 w SN2Pb4p 1 moves 6d4b 6c5c 4b5a 5c4b 3b4c 4b3b 2b2c S*4a 9a9b 7f7e 5a6b 2e2d"


cmd_lst = list(sfen_move_str.split(" "))  # スペース区切りのリスト。
if len(cmd_lst) > 2:
    sfen_str = cmd_lst[2] +" "+ cmd_lst[3] +" "+ cmd_lst[4] +" "+ cmd_lst[5]
if cmd_lst[0] != "position":
    print("sfen_move_strがpotision から始まっていません。終了します。")
    quit()

# cshogi
# 送られてきた局面まで局面をセットしなおす。
if cmd_lst[1] == "sfen":
    board = cshogi.Board(sfen_str)
    if len(cmd_lst) > 6:
        for i in range(7, len(cmd_lst)):
            board.push_usi(cmd_lst[i])  # 局面を進める
elif cmd_lst[1] == "startpos":
    board = cshogi.Board()
    if len(cmd_lst) > 2:
        for i in range(3, len(cmd_lst)):
            board.push_usi(cmd_lst[i])  # 局面を進める

move_lst = list(board.legal_moves)  # 合法手をリスト化する
print("cshogi")
print("move_lst:")
for move in move_lst:
    print(cshogi.move_to_usi(move), end=" ")
print("")
print("legal_moves: ", len(move_lst))

print("")

# wshogi
board = sfen_move_str
move_lst = legal_moves(board)  # 合法手をリスト化する
print("wshogi")
print("move_lst:")
for move in move_lst:
    print(move, end=" ")
print("")
print("legal_moves: ", len(move_lst))

print("")

# python-shogi
board = shogi.Board()
# 送られてきた局面まで局面をセットしなおす。
if cmd_lst[1] == "sfen":
    board = shogi.Board(sfen_str)
    if len(cmd_lst) > 6:
        for i in range(7, len(cmd_lst)):
            board.push(shogi.Move.from_usi(cmd_lst[i]))  # 局面を進める
elif cmd_lst[1] == "startpos":
    if len(cmd_lst) > 2:
        board = shogi.Board()
        for i in range(3, len(cmd_lst)):
            board.push(shogi.Move.from_usi(cmd_lst[i]))  # 局面を進める

move_lst = list(board.legal_moves)  # 合法手をリスト化する
print("python-shogi")
print("move_lst:")
for move in move_lst:
    print(move, end=" ")
print("")
print("legal_moves: ", len(move_lst))

# 以下、盤面表示。
print("")
print("盤面表示")
print(sfen2csa(sfen_move_str))

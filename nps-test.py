from wshogi_dll.wshogi import *
import cshogi
import shogi
import time


sfen_str = "position startpos moves 7g7f 3c3d 2g2f 4c4d 2f2e 2b3c 3i4h 8b4b 5i6h 5a6b 6h7h 3a3b 5g5f 7a7b 4i5h 3b4c 8h7g 6b7a 6g6f 4d4e 7h8h 4c5d 4h5g 6c6d 6i7h 7a8b 5h6g 4a5b 9i9h 9c9d 8h9i 9d9e 7i8h 1c1d 7h7i 1d1e 3g3f 4b2b 2i3g 2b4b 6g6h 4b4a 2h2f 3c4d 2f2g 4d3c 2e2d 2c2d 6f6e 3d3e 7g3c+ 2a3c B*2b 3e3f 2b3c+ 5d6e N*4d 5b4b 3c2b 9e9f 9g9f P*9g 9h9g B*4i 2g2d 3f3g+ P*6b 6a6b 4d3b+ 4b3b 2b3b 4a6a 3b4c 4i7f+ 1i1h P*9h 9i9h N*8f 9h9i"
next = "P*2g"

# python-shogiのみloopは1/100
loop = 10000

cnt_cshogi = 0
cnt_wshogi = 0
cnt_shogi = 0

# cshogi
print("cshogi、計測中")
start_time = time.perf_counter()

for i in range(loop):
    if i==loop/4:
        start_time = time.perf_counter()
    # 送られてきた局面まで局面をセットしなおす。
    board = cshogi.Board()
    cmd_lst = list(sfen_str.split(" "))  # スペース区切りのリスト。
    if len(cmd_lst) > 2:
        for i in range(3, len(cmd_lst)):
            board.push_usi(cmd_lst[i])  # 指す
    move_lst = list(board.legal_moves)  # 合法手をリスト化する
    cnt_cshogi += len(move_lst)
    
    board.push(move_lst[1])  # 指す
    board.pop()  # 戻す
    
    board.push_usi(next)  # 指す
    move_lst2 = list(board.legal_moves)  # 合法手をリスト化する
    cnt_cshogi += len(move_lst2)

end_time = time.perf_counter()
print("cshogi、計測終了")
print("Processing time:", end_time - start_time)
nps_cshogi = int(cnt_cshogi/(end_time - start_time))
print("nodes:", cnt_cshogi)
print("nps:", nps_cshogi)

print("")


# wshogi
print("wshogi、計測中")
start_time = time.perf_counter()

for j in range(loop):
    if j==loop/4:
        start_time = time.perf_counter()
    # 送られてきた局面まで局面をセットしなおす。
    board = sfen_str
    move_lst = legal_moves(board)  # 合法手をリスト化する
    cnt_wshogi += len(move_lst)
    
    board = wshogi_push(move_lst[1], board) #指す
    board = wshogi_pop(board)  # 戻す

    board = wshogi_push(next, board)  # 指す
    move_lst2 = legal_moves(board)  # 合法手をリスト化する
    cnt_wshogi += len(move_lst2)

end_time = time.perf_counter()
print("wshogi、計測終了")
print("Processing time:", end_time - start_time)
nps_wshogi = int(cnt_wshogi/(end_time - start_time))
print("nodes:", cnt_wshogi)
print("nps:", nps_wshogi)

print("")


# python-shogi
print("python-shogi、計測中（loopは他の1/100）")
start_time = time.perf_counter()

for k in range(loop//100):
    if k==loop//4:
        start_time = time.perf_counter()
    # 送られてきた局面まで局面をセットしなおす。
    board = shogi.Board()
    cmd_lst = list(sfen_str.split(" "))  # スペース区切りのリスト。
    if len(cmd_lst) > 2:
        for i in range(3, len(cmd_lst)):
            board.push(shogi.Move.from_usi(cmd_lst[i]))  # 指す
    move_lst = list(board.legal_moves)  # 合法手をリスト化する
    cnt_shogi += len(move_lst)
    
    board.push(move_lst[1])  # 指す
    board.pop()  # 戻す
    
    board.push(shogi.Move.from_usi(next))  # 指す
    move_lst2 = list(board.legal_moves)  # 合法手をリスト化する
    cnt_shogi += len(move_lst2)
end_time = time.perf_counter()
print("python-shogi、計測終了")
print("Processing time:", end_time - start_time)
nps_shogi = int(cnt_shogi/(end_time - start_time))
print("nodes:", cnt_shogi)
print("nps:", nps_shogi)


print("")
print("処理速度の速い順：cshogi > wshogi > python-shogi")
print("")
print("wshogiの速さ：")
print("cshogiより", int(nps_cshogi/nps_wshogi), "倍ほど遅い。")
print("python-shogiより", int(nps_wshogi/nps_shogi), "倍ほど速い。")

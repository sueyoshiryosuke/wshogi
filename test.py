# wshogiで使える関数の説明とその動作テスト。

# 以下の行が必要。
from wshogi_dll.wshogi import *

sfen_str = "position startpos moves 2g2f 3c3d 7g7f"
sfen_str2 = "position sfen lr5+Sl/3kg2+R1/p1nsp4/2pp1p2p/6P2/1PP6/PGS+s1+p2P/1K1B5/LN1+n4L b BGN6Pg 1 moves G*4c 6g6h 4c5b 6c5b P*2d 7c6e 2d2c+"
sfen_str_mate = "position startpos moves 1g1f 4a3b 4g4f 5a4b 5i6h 9a9b 5g5f 3c3d 2h5h 5c5d 5f5e 2a3c 9g9f 7a6b 6g6f 6a5b 9f9e 8c8d 4i5i 2c2d 5e5d 6b5c 8h9g 5b6b 1i1h 4b4a 1h1g 6b6a 9g7e 3d3e 5h5e 5c5d 9i9g 8b6b 3i4h P*5h 6h5g 1c1d 6i6h 5d6e P*5f 6c6d 1f1e 6e5d 4h3i 7c7d 3g3f 5h5i+ 5e5d G*7h 7e8d 6a7a 8d6b+ 4a4b 5d5a+"


"""
関数legal_moves()
  文字列sfenから合法手を生成する関数。詰みの場合は、空文字""が返ってくる。
  
  引数：文字列sfen
    例："position startpos moves 2g2f 3c3d 7g7f"
  戻り値：合法手のリスト
    例：['9a9b', '7a6b', '7a7b', '6a5b', '6a6b', '6a7b', '5a4b', '5a5b', '5a6b', '4a3b', '4a4b', '4a5b', '3a3b', '3a4b', '2a3c', '1a1b', '8b7b', '8b6b', '8b5b', '8b4b', '8b3b', '8b9b', '2b3c', '2b4d', '2b5e', '2b6f', '2b7g+', '2b7g', '2b8h+', '2b8h', '9c9d', '8c8d', '7c7d', '6c6d', '5c5d', '4c4d', '2c2d', '1c1d', '3d3e']
"""
result = legal_moves(sfen_str)
print("len(result):", len(result))
print("result:", result)
print("")


"""
関数sfen2csa()
  文字列sfenから表示用のcsaを返す関数。盤面表示などにどうぞ。
  
  引数：文字列sfen
    例："position sfen lr5+Sl/3kg2+R1/p1nsp4/2pp1p2p/6P2/1PP6/PGS+s1+p2P/1K1B5/LN1+n4L b BGN6Pg 1 moves G*4c 6g6h 4c5b 6c5b P*2d 7c6e 2d2c+"
  戻り値：盤面csa形式の文字列
"""
result1 = sfen2csa(sfen_str2)
print("result1:")
print(sfen_str2)
print(result1)
print("")


"""
関数wshogi_push()
  手を指す関数。sfenに1手追加する。
  
  引数1：文字列usi
    例："1a1b"
  引数2：文字列sfen
    例："position startpos moves 2g2f 3c3d 7g7f"
  戻り値：文字列sfen
    例："position startpos moves 2g2f 3c3d 7g7f 1a1b"
"""
result2 = wshogi_push("1a1b", sfen_str)
print("result2:", result2)
print("")


"""
関数wshogi_pop()
  sfenから1手を戻す関数。1手分を消す。
  
  引数：文字列sfen
    例："position startpos moves 2g2f 3c3d 7g7f"
  戻り値：文字列sfen
    例："position startpos moves 2g2f 3c3d"
"""
result3 = wshogi_pop(sfen_str)
print("result3:", result3)
print("")


"""
関数wshogi_turn()
  sfenの文字列から手番を返す。
  
  引数：文字列sfen
    例："position sfen lr5+Sl/3kg2+R1/p1nsp4/2pp1p2p/6P2/1PP6/PGS+s1+p2P/1K1B5/LN1+n4L b BGN6Pg 1 moves G*4c 6g6h 4c5b 6c5b P*2d 7c6e 2d2c+"
  戻り値：文字列。"BLACK"（先手）か、"WHITE"（後手）
"""
result4 = wshogi_turn(sfen_str2)
print("result4:")
print(sfen_str2)
print(result4)
print("")

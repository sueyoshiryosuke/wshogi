# wshogi
Windowの将棋プログラム用dllライブラリ

## ソフトについて  
【ソフト名】　　　wshogi  
【バージョン】　　Ver.20230710  
【著作権者】　　　末吉 竜介  
【種　別】　　　　フリーソフトウェア  
【ソースコードのライセンス】　　MIT Licence  
【連絡先】　　　　[末吉のTwitter](https://twitter.com/16shiki168)  
【配布元ページ】　https://github.com/sueyoshiryosuke/wshogi  
【動作確認環境】　Windows11  
　　　　　　　　　gcc version 10.3.0 (tdm64-1)  
　　　　　　　　　go version go1.20.4 windows/amd64  
　　　　　　　　　Python 3.10.4  
【使用ソフト】  
　TDM-GCC　　　　 https://jmeubank.github.io/tdm-gcc/  
　Go言語　　　　　https://go.dev/  
　WinPython　　　https://winpython.github.io/  
【参考ソフト】  
　shogi686micro 　https://github.com/merom686/shogi686micro  
　cshogi　　　　　https://github.com/TadaoYamaoka/cshogi  
　python-shogi　　https://github.com/gunyarakun/python-shogi  
  
―――――――――――――――――――――――――――――――――――――  
## 著作権および免責事項  
  
　本ソフトはフリーソフトです。個人／団体／社内利用を問わず、ご自由にお使い  
下さい。  
　なお，著作権は上の【著作権者】に記載している者が保有しています。  
  
　このソフトウェアを使用したことによって生じたすべての障害・損害・不具合等に  
関しては、著作権者と著作権者の関係者および著作権者の所属するいかなる  
団体・組織とも、一切の責任を負いません。各自の責任においてご使用ください。  
  
## はじめに  
　　このソフトは、Windowの将棋プログラム開発用dllライブラリです。  
　Pythonから呼び出すことを想定していますが、C言語の呼び出し規則に  
　従っているので、他のプログラム言語からも呼び出し可能なdllです。  
　合法手生成をメインに、必要最低限レベルの関数を用意しました。  
　使い方はcshogiやpython-shpgiに似せましたが異なります。  
　ちなみに、将棋の合法手（打ち歩詰め含まず、連続王手の千日手を含む）の  
　生成速度は、python-shpgiよりも速く、cshogiよりも遅いです。  
  
## ファイル構成  
　test.py  
　　→wshogiで使える関数の説明が書いてある動作テストのソースコードです。  
　test.bat  
　　→test.pyを実行します。
　test_batの結果.txt  
　　→test.pyの実行結果です。  
　sfen_check.py  
　　→cshogiやpython-shogiで生成された合法手と一致するか比較する  
　　　ソースコードです。  
　　　実行にはcshogiやpython-shogiが必要です。  
　sfen_check.bat  
　　→sfen_check.pyを実行します。  
　nps-test.py  
　　→合法手の生成速度について、cshogiやpython-shogiと比較ソースコードです。  
　nps-test.bat  
　　→nps-test.pyを実行します。  
　README.md  
　　→この説明ファイルです。  
　LICENSE  
　　→このソフトのライセンスの内容が書かれています。  
　srcフォルダ  
　　→C++とGo言語のソースコードやコンパイル方法を書いたテキストファイルが  
　　　あります。  
  
## ダウンロード方法  
　[Releases](https://github.com/sueyoshiryosuke/wshogi/releases/)から  
　バイナリも含めた一式をダウンロードできます。  
  
## 使用方法  
　test.pyを参考に、お察しください。  
  
## 使用例  
```python
from wshogi_dll.wshogi import *

# sfen棋譜「position startpos moves 2g2f 3c3d 7g7f」の
# 合法手をリストresultに格納します。
sfen_str = "position startpos moves 2g2f 3c3d 7g7f"
result = legal_moves(sfen_str)
print(result)
print(result[0])
```
  
## ライセンス  
　MIT licenseです。  
  
## 謝辞  
　　merom686さん、ソースコードを流用させてもらいました。  
　　高速な合法手生成、大変ありがたいです。  
　　山岡忠夫さん、gunyarakunさん、どういう関数があればよいかや  
　　合法手生成の速度など、大変参考にさせていただきました。  
　　皆様に感謝、感謝、感謝です。  
　　そして、本ライブラリが将棋AI開発の手助けに少しでもなれば幸いです。  
  
## 更新履歴  
　Ver.20230710　　2023/07/10  
　　初版公開。  
  
--以上--  

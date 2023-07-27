/**
 * @file wshogi_cpp.cpp
 * @brief 合法手の生成部分だけdll化したもの
 * @author SUEYOSHI Ryosuke
 * @date 2023-07-10
 * 
 * 将棋エンジン「shogi686micro」のソースコードを元に
 * 合法手の生成部分だけ抜粋してdll化したもの。
 * 
 * merom686/shogi686micro: ソースファイル1個で将棋の思考エンジン
 * https://github.com/merom686/shogi686micro
 * 
 * dll生成のコンパイルコマンドは以下。
 * --
 * g++ -O2 -mtune=generic -s -shared -o wshogi_cpp.dll wshogi_cpp.cpp
 * --
 */

#include <iostream>
#include <regex>
#include <chrono>
#include <cstring>

//#define _assert(x) ((void)0)
#define _assert(x) \
if (!(x)){ std::cout << "info string error file:" << __FILE__ << " line:" << __LINE__ << std::endl; throw; }

enum Piece_Turn {
	King, Rook, Bishop, Gold, Silver, Knight, Lance, Pawn,
	HandTypeNum, Dragon, Horse, ProSilver = 12, ProKnight, ProLance, ProPawn,
	PieceTypeNum, BlackMask = PieceTypeNum, WhiteMask = BlackMask << 1,
};

enum Color {
	Black, White, ColorNum
};

const int MaxMove = 593, MaxPly = 32, MaxGamePly = 1024;
const std::string sfenPiece = "KRBGSNLPkrbgsnlp";

const short ScoreInfinite = INT16_MAX;
const short ScoreMatedInMaxPly = ScoreInfinite - MaxPly;

volatile bool stop;
uint64_t nodes;

struct Stack;

//指し手
class Move
{
	int value;

public:
	static const int MoveNone = 0;

	int from() const {
		return value & 0xff;
	}
	int to() const {
		return value >> 8 & 0xff;
	}
	int piece() const {
		return value >> 16 & 0xf;
	}
	int promote() const {
		return value >> 20 & 0x1;
	}
	int cap() const {
		return value >> 21 & 0xf;
	}
	//移動後の駒
	int pieceTo() const {
		return piece() | promote() << 3;
	}
	bool isNone() const {
		return (value == MoveNone);
	}
	//USI形式に変換
	std::string toUSI() const;
	//移動元(駒台のときは0)8bit, 移動先8bit, 移動「前」の駒4bit, 成ったか1bit, 取った駒4bit
	Move(int from, int to, int piece, int promote, int cap){
		value = from | to << 8 | piece << 16 | promote << 20 | cap << 21;
	}
	Move(int v) : value(v){}
	Move(){}
};

//局面
struct Position
{
	static const int FileNum = 9, RankNum = 9, PromotionRank = 3;
	static const int Stride = FileNum + 1, Origin = Stride * 3, SquareNum = Origin + Stride * (RankNum + 2);

	//piece_turn: 駒の種類3bit, 成1bit, 先手の駒1bit, 後手の駒1bit 以上6bit;壁は全8bitが立っている
	uint8_t piece_turn[SquareNum], hand[ColorNum][HandTypeNum], turn;
	uint8_t king[ColorNum];//玉の位置
	uint8_t continuousCheck[ColorNum];//現在の連続王手回数

	//探索を始めた時刻と終了する時刻
	std::chrono::system_clock::time_point timeStart, timeEnd;
	int ply, gamePly;//Rootからの手数、開始局面からの手数
	Stack *ss;//RootのStack位置

	//局面を比較する 同一なら0を返す
	static int compare(const Position &p1, const Position &p2){
		bool bp = std::equal(p1.piece_turn + Origin, p1.piece_turn + Origin + Stride * RankNum, p2.piece_turn + Origin);
		bool bh = std::equal(p1.hand[Black], p1.hand[White], p2.hand[Black]);//先手の持ち駒だけ比較すればよい
		bool bt = (p1.turn == p2.turn);
		return (bp && bh && bt) ? 0 : 1;
	}
	//square(0, 0)は盤の左上隅を表す(右上隅ではない)
	static int square(int x, int y){
		return Origin + Stride * y + x;
	}
	static int turnMask(int turn){
		return (turn == Black) ? BlackMask : WhiteMask;
	}

	int turnMask() const {
		return turnMask(turn);
	}
	//升が手番にとって敵陣rank段目までにあるか
	template <int rank = PromotionRank>
	bool promotionZone(int sq) const {
		if (turn == Black){
			return sq < square(0, rank);
		} else {
			return sq >= square(0, RankNum - rank);
		}
	}
	//turnの玉に王手がかかっているか
	bool inCheck(const int turn) const;
	//手を進める
	void doMove(Stack *const ss, const Move move);
	//王手放置(自殺手)でないか確かめる
	bool isLegal(Stack *const ss, const Move pseudoLegalMove);
	void clear(){
		std::memset(this, 0, sizeof *this);
		std::fill_n(piece_turn, SquareNum, 0xff);//壁で埋める
		for (int y = 0; y < RankNum; y++){
			std::fill_n(&piece_turn[square(0, y)], FileNum, 0);//y段目を全部空き升に
		}
	}
};

struct Stack
{
	Move pv[MaxPly];//読み筋を記録する
	Move currentMove;//いま読んでいる手
	bool checked;//手番の玉に王手がかかっているか
	Position pos;//局面を保存して、千日手検出や手を戻すのに使う
};

//指定された駒の利きがある全ての升に対して、trueを返すまでfを実行する
template <class F>
inline bool forAttack(const uint8_t *pt, const int sq, const int p, const int turn, F f)
{
	static const int8_t n = Position::Stride;
	static const int8_t att[PieceTypeNum][10] = {
		{ -n - 1, -n, -n + 1, -1, 1, n - 1, n, n + 1, 0, 0 },//玉
		{ 0, -n, -1, 1, n, 0, 0, 0, 0, 0 },
		{ 0, -n - 1, -n + 1, n - 1, n + 1, 0, 0, 0, 0, 0 },
		{ -n - 1, -n, -n + 1, -1, 1, n, 0, 0, 0, 0 },
		{ -n - 1, -n, -n + 1, n - 1, n + 1, 0, 0, 0, 0, 0 },
		{ -n * 2 + 1, -n * 2 - 1, 0, 0, 0, 0, 0, 0, 0, 0 },
		{ 0, -n, 0, 0, 0, 0, 0, 0, 0, 0 },
		{ -n, 0, 0, 0, 0, 0, 0, 0, 0, 0 },//歩
		{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 },
		{ -n - 1, -n + 1, n - 1, n + 1, 0, -n, -1, 1, n, 0 },//竜
		{ -n, -1, 1, n, 0, -n - 1, -n + 1, n - 1, n + 1, 0 },
		{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 },
		{ -n - 1, -n, -n + 1, -1, 1, n, 0, 0, 0, 0 },
		{ -n - 1, -n, -n + 1, -1, 1, n, 0, 0, 0, 0 },
		{ -n - 1, -n, -n + 1, -1, 1, n, 0, 0, 0, 0 },
		{ -n - 1, -n, -n + 1, -1, 1, n, 0, 0, 0, 0 },
	};

	const int sgn = (turn == Black) ? 1 : -1;
	const int8_t *a = att[p];
	int i;
	for (i = 0; a[i] != 0; i++){
		if (f(sq + a[i] * sgn)) return true;
	}
	for (i++; a[i] != 0; i++){
		for (int d = a[i];; d += a[i]){
			if (f(sq + d * sgn)) return true;
			if (pt[sq + d * sgn] != 0) break;
		}
	}
	return false;
}

std::string Move::toUSI() const
{
	std::string s;
	auto add = [&](int sq){
		sq -= Position::Origin;
		s += '1' + Position::FileNum - 1 - sq % Position::Stride;
		s += 'a' + sq / Position::Stride;
	};

	if (from() == 0){
		s += sfenPiece[piece()];
		s += '*';
		add(to());
	} else {
		add(from());
		add(to());
		if (promote()) s += '+';
	}
	return s;
}

inline bool Position::inCheck(const int turn) const
{
	for (int p = King; p < PieceTypeNum; p++){
		const int pt = p | Position::turnMask(turn ^ 1);
		bool ret = forAttack(piece_turn, king[turn], p, turn, [&](int sq){
			return piece_turn[sq] == pt;
		});
		if (ret) return true;
	}
	return false;
}

inline void Position::doMove(Stack *const ss, const Move move)
{
	if (move.from() == 0){
		//打つ
		hand[turn][move.piece()]--;
		piece_turn[move.to()] = move.piece() | turnMask();
	} else {
		//移動
		if (move.cap()){
			//取る
			hand[turn][move.cap() % HandTypeNum]++;
		}
		piece_turn[move.from()] = 0;
		piece_turn[move.to()] = move.pieceTo() | turnMask();
		if (move.piece() == King) king[turn] = move.to();
	}
	turn ^= 1;
	ply++;
	gamePly++;

	//いま指した手
	ss->currentMove = move;
	//いま指した手が王手だったか
	(ss + 1)->checked = inCheck(turn);
	//連続王手の回数を更新
	if ((ss + 1)->checked){
		continuousCheck[ss->pos.turn]++;
	} else {
		continuousCheck[ss->pos.turn] = 0;
	}
}

inline bool Position::isLegal(Stack *const ss, const Move pseudoLegalMove)
{
	doMove(ss, pseudoLegalMove);
	bool illegal = inCheck(turn ^ 1);
	*this = ss->pos;//手を戻す
	return !illegal;
}

//全ての合法手(王手放置を含む)を生成し、生成した指し手の個数を返す
int generateMoves(Move *const moves, const Position &pos)
{
	Move *m = moves;
	int pawn = 0;//二歩検出用のビットマップ
	const int t = pos.turnMask();
	//移動
	for (int y = 0; y < Position::RankNum; y++){
		for (int x = 0; x < Position::FileNum; x++){
			int from = Position::square(x, y);
			int pt = pos.piece_turn[from];
			if (pt & t){
				int p = pt % PieceTypeNum;
				if (p == Pawn) pawn |= 1 << x;
				forAttack(pos.piece_turn, from, p, pos.turn, [&](int to){
					int cap = pos.piece_turn[to];
					if (!(cap & t)){//自分の駒と壁以外(空升と相手の駒)へなら移動できる
						if (p < HandTypeNum && p != King && p != Gold
							&& (pos.promotionZone(from) || pos.promotionZone(to))){
							*m++ = Move{ from, to, p, 1, cap % PieceTypeNum };
						}
						if (!((p == Pawn || p == Lance) && pos.promotionZone<1>(to))
							&& !(p == Knight && pos.promotionZone<2>(to))){
							*m++ = Move{ from, to, p, 0, cap % PieceTypeNum };
						}
					}
					return false;
				});
			}
		}
	}
	//打つ
	for (int p = Rook; p < HandTypeNum; p++){
		if (!pos.hand[pos.turn][p]) continue;
		for (int y = 0; y < Position::RankNum; y++){
			for (int x = 0; x < Position::FileNum; x++){
				int to = Position::square(x, y);
				int pt = pos.piece_turn[to];
				if (pt == 0 && !(p == Pawn && (pawn & 1 << x))){
					if (!((p == Pawn || p == Lance) && pos.promotionZone<1>(to))
						&& !(p == Knight && pos.promotionZone<2>(to))){
						*m++ = Move{ 0, to, p, 0, 0 };
					}
				}
			}
		}
	}
	return (int)(m - moves);
}


//SFENの局面をposとssにセットする
void setPosition(Position &pos, Stack *ss, std::istringstream &iss)
{
	//startposの処理めんどい
	//変更した。
	std::string input = iss.str().substr((size_t)iss.tellg() + 1);
	if (input.find("startpos") == 0) {
		input.replace(0, 8, "sfen lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1");
	}
	iss.str(input);

	std::string sfenPos, sfenTurn, sfenHand, sfenCount, sfenMove;

	//局面初期化
	pos.clear();

	//盤面
	iss >> sfenPos;//"sfen"
	iss >> sfenPos;
	int x = 0, y = 0, pro = 0;
	for (auto c : sfenPos){
		if ('0' < c && c <= '9'){
			x += c - '0';
		} else if (c == '+'){
			pro = 1;
		} else if (c == '/'){
			x = 0;
			y++;
		} else {
			auto i = sfenPiece.find(c);
			_assert(i != std::string::npos && i < PieceTypeNum);
			int turn = (int)i / HandTypeNum;
			int p = i % HandTypeNum | pro << 3;
			int sq = Position::square(x, y);
			pos.piece_turn[sq] = p | Position::turnMask(turn);
			pro = 0;
			x++;
			if (p == King) pos.king[turn] = sq;//玉の位置はここでセットする
		}
	}

	//手番
	iss >> sfenTurn;
	pos.turn = (sfenTurn == "b") ? Black : White;

	//持ち駒
	iss >> sfenHand;
	int n = 0;
	for (auto c : sfenHand){
		if (c == '-'){
			break;
		} else if ('0' <= c && c <= '9'){
			n = n * 10 + (c - '0');
		} else {
			auto i = sfenPiece.find(c);
			_assert(i != std::string::npos && i < PieceTypeNum);
			pos.hand[i / HandTypeNum][i % HandTypeNum] = (n == 0) ? 1 : n;
			n = 0;
		}
	}

	//手数(使わない)
	iss >> sfenCount;

	//Stackはここで
	ss->checked = pos.inCheck(pos.turn);
	ss->pos = pos;
	pos.ss = ss;

	//指し手
	iss >> sfenMove;
	if (sfenMove != "moves") return;

	while (iss >> sfenMove){
		//全ての合法手を生成して一致するものを探す
		Move moves[MaxMove];
		int n = generateMoves(moves, pos);

		auto it = std::find_if(moves, moves + n, [&](Move move){
			return sfenMove == move.toUSI();
		});
		_assert(it < moves + n);
		pos.doMove(pos.ss++, *it);
		pos.ss->pos = pos;
	}
}


/**
 * @fn int legalMoves2(const std::string)
 * 将棋のsfen形式の文字列を読み込むと、合法手の数を返す。
 * 王手放置は含まず、打ち歩詰めの手を含む。
 * @param cmd sfen形式の文字列。デフォルトは平手局面。
 * @return int 合法手の数。
 * @detail
 * 引数の例：
 *   "position startpos moves 7g7f"
 * 戻り値の例：
 *   30
 */
int legalMoves2(const std::string cmd = std::string("position startpos")) {
    Position pos;
    std::vector<Stack> vss{ MaxGamePly + 2 };
    Stack *const ss = &vss[0];

    std::string token;

    std::istringstream iss(cmd);
    iss >> token;

	std::memset(ss, 0, vss.size() * sizeof *ss);
	setPosition(pos, ss + 1, iss);

	Move moves[MaxMove];  // 合法手を格納するための配列
	int n = generateMoves(moves, pos);  //すべての合法手(王手放置を含む)の個数
	ss->pos = pos;  //現在の局面の保存。

	int trn_num = 0;  // 戻り値となる合法手の数
	for (int i = 0; i < n; i++){
		Move move = moves[i];
		if (!pos.isLegal(ss, move)) continue;  //王手放置を除く
		
		// 合法手の数をカウントする。
		trn_num++;
	}
	
	return trn_num;
}


/**
 * @fn char* legal_moves(const char*)
 * 将棋のsfen形式の文字列を読み込むと、合法手をUSI形式で返す。
 * 王手放置も、打ち歩詰めの手も含まない。
 * @param cmd sfen形式の文字列。デフォルトは平手局面。
 * @return char* 引数の局面の合法手をUSI形式で列挙する。
 * @detail 戻り値は見つからなければ「""」を返す。
 * 引数の例：
 *   "position startpos moves 7g7f"
 * 戻り値の例：
 *   "9a9b 7a6b 7a7b（以下略）"
 */
extern "C" char* legal_moves(const char* cmd = "position startpos") {
    const std::string cmd_str(cmd);
    Position pos;
    std::vector<Stack> vss{ MaxGamePly + 2 };
    Stack *const ss = &vss[0];

    std::string token, rtn_str;
	rtn_str.reserve(MaxMove * 6);  // 合法手1つは5文字+スペース1文字。
	std::string checkSfen;  //打ち歩詰めのチェック用のsfen文字列。

    std::istringstream iss(cmd_str);
    iss >> token;

	std::memset(ss, 0, vss.size() * sizeof *ss);
	setPosition(pos, ss + 1, iss);

	Move moves[MaxMove];  // 合法手を格納するための配列
	int n = generateMoves(moves, pos);  //すべての合法手(王手放置を含む)の個数
	ss->pos = pos;  //現在の局面の保存。

	for (int i = 0; i < n; i++){
		Move move = moves[i];
		if (!pos.isLegal(ss, move)) continue;  //王手放置を除く

		//手を進める
		pos.doMove(ss, move);
		// 歩を打つ手で、相手に王手がかかるか。
		if (move.from() == 0 && move.piece() == Pawn && pos.inCheck(pos.turn)){
			//打ち歩詰めのチェック。
			checkSfen = cmd_str + " " + move.toUSI();
			//相手番でのすべての合法手(王手放置を含まない)の候補手がない=打ち歩詰め。
			if (legalMoves2(checkSfen) == 0) {
				//手を戻す
				pos = ss->pos;
				continue;  //打ち歩詰めを除く
			} 
		}
		//手を戻す
		pos = ss->pos;
		
		// 候補主として記録する。
		rtn_str += moves[i].toUSI() + " ";
	}
	// C言語で扱えるようにしておく。
	return strdup(rtn_str.c_str());
}


int main() {
	/*
	std::string moves_str = legal_moves("position startpos moves 7g7f 3c3d 2g2f 4c4d 2f2e 2b3c 3i4h 8b4b 5i6h 5a6b 6h7h 3a3b 5g5f 7a7b 4i5h 3b4c 8h7g 6b7a 6g6f 4d4e 7h8h 4c5d 4h5g 6c6d 6i7h 7a8b 5h6g 4a5b 9i9h 9c9d 8h9i 9d9e 7i8h 1c1d 7h7i 1d1e 3g3f 4b2b 2i3g 2b4b 6g6h 4b4a 2h2f 3c4d 2f2g 4d3c 2e2d 2c2d 6f6e 3d3e 7g3c+ 2a3c B*2b 3e3f 2b3c+ 5d6e N*4d 5b4b 3c2b 9e9f 9g9f P*9g 9h9g B*4i 2g2d 3f3g+ P*6b 6a6b 4d3b+ 4b3b 2b3b 4a6a 3b4c 4i7f+ 1i1h P*9h 9i9h N*8f 9h9i");
    std::cout << "moves_str" << std::endl;
    std::cout << moves_str << std::endl;
	*/

    return 0;
}

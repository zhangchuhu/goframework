namespace go rank

struct TRank {
  1: i64 uid;
  2: i64 value;
  3: i64 rank;
}


service TRankService {
  list<TRank> queryRank(1: string code, 2: string rtype, 3: string ctype, 4: bool latest, 5: i32 size, 6: i32 timeParam);
}
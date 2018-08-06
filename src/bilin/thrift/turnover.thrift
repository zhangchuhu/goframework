namespace java.swift com.yy.turnover
namespace cpp com.yy.turnover
namespace py turnover


enum TUserType {
  Anchor=1, OW=2, Visitor=3, Guest=4, Sender=5, AnchorAndOw=6
}

enum UsedChannelType {
  Client=0, Web=10000, IOS=10001, Android=10002, IOS_Cracked=10003, WechatOfficialAccount=10004, TurnoverWeb=1, DatingCom=2, YLPhone=3, YLServer=4, BaiduTieba=5, FinanceApp=6, DatingAppIOS=7, DatingAppAndroid=8, YyLoveAppIOS=9, YyLoveAppAndroid=10, YyShiTingAppAndroid=11, YyShiTingAppIOS=12, FinanceSC=13, DatingAppIOSCracked=14, YyLoveAppIOSCracked=15, MEXiaomi=16, MEGameStore=17, MEBaiduBrowser=18, MEBilin=19, DatingBlindIOS=20, DatingBlindAndroid=21, VipPeiLiaoAndroid=22, VipPeiLiaoIOS=23, MEMidas=24, MEBaiduPic=25, MEJD=26, MEEmulatorApp=27, MEKuaikan=28, YYLiveIOS=29, YYLiveAndroid=30, FinancePinAnAndroid=31, FinancePinAnIOS=32, YouXiDaTing=33, WolfKillDatingIOS=34, WolfKillBindingIOS=35, WolfKillDatingAndroid=36, WolfKillBindingAndroid=37, VipNianNianIOS=38, VipNianNianAndroid=39, WolfKillPKIOS=40, WolfKillPKAndroid=41, WolfKillImIOS=42, WolfKillImAndroid=43, WolfKillPKGameIOS=44, WolfKillPKGameAndroid=45, WolfkillXiaochengxu=46, WolfkillExternal=47, XunhuanPkIOS=48, XunhuanPkAndroid=49
}

enum TAccountOperateType {
  Withdraw=1, Exchange=2, ConsumeProps=3, PropsRevenue=4, BuyVirtOverage=5, BuyVirtOverageFail=6, ActivityAutoInc=7, AutoMonthSettle=8, AccountFreeze=9, AccountUnfreeze=10, ChannelRealToPersonVirt=11, IssueSycee=12, Transfer=13, SystemOper=14, SystemCompensate=15, ExternalModification=16, GiftBagLottery=17, ChargeCurrency=18, ChargeCurrencyPresent=19, ChargeCurrencyDiscount=20, DatingBackupGroup=21, DatingBackupGroupBackFee=22, DatingBackupGroupFinish=23, BuySkin=24, BuySeal=25, BuySealBackFee=26, RedPacketIssue=27, RedPacketCharge=28, RedPacketGrab=29, RedPacketClose=30, NobleOpen=31, NobleRenew=32, NobleUpgrade=33, UPGRADE_PROPS=34, REVERT_PAY=35, WithdrawBack=36, NobleRenewExchangeVirt=37, OfficialIssue=38, PayVipRoom=39, ConsumePropsForOther=40, VirtLottery=41, BuySpoofHanging=42, LuckyTreasures=43, ProductConsume=44, ProductConsumeRevert=45, BuyLotteryChance=46, HatKing=47, GuardOpen=48, GuardRenew=49, ModifyGuardIntimateAccount=50, ModifyHatKingPond=51, HatKingReward=52, PropsExchange=53, ComboBonus=54, ExtraAccountEntry=55, RedPacketExpireBack=56, NobleDowngrade=57, ActRevenueSubsidy=58
}

enum TAccountSrcType {
  None=0, OpenNoble=1, BuyNobleCouponGiftbag=2, RenewNobleByOpenCoupon=3, RenewNobleByCouponGiftbag=4, RenewByNobleCoupon3=5, RenewByNobleCoupon6=6, RenewByNobleCoupon9=7, RenewByNobleActivityCoupon=8, Exchange=9, PointExchange=10
}

enum TAppId {
  Finance=1, Dating=2, Hundred=3, FreeShow=4, GameGuild=5, Ktv=6, Blackjack=7, Spy=8, SlaveSales=9, ScratchOff=10, Niuniu=11, MedicalTreatment=12, Sport=13, VipPk=14, HelloApp=15, FinanceForceRelieveContract=16, GameSpot=17, Bilin=18, XunHuan=19, WeiFang=20, TinyTime=21, YoMall=22, GameTemplate=23, MEPlus=24, WerewolfKill=25, TinyVideo=26, MGameVoice=27, DianHu=28
}

enum TCurrencyType {
  Virt=1, Real=2, Activity=3, Yb=4, Time=5, Commission=6, Sycee=7, Golden=8, Silver=9, Copper=10, RMB=11, SilverShell=12, Hello_Golden=13, Hello_Diamond=14, Hello_AppleDiamond=15, Hello_RedDiamond=16, SuperPurpleDiamond=17, RedPacket=18, Xh_Golden=19, Xh_Diamond=20, Xh_Ruby=21, Bilin_Whale=22, Bilin_Profit=23, TinyTime_MiBi=24, TinyTime_MiDou=25, TinyTime_Profit=26, TinyTime_EDou=27, YoMall_Salary=28, ME_Midas_Mibi=29, GameTemplate_Diamond=30, GameTemplate_Ruby=31, YYLive_RedDiamond=32, HappyCoin=33, HappyDiamond=34, MGameDiamond=35, MGameBlackCoin=36, AlipayRedPacket=37, AlipayRedPacketFreeze=38, HappyDrill=39, Bilin_Whale_NEW=40, Xh_Diamond_NEW=41, TB_tdou=100
}

enum TSortType {
  NoSort=0, ASC=1, DESC=2, DEFAULT=3
}

enum TRevenueSrcType {
  Props=1, Tutor=2, ExternalCharge=3, PropsWithoutDaySettleLevel=4, VipPkRevenue=5, FinanceStrategy=6, ForceRelieveContract=7, YoMallSalaryRevenue=8, XunhuanRoomRevenue=9, XunhuanPropsAndRoomRevenue=10, NiuWan=11, Ask=12, DATING_HATKING=13, IssueBonus=14, DATING_6P=15, DATING_PC6P=16, DATING_PC6P_VIRT=17, MGVRandom=18, DATING_6PFACE=19, XunhuanActSubsidy=20, DatingSeal=21
}

exception TServiceException {
  1: i32 code;
  2: string message;
}

struct TPayResult {
  1: i32 code;
  2: string message;
  3: string data;
}

struct TMonthRevenueRecord {
  1: i64 id;
  2: i64 settleDate;
  3: i64 uid;
  4: i32 appid;
  5: i32 contractType;
  6: i32 userType;
  7: i64 accountId;
  8: i32 accountAmount;
  9: i32 accountCurrencyType;
  10: i32 descCurrencyType;
  11: i32 exchangeLevel;
  12: i32 exchangeIncome;
  13: i64 settleTime;
  14: i32 contractWeight;
  15: i64 contractOwuid;
  16: i64 contractAnchorUid;
}

struct TUserPropsTransfer {
  1: i64 transedUid;
  2: i32 count;
  3: bool isMaster;
}

struct PinkDiamondSummary {
  1: i32 totalMonth;
  2: i64 commission;
  3: i64 anchorUid;
  4: i64 sid;
  5: bool forceRelieve;
  6: i64 contractDuration;
  7: i64 owUid;
  8: i32 appid;
}

struct TChannelAccount {
  1: i64 uid;
  2: TCurrencyType currencyType;
  3: i64 amount;
  4: i64 freezed;
  5: i32 appid;
  6: i64 sid;
}

struct TDaySettleLevelConfig {
  1: i32 appid;
  2: i32 contractType;
  3: string level;
  4: i32 greaterAndEqual;
  5: i32 lessThan;
  6: i32 rate;
  7: i32 notSuperAnchorMaxRate;
  8: i32 additionRate;
  9: i32 canAdditionRate;
  10: i32 bonusRate;
  11: i32 overflowLine;
  12: i32 overflowRate;
  13: i32 owRate;
  14: i32 overflowOwRate;
  15: i32 canExceed;
  16: i32 rewardType;
}

struct TCompereStandings {
  1: i64 id;
  2: i64 statDate;
  3: i64 uid;
  4: i64 sid;
  5: i32 appid;
  6: i64 signTime;
  7: i64 relieveTime;
  8: i64 revenueSum;
  9: double monthSettle;
  10: i64 virtExchange;
}

struct TChargeCurrencyConfig {
  1: i32 id;
  2: string name;
  3: i32 appid;
  4: i32 usedChannelType;
  5: i32 destCurrencyType;
  6: i32 chargeRate;
  7: i32 offersType;
  8: i32 offersRate;
  9: i32 srcAmount;
  10: i32 destAmount;
  11: i64 effectStartTime;
  12: i64 effectEndTime;
  13: i32 status;
  14: i32 weight;
  15: string productId;
  16: bool offersCurrencySame;
  17: i32 offersCurrencyType;
  18: string expand;
}

struct TExchangeCurrencyConfig {
  1: i64 id;
  2: i32 appid;
  3: i32 srcCurrencyType;
  4: i32 destCurrencyType;
  5: i32 srcAmount;
  6: i32 destAmount;
  7: i32 exchangeRate;
  8: i32 weight;
  9: i64 startTime;
  10: i64 endTime;
}

struct TUserAccountHistory {
  1: i64 uid;
  2: i64 accountId;
  3: TCurrencyType currencyType;
  4: i64 amountOrig;
  5: i64 amountChange;
  6: i64 freezedOrig;
  7: i64 freezedChange;
  8: i64 optTime;
  9: string description;
  10: TAccountOperateType optType;
  11: i32 appid;
  12: i64 id;
}

struct TDaySettleAdditionRate {
  1: i64 uid;
  2: i32 appid;
  3: i32 additionRate;
  4: i64 additionMonth;
  5: string memo;
  6: i32 legendary;
}


struct TDaySettleAdditionConfig {
  1: i32 appid;
  2: i32 months;
  3: i32 days;
  4: i32 canexceed;
}

struct TRank {
  1: i64 uid;
  2: i64 value;
  3: i64 rank;
}

struct TRevertModifyAccountOrder {
  1: i64 uid;
  2: string seqid;
  3: i32 amount;
  4: i32 status;
  5: TCurrencyType currencyType;
  6: i64 optTime;
}

struct PinkDiamondRec {
  1: i32 id;
  2: i64 uid;
  3: i64 liveUid;
  4: i64 owid;
  5: i64 sid;
  6: i64 ssid;
  7: i32 appId;
  8: i32 yb;
  9: i32 period;
  10: i32 durationType;
  11: i64 activateTime;
  12: i32 usedChannel;
  13: bool monthRewarded;
  14: bool ybRewarded;
  15: i64 contractSid;
  16: i64 payRequestId;
}

struct TDaySettleAdditionResultInfo {
  1: i64 uid;
  2: TAppId appid;
  3: i64 month;
  4: i32 curLevel;
  5: i32 nextLevel;
  6: i32 nextNeedDays;
  7: bool canLevelUp;
  8: bool legendary;
  9: bool additionEnable;
  10: list<i32> additionRates;
  11: i32 curLevelIndex;
  12: i32 nextLevelIndex;
  13: bool isMaxLevel;
  14: list<i32> otherMaxThanCurLevelIndexes;
  15: list<i32> otherMaxThanCurLevels;
  16: list<i32> otherMaxThanCurLevelNeedDays;
  17: list<bool> canLevelUpToOtherMaxThanCurLevels;
  18: i32 days;
  19: string curLevelString;
  20: list<string> additionRateLevels;
  21: list<i32> additionRateGEDays;
  22: list<i32> additionRateDayCounts;
}

struct TQueryPageInfo {
  1: i32 page;
  2: i32 pagesize;
  3: i32 totalElement;
  4: i32 totalPage;
  5: list<map<string, string>> content;
  6: map<string, string> extend;
  7: i32 code;
  8: string message;
}

struct TExtraUserAccount {
  1: i64 id;
  2: i64 uid;
  3: i32 appid;
  4: i64 actId;
  5: i32 currencyType;
  6: i64 amount;
  7: i64 freezed;
  8: i64 validTime;
  9: i64 version;
  10: byte isExchange;
  11: i64 exchangeTime;
}

struct TChannelAccountCumulative {
  1: i64 id;
  2: i64 uid;
  3: i64 sid;
  4: i32 appid;
  5: i32 currencyType;
  6: i64 amount;
}

struct TUserAccountPeriod {
  1: i64 id;
  2: i64 uid;
  3: i32 currencyType;
  4: i64 totalAmount;
  5: i64 amount;
  6: i32 isFreezed;
  7: i32 appid;
  8: i64 createTime;
  9: i64 startTime;
  10: i64 endTime;
  11: i64 version;
}

struct TExtraUserAccountHistory {
  1: i64 id;
  2: i64 accountId;
  3: i64 uid;
  4: i32 currencyType;
  5: i64 amountOrig;
  6: i64 amountChange;
  7: i64 freezedOrig;
  8: i64 freezedChange;
  9: i32 optType;
  10: i64 optTime;
  11: string description;
  12: string userIp;
  13: i64 actId;
  14: i32 appid;
  15: string seqId;
  16: i64 validTime;
  17: string platform;
  18: string device;
}

struct TRevenueRecord {
  1: i64 id;
  2: i64 uid;
  3: i64 contributeUid;
  4: i64 sid;
  /* 收礼心值 */
  5: double income;
  /* 分成比例 */
  6: double incomeRate;
  /* 收入，单位是收益币，10000收益币=1元，公会抽成收入就是income-realIncome */
  7: double realIncome;
  8: i64 optTime;
  9: i64 revenueDate;
  10: i32 revenueType;
  11: i32 exchageLevel;
  12: i32 appid;
  13: i32 srcType;
  14: i32 additionRate;
  15: double allIncome;
  16: string level;
  17: i32 bonusRate;
  18: double bonusIncome;
  19: i32 bonusStatus;
}

struct TRevenueData {
  1: i64 income;
  2: i64 realIncome;
}

struct TCurrencyIssue {
  1: i64 uid;
  2: i64 amount;
  3: i32 currencyType;
  4: i32 appid;
  5: string seq;
  6: i64 createTime;
  7: i64 finishtime;
  8: i32 status;
}

struct TMonthSettleApply {
  1: i64 id;
  2: i64 uid;
  3: i64 settleDate;
  4: i64 applyTime;
  5: double applyRealAmount;
  6: double exchangeSalaryAmount;
  7: i32 result;
  8: i32 appid;
  9: i32 destCurrencyType;
  10: TUserType userType;
  11: double compensationAmount;
  12: i32 reOrderFlag;
  13: i32 settleType;
  14: i32 withdrawAccountType;
  15: string withdrawAccount;
  16: i64 contractSid;
  17: i32 violationPunishLevel;
  18: double violationPunishAmount;
  19: string resultMsg;
}

struct TUserAccount {
  1: i64 uid;
  2: TCurrencyType currencyType;
  3: i64 amount;
  4: i64 freezed;
  5: i32 appid;
}

service TCurrencyService {

  /**
     * 查询频道指定特定时间范围内的收入情况 timeGreaterThan与timeLessThan不能同时为null
     * 
     * @param sid
     * @param appid
     * @param timeGreaterThan 日期格式：yyyyMMddHHmmss
     * @param timeLessThan 日期格式：yyyyMMddHHmmss
     * @param revenueUserType 1-主播，2-OW，0-所有
     * @param page 从1开始
     * @param pagesize
     * @param anchorUid 当revenueUserType=1时，需要传入主播UID
     * @return
     */
  list<TRevenueRecord> queryRevenueRecord(1: i64 sid, 2: i32 appid, 3: string timeGreaterThan, 4: string timeLessThan, 5: i32 revenueUserType, 6: i32 page, 7: i32 pagesize, 8: i64 anchorUid);
  
  
  
   /**
     * 提现记录查询
     * 
     * @param uid
     * @param timeGreaterThan 
     * @param timeLessThan
     * @param appid
     * @param page
     * @param pagesize
     * @return
     */
  list<TMonthSettleApply> queryUserMonthSettleApply(1: i64 uid, 2: i64 timeGreaterThan, 3: i64 timeLessThan, 4: i32 appid, 5: i32 page, 6: i32 pagesize);
  
  /*
  * 查看ow账户的可提现金额
  * currencyType: Bilin_Profit 10000收益币=1元 
  */
  TChannelAccount getChannelAccountByUidAndType(1: i64 uid, 2: i64 sid, 3: TAppId appid, 4: TCurrencyType currencyType);

  TQueryPageInfo queryRevenueRecordPaging(1: i64 uid, 2: TAppId appid, 3: i32 pagesize, 4: i32 page, 5: TUserType revenueUserType, 6: i64 startDate, 7: i64 endDate, 8: i64 anchorUid, 9: i64 sid, 10: TRevenueSrcType srcType);

}

service TUserAccountService {

  /*
   * 查询可提现金额
   * currencyType: Bilin_Profit 10000收益币=1元 
  */
  TUserAccount getUserAccountByUidAndType(1: i64 uid, 2: TAppId appid, 3: TCurrencyType currencyType);
  
  /*
  * 查询总心值
  */
  i64 bilinCumulativeProfit(1: i64 uid)
 
}

struct TContract {
  1: i64 liveUid;
  2: string groupName;
  3: i64 sid;
  4: i64 owUid;
  5: i32 weight;
  6: i64 signTime;
  7: i32 appid;
  8: i32 companySign;
  9: i64 finishTime;
  10: i32 months;
  11: i32 superAnchorSign;
  12: i32 templateId;
}

struct TWeekPropsRecvInfo {
  1: i64 uid;
  2: i32 propId;
  3: string propName;
  4: i32 pricingId;
  5: double amount;
  6: i64 usedTime;
  7: i64 sid;
  8: i32 propCnt;
  9: string guestUid;
  10: i64 anchorUid;
  11: double sumAmount;
  12: i64 id;
  13: TCurrencyType currencyType;
  14: i32 appid;
  15: i32 playType;
  16: string expand;
}
struct TWeekPropsRecvInfoQueryPage {
  1: list<TWeekPropsRecvInfo> content;
  2: i32 page;
  3: i32 pagesize;
  4: i32 totalElement;
  5: i32 totalPage;
  6: map<string, string> extend;
}

service TPingService {
  i64 ping(1: i64 seq);
  void ping2();
}

service TContractService  {

  i32 addContractInfoExternal(1: i64 uid, 2: TAppId appid, 3: i64 sid, 4: i64 owuid, 5: i32 weight, 6: i32 templateId);

  TContract queryContractByAnchor(1: i64 uid, 2: TAppId appid); 
}

service TPropsService {
 
  TWeekPropsRecvInfoQueryPage queryAnchorWeekPropsRecieve(1: i64 uid, 2: TAppId appid, 3: i64 startTime, 4: i64 endTime, 5: i32 page, 6: i32 pagesize, 7: list<i32> propIds, 8: list<i32> playTypes);
  TWeekPropsRecvInfoQueryPage queryChannelWeekPropsRecieve(1: i64 sid, 2: TAppId appid, 3: i64 startTime, 4: i64 endTime, 5: i32 page, 6: i32 pagesize, 7: i64 usedUid);

}

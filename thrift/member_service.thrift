namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.member_service
include "ttype.thrift"


//会员服务
service MemberService{
    /**
     * 注册会员
     * @param member 会员信息
     * @param profile 资料
     * @param mchId 商户编号
     * @param cardId 会员卡号
     * @param inviteCode 邀请码
     **/
    ttype.Result RegisterMemberV1(1:SMember member,2:SProfile profile,3:i32 mchId,4:string cardId,5:string inviteCode)

    // 登录，返回结果(Result)和会员编号(Id);
    // Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
    ttype.Result CheckLogin(1:string user,2:string pwd,3:bool update)
    // 检查交易密码
    ttype.Result CheckTradePwd(1:i64 id,2:string tradePwd)
    // 等级列表
    list<SLevel> LevelList()
    // 获取实名信息
    STrustedInfo GetTrustInfo(1:i64 id)
    // 获取等级信息
    SLevel GetLevel(1:i32 id)
    // 根据SIGN获取等级
    SLevel GetLevelBySign(1:string sign)
    // 根据会员编号获取会员信息
    SMember GetMember(1:i64 id)
    // 根据用户名获取会员信息
    SMember GetMemberByUser(1:string user)
    // 根据会员编号获取会员资料
    SProfile GetProfile(1:i64 id)
    /** 锁定/解锁会员 */
    ttype.Result ToggleLock(1:i64 memberId)
    // 获取会员汇总信息
    SComplexMember Complex(1:i64 memberId)
    // 检查资料是否完成
    ttype.Result CheckProfileComplete(1:i64 memberId)
    // 更改会员等级
    ttype.Result UpdateLevel(1:i64 memberId,2:i32 level,3:bool review,4:i64 paymentOrderId)
    /* 更改手机号码，不验证手机格式 */
    ttype.Result ChangePhone(1:i64 memberId,2:string phone)
     /* 更改用户名 */
    ttype.Result ChangeUsr(1:i64 memberId,2:string usr)
    // 升级为高级会员
    ttype.Result Premium(1:i64 memberId,2:i32 v,3:i64 expires)
    // 获取会员的会员Token,reset表示是否重置token
    string GetToken(1:i64 memberId,2:bool reset)
    // 检查会员的会话Token是否正确，如正确返回: 1
    bool CheckToken(1:i64 memberId,2:string token)
    // 移除会员的Token
    void RemoveToken(1:i64 memberId)
    // 获取会员的收货地址
    list<SAddress> GetAddressList(1:i64 memberId)
    // 获取地址，如果addrId为0，则返回默认地址
    SAddress GetAddress(1:i64 memberId,2:i64 addrId)
    // 获取会员账户信息
    SAccount GetAccount(1:i64 memberId)
    // 获取自己的邀请人会员编号数组
    list<i64> InviterArray(1:i64 memberId,2:i32 depth)
    // 按条件获取荐指定等级会员的数量
    i32 GetInviterQuantity(1:i64 memberId,2:map<string,string> data)
    // 按条件获取荐指定等级会员的列表
    list<i64> GetInviterArray(1:i64 memberId,2:map<string,string> data)
    // 账户充值
    ttype.Result ChargeAccount(1:i64 memberId ,2:i32 account,3:i32 kind,
      4:string title,5:string outerNo,6:double amount,7:i64 relateUser)
    // 抵扣账户
    ttype.Result DiscountAccount(1:i64 memberId,2:i32 account,3:string title,
      4:string outerNo,5:double amount,6:i64 relateUser,7:bool mustLargeZero)
    // 调整账户
    ttype.Result AdjustAccount(1:i64 memberId,2:i32 account,3:double value,4:i64 relateUser,5:string remark,)

    // !银行四要素认证
    ttype.Result B4EAuth(1:i64 memberId,2:string action,3:map<string,string> data)
}


/** 等级 */
struct SLevel {
    1: i32 ID
    2: string Name
    3: i32 RequireExp
    4: string ProgramSignal
    5: i32 IsOfficial
    6: i32 Enabled
}

/** 会员 */
struct SMember {
    1: i64 ID
    2: string Usr
    3: string Pwd
    4: string TradePwd
    5: i32 Exp
    6: i32 Level
    7: string InvitationCode
    // 高级用户类型
    8:i32   PremiumUser
    // 高级用户过期时间
    9:i64   PremiumExpires
    10: string RegFrom
    11: string RegIp
    12: i64 RegTime
    13: string CheckCode
    14: i64 CheckExpires
    15: i32 Flag
    16: i32 State
    17: i64 LoginTime
    18: i64 LastLoginTime
    19: i64 UpdateTime
    20: string DynamicToken
    21: i64 TimeoutTime
}

/** 资料 */
struct SProfile {
    1: i64 MemberId
    2: string Name
    3: string Avatar
    4: i32 Sex
    5: string BirthDay
    6: string Phone
    7: string Address
    8: string Im
    9: string Email
    10: i32 Province
    11: i32 City
    12: i32 District
    13: string Remark
    14: string Ext1
    15: string Ext2
    16: string Ext3
    17: string Ext4
    18: string Ext5
    19: string Ext6
    20: i64 UpdateTime
}

/** 账户 */
struct SAccount {
    1: i64 MemberId
    2: i64 Integral
    3: i64 FreezeIntegral
    4: double Balance
    5: double FreezeBalance
    6: double ExpiredBalance
    7: double WalletBalance
    8: double FreezeWallet
    9: double ExpiredPresent
    10: double TotalPresentFee
    11: double FlowBalance
    12: double GrowBalance
    13: double GrowAmount
    14: double GrowEarnings
    15: double GrowTotalEarnings
    16: double TotalExpense
    17: double TotalCharge
    18: double TotalPay
    19: i64 PriorityPay
    20: i64 UpdateTime
}

struct SComplexMember {
    1: i64 MemberId
    2: string Usr
    3: string Name
    4: string Avatar
    5: i32 Exp
    6: i32 Level
    7: string LevelName
    8: string LevelSign
    9: i32 LevelOfficial
    10: i32	PremiumUser
    11: i64	PremiumExpires
    12: string InvitationCode
    13: i32 TrustAuthState
    14: i32 State
    15: i64 Integral
    16: double Balance
    17: double WalletBalance
    18: double GrowBalance
    19: double GrowAmount
    20: double GrowEarnings
    21: double GrowTotalEarnings
    22: i64 UpdateTime
}

struct SMemberRelation {
    1: i64 MemberId
    2: string CardId
    3: i64 InviterId
    4: string InviterStr
    5: i32 RegisterMchId
}


struct STrustedInfo {
    1: i64 MemberId
    2: string RealName
    3: string CardId
    4: string TrustImage
    5: i32 ReviewState
    6: i64 ReviewTime
    7: string Remark
    8: i64 UpdateTime
}


struct SAddress {
    1: i64 ID
    2: i64 MemberId
    3: string RealName
    4: string Phone
    5: i32 Province
    6: i32 City
    7: i32 District
    8: string Area
    9: string Address
    10: i32 IsDefault
}

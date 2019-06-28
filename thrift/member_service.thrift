namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.member_service
include "ttype.thrift"


//会员服务
service MemberService{
    /**
     * 注册会员
     * @param user 登陆用户名
     * @param pwd 登陆密码,md5运算后的字符串
     * @param flag 用户自定义标志
     * @param phone 手机号码
     * @param email 邮箱
     * @param avatar 头像
     * @param extend 扩展数据
     * @return 注册结果，返回user_code
     */
    ttype.Result RegisterMemberV2(1:string user,2:string pwd,3:i32 flag,4:string name,
        5:string phone,6:string email,7:string avatar,8:map<string,string> extend)

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
    /** 交换会员编号 */
    i64 SwapMemberId(1:ECredentials cred,2:string value)
    // 根据会员编号获取会员信息
    SMember GetMember(1:i64 id)
    // 根据用户名获取会员信息
    SMember GetMemberByUser(1:string user)
    // 根据会员编号获取会员资料
    SProfile GetProfile(1:i64 id)
    /** 激活会员 */
    ttype.Result Active(1:i64 memberId)
    /** 锁定/解锁会员 */
    ttype.Result Lock(1:i64 memberId,2:bool lock,3:string remark)
    /** 标志赋值, 如果flag小于零, 则异或运算 */
    ttype.Result GrantFlag(1:i64 memberId,2:i32 flag)
    /** 获取会员汇总信息 */
    SComplexMember Complex(1:i64 memberId)
    /** 发送会员验证码消息, 并返回验证码, 验证码通过data.code获取 */
    ttype.Result SendCode(1:i64 memberId ,2:string op,3:i32 msgType)
    /** 比较验证码是否正确 */
    ttype.Result CompareCode(1:i64 memberId ,2:string code)
    /** 获取收款码 */
    list<SCollectsCode> GetCollectsCodes(1:i64 memberId)
    /** 保存收款码 */
    ttype.Result SaveCollectsCode(1:i64 memberId,2:SCollectsCode code)
    // 检查资料是否完成
    ttype.Result CheckProfileComplete(1:i64 memberId)
    /** 获取会员等级信息 */
    SMemberLevelInfo MemberLevelInfo(1:i64 memberId)
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
    // 获取邀请会员的数量
    i32 InviteMembersQuantity(1:i64 memberId,2:i32 depth)
    // 按条件获取荐指定等级会员的数量
    i32 QueryInviteQuantity(1:i64 memberId,2:map<string,string> data)
    // 按条件获取荐指定等级会员的列表
    list<i64> QueryInviteArray(1:i64 memberId,2:map<string,string> data)
    // 账户充值,amount精确到分
    ttype.Result AccountCharge(1:i64 memberId ,2:i32 account,3:string title,
      4:i32 amount,5:string outerNo,6:string remark)
    // 账户消耗,amount精确到分
    ttype.Result AccountConsume(1:i64 memberId,2:i32 account,3:string title,
      4:i32 amount, 5:string outerNo,6:string remark)
    // 账户抵扣,amount精确到分
    ttype.Result AccountDiscount(1:i64 memberId,2:i32 account,3:string title,
      4:i32 amount, 5:string outerNo,6:string remark)
    // 账户退款,amount精确到分
    ttype.Result AccountRefund(1:i64 memberId,2:i32 account,3:string title,
        4:i32 amount, 5:string outerNo,6:string remark)
    // 账户人工调整
    ttype.Result AccountAdjust(1:i64 memberId,2:i32 account,3:i32 value,4:i64 relateUser,5:string remark)

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

    /**  */
    1:i64 Id
    /** 用户名 */
    2:string User
    /**  */
    3:string Pwd
    /**  */
    4:string TradePwd
    /**  */
    5:i64 Exp
    /**  */
    6:i32 Level
    /** 高级用户级别 */
    7:i32 PremiumUser
    /** 高级用户过期时间 */
    8:i64 PremiumExpires
    /**  */
    9:string InvitationCode
    /**  */
    10:string RegIp
    /**  */
    11:string RegFrom
    /**  */
    12:i32 State
    /** 会员标志 */
    13:i32 Flag
    /**  */
    14:string Code
    /**  */
    15:string Avatar
    /**  */
    16:string Phone
    /**  */
    17:string Email
    /** 昵称 */
    18:string Name
    /* 用户会员密钥 */
    19:string DynamicToken
    /** 注册时间 */
    20:i64 RegTime
    /** 最后登录时间 */
    21:i64 LastLoginTime
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
    9: double ExpiredWallet
    10: double TotalWalletAmount
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
    1: string Name
    2: string Avatar
    3: string Phone
    4: i32 Exp
    5: i32 Level
    6: string LevelName
    7: string InvitationCode
    8: i32 TrustAuthState
    9: i32 PremiumUser
    10: i32 Flag
    11: i64 UpdateTime
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

/** 收款码 */
struct SCollectsCode{
    /** 编号 */
    1:i32 Id
    /** 账户标识,如:alipay */
    2:string Identity
    /** 账户名称 */
    3:string Name
    /** 账号 */
    4:string AccountId
    /** 收款码地址 */
    5:string CodeUrl
    /** 是否启用 */
    6:i32 State
}

/* 会员等级信息 */
struct SMemberLevelInfo{
    /** 等级 */
    1:i32 Level
    /** 等级名称 */
    2:string LevelName
    /** 经验值 */
    3:i32 Exp
    /** 编程符号 */
    4:string ProgramSignal
    /** 下一级等级,返回-1表示最高级别 */
    5:i32 NextLevel
    /** 下一等级名称 */
    6:string NextLevelName
    /** 编程符号 */
    7:string NextProgramSignal
    /** 需要经验值 */
    8:i32 RequireExp
}


/** 凭据 */
enum ECredentials{
    /** 用户名 */
    User = 1
    /** 用户代码 */
    Code = 2
    /** 邮箱 */
    Email = 3
    /** 手机号码 */
    Phone = 4
}
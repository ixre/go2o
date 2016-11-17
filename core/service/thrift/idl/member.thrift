namespace go define

struct Member {
    1: i64 Id
    2: string Usr
    3: string Pwd
    4: string TradePwd
    5: i64 Exp
    6: i64 Level
    7: string InvitationCode
    8: string RegFrom
    9: string RegIp
    10: i64 RegTime
    11: string CheckCode
    12: i64 CheckExpires
    13: i64 State
    14: i64 LoginTime
    15: i64 LastLoginTime
    16: i64 UpdateTime
    17: string DynamicToken
    18: i64 TimeoutTime
}

struct Profile {
    1: i64 MemberId
    2: string Name
    3: string Avatar
    4: i64 Sex
    5: string BirthDay
    6: string Phone
    7: string Address
    8: string Im
    9: string Email
    10: i64 Province
    11: i64 City
    12: i64 District
    13: string Remark
    14: string Ext1
    15: string Ext2
    16: string Ext3
    17: string Ext4
    18: string Ext5
    19: string Ext6
    20: i64 UpdateTime
}

//会员服务
service MemberService{
    // 登陆，返回结果(Result)和会员编号(Id);
    // Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
    map<string,i32> Login(1:string user,2:string pwd,3:bool update),
    // 根据会员编号获取会员信息
    Member GetMember(1:i32 id),
    // 根据用户名获取会员信息
    Member GetMemberByUser(1:string user),
    // 根据会员编号获取会员资料
    Profile GetProfile(1:i32 id),
}

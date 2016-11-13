namespace go define

//会员服务
service MemberService{
    // 登陆，返回结果(Result)和会员编号(Id);
    // Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
    map<string,i32> Login(1:string user,2:string pwd,3:bool update),
}

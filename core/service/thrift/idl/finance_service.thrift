namespace go define

include "ttype.thrift"

// 财务服务
service FinanceService{
    // 转入(业务放在service,是为person_finance解耦)
    ttype.DResult RiseTransferIn(1:i64 personId,2:i32 transferWith,3:double amount)
}
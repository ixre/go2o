package dao




import(
    "go2o/core/dao/model"
)

type IRbacDao interface{
    // auto generate by gof
    // Get 部门
    GetPermDept(primary interface{}) *model.PermDept
    // GetBy 部门
    GetPermDeptBy(where string, v ...interface{}) *model.PermDept
    // Count 部门 by condition
    CountPermDept(where string, v ...interface{}) (int, error)
    // Select 部门
    SelectPermDept(where string, v ...interface{}) []*model.PermDept
    // Save 部门
    SavePermDept(v *model.PermDept) (int, error)
    // Delete 部门
    DeletePermDept(primary interface{}) error
    // Batch Delete 部门
    BatchDeletePermDept(where string, v ...interface{}) (int64, error)
    // Query paging data
    PagingQueryPermDept(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})
}
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

    // Get 岗位
    GetPermJob(primary interface{}) *model.PermJob
    // GetBy 岗位
    GetPermJobBy(where string, v ...interface{}) *model.PermJob
    // Count 岗位 by condition
    CountPermJob(where string, v ...interface{}) (int, error)
    // Select 岗位
    SelectPermJob(where string, v ...interface{}) []*model.PermJob
    // Save 岗位
    SavePermJob(v *model.PermJob) (int, error)
    // Delete 岗位
    DeletePermJob(primary interface{}) error
    // Batch Delete 岗位
    BatchDeletePermJob(where string, v ...interface{}) (int64, error)
    // Query paging data
    PagingQueryPermJob(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})
 }
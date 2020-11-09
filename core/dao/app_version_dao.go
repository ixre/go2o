package dao




import(
    "go2o/core/dao/model"
)

type IAppProdDao interface{
    // auto generate by gof
    // Get APP产品
    Get(primary interface{})*model.AppProd
    // GetBy APP产品
    GetBy(where string,v ...interface{})*model.AppProd
    // Count APP版本 by condition
    Count(where string, v ...interface{}) (int, error)
    // Select APP产品
    Select(where string,v ...interface{})[]*model.AppProd
    // Save APP产品
    Save(v *model.AppProd)(int,error)
    // Delete APP产品
    Delete(primary interface{}) error
    // Batch Delete APP产品
    BatchDelete(where string,v ...interface{})(int64,error)

    // Get APP版本
    GetVersion(primary interface{})*model.AppVersion
    // GetBy APP版本
    GetVersionBy(where string,v ...interface{})*model.AppVersion
    // Select APP版本
    SelectVersion(where string,v ...interface{})[]*model.AppVersion
    // Save APP版本
    SaveVersion(v *model.AppVersion)(int,error)
    // Delete APP版本
    DeleteVersion(primary interface{}) error
    // Batch Delete APP版本
    BatchDeleteVersion(where string,v ...interface{})(int64,error)
}
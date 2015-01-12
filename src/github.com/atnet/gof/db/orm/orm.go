package orm

type Orm interface {
	//Set orm output information
	SetTrace(b bool)

	CreateTableMap(v interface{}, tableName string)

	Get(primaryVal interface{}, entity interface{}) error

	//get entity by condition
	GetBy(entity interface{}, where string, args ...interface{}) error

	//get entity by sql query result
	GetByQuery(entity interface{}, sql string, args ...interface{}) error

	//Select more than 1 entity list
	//@to : refrence to queryed entity list
	//@params : query condition
	//@where : other condition
	Select(to interface{}, where string, args ...interface{}) error

	SelectByQuery(to interface{}, sql string, args ...interface{}) error

	//delete entity and effect to database
	Delete(entity interface{}, where string, args ...interface{}) (effect int64, err error)

	//delete entity by primary key
	DeleteByPk(entity interface{}, primary interface{}) (err error)

	Save(primary interface{}, entity interface{}) (rows int64, lastInsertId int64, err error)
}

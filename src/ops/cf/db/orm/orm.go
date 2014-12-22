package orm

type Orm interface {
	//Set orm output information
	SetTrace(b bool)
	Get(entity interface{}, primaryVal interface{}) error
	//get entity by condition
	GetBy(entity interface{}, where string) error
	//get entity by sql query result
	GetByQuery(entity interface{}, sql string) error

	CreateTableMap(v interface{}, tableName string)

	//Select more than 1 entity list
	//@to : refrence to queryed entity list
	//@params : query condition
	//@where : other condition
	Select(to interface{}, params interface{}, where string) error

	SelectByQuery(to interface{}, entity interface{}, sql string) error

	//delete entity and effect to database
	Delete(entity interface{}, where string) (effect int64, err error)

	//delete entity by primary key
	DeleteByPk(entity interface{}, primary interface{}) (err error)

	Save(primary interface{}, entity interface{}) (rows int64, lastInsertId int64, err error)
}

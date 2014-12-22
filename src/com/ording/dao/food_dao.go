package dao

import (
	"com/ording/entity"
	"database/sql"
	"fmt"
	"ops/cf/db"
)

type itemDao struct {
	db.Connector
}

func (this *itemDao) GetFoodItemById(partnerId, id int) (e entity.FoodItem) {
	//glob.CurrContext().ORM.Get(&e,id)
	this.Connector.GetOrm().GetByQuery(&e, fmt.Sprintf(`select * FROM it_item
			INNER JOIN it_category c ON c.id = it_item.cid WHERE it_item.id=%d
			AND c.ptid=%d`, id, partnerId))

	//
	//	var es []entity.FoodItem
	//
	//	glob.CurrContext().ORM.Select(&es,e,"1=1")
	//	fmt.Println(es)

	//	e.Id = 0
	//	e.Name = e.Name+"-2"
	//	SaveFoodItem(&e)

	//	var e2 []entity.FoodItem
	//	glob.CurrContext().ORM.SelectByQuery(&e2,fmt.Sprintf(`select * FROM it_item
	//			INNER JOIN it_category c ON c.id = it_item.cid WHERE it_item.id=%d
	//			AND c.ptid=%d`,id,partnerId))
	//
	//	fmt.Println(e2)
	return e
}

//获取食物数量
func (this *itemDao) FoodItemsCount(partnerId, cid int) (count int) {
	this.Connector.QueryRow(`
		SELECT COUNT(0) FROM it_item f
	INNER JOIN it_category c ON f.cid = c.id
	 where c.ptid = ?
	AND (cid == -1 OR cid = ?)
	`, func(r *sql.Row) {
		r.Scan(count)
	}, partnerId, partnerId)
	return count
}

func (this *itemDao) DelFoodItem(partnerId, id int) (row int, err error) {
	r, _, _ := this.Connector.Exec(`
		DELETE f,f2 FROM it_item AS f
		INNER JOIN it_category AS c ON f.cid=c.id
		INNER JOIN it_itemprop as f2 ON f2.id=f.id
		WHERE f.id=? AND c.ptid=?`, id, partnerId)
	return r, nil
}

func (this *itemDao) SaveFoodItem(item *entity.FoodItem) (int, error) {
	orm := this.Connector.GetOrm()
	if item.Id <= 0 {
		//多行字符用
		_, id, err := orm.Save(nil, item)
		return int(id), err
	} else {
		_, _, err := orm.Save(item.Id, item)
		return item.Id, err
	}
}

func (this *itemDao) GetItemsByCid(partnerId, categoryId, num int) (e []entity.FoodItem) {
	var sql string
	if num <= 0 {
		sql = fmt.Sprintf(`SELECT * FROM it_item INNER JOIN it_category ON it_item.cid=it_category.id
		WHERE ptid=%d AND it_category.id=%d`, partnerId, categoryId)
	} else {
		sql = fmt.Sprintf(`SELECT * FROM it_item INNER JOIN it_category ON it_item.cid=it_category.id
		WHERE ptid=%d AND it_category.id=%d LIMIT 0,%d`, partnerId, categoryId, num)
	}
	e = []entity.FoodItem{}
	err := this.Connector.GetOrm().SelectByQuery(&e, entity.FoodItem{}, sql)
	if err != nil {
		return nil
	}
	return e
}

/*
def getfooditems(partnerid,page,size,cid=-1):
   '''获取食物'''
   return newdb(True).fetchall('''SELECT f.*,c.name as cname FROM it_item f
    INNER JOIN it_category c ON f.cid=c.id where c.ptid=%(ptid)s'''
    +(' AND cid=%(cid)s' if cid!=-1 else '')+
    ''' ORDER BY updatetime desc limit %(s)s,%(e)s''',
                         {
                          'ptid':partnerid,
                          's':(page-1)*size,
                          'e':size,
                          'cid':cid
                          }
                         )

def getfoods(partnerid,cid):
    return newdb(True).fetchall('''SELECT f.*,c.name as cname FROM it_item f
    INNER JOIN it_category c ON f.cid=c.id where c.ptid=%(ptid)s
      AND cid=%(cid)s ORDER BY id''',
                         {
                          'ptid':partnerid,
                          'cid':cid
                          }
                                )

def getcupornfoods(partnerid):
    return newdb(True).fetchall('''SELECT f.*,c.name as cname FROM it_item f
    INNER JOIN it_category c ON f.cid=c.id where c.ptid=%(ptid)s
       AND `percent`<1 ORDER BY id''',
                         {
                          'ptid':partnerid
                          }
                                )
*/

package entity

import "time"

//门店
type Shop struct {
	Id         int       `db:"id" pk:"yes" auto:"yes"`
	PartnerId  int       `db:"pt_id"`
	Name       string    `db:"name"`
	Address    string    `db:"address"`
	Phone      string    `db:"phone"`
	OrderIndex int       `db:"order_index"`
	State      int       `db:"state"`
	CreateTime time.Time `db:"create_time"`
}

/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"go2o/src/core/infrastructure"
	"go2o/src/core/infrastructure/log"
	"go2o/src/core/variable"
	"regexp"
)

type PartnerQuery struct {
	db.Connector
	gof.Storage
}

func NewPartnerQuery(c gof.App) *PartnerQuery {
	return &PartnerQuery{
		Connector: c.Db(),
		Storage:   c.Storage(),
	}
}

/*


def getpage(partnerid,type):
    '获取页面内容'
    row = newdb().fetchone('SELECT `content` FROM pt_page WHERE `ptid`=%(ptid)s and `type`=%(type)s',
                            {
                             'ptid':partnerid,
                             'type':type
                            })

    return '' if row==None else row[0]


def updatepage(partnerid,type,content):
    '更新页面内容'
    data={
          'ptid':partnerid,
          'type':type,
          'content':content,
          'updatetime':utility.timestr()
          }

    row=newdb().query('UPDATE pt_page SET `content`=%(content)s,`updatetime`=%(updatetime)s WHERE `ptid`=%(ptid)s AND `type`=%(type)s',data)
    if row==0:
        newdb().query('INSERT INTO pt_page (`ptid`,`type`,`content`,`updatetime`) VALUES(%(ptid)s,%(type)s,%(content)s,%(updatetime)s)',data)

*/

var (
	commHostRegexp *regexp.Regexp
)

func getHostRegexp() *regexp.Regexp {
	if commHostRegexp == nil {
		commHostRegexp = regexp.MustCompile("([^.]+)." +
			infrastructure.GetApp().Config().GetString(variable.ServerDomain))
	}
	return commHostRegexp
}

// 根据主机查询商户编号
func (this *PartnerQuery) QueryPartnerIdByHost(host string) int {
	//  $ 获取合作商ID
	// $ hostname : 域名
	// *.wdian.net  二级域名
	// www.dc1.com  顶级域名

	var partnerId int
	var err error

	reg := getHostRegexp()
	if reg.MatchString(host) {
		matches := reg.FindAllStringSubmatch(host, 1)
		usr := matches[0][1]
		err = this.Connector.ExecScalar(`SELECT id FROM pt_partner WHERE usr=?`, &partnerId, usr)
	} else {
		err = this.Connector.ExecScalar(
			`SELECT id FROM pt_partner INNER JOIN pt_siteconf
					 ON pt_siteconf.partner_id = pt_partner.id
					 WHERE host=?`, &partnerId, host)
	}

	if err != nil {
		log.PrintErr(err)
	}
	return partnerId
}

// 验证用户密码并返回编号
func (this *PartnerQuery) Verify(usr, pwd string) int {
	var id int = -1
	this.Connector.ExecScalar("SELECT id FROM pt_partner WHERE usr=? AND pwd=?", &id, usr, pwd)
	return id
}

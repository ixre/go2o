#-*- coding:utf-8 -*-
'''万绿园商城会员同步'''
import MySQLdb
import org.conf as conf
import org.utility
import hashlib

def adduser(request):
    conn=MySQLdb.connect(
                          host=conf.db_host,
                          user='root',
                          passwd='$Newmin',
                          db='wly_mall',
                          port=conf.db_port,
                          charset='utf8'
                          )
    cursor=conn.cursor()
    cursor.execute('''INSERT INTO  mall_users
                    (user_name, email, sex, birthday, reg_time, password)
                    VALUES(%(user_name)s,%(email)s,%(sex)s,%(birthday)s,%(reg_time)s,%(password)s)''',
                    {
                      'user_name':request.username,
                      'email':request.email,
                      'sex':request.sex,
                      'birthday':request.birthday,
                      'reg_time':org.utility.timestr(),
                      'password':hashlib.md5(request.password).hexdigest()
                     })
    
    conn.commit()
    cursor.close()
    conn.close()
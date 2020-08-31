package orm

import (
	"fmt"
	"strings"

	"xorm.io/xorm"
)

type XormSessionCreator interface {
	GetXormSession() *xorm.Session
}
type XormBase struct {
	XormSessionCreator
}

func (c *XormBase) XormGetByMap(q map[string]interface{}, columns []string, result interface{}) (exist bool, err error) {
	sess := c.GetXormSession()
	if q != nil {
		var query []string
		var args []interface{}
		for k, v := range q {
			query = append(query, fmt.Sprintf("%s=? ", k))
			args = append(args, v)
		}
		sess = sess.Where(strings.Join(query, " and "), args...)
	}
	if len(columns) > 0 {
		sess = sess.Cols(columns...)
	}
	return sess.Get(result)
}

func (c *XormBase) XormFindByMap(q map[string]interface{}, columns []string, result interface{}) (err error) {
	sess := c.GetXormSession()
	if q != nil {
		var query []string
		var args []interface{}
		for k, v := range q {
			query = append(query, fmt.Sprintf("%s=? ", k))
			args = append(args, v)
		}
		sess = sess.Where(strings.Join(query, " and "), args...)
	}
	if len(columns) > 0 {
		sess = sess.Cols(columns...)
	}
	return sess.Find(result)
}

func (c *XormBase) XormUpdateByMap(q map[string]interface{}, columns []string, data interface{}) (rowsAffected int64, err error) {
	sess := c.GetXormSession()
	if q != nil {
		var query []string
		var args []interface{}
		for k, v := range q {
			query = append(query, fmt.Sprintf("%s=? ", k))
			args = append(args, v)
		}
		sess = sess.Where(strings.Join(query, " and "), args...)
	}
	if len(columns) > 0 {
		sess = sess.Cols(columns...)
	}
	return sess.UseBool().Update(data)
}



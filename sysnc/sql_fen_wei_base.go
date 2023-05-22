package sysnc

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

const WENWEITableName = "BASE_FW"

type SqlfenweiBase struct {
	FWBM string    `json:"FWBM"`
	FWMC string    `json:"FWMC"`
	FWXS float64   `json:"FWXS"`
	FWJC string    `json:"FWJC"`
	CJSJ time.Time `json:"CJSJ"`
	XGSJ time.Time `json:"XGSJ"`
}

type CpFenWeiList []SqlfenweiBase

func (rList *CpFenWeiList) GetListPage(size, offset int, d *gorm.DB) (err error) {
	if d == nil {
		d = db.SqlDB
	}
	query := d.Table(WENWEITableName)

	err = query.Offset(offset).Limit(size).Order("FWBM desc").Find(&rList).Error
	return

}

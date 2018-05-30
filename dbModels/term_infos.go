package dbModels

type DbTermInfo struct {
	MchtCd    string `gorm:"type:varchar(15);primary_key"`
	TermId    string `gorm:"type:varchar(8);primary_key"`
	MchtNm    string `gorm:"type:varchar(40)"`
	ActiveFlg string `gorm:"type:varchar(1);default:'0'"`
	SignFlg   string `gorm:"type:varchar(1);default:'0'"`
	//批次号
	Term_batch string `gorm:"type:varchar(6);default:'000001'"`
	Term_seq   string `gorm:"type:varchar(6);default:'000001'"`
	//term param
	TermProd   string `gorm:"type:varchar(10)"`
	TermModel  string `gorm:"type:varchar(28)"`
	BrandKsn   string `gorm:"type:varchar(64)"`
	BrandSn    string `gorm:"type:varchar(20)"`
	ActiveCode string `gorm:"type:varchar(10)"`
	//term key info
	KeyDownFlg   string `gorm:"type:varchar(10)"`
	SessionKey   string `gorm:"type:varchar(32)"`
	TermPubKey   string `gorm:"type:text(5000)"`
	ServerPubKey string `gorm:"type:text(5000)"`
	ServerPriKey string `gorm:"type:text(5000)"`
	TermTmk      string `gorm:"type:varchar(32)"`
	TermTmkChk   string `gorm:"type:varchar(16)"`
	TermPik      string `gorm:"type:varchar(32)"`
	TermPikChk   string `gorm:"type:varchar(16)"`
	TermTrk      string `gorm:"type:varchar(32)"`
	TermTrkChk   string `gorm:"type:varchar(16)"`
	PbocIcKey    string `gorm:"type:text(5000)"`
	PbocIcParam  string `gorm:"type:text(5000)"`
	DbBase
}

func (t DbTermInfo) TableName() string {
	return "term_infos"
}


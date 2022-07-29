package LDB

type IDataTables interface {
	Register(fileName string, protoType interface{})
	Load(resPath string) error
	HotUpdateLoadFile(fileName string) error
	Get(fileName string) interface{}
}

func NewDataTables() IDataTables {
	ldb := &_DataTables{}
	ldb.init()
	return ldb
}

package LDB

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

type _DataTables struct {
	dic     map[string]interface{}
}

func (dt *_DataTables) init() {
	dt.dic = make(map[string]interface{})
}

func (dt *_DataTables) Register(fileName string, protoType interface{}) {
	dt.dic[fileName] = protoType
}

//
func (dt *_DataTables) loadFile(fileName string) error {

	baseName := filepath.Base(fileName)
	name := strings.TrimSuffix(baseName, ".bytes")

	pt, ok := dt.dic[name]
	if !ok {
		log.Debug("couldn't find ", name, "  prototype")
		return nil
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Debug("read file error ", fileName, "  ", err)
		return err
	}

	err = proto.Unmarshal(data, pt.(proto.Message))
	if err != nil {
		log.Debug("unmarshal data error ", fileName, reflect.TypeOf(pt).Name())
		return err
	}
	return nil
}

func (dt *_DataTables) HotUpdateLoadFile(fileName string) error {

	baseName := filepath.Base(fileName)
	name := strings.TrimSuffix(baseName, ".bytes")

	pt, ok := dt.dic[name]
	if !ok {
		return errors.New(fmt.Sprint("couldn't find ", name, "  prototype"))
	}

	_dic := make(map[string]interface{})
	for key, val := range dt.dic {
		if key == name {
			pt = reflect.New(reflect.TypeOf(val).Elem()).Interface()
			_dic[key] = pt
		} else {
			_dic[key] = val
		}
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error("read file error ", fileName, "  ", err)
		return err
	}

	err = proto.Unmarshal(data, pt.(proto.Message))
	if err != nil {
		log.Error("unmarshal data error ", fileName, reflect.TypeOf(pt).Name())
		return err
	}

	dt.dic = _dic
	return nil
}

func (dt *_DataTables) Load(resPath string) error {

	if !strings.HasSuffix(resPath, "/") && !strings.HasSuffix(resPath, "\\") {
		resPath += "/"
	}
	rd, err := ioutil.ReadDir(resPath)
	if err != nil {
		return err
	}

	for _, fi := range rd {

		if !fi.IsDir() && filepath.Ext(fi.Name()) == ".bytes" {
			err := dt.loadFile(resPath + fi.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (dt *_DataTables) Get(fileName string) interface{} {
	pt, ok := dt.dic[fileName]
	if !ok {
		return nil
	}
	return pt
}

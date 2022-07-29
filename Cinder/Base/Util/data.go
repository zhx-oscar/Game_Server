package Util

import "encoding/json"

type IData interface {
	SetBool(v bool)
	GetBool() bool

	Set(v interface{})
	Get() interface{}

	SetInt(v int)
	GetInt() int

	SetString(v string)
	GetString() string

	TraversalPairs(func(key, value interface{}) bool)
	IsPairExist(key string) bool

	AddPair(key string, value interface{})
	GetValue(key string) interface{}
	GetStringValue(key string) string
	GetIntValue(key string) int
	GetFloat32Value(key string) float32
	GetFloat64Value(key string) float64
	GetJsonData() string
}

type _Data struct {
	Data    interface{}
	DataMap map[string]interface{}
}

func NewData() IData {
	return &_Data{
		Data:    nil,
		DataMap: make(map[string]interface{}),
	}
}

func NewStrData(strVal string) IData {
	return &_Data{
		Data:    strVal,
		DataMap: make(map[string]interface{}),
	}
}

func NewIntData(intVal int) IData {
	return &_Data{
		Data:    float64(intVal),
		DataMap: make(map[string]interface{}),
	}
}

func NewDataFromBytes(data string) IData {

	d := NewData()

	err := json.Unmarshal([]byte(data), d)
	if err != nil {
		return NewData()
	}

	return d
}

func (d *_Data) SetBool(v bool) {
	d.Data = v
}

func (d *_Data) GetBool() bool {
	return d.Data.(bool)
}

func (d *_Data) Set(v interface{}) {
	d.Data = v
}

func (d *_Data) Get() interface{} {
	return d.Data
}

func (d *_Data) SetInt(v int) {
	d.Data = float64(v)
}

func (d *_Data) GetInt() int {
	return int(d.Data.(float64))
}

func (d *_Data) SetString(v string) {
	d.Data = v
}

func (d *_Data) GetString() string {
	return d.Data.(string)
}

func (d *_Data) AddPair(key string, value interface{}) {
	d.DataMap[key] = value
}

func (d *_Data) GetValue(key string) interface{} {
	return d.DataMap[key]
}

func (d *_Data) GetStringValue(key string) string {
	return d.DataMap[key].(string)
}

func (d *_Data) GetIntValue(key string) int {
	return int(d.DataMap[key].(float64))
}

func (d *_Data) GetFloat32Value(key string) float32 {
	return float32(d.DataMap[key].(float64))
}

func (d *_Data) GetFloat64Value(key string) float64 {
	return d.DataMap[key].(float64)
}

func (d *_Data) GetJsonData() string {

	bs, err := json.Marshal(d)
	if err != nil {
		return ""
	}

	return string(bs)
}

func (d *_Data) TraversalPairs(f func(key, value interface{}) bool) {
	for k, v := range d.DataMap {
		if !f(k, v) {
			break
		}
	}
}

func (d *_Data) IsPairExist(key string) bool {
	_, ok := d.DataMap[key]
	return ok
}

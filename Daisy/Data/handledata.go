package Data

import "Daisy/DataTables"

func initHandleData() {
	SkillData = make(map[uint32]map[uint32]*DataTables.SkillValue_Config, len(GetSkillConfig().SkillMain_ConfigItems))

	hotUpdateHandleData(false)
}

// hotUpdateHandleData 热更 handleData
func hotUpdateHandleData(isHotUpdate bool) {
	handleData_SkillData()
	handleData_SeasonLevelData()
}

var SkillData map[uint32]map[uint32]*DataTables.SkillValue_Config

func handleData_SkillData() {
	exceldata := GetSkillConfig()
	if exceldata == nil {
		return
	}
	_SkillData := make(map[uint32]map[uint32]*DataTables.SkillValue_Config, len(exceldata.SkillMain_ConfigItems))
	for _, val := range exceldata.SkillValue_ConfigItems {
		temp, ok := _SkillData[val.SkillID]
		if !ok {
			temp = map[uint32]*DataTables.SkillValue_Config{
				val.Level: val,
			}
		} else {
			temp[val.Level] = val
		}

		_SkillData[val.SkillID] = temp
	}

	SkillData = _SkillData
}

var SeasonLevelData map[uint32]map[uint32]*DataTables.SeasonLevel_Config
func handleData_SeasonLevelData() {
	exceldata := GetSeasonConfig()
	if exceldata == nil {
		return
	}
	_data := make(map[uint32]map[uint32]*DataTables.SeasonLevel_Config)
	for _, val := range exceldata.SeasonLevel_ConfigItems {
		_, ok := _data[val.SeasonID]
		if !ok {
			_data[val.SeasonID] = make(map[uint32]*DataTables.SeasonLevel_Config)
		}
		_data[val.SeasonID][val.Seasonlevel] = val
	}
	SeasonLevelData = _data
}
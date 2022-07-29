package conf

import (
	"Daisy/DataTables"
	"Daisy/Fight/internal/log"
	"encoding/json"
)

type BlackBoardKeys struct {
	Data map[string]interface{}
}

// AIInfo 行为树配置
type AIInfo struct {
	ID             uint32         //AI ID
	TreeName       string         //行为树路径
	BlackBoardKeys BlackBoardKeys //行为树黑板
}

//loadAIInfo 更新AIdata
func loadAIInfo(data *DataTables.AIData_Config_Data) map[uint32]*AIInfo {
	aiDdata := map[uint32]*AIInfo{}

	if data == nil {
		log.Error("LoadAIInfo data is nil")
		return nil
	}

	for id, val := range data.AIData_ConfigItems {
		temp := &AIInfo{
			ID:       id,
			TreeName: val.TreeName,
		}
		if len(val.BlackBoardKeys) > 0 {
			err := json.Unmarshal([]byte(val.BlackBoardKeys), &temp.BlackBoardKeys.Data)
			if err != nil {
				log.Error("BlackBoardKeys Unmarshal fail ", err, id, val.BlackBoardKeys)
				continue
			}
		}

		aiDdata[id] = temp
	}

	return aiDdata
}

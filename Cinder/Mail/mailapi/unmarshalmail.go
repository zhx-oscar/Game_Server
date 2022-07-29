package mailapi

import (
	"Cinder/Mail/mailapi/types"
	"encoding/json"
)

// UnmarshalMail 解包 Mail json 串
func UnmarshalMail(mailJsonData []byte) (types.Mail, error) {
	m := types.Mail{}
	if err := json.Unmarshal(mailJsonData, &m); err != nil {
		return m, err
	}
	return m, nil
}

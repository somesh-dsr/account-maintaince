package generic

import (
	"encoding/hex"
	"encoding/json"
	"crypto/rand"
)

func UnMarshalData(data[]byte,dataStruct interface{})error{
	err := json.Unmarshal(data,dataStruct)
	return err
}
func MarshalData(data []byte)([]byte,error){
	resJson,err := json.Marshal(data)
	return resJson,err
}
func GetKey()(string,error){
	bytes := make([]byte, 10)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
package tempFile

//import (
//	"encoding/json"
//	"os"
//)
//
//func SaveFile(data interface{}) error {
//	f, err := os.OpenFile("", os.O_CREATE|os.O_APPEND, 0777)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	s, _ := json.Marshal(data)
//	f.WriteString(string(s))
//	return nil
//}

package bdsr

import "net/http"
import "bytes"
import "encoding/base64"
import "strconv"
import "errors"
import "io/ioutil"
import "encoding/json"

func Request(access_token string, rawWAV []byte) (string,error) {
	reqData := make(map[string]interface{})
	reqData["format"]="wav"
	reqData["rate"]=8000
	reqData["channel"]=1
	reqData["token"]=access_token
	reqData["cuid"]="bdsr_api"
	reqData["len"]=len(rawWAV)
	reqData["speech"]=base64.StdEncoding.EncodeToString(rawWAV)

	reqJSON,err := json.Marshal(reqData)
	if err!=nil {
		return "",err
	}

	res,err := http.Post("http://vop.baidu.com/server_api","application/json",bytes.NewBuffer([]byte(reqJSON)))
	if err!=nil {
		return "",err
	}
	result,err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err!=nil {
		return "",err
	}

	decoded := make(map[string]interface{})

	err = json.Unmarshal(result,&decoded)
	if err!=nil {
		return "",err
	}

//!
	err_no_i,ok := decoded["err_no"]
	if !ok {
		return "",errors.New("Bad return format")
	}
	err_no := int(err_no_i.(float64))
	if err_no!=0 {
		return "",errors.New("API error: "+strconv.Itoa(err_no))
	}

	return decoded["result"].([]string)[0],nil
}


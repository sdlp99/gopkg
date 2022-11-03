package str

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

func BasePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}

	return dir + string(os.PathSeparator)
}

func ExistsFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func B2s(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

func S2b(s string) []byte {
	return []byte(s)
}

// 返回右边和左边不一样的文本
func Jsondiff(left, right string, ignoreColumn string) string {
	jmap := gjson.Parse(left)
	ret := ""
	jmap.ForEach(func(key, line gjson.Result) bool {
		r := gjson.Get(right, key.String())
		//不在忽略字段上面
		if find := strings.Contains(";"+ignoreColumn+";", ";"+key.String()+";"); !find {

			if r.String() != line.String() {
				ret += "|" + key.String() + ":" + line.String() + "," + r.String()
			}
		}
		return true
	})
	return ret
}

func Map2Json(object map[string]interface{}) string {
	jsonStr, _ := json.Marshal(object)
	return string(jsonStr)
}

func MapList2Json(object []map[string]interface{}) string {
	jsonStr, _ := json.Marshal(object)
	return string(jsonStr)
}

func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func ConvertByte2String(byte []byte) string {
	var str string
	var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
	str = string(decodeBytes)
	return str
}
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		if start == "" {
			n = 0
		} else {
			n = n + 1
		}
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	} else {
		if end == "" {
			m = len(str)
		}
	}
	str = string([]byte(str)[:m])
	return str
}

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	if strings.Contains(dir, GetTmpDir()) {
		return GetCurrentAbPathByCaller()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func GetTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
		abPath += "/.."
	}
	return abPath
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

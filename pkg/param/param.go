package param

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Params map[string]interface{}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func GetQueryString(params Params) string {
	var data []string
	for k, v := range params {
		if v != "" && k != "sign" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	return strings.Join(data, "&")
}

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (m Params) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}
	start.Name = xml.Name{"", "xml"}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: strval(v)})
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML unmarshals the XML into a map of string to strings,
// creating a key in the map for each tag and setting it's value to the
// tags contents.
//
// The fact this function is on the pointer of Map is important, so that
// if m is nil it can be initialized, which is often the case if m is
// nested in another xml structurel. This is also why the first thing done
// on the first line is initialize it.
func (m *Params) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Params{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func strval(value interface{}) string {
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


func (m Params) FormatURLParam() string {
	v := url.Values{}
	for key, value := range m {
		v.Add(key, strval(value))
	}
	return v.Encode()
}
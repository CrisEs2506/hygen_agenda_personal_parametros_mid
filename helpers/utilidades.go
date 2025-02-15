package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/xray"
)

const (
	JSON_error          string = "Error en el archivo JSON"
	ErrorParametros     string = "Error en los parametros de ingreso"
	ErrorBody           string = "Cuerpo de la peticion invalido"
	AppJson             string = "application/json"
	Calibri             string = "Calibri"
	CalibriBold         string = "Calibri-Bold"
	MinionProBoldCn     string = "MinionPro-BoldCn"
	MinionProMediumCn   string = "MinionPro-MediumCn"
	MinionProBoldItalic string = "MinionProBoldItalic"
)

// Envia una petición con datos al endpoint indicado y extrae la respuesta del campo Data para retornarla
func SendRequestNew(endpoint string, route string, trequest string, target interface{}, datajson interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + beego.AppConfig.String(endpoint) + route
	var response map[string]interface{}
	var err error
	err = SendJson(url, trequest, &response, &datajson)
	err = ExtractData(response, target, err)
	return err
}

// Envia una petición con datos a endpoints que responden con el body sin encapsular
func SendRequestLegacy(endpoint string, route string, trequest string, target interface{}, datajson interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + beego.AppConfig.String(endpoint) + route
	if err := SendJson(url, trequest, &target, &datajson); err != nil {
		return err
	}
	return nil
}

// Envia una petición al endpoint indicado y extrae la respuesta del campo Data para retornarla
func GetRequestNew(endpoint string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + beego.AppConfig.String(endpoint) + route
	var response map[string]interface{}
	var err error
	err = GetJson(url, &response)
	err = ExtractData(response, &target, err)
	return err
}

// Envia una petición a endpoints que responden con el body sin encapsular
func GetRequestLegacy(endpoint string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + beego.AppConfig.String(endpoint) + route
	if err := GetJson(url, target); err != nil {
		return err
	}
	return nil
}

func GetRequestWSO2(service string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") +
		beego.AppConfig.String("AgendaPersonalParametrosApiMidUrlWso2") +
		beego.AppConfig.String(service) + "/" + route
	if response, err := GetJsonWSO2Test(url, &target); response == 200 && err == nil {
		return nil
	} else {
		return err
	}
}

// Esta función extrae la información cuando se recibe encapsulada en una estructura
// y da manejo a las respuestas que contienen arreglos de objetos vacíos
func ExtractData(respuesta map[string]interface{}, v interface{}, err2 error) error {
	var err error

	if err2 != nil {
		return err2
	}
	if respuesta["Success"] == false {
		err = errors.New(respuesta["Message"].(string))
		panic(err)
	}
	datatype := fmt.Sprintf("%v", respuesta["Data"])
	switch datatype {
	case "map[]", "[map[]]": // response vacio
		break
	default:
		err = formatdata.FillStruct(respuesta["Data"], &v)
		respuesta = nil
	}
	return err
}

func JsonDebug(i interface{}) {
	formatdata.JsonPrint(i)
	fmt.Println()
}

func iguales(a interface{}, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func SendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	req, _ := http.NewRequest(trequest, url, b)
	req.Header.Set("Accept", AppJson)
	req.Header.Add("Content-Type", AppJson)
	seg2 := xray.BeginSegmentSec(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	xray.UpdateSegment(resp, err, seg2)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return json.NewDecoder(resp.Body).Decode(target)
}

func GetJsonTest(url string, target interface{}) (status int, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	seg2 := xray.BeginSegmentSec(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	xray.UpdateSegment(resp, err, seg2)
	if err != nil {
		return resp.StatusCode, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(target)
}

func GetJson(url string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	seg2 := xray.BeginSegmentSec(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	xray.UpdateSegment(resp, err, seg2)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return json.NewDecoder(resp.Body).Decode(target)
}

func GetXml(url string, target interface{}) error {

	req, _ := http.NewRequest("GET", url, nil)
	seg2 := xray.BeginSegmentSec(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	xray.UpdateSegment(resp, err, seg2)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return xml.NewDecoder(resp.Body).Decode(target)
}

func GetJsonWSO2(urlp string, target interface{}) error {
	b := new(bytes.Buffer)
	req, _ := http.NewRequest("GET", urlp, b)
	req.Header.Set("Accept", AppJson)
	seg2 := xray.BeginSegmentSec(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	xray.UpdateSegment(resp, err, seg2)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return json.NewDecoder(resp.Body).Decode(target)
}

func GetJsonWSO2Test(urlp string, target interface{}) (status int, err error) {

	b := new(bytes.Buffer)
	req, _ := http.NewRequest("GET", urlp, b)
	req.Header.Set("Accept", AppJson)
	seg2 := xray.BeginSegmentSec(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	xray.UpdateSegment(resp, err, seg2)
	if err != nil {
		beego.Error("error", err)
		return resp.StatusCode, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			beego.Error(nil, err)
		}
	}()
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(target)
}

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	oneDay := time.Hour * 5
	a = a.Add(oneDay)
	b = b.Add(oneDay)
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	// Normalize negative values

	if day < 0 {
		// days in month: p
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func FormatMoney(value interface{}, Precision int) string {
	formattedNumber := FormatNumber(value, Precision, ",", ".")
	return FormatMoneyString(formattedNumber, Precision)
}

func FormatMoneyString(formattedNumber string, Precision int) string {
	var format string

	zero := "0"
	if Precision > 0 {
		zero += "." + strings.Repeat("0", Precision)
	}

	format = "%s%v"
	result := strings.Replace(format, "%s", "$", -1)
	result = strings.Replace(result, "%v", formattedNumber, -1)

	return result
}

func FormatNumber(value interface{}, precision int, thousand string, decimal string) string {
	v := reflect.ValueOf(value)
	var x string
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x = fmt.Sprintf("%d", v.Int())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x = fmt.Sprintf("%d", v.Uint())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Float32, reflect.Float64:
		x = fmt.Sprintf(fmt.Sprintf("%%.%df", precision), v.Float())
	case reflect.Ptr:
		switch v.Type().String() {
		case "*big.Rat":
			x = value.(*big.Rat).FloatString(precision)

		default:
			panic("Unsupported type - " + v.Type().String())
		}
	default:
		panic("Unsupported type - " + v.Kind().String())
	}

	return formatNumberString(x, precision, thousand, decimal)
}

func formatNumberString(x string, precision int, thousand string, decimal string) string {
	lastIndex := strings.Index(x, ".") - 1
	if lastIndex < 0 {
		lastIndex = len(x) - 1
	}

	var buffer []byte
	var strBuffer bytes.Buffer

	j := 0
	for i := lastIndex; i >= 0; i-- {
		j++
		buffer = append(buffer, x[i])

		if j == 3 && i > 0 && !(i == 1 && x[0] == '-') {
			buffer = append(buffer, ',')
			j = 0
		}
	}

	for i := len(buffer) - 1; i >= 0; i-- {
		strBuffer.WriteByte(buffer[i])
	}
	result := strBuffer.String()

	if thousand != "," {
		result = strings.Replace(result, ",", thousand, -1)
	}

	extra := x[lastIndex+1:]
	if decimal != "." {
		extra = strings.Replace(extra, ".", decimal, 1)
	}

	return result + extra
}

// Valida que el body recibido en la petición tenga contenido válido
func ValidarBody(body []byte) (valid bool, err error) {
	var test interface{}
	if err = json.Unmarshal(body, &test); err != nil {
		return false, err
	} else {
		content := fmt.Sprintf("%v", test)
		fmt.Println(content)
		switch content {
		case "map[]", "[map[]]": // body vacio
			return false, nil
		}
	}
	return true, nil
}

// Quita el formato de moneda a un string y lo convierte en valor flotante
func DeformatNumber(formatted string) (number float64) {
	formatted = strings.ReplaceAll(formatted, ",", "")
	formatted = strings.Trim(formatted, "$")
	number, _ = strconv.ParseFloat(formatted, 64)
	return
}

// Obtiene los datos del usuario autenticado
func GetUsuario(usuario string) (nombreUsuario map[string]interface{}, err error) {
	if len(usuario) > 0 {
		var decData map[string]interface{}
		if data, err6 := base64.StdEncoding.DecodeString(usuario); err6 != nil {
			return nombreUsuario, err6
		} else {
			if err7 := json.Unmarshal(data, &decData); err7 != nil {
				return nombreUsuario, err7
			}
		}
		nombreUsuario = decData["user"].(map[string]interface{})
	}
	return nombreUsuario, err
}

// Manejo único de errores para controladores sin repetir código
func ErrorController(c beego.Controller, controller string) {
	if err := recover(); err != nil {
		logs.Error(err)
		localError := err.(map[string]interface{})
		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + controller + "/" + (localError["funcion"]).(string))
		c.Data["data"] = (localError["err"])
		xray.EndSegmentErr(http.StatusBadRequest, localError["err"])
		if status, ok := localError["status"]; ok {
			c.Abort(status.(string))
		} else {
			c.Abort("500")
		}
	}
}
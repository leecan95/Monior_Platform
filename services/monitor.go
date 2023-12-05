package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
)

func MonitorKpiModule(db *sqlx.DB) error {
	//var response PodReponse
	//var sum []config.KpiData
	//url := config.PrometheusUrl
	//params := "?query=sum(users_api_request_error_count)"
	//resp, err := http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//evalue, eok := response.Data.Result[0].Value[1].(string)
	//
	//params = "?query=sum(users_api_request_count)"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//
	//value, ok := response.Data.Result[0].Value[1].(string)
	//if ok && eok {
	//	data := config.KpiData{
	//		Pod:   "user",
	//		Req:   value,
	//		Error: evalue,
	//	}
	//	sum = append(sum, data)
	//}
	//
	//params = "?query=sum(vtdevices_api_request_error_count)"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//
	//evalue, eok = response.Data.Result[0].Value[1].(string)
	//
	//params = "?query=sum(vtdevices_api_request_count)"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//value, ok = response.Data.Result[0].Value[1].(string)
	//if ok && eok {
	//	data := config.KpiData{
	//		Pod:   "device",
	//		Req:   value,
	//		Error: evalue,
	//	}
	//	sum = append(sum, data)
	//}
	//params = "?query=sum(organizations_api_request_error_count)"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//
	//evalue, eok = response.Data.Result[0].Value[1].(string)
	//
	//params = "?query=sum(organizations_api_request_count)"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//value, ok = response.Data.Result[0].Value[1].(string)
	//if ok && eok {
	//	data := config.KpiData{
	//		Pod:   "organization",
	//		Req:   value,
	//		Error: evalue,
	//	}
	//	sum = append(sum, data)
	//}
	//
	//params = "?query=sum(attributes_api_request_error_count)"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//
	//evalue, eok = response.Data.Result[0].Value[1].(string)
	//
	//params = "?query=sum(attributes_api_request_count)"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	return err
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	return err
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	return err
	//}
	//value, ok = response.Data.Result[0].Value[1].(string)
	//if ok && eok {
	//	data := config.KpiData{
	//		Pod:   "attribute",
	//		Req:   value,
	//		Error: evalue,
	//	}
	//	sum = append(sum, data)
	//}
	model.KpiModule(db, nil)

	return nil
}

func MonitorKpiApi(db *sqlx.DB) error {
	var response PodReponse
	var sum []config.KpiData
	url := config.PrometheusUrl
	params := "?query=users_api_request_error_count{method=\"user_vtracking_login\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}
	evalue, eok := response.Data.Result[0].Value[1].(string)

	params = "?query=users_api_request_count{method=\"user_vtracking_login\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}

	value, ok := response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := config.KpiData{
			Pod:    "user",
			Method: "user_vtracking_login",
			Req:    value,
			Error:  evalue,
		}
		sum = append(sum, data)
	}

	params = "?query=users_api_request_error_count{method=\"get_list_users\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}

	evalue, eok = response.Data.Result[0].Value[1].(string)

	params = "?query=users_api_request_count{method=\"get_list_users\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}
	value, ok = response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := config.KpiData{
			Pod:    "users",
			Method: "get_list_user",
			Req:    value,
			Error:  evalue,
		}
		sum = append(sum, data)
	}
	params = "?query=vtdevices_api_request_error_count{method=\"GetFilteredVehicles\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}

	evalue, eok = response.Data.Result[0].Value[1].(string)

	params = "?query=vtdevices_api_request_count{method=\"GetFilteredVehicles\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}
	value, ok = response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := config.KpiData{
			Pod:    "devices",
			Method: "Get Vehicles",
			Req:    value,
			Error:  evalue,
		}
		sum = append(sum, data)
	}

	params = "?query=vtdevices_api_request_error_count{method=\"GetFilteredDevices\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}

	evalue, eok = response.Data.Result[0].Value[1].(string)

	params = "?query=vtdevices_api_request_count{method=\"GetFilteredDevices\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}
	value, ok = response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := config.KpiData{
			Pod:    "attribute",
			Method: "Get List Devices",
			Req:    value,
			Error:  evalue,
		}
		sum = append(sum, data)
	}

	params = "?query=organizations_api_request_error_count{method=\"get_tree_org_ids\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}

	evalue, eok = response.Data.Result[0].Value[1].(string)

	params = "?query=organizations_api_request_count{method=\"get_tree_org_ids\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		return err
	}
	value, ok = response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := config.KpiData{
			Pod:    "organizations",
			Method: "Get list org",
			Req:    value,
			Error:  evalue,
		}
		sum = append(sum, data)
	}
	model.KpiModule(db, nil)
	return nil
}

package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type KafkaReponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
type Systemreponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
				Device   string `json:"device"`
				Job      string `json:"job"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type CmdReponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Cmd string `json:"cmd"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type KafkaData struct {
	ResultType string `json:"resultType"`
	Result     []struct {
		Metric struct {
			Instance string `json:"instance"`
		} `json:"metric"`
		Value []interface{} `json:"value"`
	} `json:"result"`
}

type PodReponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Pod    string `json:"pod"`
				Method string `json:"method"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type PodData struct {
	Pod   string `json:"pod"`
	Value string `json:"value"`
}

type UrlData struct {
	Url   string `json:"url"`
	Value string `json:"value"`
}

type CmdData struct {
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
}

type SysData struct {
	Device string `json:"device"`
	Url    string `json:"url"`
	Value  string `json:"value"`
}

type KpiData struct {
	Pod   string `json:"pod"`
	Req   string `json:"request"`
	Error string `json:"error"`
}

func SuccessResponse(c *gin.Context, status int, message interface{}, data interface{}) {
	c.AbortWithStatusJSON(status, gin.H{"message": message, "data": data})
}
func Success(c *gin.Context, message interface{}, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

func ErrorResponse(c *gin.Context, status int, message interface{}, error interface{}) {
	c.AbortWithStatusJSON(status, gin.H{"message": message, "errors": error})
}
func ValidationError(c *gin.Context, message interface{}, error interface{}) {
	ErrorResponse(c, http.StatusUnprocessableEntity, message, error)
}

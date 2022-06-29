package elastic

import (
	"encoding/json"
	"github.com/oaago/cloud/config"
	elastics "github.com/olivere/elastic/v7"
)

type elasticType struct {
	*elastics.Client
}

type ElasticType struct {
	URL      string `yaml:"url"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
}

var Elastic = &elasticType{}
var ElasticOptions = &ElasticType{}

func init() {
	elasticStr := config.Op.GetString("elastic")
	if len(elasticStr) > 0 {
		json.Unmarshal([]byte(elasticStr), ElasticOptions)
	}
}

func NewClient() *elasticType {
	var errConnect error
	Elastic.Client, errConnect = elastics.NewClient(
		elastics.SetURL(ElasticOptions.URL),
		elastics.SetGzip(true),
		elastics.SetSniff(false),
		elastics.SetBasicAuth(ElasticOptions.UserName, ElasticOptions.Password),
	)
	if errConnect != nil {
		panic(errConnect)
	}
	return Elastic
}

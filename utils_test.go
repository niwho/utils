package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitEsClient([]string{"http://10.116.27.15:9600"}, "test_wt_idx", "test_type", 2, 0, 60)
	os.Exit(m.Run())
}

func TestInsert2Es(t *testing.T) {
	esversion, err := EsClient.ElasticsearchVersion("http://10.116.27.15:9600")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	err = Insert2Es(`{"a":1, "b":2, "c":"lala"}`)
	fmt.Printf("Insert2Es err=", err)
}

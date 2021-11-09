package elastic

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

type Index struct {
	client    *elasticsearch.Client
	indexName string
}

func NewCreateIndex(addres []string) *Index {
	elastic, err := elasticsearch.NewClient(elasticsearch.Config{
		
		Addresses: addres,
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	return &Index{
		client: elastic,
	}

}

func (i *Index) CreateIndex(index_name string) error {
	index_name = fmt.Sprintf("funding_%s", index_name)
	i.indexName = index_name
	res, err := i.client.Indices.Exists([]string{index_name})

	if err != nil {
		return err
	}

	if res.StatusCode != 404 {
		return fmt.Errorf("error in index existence response %s", res.String())
	}

	res, err = i.client.Indices.Create(i.indexName)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf(res.String())
	}

	return nil

}
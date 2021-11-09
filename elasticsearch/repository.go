package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"funding/model"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type repository struct {
	elastic Index
	timeout time.Duration
}

func NewElasticRepo(elastic Index, timeduration time.Duration) *repository {
	return &repository{elastic: elastic, timeout: timeduration}
}

type Storage struct {
	Source interface{} `json:"_source"`
}

type ElasticRepository interface {
	CreateCampaign(ctx context.Context, campaignye model.Campaign) error
	UpdateCampaign(ctx context.Context, campanye model.Campaign) (model.Campaign, error)
}

func (r *repository) CreateCampaign(ctx context.Context, campaignye model.Campaign) error {
	reqBody, err := json.Marshal(campaignye)
	if err != nil {
		return err
	}

	req := esapi.CreateRequest{
		Index:      r.elastic.indexName,
		DocumentID: fmt.Sprint(campaignye.ID),
		Body:       bytes.NewBuffer(reqBody),
	}

	ctx1, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := req.Do(ctx1, r.elastic.client)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 400 {
		return fmt.Errorf("status code is not 400")
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	fmt.Println(res.Body)
	return nil

}

func (r *repository) UpdateCampaign(ctx context.Context, campanye model.Campaign) (model.Campaign, error) {
	reqBody, err := json.Marshal(campanye)
	if err != nil {
		return campanye, err
	}

	req := esapi.UpdateRequest{
		Index:      r.elastic.indexName,
		DocumentID: fmt.Sprint(campanye.ID),
		Body:       bytes.NewBuffer(reqBody),
	}

	res, err := req.Do(ctx, r.elastic.client)
	if err != nil {
		return campanye, err
	}

	defer res.Body.Close()
	if res.StatusCode == 404 {
		return campanye, errors.New("value not found")
	}

	if res.IsError() {
		return campanye, fmt.Errorf("update: response: %s", res.String())
	}

	var (
		body    model.Campaign
		storage Storage
	)

	storage.Source = &body

	if err := json.NewDecoder(res.Body).Decode(&storage); err != nil {
		return campanye, fmt.Errorf("find one: decode: %w", err)
	}

	return body, nil
}

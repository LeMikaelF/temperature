package metrics

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"log"
	"time"
)

type Exporter struct {
	metricsChan  <-chan float64
	client       *cloudwatch.Client
	lastReceived *float64
}

func NewExporter(metricsChan <-chan float64) *Exporter {
	return &Exporter{metricsChan: metricsChan, lastReceived: aws.Float64(0)}
}

func (e Exporter) Start() error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("loading config: %v", err)
	}

	e.client = cloudwatch.NewFromConfig(cfg)

	tick := time.Tick(10 * time.Second)

	for {
		select {
		case _ = <-tick:
			var metrics []float64
			for len(e.metricsChan) > 0 {
				metric := <-e.metricsChan
				metrics = append(metrics, metric)
				*e.lastReceived = metric
			}
			if len(metrics) > 0 {
				e.pushMetrics(metrics)
			}
		}
	}
}

func (e Exporter) LastReceived() float64 {
	return *e.lastReceived
}

func (e Exporter) pushMetrics(metrics []float64) {
	log.Printf("pushing metrics %v", metrics)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	_, err := e.client.PutMetricData(ctx,
		&cloudwatch.PutMetricDataInput{
			MetricData: []types.MetricDatum{{
				MetricName: aws.String("temperature"),
				Dimensions: []types.Dimension{{
					Name:  aws.String("room"),
					Value: aws.String("living room"),
				}},
				Timestamp: aws.Time(time.Now()),
				Values:    metrics,
			}},
			Namespace: aws.String("Buildings"),
		})
	if err != nil {
		log.Printf("pushMetrics: %v", err)
	}
}

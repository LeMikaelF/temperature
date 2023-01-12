package metrics

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"log"
	"time"
)

type Exporter struct {
	metricsChan <-chan float64
	client      *cloudwatch.Client
	lastPushed  *float64
}

func NewExporter(metricsChan <-chan float64) *Exporter {
	return &Exporter{metricsChan: metricsChan, lastPushed: aws.Float64(0)}
}

func (e Exporter) Start() {
	e.client = cloudwatch.New(cloudwatch.Options{})

	tick := time.Tick(1 * time.Minute)

	for {
		select {
		case _ = <-tick:
			var metrics []float64
			for metric := range e.metricsChan {
				*e.lastPushed = metric
				metrics = append(metrics, metric)
			}
			e.pushMetrics(metrics)
		}
	}
}

func (e Exporter) LastPushed() float64 {
	if e.lastPushed == nil {
		return 0
	}
	return *e.lastPushed
}

func (e Exporter) pushMetrics(metrics []float64) {
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
				StorageResolution: aws.Int32(60),
				Timestamp:         aws.Time(time.Now()),
				Values:            metrics,
			}},
			Namespace: aws.String("Buildings"),
		})
	if err != nil {
		log.Printf("pushMetrics: %v", err)
	}
}

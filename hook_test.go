package elastic_logrus

import (
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestHook(t *testing.T) {
	
	logger := setupLogger()
	logger.Info("All done")
	logger.WithFields(logrus.Fields{"test1" : "a" , "test2": "b" , "test3": "c"})
	
	
}

func BenchmarkHook(b *testing.B) {
	
	logger := setupLogger()
	
	for i := 0; i < b.N; i++ {
		logger.WithFields(logrus.Fields{"test1" : "a" , "test2": "b" , "test3": "c", "test4": i} )
	}
}


func setupLogger() *logrus.Logger {
	logger := logrus.New()
	client, err := elastic.NewClient(elastic.SetSniff(false),  elastic.SetURL("http://localhost:9200"))
	if err != nil {
		logger.WithError(err).Fatal("Failed to construct elasticsearch client")
	}
	
	// Create logger with 15 seconds flush interval
	hook, err := NewElasticHook(client, "localhost", logrus.DebugLevel, func() string {
		return fmt.Sprintf("%s-%s", "some-index", time.Now().Format("2006-01-02"))
	}, time.Second * 15)
	
	if err != nil {
		logger.WithError(err).Fatal("Failed to create elasticsearch hook for logger")
	}
	
	logger.Hooks.Add(hook)
	return logger

}
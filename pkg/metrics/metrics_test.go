package metrics

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestController_Run_StopsWhenContextIsCanceled(t *testing.T) {

	log := logrus.NewEntry(logrus.New())
	c := New(log)

	timeout := time.After(2 * time.Second)
	done := make(chan bool)
	go func() {
		if err := c.Run("127.0.0.1:0"); err != nil {
			t.Errorf("Controller.Run() error = %v", err)
		}
		done <- true
		c.Shutdown()
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}

}

package soul

import (
	"sync"
	"testing"
	"time"

	"github.com/soul-go/soul/context"
	"github.com/spf13/cobra"

	"github.com/soul-go/soul/engines"
	. "github.com/smartystreets/goconvey/convey"
)

type cronJobMock struct {
	wg     *sync.WaitGroup
	called bool
}

func (c *cronJobMock) Run() {
	c.called = true
	c.wg.Done()
}

func TestListen(t *testing.T) {
	Convey("Given a service with routes", t, func() {
		mock := &engines.Mock{}
		service := New("test")
		router := service.NewRouter(mock)
		router.Include(func(c *context.Context) {})
		router.Listen("logs.warning", func(c *context.Context) {})

		Convey("When the routes are listened", func() {
			service.listen(&cobra.Command{}, []string{"logs.warning"})
			time.Sleep(1 * time.Millisecond)

			Convey("The engine should be connected", func() {
				So(mock.IsConnected, ShouldBeTrue)
			})

			Convey("The listened routing key should be 'logs.warning'", func() {
				So(mock.RoutingKey, ShouldEqual, "logs.warning")
			})

			Convey("The number of handlers should be 2", func() {
				So(mock.Handlers, ShouldHaveLength, 2)
			})
		})
	})
}

func TestCronJob(t *testing.T) {
	Convey("Given a service with a cronjob", t, func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		mock := &cronJobMock{wg: wg}
		service := New("test")
		service.CronJob("test", "@every 2ms", mock.Run)

		Convey("When the cronjob starts", func() {
			service.cronjob(&cobra.Command{}, []string{"test"})
			wg.Wait()

			Convey("The cronjob should be called", func() {
				So(mock.called, ShouldBeTrue)
			})
		})
	})
}

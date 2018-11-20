package features

import (
	"os"
	"testing"

	"github.com/sofmon/wedit/builder"
	"github.com/sofmon/wedit/service"
)

func TestMain(m *testing.M) {

	err := service.LoadConfig("wedit_test.json")
	if err != nil {
		panic(err)
	}

	err = builder.LoadConfig("wedit_test.json")
	if err != nil {
		panic(err)
	}

	go func() {
		err := service.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	os.Exit(m.Run())

}

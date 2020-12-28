package imports

import (
	"os"
	"testing"

	"github.com/rinchsan/gosimports/internal/testenv"
)

func TestMain(m *testing.M) {
	testenv.ExitIfSmallMachine()
	os.Exit(m.Run())
}

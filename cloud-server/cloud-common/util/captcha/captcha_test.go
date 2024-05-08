package captcha

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(Tencent.Check("randStr", "ticket"))
}

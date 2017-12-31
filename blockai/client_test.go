package blockai

import (
	"fmt"
	"testing"
	"time"
)

func TestFmtImg(t *testing.T) {
	img := &Img{
		URL: "http://s.cn.bing.net/az/hprichbg/rb/CliffPalaceLuminara_ZH-CN10279459718_1920x1080.jpg",
	}
	img.FmtName()
	if img.Name != "Warmiapoland Zh Cn13324541925 1366x768" {
		t.Errorf("%+v", img)
	}
}

func TestAccount(t *testing.T) {
	name := fmt.Sprint(time.Now().Unix())
	email := name + "@yahoo.com"
	passwd := name + "pass"
	cli, err := NewClient(name, email, passwd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("account: %+v", cli)

	img := NewImg("http://s.cn.bing.net/az/hprichbg/rb/CliffPalaceLuminara_ZH-CN10279459718_1920x1080.jpg")

	err = cli.Upload(img)
	if err != nil {
		t.Fatal(err)
	}
	t.Error("...")
}

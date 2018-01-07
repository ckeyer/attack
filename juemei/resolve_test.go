package juemei

import (
	"testing"

	"github.com/ckeyer/ckeyer/httpcli"
	"gopkg.in/redis.v4"
)

func TestGetPage(t *testing.T) {
	return
	url := "http://www.juemei.com/mm/201612/6386.html"
	cli := httpcli.NewClient()
	doc, err := LoadDoc(cli, url)
	if err != nil {
		t.Fatal(err)
	}

	ret, err := Resolve(cli, doc)
	if err != nil {
		t.Fatal(err)
	}

	if ret.URL != url {
		t.Fatal("url not equie")
	}
}

func TestRedisStore(t *testing.T) {
	return
	url := "http://www.juemei.com/mm"
	cli := httpcli.NewClient()
	rcli := redis.NewClient(&redis.Options{
		Addr: "u5.mj:6379",
	})
	if err := rcli.Ping().Err(); err != nil {
		t.Fatal(err)
	}

	store := NewRedisStore(rcli)

	for i := 0; i < 10; i++ {
		t.Error(".....", url)
		doc, err := LoadDoc(cli, url)
		if err != nil {
			t.Fatal(err)
		}

		ret, err := Resolve(cli, doc)
		if err != nil {
			t.Fatal(err)
		}

		store.Done(ret)
		us, err := store.Jobs(1)
		if err != nil {
			t.Fatal(err)
		}
		if len(us) != 1 {
			t.Fatal("urls.length should be 1")
		}
		url = us[0]
	}

}

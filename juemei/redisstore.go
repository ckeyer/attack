package juemei

import (
	"gopkg.in/redis.v4"
)

const (
	RedisKeyPrefix = "juemei_" // redis key 前缀

	KeySetToDo     = RedisKeyPrefix + "todo"      // 未爬取的页面链接
	KeySetDoing    = RedisKeyPrefix + "doing"     // 正在处理的页面
	KeySetDone     = RedisKeyPrefix + "done"      // 已爬取完成的页面
	KeySetImgs     = RedisKeyPrefix + "imgs"      // 爬取的图片链接
	KeySetOutLinks = RedisKeyPrefix + "out_links" // 外链
)

type RedisStore struct {
	*redis.Client
}

func NewRedisStore(cli *redis.Client) *RedisStore {
	return &RedisStore{
		Client: cli,
	}
}

func (r *RedisStore) Init() error {
	for {
		Urls, _, err := r.SScan(KeySetDoing, 0, "", 100).Result()
		if err != nil {
			return err
		}
		if len(Urls) == 0 {
			return nil
		}
		for _, Url := range Urls {
			err = r.SMove(KeySetDoing, KeySetToDo, Url).Err()
			if err != nil {
				return err
			}
		}
	}
}

func (r *RedisStore) Done(rets ...*ResolveResult) (*StoreResult, error) {
	sret := &StoreResult{
		URLs: []string{},
	}
	for _, rr := range rets {
		_, err := r.setDone(rr.URL)
		if err != nil {
			return sret, err
		}
		sret.URLs = append(sret.URLs, rr.URL)

		imgN, err := r.addImgs(rr.Imgs)
		if err != nil {
			return sret, err
		}
		sret.Imgs += imgN

		linkN, err := r.addLinks(rr.Links)
		if err != nil {
			return sret, err
		}
		sret.Links += linkN

		outN, err := r.addOutLinks(rr.OutLinks)
		if err != nil {
			return sret, err
		}
		sret.OutLinks += outN
	}

	return sret, nil
}

func (r *RedisStore) setDone(done string) (bool, error) {
	return r.SMove(KeySetDoing, KeySetDone, done).Result()
}

func (r *RedisStore) addImgs(imgs []string) (int64, error) {
	return r.SAdd(KeySetImgs, stringsToInterface(imgs)...).Result()
}

func (r *RedisStore) addLinks(links []string) (int64, error) {
	var count int64
	for _, link := range links {
		// Note: this is not safety
		if r.SIsMember(KeySetDone, link).Val() || r.SIsMember(KeySetDoing, link).Val() {
			continue
		}

		n, err := r.SAdd(KeySetToDo, link).Result()
		if err != nil {
			return 0, err
		}
		count += n
	}

	return count, nil
}

func (r *RedisStore) addOutLinks(outs []string) (int64, error) {
	return r.SAdd(KeySetOutLinks, stringsToInterface(outs)...).Result()
}

func (r *RedisStore) Jobs(count int) ([]string, error) {
	ss, _, err := r.SScan(KeySetToDo, 0, "", int64(count)).Result()
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, Url := range ss {
		ok, err := r.SMove(KeySetToDo, KeySetDoing, Url).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			ret = append(ret, Url)
		}
	}

	return ret, nil
}

func stringsToInterface(ss []string) []interface{} {
	ret := []interface{}{}
	for _, s := range ss {
		ret = append(ret, s)
	}
	return ret
}

package Cache

func SetPropCache(typ, id string, data []byte) error {
	if typ == "" || id == "" || data == nil {
		return ErrInvalidParam
	}

	return RedisDB.Set(genPropCacheKey(typ, id), data, 0).Err()
}

func SetPropCacheList(typs, ids []string, datas [][]byte) error {
	if len(typs) == 0 {
		return nil
	}

	if len(typs) != len(ids) || len(typs) != len(datas) {
		return ErrInvalidParam
	}

	args := make([]interface{}, 0, 1)
	for i := range typs {
		if typs[i] == "" || ids[i] == "" || datas[i] == nil {
			return ErrInvalidParam
		}

		args = append(args, genPropCacheKey(typs[i], ids[i]))
		args = append(args, datas[i])
	}

	return RedisDB.MSet(args...).Err()
}

func GetPropCache(typ, id string) ([]byte, error) {
	if typ == "" || id == "" {
		return nil, ErrInvalidParam
	}

	return RedisDB.Get(genPropCacheKey(typ, id)).Bytes()
}

func GetPropCacheList(typs, ids []string) ([][]byte, error) {
	if len(typs) == 0 {
		return nil, ErrInvalidParam
	}

	if len(typs) != len(ids) {
		return nil, ErrInvalidParam
	}

	args := make([]string, 0, 1)
	for i := range typs {
		args = append(args, genPropCacheKey(typs[i], ids[i]))
	}

	results, err := RedisDB.MGet(args...).Result()
	if err != nil {
		return nil, err
	}

	datas := make([][]byte, 0, 1)
	for _, r := range results {
		if r == nil {
			datas = append(datas, nil)
		} else {
			datas = append(datas, []byte(r.(string)))
		}
	}

	return datas, nil
}

func genPropCacheKey(typ, id string) string {
	return "PropCache:" + typ + ":" + id
}

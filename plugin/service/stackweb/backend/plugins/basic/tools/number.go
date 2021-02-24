package tools

func Int32ArrayTo64(in []int32) []int64 {
	ret := make([]int64, len(in))
	for _, v := range in {
		ret = append(ret, int64(v))
	}

	return ret
}

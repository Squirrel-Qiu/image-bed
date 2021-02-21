package id

//type Generate struct {
//	IdType string
//	IdList []string
//}
//
//func Generator(g *Generate) (idList []string, err error) {
//	var db dbb.DBApi
//	idValue, err := db.GetIdValue(g.IdType) // 原子操作
//	if err != nil {
//		return nil, err
//	}
//
//	for i := idValue-10; i <= idValue; i++ {
//		buff := make([]byte, 8)
//		binary.BigEndian.PutUint64(buff, uint64(i))
//		m := fmt.Sprintf("%x", md5.Sum(buff))
//		g.IdList = append(g.IdList, m[:10])
//	}
//	return g.IdList, nil
//}

package dao

import (
	C "com.mutantcat.cloud_step/collection"
	"com.mutantcat.cloud_step/entity"
)

func GetAllCollections() []entity.Collection {
	collections := make([]entity.Collection, 0)
	err := PublicEngine.Find(&collections)
	if err != nil {
		return nil
	}
	return collections
}

// 检查集合名是否存在
func CheckCollectionNameExist(name string) bool {
	count, err := PublicEngine.Where("name = ?", name).Count(&entity.Collection{})
	if err != nil {
		return false
	}
	return count > 0
}

func AddCollection(collectionName string, urls []entity.Url) (bool, int) {
	// 开启事务
	session := PublicEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return false, 0
	}
	collection := entity.Collection{}
	collection.Name = collectionName
	_, err = session.Insert(&collection)
	if err != nil {
		session.Rollback()
		return false, 0
	}
	for i, url := range urls {
		url.Parent = collectionName
		url.Alive = true
		url.Retry = 0
		_, err = session.Insert(&url)
		if err != nil {
			session.Rollback()
			return false, 0
		}
		urls[i].Id = url.Id
	}
	err = session.Commit()
	if err != nil {
		session.Rollback()
		return false, 0
	}
	// 数据库中添加无误之后 添加进缓存中
	C.MWorkCllection.Lock()
	defer C.MWorkCllection.Unlock()
	C.WorkCllection[collection.Name] = urls

	// 获得刚才添加的集合的id
	collection = entity.Collection{}
	_, err = PublicEngine.Where("name = ?", collectionName).Get(&collection)
	if err != nil {
		return false, 0
	}

	return true, collection.Id
}

// 删除集合 （还需要删除parent是这个集合的url）
func DeleteCollectionById(id int) bool {
	session := PublicEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return false
	}
	collection := entity.Collection{}
	_, err = session.ID(id).Get(&collection)
	if err != nil {
		session.Rollback()
		return false
	}
	_, err = session.ID(id).Delete(&collection)
	if err != nil {
		session.Rollback()
		return false
	}
	_, err = session.Where("parent = ?", collection.Name).Delete(&entity.Url{})
	if err != nil {
		session.Rollback()
		return false
	}
	err = session.Commit()
	if err != nil {
		session.Rollback()
		return false
	}
	// 删除缓存中的集合
	C.MWorkCllection.Lock()
	defer C.MWorkCllection.Unlock()
	delete(C.WorkCllection, collection.Name)
	return true
}

func GetCollectionNameById(id int) string {
	collection := entity.Collection{}
	_, err := PublicEngine.ID(id).Get(&collection)
	if err != nil {
		return ""
	}
	return collection.Name
}

// 判断是否有自助模式或者代理模式依赖这个集合
func CheckCollectionDepend(id int) bool {
	name := GetCollectionNameById(id)
	proxy := entity.Proxy{}
	has, err := PublicEngine.Where("point = ?", name).Get(&proxy)
	if err != nil {
		return true
	}
	if has {
		return true
	}
	selfHelp := entity.SelfHelp{}
	has, err = PublicEngine.Where("point = ?", name).Get(&selfHelp)
	if err != nil {
		return true
	}
	if has {
		return true
	}
	return false
}

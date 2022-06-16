集成官方mongo驱动， 并且库中提供了基本的客户端管理。

---

### 使用示例

```
import (
    mongoOfficial "github.com/yunxiaoyang01/open_sdk/mongoc_official"
)

// 库中提供了基本的客户端管理。
// 新增一个mongo客户端
mongoOfficial.AddClientAddress(ctx, "", conf.MongoURI)

// 原始官方驱动相比mgo不那么友好，为了更容易使用，封装了驱动的原始方法(mongoOfficial.Base)。
// 请基于以下封装来实现
type UserMongo struct {
	*mongoOfficial.Base
}

func (db *UserMongo) Find(ctx context.Context, name string) (*user.UserInfo, error) {
	user := &user.UserInfo{}
	err := db.FindOne(ctx, map[string]interface{}{"name": name}, user)
	return user, err
}

func (db *UserMongo) Insert(ctx context.Context, user *user.UserInfo) error {
	_, err := db.InsertOne(ctx, user)
	return err
}

func (db *TUserMongo) Upsert(ctx context.Context, filter bson.D, user *user.UserInfo) error {
	_, err := db.Base.Upsert(ctx, filter, user)
	return err
}

UserMongo = &UserMongo{Base: mongoOfficial.NewBaseModel("", "testdb", "user")}


```

- 链接URI  
mongodb://username:password@example1.com,example2.com,example3.com/?replicaSet=test&w=majority&wtimeoutMS=5000

https://docs.mongodb.com/manual/reference/connection-string/

| 参数             | 描述                                                                         |
| ---------------- | ---------------------------------------------------------------------------- |
| replicaSet       | repset name                                                                  |
| appName          | 建立连接的时候指定一个标签，可用于筛选日志 & 慢日志分析                      |
| w                | write concern, https://docs.mongodb.com/manual/reference/write-concern/#wc-w |
| readConcernLevel | read concern, https://docs.mongodb.com/manual/reference/read-concern/        |
| compressors      | 数据压缩， zstd(>=4.2 & enable cgo) ,zlib (>= 3.6), snappy(>= 3.4)           |
| connectTimeoutMS | 毫秒，连接超时                                                               |
| socketTimeoutMS  | 毫秒，读写超时                                                               |
| maxPoolSize      | 连接池最大连接数， default: 100                                              |
| minPoolSize      | 连接池最小连接数， default: 0                                                |
| maxIdleTimeMS    | 连接空闲多长时间被回收                                                       |
| readPreference   | 读优先级，https://docs.mongodb.com/manual/core/read-preference/              |

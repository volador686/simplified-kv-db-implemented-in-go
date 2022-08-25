simplified key-value database implemented in go:

介绍文档

一、主要功能：

一个简化的使用go语言实现的分布式的key-value数据库

可以对数据库中新添加的数据增量更新到旧的备份中

主要结构：

a)config 用于定义配置文件结构体

b)db用于创建与数据库的连接，定义数据库相关功能函数

c)web 用于处理客户请求，对相关功能请求进行回应

二、相关包/依赖项：
a)net/http:（相关使用）

i.功能监听

```
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	
	r.ParseForm()
	key := r.Form.Get("key")

	shard := s.getShard(key)
	value, err := s.db.GetKey(key)

	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "shard = %d, current shard = %d, addr = %q, value = %q, error = %v\n", shard, s.shardIdx, s.addrs[shard], value, err)
}


func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	shard := s.getShard(key)

	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}
	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "error = %v, shardIdx = %d, current shard = %d\n", err, shard, s.shardIdx)
}

func (s *Server) ReplicaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	replica_addr := r.Form.Get("addr")
	s.db.SendReplica(replica_addr)
	fmt.Fprintf(w, "replica_addr = %v\n", replica_addr)
}
```
b)boltDB:（数据库）
i.读写功能的实现

1.实现过程

ii.增量更新的实现

1.实现过程

通过创建一个replicationbucket来存储需要进行复制的key-value对，当接受到进行复制的指令时，向对应的socket发送key-value对，并且删除已经发送的key-value对。

func (d *Database) SendReplica(replica_addr string) (string, string) {
	var res1 []byte
	var res2 []byte
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(replicaBucket)
		c := b.Cursor()
		res1, res2 = c.First()
		for key, value := c.First(); key != nil; key, value = c.Next() {
			log.Printf("key = %v, value = %v\n", string(key), string(value))
			cli := goz.NewClient()
			url := "http://" + replica_addr + "/set"
			resp, err := cli.Post(string(url), goz.Options{
				FormParams: map[string]interface{}{
					"key":   string(key),
					"value": string(value),
				},
			})
			if err != nil {
				log.Printf("replication error:%v\n", err)
			}
			body, _ := resp.GetBody()
			log.Println(body)
		}
		// for key, _ := c.First(); key != nil; key, _ = c.Next() {
		// 	b := tx.Bucket(replicaBucket)
		// 	b.Delete(key)
		// }
		return nil
	})
	return string(res1), string(res2)
}
2.通过这种方式实现增量更新的问题

a)问题阐述

i.以当前方式无法实现对多个数据库进行增量更新

ii.可扩展性差

b)改进方式（Todo）

i.更改增量更新的变量，使其可以接受多个指定socket作为目的数据库

ii.可以对defaultbucket执行复制但是不删除，可以使指定数据库中的内容与目标数据库一致

c)hash:（哈希）

通过对提供的key进行哈希操作，从而确定对应相关操作的数据库，同时将用户提供的相关信息和操作类型重定向到相应数据库

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q)\n", s.shardIdx, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}
三、随机读写测试：
a)测试结果
i.随机写

goos: linux

goarch: amd64

pkg: modules/benchmark

cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

BenchmarkWrite-12  49  24769037 ns/op  14965 B/op   116 allocs/op


ii.随机读
goos: linux

goarch: amd64

pkg: modules/benchmark

cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

BenchmarkRead-12  4815 233616 ns/op  14687 B/op	 112 allocs/op

iii.结果分析
从单位时间内对随机读写的次数比较可以得知

b)分析读写差距原因

i.程序设计中所采用的boltDB数据库对读写方式的控制方式不同：

在boltDB中，读操作支持并发，写操作不支持

ii.每次进行写操作和读操作时，未采用缓存，造成速度下降

c)可以改进的方式

i.添加缓存（Todo）

可以规定需要缓存的数据超过一定大小时，再将缓存中的数据写入数据库，从而平缓数据库的写入。

但是，会因此增加程序复杂度，并且相较于不采用缓存，如果程序意外终止，缓存中的数据可能因此丢失，不同数据库的写入顺序/速度可能不同，还会因此造成数据库的不一致

ii.增加数据库的数量

在测试中数据库的数量为4（不包括备份数据库），在面对大量的数据时，可以适当增加数据库的数量，从而降低单一数据库写操作的数量，进而提升数据库的整体性能

iii.换用读写速度更高的数据库

可以换用读写速度更高的数据库（redis等）

四、结语

a)参考

i.boltDB -- https://github.com/boltdb/bolt

ii.go-curl -- https://github.com/idoubi/goz

iii.implementation -- https://www.youtube.com/watch?v=EdPkmJrtTWQ


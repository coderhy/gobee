package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gobee/pkg/common"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus" //定义为log 方便后面可插拔
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoModel 模型结构
type MongoModel struct {
	Client     *mongo.Database //数据库DB
	Connection string          //yaml配置
	PoolID     string          //连接池数组（dbLinkPool） 里面的key,为了删除使用

	// TablePrefix    string            //表前缀
	TableName      string            //表名
	SubTableMod    []int64           //分表模型
	Fields         []string          //表字段
	WhereFields    []string          //表字段
	TransferFields map[string]string //新-老字段流转关系
	TransferOrders map[string]string //新-老排序流转关系

	Transaction *sql.Tx   //事务
	Options     MgOptions //参数选项
	SQL         string    //当前SQL
	DbError     error     //错误
	Debug       bool
}

// MgOptions 参数选项结构
type MgOptions struct {
	Distinct bool        //是否去重
	Fields   []string    //字段
	Force    string      //强制使用索引
	Where    interface{} //参数表达式

	Order  string //排序
	Join   string //JOIN 需要完善
	Group  string //分组
	Having string //弥补了WHERE关键字不能与聚合函数联合使用的不足
	Union  string //选取不同的值 没有重复的  需要完善
	// Lock    string //锁表
	Replace bool //是否替换插入
	Offset  int  //从第N个开始
	Limit   int  //获取X条数据
}

// db连接池
var mongoLinkPool = map[string]interface{}{}

// 初始化链接  => ok
func getMgConnect(option map[interface{}]interface{}) (*mongo.Database, string, error) {
	dbDSN := fmt.Sprintf("mongodb://%s:%s@%s:%s", option["DbUser"], option["DbPwd"], option["DbHost"], option["DbPort"])
	// dbDSN := fmt.Sprintf("mongodb://root:123456@localhost:27017")
	// fmt.Println("dbDSN", dbDSN)
	PoolID := common.Md5(dbDSN)
	if db, ok := mongoLinkPool[PoolID]; ok {
		return db.(*mongo.Database), PoolID, nil
	}
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(dbDSN)
	clientOptions.SetMaxPoolSize(uint64(option["MaxOpenConns"].(int))) //用于设置最大打开的连接数，默认值为100
	clientOptions.SetMinPoolSize(uint64(option["MaxIdleConns"].(int))) //用于设置最小打开的连接数，默认值为0
	// clientOptions.SetMaxConnIdleTime(10000) //指定连接池中连接保持空闲的最长时间 分钟
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err.Error()}).Warn("mongo link error")
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err.Error()}).Warn("mongo link error1")
	}
	db := client.Database(option["DbName"].(string))
	mongoLinkPool[PoolID] = db
	return db, PoolID, err
}

// NewMongo 初始化 OK
func NewMongo(modelType map[string]interface{}) MongoModel {
	setModel := MongoModel{}
	dbConfig := GetConfig("mongo")
	// options := dbConfig.Get("MG_ALL_DEMO")
	options := dbConfig.Get(modelType["connection"].(string))
	// fmt.Println("options", options)
	if options != nil {
		options := options.(map[interface{}]interface{})
		// Client, PoolID, DbError := getConnect(options) //连接
		Client, PoolID, DbError := getMgConnect(options)

		setModel.DbError = DbError
		setModel.Client = Client
		setModel.PoolID = PoolID
		if _, ok := options["Debug"]; ok {
			setModel.Debug = options["Debug"].(bool)
		}
	} else {
		setModel.DbError = errors.New("配置不存在或异常")
		// log.WithFields(log.Fields{"isPush": true, "error": modelType["connection"].(string) + " 配置不存在或异常"}).Warn("mongo error")
	}

	//方式1：map转成(映射到)struct结构
	mapstructure.Decode(modelType, &setModel)

	// fmt.Println("setModel ", setModel.Debug)
	//方式2：性能优化 start
	/* if _, ok := modelType["connection"]; ok {
		setModel.Connection = modelType["connection"].(string) //数据库DB配置
	} else {
		setModel.DbError = errors.New("connection异常")
	}
	if _, ok := modelType["tableName"]; ok {
		setModel.TableName = modelType["tableName"].(string) //表名
	} else {
		setModel.DbError = errors.New("tableName异常")
	}
	if _, ok := modelType["subTableMod"]; ok {
		setModel.SubTableMod = modelType["subTableMod"].([]int64) //分表模型
	}

	if _, ok := modelType["fields"]; ok {
		setModel.Fields = modelType["fields"].([]string) //表字段
	} else {
		setModel.DbError = errors.New("fields异常")
	}
	if _, ok := modelType["whereFields"]; ok {
		setModel.WhereFields = modelType["whereFields"].([]string) //表字段
	} else {
		setModel.DbError = errors.New("whereFields异常")
	}
	if _, ok := modelType["transferFields"]; ok {
		setModel.TransferFields = modelType["transferFields"].(map[string]string) //新-老字段流转关系
	} else {
		setModel.DbError = errors.New("transferFields异常")
	}
	if _, ok := modelType["transferOrders"]; ok {
		setModel.TransferOrders = modelType["transferOrders"].(map[string]string) //新-老排序流转关系
	} else {
		setModel.DbError = errors.New("transferOrders异常")
	} */
	//方式2：性能优化 end

	return setModel
}

// Close 关闭释放
func (m *MongoModel) Close() error {
	if _, ok := mongoLinkPool[m.PoolID]; ok {
		delete(mongoLinkPool, m.PoolID)
	}
	return m.Client.Client().Disconnect(context.TODO())
}

// Add 添加数据  => ok (replace = false)
func (m *MongoModel) Add(data interface{}) (interface{}, error) {
	insertResult, err := m.Client.Collection(m.TableName).InsertOne(context.TODO(), data)
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo add error")
	}
	if m.Debug { //组装sql debug用
		strjson, _ := json.Marshal(data)
		_sql := "insert into " + m.TableName + " " + string(strjson)
		m.SQL = _sql
	}
	return insertResult.InsertedID, err
}

// AddAll 添加数据  => ok (replace = false)
func (m *MongoModel) AddAll(data []interface{}) (interface{}, error) {
	insertResult, err := m.Client.Collection(m.TableName).InsertMany(context.TODO(), data)
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo addall error")
	}
	if m.Debug { //组装sql debug用
		strjson, _ := json.Marshal(data)
		_sql := "insert into " + m.TableName + " " + string(strjson)
		m.SQL = _sql
	}
	return insertResult.InsertedIDs, err
}

//Save 保存、更新数据
func (m *MongoModel) Save(update map[string]interface{}) (int64, error) {
	var (
		updateResult *mongo.UpdateResult
		err          error
	)
	where, err := mgParseQuery(m.Options.Where.(map[string]interface{}))
	save := map[string]interface{}{"$set": map[string]interface{}{}}
	if len(update) > 0 {
		temp := make(map[string]interface{}, 0)
		for k, v := range update {
			temp[k] = v
		}
		save["$set"] = temp
	}
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo Save error")
	}
	if m.Options.Limit == 1 {
		updateResult, err = m.Client.Collection(m.TableName).UpdateOne(context.TODO(), where, save)
	} else {
		updateResult, err = m.Client.Collection(m.TableName).UpdateMany(context.TODO(), where, save)
	}
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo Save error1")
	}
	if m.Debug { //组装sql debug用
		wherejson, _ := json.Marshal(where)
		savejson, _ := json.Marshal(save)
		_sql := "update " + m.TableName + " set " + string(savejson) + " where " + string(wherejson)
		m.SQL = _sql
	}
	return updateResult.ModifiedCount, err
}

//Delete 删除数据
func (m *MongoModel) Delete() (int64, error) {
	var (
		deleteResult *mongo.DeleteResult
		err          error
	)
	where, err := mgParseQuery(m.Options.Where.(map[string]interface{}))
	if m.Options.Limit == 1 {
		// 删除1条
		deleteResult, err = m.Client.Collection(m.TableName).DeleteOne(context.TODO(), where)
	} else {
		// 删除所有
		deleteResult, err = m.Client.Collection(m.TableName).DeleteMany(context.TODO(), where)
	}
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo delete error")
	}
	if m.Debug { //组装sql debug用
		wherejson, _ := json.Marshal(where)
		_sql := "delete from " + m.TableName + " where " + string(wherejson)
		m.SQL = _sql
	}
	return deleteResult.DeletedCount, err
}

//Find 查询数据单条 => ok
func (m *MongoModel) Find() (map[string]interface{}, error) {
	m.Options.Limit = 1 //只取一条
	result, err := mgParseQuerySelect(m, "find")
	return result.(map[string]interface{}), err
	// if m.DbError != nil {
	// 	log.WithFields(log.Fields{"isPush": true, "error": m.DbError}).Warn("mongo find error")
	// 	return nil, m.DbError
	// }
	// // 将选项传递给FindOne()
	// findOptions := options.FindOne()
	// findOptions.SetSort(m.Options.Order)
	// var result map[string]interface{}
	// err := m.Client.Collection(m.TableName).FindOne(context.TODO(), m.Options.Where, findOptions).Decode(&result)
	// id := result["_id"].(primitive.ObjectID)
	// result["_id"] = hex.EncodeToString(id[:])
	// return result, err
}

//Select 查询多条数据 => ok
func (m *MongoModel) Select() ([]map[string]interface{}, error) {
	result, err := mgParseQuerySelect(m, "select")
	return result.([]map[string]interface{}), err
}

// Count 统计
func (m *MongoModel) Count() (int64, error) {
	if m.Options.Group != "" {
		var pipeline []map[string]interface{}
		where := map[string]interface{}{ //查询条件
			"$match": map[string]interface{}{},
		}
		group := map[string]interface{}{ //分组条件
			"$group": map[string]interface{}{},
		}
		where["$match"], _ = mgParseQuery(m.Options.Where.(map[string]interface{}))
		group["$group"] = mgParseGroup(m.Options.Group)
		pipeline = append(pipeline, where)
		if m.Options.Group != "" {
			pipeline = append(pipeline, group)
		}
		cur, err := m.Client.Collection(m.TableName).Aggregate(context.TODO(), pipeline)
		if err != nil {
			log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo Count error")
		}
		return int64(cur.RemainingBatchLength()), err
	}
	where, _ := mgParseQuery(m.Options.Where.(map[string]interface{}))
	cur, err := m.Client.Collection(m.TableName).CountDocuments(context.TODO(), where)
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo Count error1")
	}
	if m.Debug { //组装sql debug用
		wherejson, _ := json.Marshal(where)
		_sql := "select count(*) from " + m.TableName + " where " + string(wherejson)
		m.SQL = _sql
	}
	return cur, err
}

//SetInc 字段值增长 array("fieldName",1)
//@param  data 字段名,增长值
//@return id,error
func (m *MongoModel) SetInc(update map[string]interface{}) (int64, error) {
	var (
		updateResult *mongo.UpdateResult
		err          error
	)
	where, err := mgParseQuery(m.Options.Where.(map[string]interface{}))
	save := map[string]interface{}{"$inc": map[string]interface{}{}}
	if len(update) > 0 {
		temp := make(map[string]interface{}, 0)
		for k, v := range update {
			temp[k] = v
		}
		save["$inc"] = temp
	}
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo SetInc error")
	}
	if m.Options.Limit == 1 {
		updateResult, err = m.Client.Collection(m.TableName).UpdateOne(context.TODO(), where, save)
	} else {
		updateResult, err = m.Client.Collection(m.TableName).UpdateMany(context.TODO(), where, save)
	}
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo SetInc error1")
	}
	if m.Debug { //组装sql debug用
		wherejson, _ := json.Marshal(where)
		savejson, _ := json.Marshal(save)
		_sql := "update " + m.TableName + " set " + string(savejson) + " where " + string(wherejson)
		m.SQL = _sql
	}
	return updateResult.ModifiedCount, err
}

//SetDec 字段值增长 array("fieldName",1)
//@param  data 字段名,增长值
//@return id,error
func (m *MongoModel) SetDec(update map[string]interface{}) (int64, error) {
	var (
		updateResult *mongo.UpdateResult
		err          error
	)
	where, err := mgParseQuery(m.Options.Where.(map[string]interface{}))
	save := map[string]interface{}{"$inc": map[string]interface{}{}}
	if len(update) > 0 {
		temp := make(map[string]interface{}, 0)
		for k, v := range update {
			switch reflect.TypeOf(v).String() {
			case "int":
				temp[k] = -v.(int)
			case "int32":
				temp[k] = -v.(int32)
			case "int64":
				temp[k] = -v.(int64)
			default:
				log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo SetDec 转换 error")
			}
		}
		save["$inc"] = temp
	}
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo SetDec error")
	}
	if m.Options.Limit == 1 {
		updateResult, err = m.Client.Collection(m.TableName).UpdateOne(context.TODO(), where, save)
	} else {
		updateResult, err = m.Client.Collection(m.TableName).UpdateMany(context.TODO(), where, save)
	}
	if err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo SetDec error1")
	}
	if m.Debug { //组装sql debug用
		wherejson, _ := json.Marshal(where)
		savejson, _ := json.Marshal(save)
		_sql := "update " + m.TableName + " set " + string(savejson) + " where " + string(wherejson)
		m.SQL = _sql
	}
	return updateResult.ModifiedCount, err
}

// Where 设置条件
func (m *MongoModel) Where(where interface{}) *MongoModel {
	m.Options.Where = where
	return m
}

// Limit 设置条件
func (m *MongoModel) Limit(offset int) *MongoModel {
	m.Options.Limit = offset
	return m
}

// Page 分页查询
func (m *MongoModel) Page(page int, rows int) *MongoModel {
	m.Options.Limit = rows
	m.Options.Offset = (page - 1) * m.Options.Limit //偏移量
	return m
}

// Field 设置字段
func (m *MongoModel) Field(fields []string) *MongoModel {
	m.Fields = fields
	return m
}

// Order 设置字段
func (m *MongoModel) Order(order string) *MongoModel {
	m.Options.Order = order
	return m
}

// Group 设置字段
func (m *MongoModel) Group(group string) *MongoModel {
	m.Options.Group = group
	return m
}

// GetFields 获取字段
func (m *MongoModel) GetFields() []string {
	return m.Fields
}

// GetWhereFields 获取可查询字段
func (m *MongoModel) GetWhereFields() []string {
	return m.WhereFields
}

// GetTransferFields 获取流转字段
func (m *MongoModel) GetTransferFields() map[string]string {
	return m.TransferFields
}

// GetTransferOrders 获取排序流转字段
func (m *MongoModel) GetTransferOrders() map[string]string {
	return m.TransferOrders
}

//GetTableName 得到完整的数据表名
func (m *MongoModel) GetTableName() string {
	return m.TableName
}

//GetLastSQL 返回最后执行的sql语句 ok
func (m *MongoModel) GetLastSQL() string {
	return m.SQL
}

//mgParseQuerySelect 根据条件返回单条多条数据
func mgParseQuerySelect(m *MongoModel, mode string) (interface{}, error) {
	var results interface{} //返回结果
	switch mode {
	case "find":
		results = make(map[string]interface{}) //返回结果
	case "select":
		results = []map[string]interface{}{} //返回结果
	default:
		log.WithFields(log.Fields{"isPush": true, "error_info": mode}).Warn("查询类型错误！")
		return results, errors.New("查询类型错误！")
	}
	//组装字段
	mapField := map[string]int{}
	for k := range m.Fields {
		mapField[m.Fields[k]] = 1
	}
	var (
		cur *mongo.Cursor
		err error
	)
	if m.Options.Group == "" {
		findOptions := options.Find()
		findOptions.SetProjection(mapField)
		findOptions.SetSort(mgParseOrder(m.Options.Order))
		findOptions.SetSkip(int64(m.Options.Offset))
		findOptions.SetLimit(int64(m.Options.Limit))
		where, _ := mgParseQuery(m.Options.Where.(map[string]interface{}))
		cur, err = m.Client.Collection(m.TableName).Find(context.TODO(), where, findOptions)
		if m.Debug { //组装sql debug用
			wherestrjson, _ := json.Marshal(where)
			mapFieldstrjson, _ := json.Marshal(mapField)
			orderstrjson, _ := json.Marshal(mgParseOrder(m.Options.Order))
			_sql := "select " + string(mapFieldstrjson) + " from " + m.TableName + " where " + string(wherestrjson)
			_sql = _sql + " order " + string(orderstrjson) + " limit " + strconv.Itoa(m.Options.Offset) + "," + strconv.Itoa(m.Options.Limit)
			m.SQL = _sql
		}
		if err != nil {
			log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo select-find error")
		}
	} else {
		var pipeline []map[string]interface{}
		field := map[string]interface{}{ //查询字段
			"$project": map[string]interface{}{},
		}
		where := map[string]interface{}{ //查询条件
			"$match": map[string]interface{}{},
		}
		group := map[string]interface{}{ //分组条件
			"$group": map[string]interface{}{},
		}
		order := map[string]interface{}{ //排序条件
			"$sort": map[string]interface{}{},
		}
		skip := map[string]interface{}{ //起始条数条件
			"$skip": m.Options.Offset,
		}
		limit := map[string]interface{}{ //条数条件
			"$limit": m.Options.Limit,
		}
		//组装条件
		field["$project"] = mapField
		where["$match"], _ = mgParseQuery(m.Options.Where.(map[string]interface{}))
		group["$group"] = mgParseGroup(m.Options.Group)
		order["$sort"] = mgParseOrder(m.Options.Order)
		pipeline = append(pipeline, field)
		pipeline = append(pipeline, where)
		if m.Options.Group != "" {
			pipeline = append(pipeline, group)
		}
		if m.Options.Order != "" {
			pipeline = append(pipeline, order)
		}
		if m.Options.Limit > 0 {
			pipeline = append(pipeline, skip)
			pipeline = append(pipeline, limit)
		}
		//end组装条件
		// fmt.Println("pipeline", pipeline, reflect.TypeOf(m.Options.Where).String())
		cur, err = m.Client.Collection(m.TableName).Aggregate(context.TODO(), pipeline)
		if m.Debug { //组装sql debug用
			pipelinestrjson, _ := json.Marshal(pipeline)
			_sql := "select from " + m.TableName + " where " + string(pipelinestrjson)
			m.SQL = _sql
		}
		if err != nil {
			log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo select error")
		}
	}

	// 查找多个文档返回一个光标
	// 遍历游标允许我们一次解码一个文档
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo select error1")
		}
		// id := elem["_id"].(primitive.ObjectID)
		// elem["_id"] = hex.EncodeToString(id[:])
		if mode == "find" {
			results = elem
		} else if mode == "select" {
			results = append(results.([]map[string]interface{}), elem)
		}
	}

	if err := cur.Err(); err != nil {
		log.WithFields(log.Fields{"isPush": true, "error": err}).Warn("mongo select error2")
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	return results, err
}

//order分析 后面要支持数组
//@access protected
//@param mixed $order
//@return string
func mgParseOrder(order string) map[string]interface{} {
	data := strings.Split(order, ",")
	orderMap := map[string]interface{}{}
	if len(data) > 0 {
		for key := range data {
			// fmt.Println("data[key]", data[key])
			dataSplit := strings.Split(string(data[key]), " ")
			v := strings.ToLower(dataSplit[1])
			switch v {
			case "asc":
				orderMap[dataSplit[0]] = 1
			case "desc":
				orderMap[dataSplit[0]] = -1
			}
		}
	}
	return orderMap
}

//group分析
//@param mixed group
//@return string
func mgParseGroup(group string) map[string]map[string]interface{} {
	data := strings.Split(group, ",")
	groupMap := map[string]map[string]interface{}{
		"_id": {},
		"count": {
			"$sum": 1,
		},
	}
	if len(data) > 0 {
		for k := range data {
			groupMap["_id"][data[k]] = "$" + data[k]
		}
	}
	return groupMap
}

//where 条件分析
func mgParseQuery(where map[string]interface{}) (map[string]interface{}, error) {
	//初始化查询格式
	query := make(map[string]interface{})
	dataType := reflect.TypeOf(where).String()
	switch dataType {
	case "string":
		return query, errors.New("mgParseQuery 不支持string")
	case "map[string]interface {}":
		//数组形式的where条件参数
		operate := "AND"
		//whereTemp 为了把_logic过滤掉
		whereTemp := make(map[string]interface{})
		for k, v := range where {
			if k == "_logic" {
				operate = common.Strtoupper(v.(string))
				continue
			} else {
				whereTemp[k] = v
			}
		}
		// fmt.Println("whereTemp", whereTemp, operate)
		//初始化
		switch operate {
		case "AND":
			query["$and"] = make([]map[string]interface{}, 0)
		case "OR":
			query["$or"] = make([]map[string]interface{}, 0)
		default:
			log.WithFields(log.Fields{"isPush": true, "operate": operate}).Warn("mgParseQuery operate 类型不支持")
			return query, errors.New("mgParseQuery operate:" + operate + " 类型不支持")
		}
		for k, v := range whereTemp {
			// fmt.Println("whereTemp1", k, v)
			//复合查询
			if k == "_complex" || common.IsNumeric(k) {
				// fmt.Println("_complex_complex_complex_complex")
				whereParse, err := mgParseSpecialWhere(v.(map[string]interface{}))
				// fmt.Println("whereParse", whereParse, err)
				if err != nil {
					return query, err
				}

				//如果有logic 取logic逻辑
				switch operate {
				case "AND":
					query["$and"] = append(query["$and"].([]map[string]interface{}), whereParse)
				case "OR":
					query["$or"] = append(query["$or"].([]map[string]interface{}), whereParse)
				default:
					log.WithFields(log.Fields{"isPush": true, "operate": operate}).Warn("mgParseQuery operate 类型不支持")
					return query, errors.New("mgParseQuery operate:" + operate + " 类型不支持")
				}
			} else {
				matchQuery, err := mgParseWhereItem(k, v)
				if err != nil {
					return query, err
				}
				// fmt.Println("matchQuery", matchQuery, err)
				//如果有logic 取logic逻辑
				switch operate {
				case "AND":
					if len(matchQuery) > 0 {
						query["$and"] = append(query["$and"].([]map[string]interface{}), matchQuery)
					}
				case "OR":
					if len(matchQuery) > 0 {
						query["$or"] = append(query["$or"].([]map[string]interface{}), matchQuery)
					}
				default:
					log.WithFields(log.Fields{"isPush": true, "operate": operate}).Warn("mgParseQuery operate 类型不支持")
					return query, errors.New("mgParseQuery operate:" + operate + " 类型不支持")
				}
			}
		}
	default:
		return query, errors.New("mgParseQuery 不支持" + dataType)
	}
	return query, nil
}

//mgParseSpecialWhere 特殊条件分析
func mgParseSpecialWhere(where map[string]interface{}) (map[string]interface{}, error) {
	return mgParseQuery(where)
}

//mgParseWhereItem 分析where条件
func mgParseWhereItem(key interface{}, val interface{}) (map[string]interface{}, error) {
	matchQuery := map[string]interface{}{}
	valType := reflect.TypeOf(val).String()

	// reg1, _ := regexp.Compile(`^(eq|gt|egt|lt|elt)$`)
	// reg2, _ := regexp.Compile(`^(like)$`)
	// reg3, _ := regexp.Compile(`^(AND|OR|XOR)$`)
	// reg4, _ := regexp.Compile(`^(in)$`)
	// reg5, _ := regexp.Compile(`^(between)$`)
	var eqData = []string{"eq", "neq", "gt", "egt", "lt", "elt"}
	var likeData = []string{"notlike", "like"}
	// var logicData = []string{"AND", "OR", "XOR"}
	var inData = []string{"notin", "not in", "in"}
	var betweenData = []string{"notbetween", "not between", "between"}
	switch valType {
	case "[]interface {}":
		temp := val.([]interface{})
		switch reflect.TypeOf(temp[0]).String() {
		case "string": //值的类型是string
			exp := common.Strtolower(temp[0].(string))
			if common.InArray(exp, eqData) { //eq|gt|egt|lt|elt
				temp1Type := reflect.TypeOf(temp[1]).String()
				switch exp {
				case "eq":
					switch temp1Type {
					case "float64":
						matchQuery[key.(string)] = strconv.FormatFloat(temp[1].(float64), 'f', -1, 64)
					case "int":
						matchQuery[key.(string)] = strconv.Itoa(temp[1].(int))
					case "string":
						matchQuery[key.(string)] = temp[1].(string)
					case "primitive.ObjectID":
						matchQuery[key.(string)] = temp[1]
					default:
						log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[1], "type": reflect.TypeOf(temp[1]).String()}).Warn("type类型不支持")
						return matchQuery, errors.New("type类型不支持")
					}
				case "gt":
					matchQuery[key.(string)] = map[string]interface{}{"$gt": temp[1]}
				case "egt":
					matchQuery[key.(string)] = map[string]interface{}{"$gte": temp[1]}
				case "lt":
					matchQuery[key.(string)] = map[string]interface{}{"$lt": temp[1]}
				case "elt":
					matchQuery[key.(string)] = map[string]interface{}{"$lte": temp[1]}
				case "neq":
					matchQuery[key.(string)] = map[string]interface{}{"$ne": temp[1]}
				default:
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("条件类型不支持")
					return matchQuery, errors.New("条件类型不支持")
				}
			} else if common.InArray(exp, likeData) {
				switch reflect.TypeOf(temp[1]).String() {
				case "[]interface {}":
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[1], "type": reflect.TypeOf(temp[1]).String()}).Warn("like多条件类型不支持")
					return matchQuery, errors.New("like多条件类型不支持")
				case "string":
					matchQuery[key.(string)] = map[string]interface{}{"$like": strings.Replace(temp[1].(string), "%", "*", -1)}
				default:
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[1], "type": reflect.TypeOf(temp[1]).String()}).Warn("like type类型不支持")
					return matchQuery, errors.New("like type类型不支持")
				}

			} else if "bind" == exp { // 使用表达式
				log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("bind 条件类型不支持")
				return matchQuery, errors.New("bind 条件类型不支持")
			} else if "exp" == exp { // 使用表达式
				log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("exp 条件类型不支持")
				return matchQuery, errors.New("exp 条件类型不支持")
			} else if common.InArray(exp, inData) { //in
				//三个参数不支持
				if len(temp) == 3 && "exp" == temp[2].(string) {
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("in 三个参数条件类型不支持")
					return matchQuery, errors.New("in 三个参数条件类型不支持")
				}
				zone := make([]interface{}, 0)
				switch reflect.TypeOf(temp[1]).String() {
				case "string": //逗号分割的字符串,例如：Y,ok,lock
					temp[1] = strings.Split(temp[1].(string), ",")
					for _, v := range temp[1].([]string) {
						zone = append(zone, v)
					}
				case "[]interface {}": //数组，例如：['Y','ok','lock']
					for _, v := range temp[1].([]interface{}) {
						zone = append(zone, v)
					}
				case "[]string": //数组，例如：['Y','ok','lock']
					for _, v := range temp[1].([]string) {
						zone = append(zone, v)
					}
				case "[]int64": //数组
					for _, v := range temp[1].([]int64) {
						zone = append(zone, v)
					}
				default:
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[1], "type": reflect.TypeOf(temp[1]).String()}).Warn("in 参数类型不支持")
					return matchQuery, errors.New("in 参数类型不支持")
				}
				matchQuery[key.(string)] = map[string]interface{}{"$in": zone}
			} else if common.InArray(exp, betweenData) { //between
				switch reflect.TypeOf(temp[1]).String() {
				case "string": //1,8
					var data []string
					data = strings.Split(temp[1].(string), ",")
					matchQuery[key.(string)] = map[string]interface{}{"$gte": data[0], "$lte": data[1]}
				case "[]interface {}": //['1','8']
					var data []interface{}
					for _, v := range temp[1].([]interface{}) {
						data = append(data, v)
					}
					matchQuery[key.(string)] = map[string]interface{}{"$gte": data[0], "$lte": data[1]}
				case "[]string": //['1','8']
					var data []interface{}
					for _, v := range temp[1].([]string) {
						data = append(data, v)
					}
					matchQuery[key.(string)] = map[string]interface{}{"$gte": data[0], "$lte": data[1]}
				default:
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[1], "type": reflect.TypeOf(temp[1]).String()}).Warn("between 参数类型不支持")
					return matchQuery, errors.New("between 参数类型不支持")
				}

			} else {
				log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("条件类型不支持")
				return matchQuery, errors.New("条件类型不支持")
			}
		default:
			log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[0], "type": reflect.TypeOf(temp[0]).String()}).Warn("值类型不支持")
			return matchQuery, errors.New("值类型不支持")
		}
	case "[]string":
		temp := val.([]string)
		switch reflect.TypeOf(temp[0]).String() {
		case "string": //值的类型是string
			exp := common.Strtolower(temp[0])
			if common.InArray(exp, eqData) { //eq|gt|egt|lt|elt
				// fmt.Println("exp:", exp)
				switch exp {
				case "eq":
					matchQuery[key.(string)] = temp[1]
				case "gt":
					matchQuery[key.(string)] = map[string]interface{}{"$gt": temp[1]}
				case "egt":
					matchQuery[key.(string)] = map[string]interface{}{"$gte": temp[1]}
				case "lt":
					matchQuery[key.(string)] = map[string]interface{}{"$lt": temp[1]}
				case "elt":
					matchQuery[key.(string)] = map[string]interface{}{"$lte": temp[1]}
				default:
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("条件类型不支持")
					return matchQuery, errors.New("条件类型不支持")
				}
			} else if common.InArray(exp, likeData) {
				switch reflect.TypeOf(temp[1]).String() {
				case "[]string": //like 多个条件
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp, "type": reflect.TypeOf(temp[1]).String()}).Warn("like 多条件类型不支持")
					return matchQuery, errors.New("like 多条件类型不支持")
				case "string":
					matchQuery[key.(string)] = map[string]interface{}{"$like": strings.Replace(temp[1], "%", "*", -1)}
				default:
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[1], "type": reflect.TypeOf(temp[1]).String()}).Warn("like 值类型不支持")
					return matchQuery, errors.New("like 值类型不支持")
				}

			} else if "bind" == exp { // 使用表达式
				log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("bind 条件类型不支持")
				return matchQuery, errors.New("bind 条件类型不支持")
			} else if "exp" == exp { // 使用表达式
				log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("exp 条件类型不支持")
				return matchQuery, errors.New("exp 条件类型不支持")
			} else if common.InArray(exp, inData) { //in
				//三个参数不支持
				if len(temp) == 3 && "exp" == temp[2] {
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("in 多条件类型不支持")
					return matchQuery, errors.New("in 多条件类型不支持")
				} else {
					zone := make([]interface{}, 0)
					switch reflect.TypeOf(temp[1]).String() {
					case "string":
						expTemp := common.Explode(",", temp[1])
						for _, v := range expTemp {
							zone = append(zone, v)
						}
					default:
						log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[1], "type": reflect.TypeOf(temp[1]).String()}).Warn("in 条件值类型不支持")
						return matchQuery, errors.New("in 条件值类型不支持")
					}
					matchQuery[key.(string)] = map[string]interface{}{"$in": zone}
				}
			} else if common.InArray(exp, betweenData) { //between
				switch reflect.TypeOf(temp[1]).String() {
				case "string":
					var data []string
					data = strings.Split(temp[1], ",")
					matchQuery[key.(string)] = map[string]interface{}{"$gte": data[0], "$lte": data[1]}
				default:
					log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp, "type": reflect.TypeOf(temp[1]).String()}).Warn("between 值类型不支持")
					return matchQuery, errors.New("between 值类型不支持")
				}

			} else {
				log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": exp}).Warn("条件类型不支持")
				return matchQuery, errors.New("条件类型不支持")
			}
		default:
			log.WithFields(log.Fields{"isPush": true, "error_info": temp, "val": temp[0], "type": reflect.TypeOf(temp[0]).String()}).Warn("类型不支持")
			return matchQuery, errors.New("类型不支持")
		}

	case "float64":
		matchQuery[key.(string)] = strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case "int":
		matchQuery[key.(string)] = strconv.Itoa(val.(int))
	case "string":
		matchQuery[key.(string)] = val.(string)
	case "int64":
		matchQuery[key.(string)] = strconv.FormatInt(val.(int64), 10)
	case "primitive.ObjectID":
		matchQuery[key.(string)] = val
	default:
		log.WithFields(log.Fields{"isPush": true, "error_info": val, "val": val, "type": valType}).Warn("类型不支持")
		return matchQuery, errors.New("类型不支持")
	}
	return matchQuery, nil
}

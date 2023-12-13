package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// Чтение переменной окружения REDIS_URL
	redisURL := os.Getenv("REDIS_URL")
	// Если REDIS_URL не установлена, используйте значение по умолчанию "localhost:6379"
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Разбор URL-адреса Redis
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalf("Ошибка разбора URL-адреса Redis: %v", err)
	}

	// Извлечение компонент URL-адреса
	hostname := parsedURL.Hostname()
	port := parsedURL.Port()
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	// Если хост и порт не указаны, используйте значения по умолчанию
	if hostname == "" {
		hostname = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	// Создаем клиент Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", hostname, port), // Хост и порт
		Password: password,                             // Пароль (если требуется)
		Username: username,                             // username (если требуется)
		DB:       0,                                    // Номер базы данных
	})

	defer client.Close()

	//////////////////////////////////////////////////////////////////////////////////

	// Выполнение команды FLUSHALL
	err = client.FlushAll(ctx).Err()
	if err != nil {
		log.Printf("Ошибка выполнения команды FLUSHALL: %v\n", err)
	} else {
		fmt.Println("Команда FLUSHALL выполнена успешно.")
	}

	log.Println("=================== String ======================")

	// SET и GET
	err = client.Set(ctx, "key1", "value1", 0).Err()
	if err != nil {
		log.Printf("Ошибка SET: %v\n", err)
	}

	value, err := client.Get(ctx, "key1").Result()
	if err != nil {
		log.Printf("Ошибка GET: %v\n", err)
	} else {
		fmt.Printf("key1: %s\n", value)
	}

	// SETEX (установка с TTL) и GET
	err = client.SetEX(ctx, "key1", "value2", 10*time.Second).Err()
	if err != nil {
		log.Printf("Ошибка SETEX: %v\n", err)
	}

	value, err = client.Get(ctx, "key1").Result()
	if err != nil {
		log.Printf("Ошибка GET: %v\n", err)
	} else {
		fmt.Printf("key1: %s\n", value)
	}

	// SETXX (установить только если существует) и GET
	err = client.SetXX(ctx, "key2", "value2", 10*time.Second).Err()
	if err != nil {
		log.Printf("Ошибка SETXX: %v\n", err)
	}

	value, err = client.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 не существует")
	} else if err != nil {
		log.Printf("Ошибка GET: %v\n", err)
	} else {
		fmt.Printf("key2: %s\n", value)
	}

	// SETXX (установить только если существует) и GET
	err = client.SetXX(ctx, "key1", "value2", 10*time.Second).Err()
	if err != nil {
		log.Printf("Ошибка SETXX: %v\n", err)
	}

	value, err = client.Get(ctx, "key1").Result()
	if err == redis.Nil {
		fmt.Println("key1 не существует")
	} else if err != nil {
		log.Printf("Ошибка GET: %v\n", err)
	} else {
		fmt.Printf("key1: %s\n", value)
	}

	// SETNX (установить только если не существует)
	setNXResult, err := client.SetNX(ctx, "key3", "Hello", 0).Result()
	if err != nil {
		log.Printf("Ошибка SETNX: %v\n", err)
	} else if setNXResult {
		fmt.Println("key3 установлен как 'Hello'")
	}

	setNXResult, err = client.SetNX(ctx, "key3", "World", 0).Result()
	if err != nil {
		log.Printf("Ошибка SETNX: %v\n", err)
	} else if setNXResult {
		fmt.Println("key3 установлен как 'World'")
	}

	value, err = client.Get(ctx, "key3").Result()
	if err != nil {
		log.Printf("Ошибка GET: %v\n", err)
	} else {
		fmt.Printf("key3: %s\n", value)
	}

	// GETSET
	getSetResult, err := client.GetSet(ctx, "key3", "Redis").Result()
	if err != nil {
		log.Printf("Ошибка GETSET: %v\n", err)
	} else {
		fmt.Printf("Результат GETSET: %s\n", getSetResult)
	}

	value, err = client.Get(ctx, "key3").Result()
	if err != nil {
		log.Printf("Ошибка GET: %v\n", err)
	} else {
		fmt.Printf("key3: %s\n", value)
	}

	// MGET
	keysToGet := []string{"key1", "key2", "key3"}
	mGetResult, err := client.MGet(ctx, keysToGet...).Result()
	if err != nil {
		log.Printf("Ошибка MGET: %v\n", err)
	} else {
		fmt.Printf("Результат MGET: %v\n", mGetResult)
	}

	// MSET
	keysAndValues := []interface{}{"key2", "Hello", "key3", "Redis"}
	err = client.MSet(ctx, keysAndValues...).Err()
	if err != nil {
		log.Printf("Ошибка MSET: %v\n", err)
	}

	mGetResult, err = client.MGet(ctx, keysToGet...).Result()
	if err != nil {
		log.Printf("Ошибка MGET: %v\n", err)
	} else {
		fmt.Printf("Результат MGET после MSET: %v\n", mGetResult)
	}

	// MSETNX
	keysAndValues = []interface{}{"key4", "value4", "key5", "value5"}
	mSetNXResult, err := client.MSetNX(ctx, keysAndValues...).Result()
	if err != nil {
		log.Printf("Ошибка MSETNX: %v\n", err)
	} else {
		fmt.Printf("Результат MSETNX: %v\n", mSetNXResult)
	}

	keysToGet = []string{"key4", "key5"}
	mGetResult, err = client.MGet(ctx, keysToGet...).Result()
	if err != nil {
		log.Printf("Ошибка MGET: %v\n", err)
	} else {
		fmt.Printf("Результат MGET после MSETNX: %v\n", mGetResult)
	}

	// SETEX (с TTL в секундах)
	err = client.SetEX(ctx, "time", "hello", 10*time.Second).Err()
	if err != nil {
		log.Printf("Ошибка SETEX: %v\n", err)
	}

	ttlResult, err := client.TTL(ctx, "time").Result()
	if err != nil {
		log.Printf("Ошибка TTL: %v\n", err)
	} else {
		fmt.Printf("TTL для ключа 'time': %s\n", ttlResult)
	}

	// EXPIRE (установка TTL)
	err = client.Expire(ctx, "time", 5*time.Second).Err()
	if err != nil {
		log.Printf("Ошибка EXPIRE: %v\n", err)
	}

	ttlResult, err = client.TTL(ctx, "time").Result()
	if err != nil {
		log.Printf("Ошибка TTL: %v\n", err)
	} else {
		fmt.Printf("TTL для ключа 'time' после EXPIRE: %s\n", ttlResult)
	}

	// Удаление ключа
	err = client.Del(ctx, "time").Err()
	if err != nil {
		log.Printf("Ошибка удаления ключа 'time': %v\n", err)
	}

	value, err = client.Get(ctx, "time").Result()
	if err == redis.Nil {
		fmt.Println("Ключ 'time' не существует")
	} else if err != nil {
		log.Printf("Ошибка GET: %v\n", err)
	} else {
		fmt.Printf("Значение ключа 'time' после удаления: %s\n", value)
	}

	// INCR и DECR
	err = client.Set(ctx, "count", 10, 0).Err()
	if err != nil {
		log.Printf("Ошибка SET: %v\n", err)
	}

	fmt.Printf("INCR: %d\n", client.Incr(ctx, "count").Val())
	fmt.Printf("INCRBY: %d\n", client.IncrBy(ctx, "count", 5).Val())
	fmt.Printf("DECR: %d\n", client.Decr(ctx, "count").Val())
	fmt.Printf("DECRBY: %d\n", client.DecrBy(ctx, "count", 5).Val())
	fmt.Printf("INCRBYFLOAT: %f\n", client.IncrByFloat(ctx, "count", 1.5).Val())

	// APPEND и EXISTS
	err = client.Append(ctx, "hello.world", "Hello").Err()
	if err != nil {
		log.Printf("Ошибка APPEND: %v\n", err)
	}

	err = client.Append(ctx, "hello.world", " World").Err()
	if err != nil {
		log.Printf("Ошибка APPEND: %v\n", err)
	}

	fmt.Printf("APPEND: %s\n", client.Get(ctx, "hello.world").Val())
	fmt.Printf("EXISTS 'hello.world': %d\n", client.Exists(ctx, "hello.world").Val())
	fmt.Printf("EXISTS 'hello': %d\n", client.Exists(ctx, "hello").Val())

	log.Println("=================== List/Array ======================")

	// RPUSH - добавление элементов в конец списка
	err = client.RPush(ctx, "mylist", "one").Err()
	if err != nil {
		log.Printf("Ошибка RPUSH: %v\n", err)
	}

	err = client.RPush(ctx, "mylist", []string{"two", "three"}).Err()
	if err != nil {
		log.Printf("Ошибка RPUSH: %v\n", err)
	}

	// LPUSH - добавление элемента в начало списка
	err = client.LPush(ctx, "mylist", "zero").Err()
	if err != nil {
		log.Printf("Ошибка LPUSH: %v\n", err)
	}

	// LLEN - получение длины списка
	length, err := client.LLen(ctx, "mylist").Result()
	if err != nil {
		log.Printf("Ошибка LLEN: %v\n", err)
	} else {
		fmt.Printf("Длина списка mylist: %d\n", length)
	}

	// LRANGE - получение элементов списка по диапазону
	rangeResult, err := client.LRange(ctx, "mylist", 0, -1).Result()
	if err != nil {
		log.Printf("Ошибка LRANGE: %v\n", err)
	} else {
		fmt.Printf("Элементы списка mylist: %v\n", rangeResult)
	}

	// LINDEX - получение элемента по индексу
	indexValue, err := client.LIndex(ctx, "mylist", -1).Result()
	if err != nil {
		log.Printf("Ошибка LINDEX: %v\n", err)
	} else {
		fmt.Printf("Элемент с индексом -1: %s\n", indexValue)
	}

	// LTRIM - обрезка списка по диапазону
	err = client.LTrim(ctx, "mylist", 1, -1).Err()
	if err != nil {
		log.Printf("Ошибка LTRIM: %v\n", err)
	}

	// LSET - установка значения элемента по индексу
	err = client.LSet(ctx, "mylist", 0, "Hello").Err()
	if err != nil {
		log.Printf("Ошибка LSET: %v\n", err)
	}

	// LREM - удаление элементов из списка
	err = client.LRem(ctx, "mylist", 0, "Hello").Err()
	if err != nil {
		log.Printf("Ошибка LREM: %v\n", err)
	}

	// LINSERT - вставка элемента перед или после указанного элемента
	err = client.LInsert(ctx, "mylist", "BEFORE", "two", "one").Err()
	if err != nil {
		log.Printf("Ошибка LINSERT: %v\n", err)
	}

	err = client.LInsert(ctx, "mylist", "AFTER", "1", "0").Err()
	if err != nil {
		log.Printf("Ошибка LINSERT: %v\n", err)
	}

	// RPOPLPUSH - вынимает элемент справа из одного списка и вставляет его слева в другой список
	poppedValue, err := client.RPopLPush(ctx, "mylist", "mylist2").Result()
	if err != nil {
		log.Printf("Ошибка RPopLPush: %v\n", err)
	} else {
		fmt.Printf("Значение после RPopLPush: %s\n", poppedValue)
	}

	// Вывод содержимого списков
	mylist, err := client.LRange(ctx, "mylist", 0, -1).Result()
	if err != nil {
		log.Printf("Ошибка LRange: %v\n", err)
	} else {
		fmt.Printf("Содержимое mylist: %v\n", mylist)
	}

	mylist2, err := client.LRange(ctx, "mylist2", 0, -1).Result()
	if err != nil {
		log.Printf("Ошибка LRange: %v\n", err)
	} else {
		fmt.Printf("Содержимое mylist2: %v\n", mylist2)
	}

	log.Println("=================== Hash ======================")

	// HSET - установка значения поля в хеше
	err = client.HSet(ctx, "myhash", "field1", "0").Err()
	if err != nil {
		log.Printf("Ошибка HSET: %v\n", err)
	}

	err = client.HSet(ctx, "myhash", "field2", "value2").Err()
	if err != nil {
		log.Printf("Ошибка HSET: %v\n", err)
	}

	// HGETALL - получение всех полей и значений хеша
	values, err := client.HGetAll(ctx, "myhash").Result()
	if err != nil {
		log.Printf("Ошибка HGETALL: %v\n", err)
	} else {
		fmt.Println("HGETALL", values)
		fmt.Println("field1", values["field1"])
	}

	// HGET - получение значения поля из хеша
	value, err = client.HGet(ctx, "myhash", "field2").Result()
	if err != nil {
		log.Printf("Ошибка HGET: %v\n", err)
	} else {
		fmt.Println("field2", value)
	}

	// HKEYS - получение всех ключей хеша
	keys, err := client.HKeys(ctx, "myhash").Result()
	if err != nil {
		log.Printf("Ошибка HKEYS: %v\n", err)
	} else {
		fmt.Println("HKEYS", keys)
	}

	// HVALS - получение всех значений хеша
	valuesSlice, err := client.HVals(ctx, "myhash").Result()
	if err != nil {
		log.Printf("Ошибка HVALS: %v\n", err)
	} else {
		fmt.Println("HVALS", valuesSlice)
	}

	// HDEL - удаление поля из хеша
	err = client.HDel(ctx, "myhash", "field2").Err()
	if err != nil {
		log.Printf("Ошибка HDEL: %v\n", err)
	}

	// HEXISTS - проверка существования поля в хеше
	exists, err := client.HExists(ctx, "myhash", "field2").Result()
	if err != nil {
		log.Printf("Ошибка HEXISTS: %v\n", err)
	} else {
		fmt.Println("HEXISTS", exists)
	}

	exists, err = client.HExists(ctx, "myhash", "field1").Result()
	if err != nil {
		log.Printf("Ошибка HEXISTS: %v\n", err)
	} else {
		fmt.Println("HEXISTS", exists)
	}

	// HLEN - получение количества полей в хеше
	length, err = client.HLen(ctx, "myhash").Result()
	if err != nil {
		log.Printf("Ошибка HLEN: %v\n", err)
	} else {
		fmt.Println("HLEN", length)
	}

	// HSETNX - установка значения поля в хеше, если оно не существует
	err = client.HSetNX(ctx, "myhash", "field2", "Hello").Err()
	if err != nil {
		log.Printf("Ошибка HSETNX: %v\n", err)
	}

	err = client.HSetNX(ctx, "myhash", "field2", "World").Err()
	if err != nil {
		log.Printf("Ошибка HSETNX: %v\n", err)
	}

	// Повторная проверка значения поля
	value, err = client.HGet(ctx, "myhash", "field2").Result()
	if err != nil {
		log.Printf("Ошибка HGET: %v\n", err)
	} else {
		fmt.Println("HSETNX", value)
	}

	// HINCRBY - увеличение значения поля на заданную величину
	newValue, err := client.HIncrBy(ctx, "myhash", "field1", 5).Result()
	if err != nil {
		log.Printf("Ошибка HINCRBY: %v\n", err)
	} else {
		fmt.Println("HINCRBY", newValue)
	}

	// Повторная проверка значения поля
	value, err = client.HGet(ctx, "myhash", "field1").Result()
	if err != nil {
		log.Printf("Ошибка HGET: %v\n", err)
	} else {
		fmt.Println("HINCRBY", value)
	}

	// HINCRBYFLOAT - увеличение значения поля на заданную величину (с плавающей запятой)
	valueFloat, err := client.HIncrByFloat(ctx, "myhash", "field1", -2.6).Result()
	if err != nil {
		log.Printf("Ошибка HINCRBYFLOAT: %v\n", err)
	} else {
		fmt.Println("HINCRBYFLOAT", valueFloat)
	}

	fmt.Println("=================Transactions=============================")

	// Устанавливаем значение ключа 'another-key'
	err = client.Set(ctx, "another-key", "another-value", 0).Err()
	if err != nil {
		log.Printf("Ошибка установки ключа 'another-key': %v\n", err)
		return
	}

	// Начинаем мультикоманду (транзакцию)
	pipe := client.TxPipeline()
	setKeyCmd := pipe.Set(ctx, "other-key", "other-value", 0)
	getKeyCmd := pipe.Get(ctx, "another-key")

	// Выполняем все команды транзакции
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Ошибка выполнения транзакции: %v\n", err)
		return
	}

	// Получаем результаты команд
	setKeyReply := setKeyCmd.Val()
	otherKeyValue := getKeyCmd.Val()

	fmt.Println("Transactions", setKeyReply, otherKeyValue)

	fmt.Println("===================scanIterator===========================")

	// Итерация по ключам с использованием SCAN
	var cursor uint64 = 0
	keysCount := 0

	for {
		// Получаем следующую порцию ключей
		keys, newCursor, err := client.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			log.Printf("Ошибка SCAN: %v\n", err)
			break
		}

		for _, key := range keys {
			// Исключаем определенные ключи (ваша логика)
			if key != "myhash" && key != "mylist" && key != "mylist2" {
				// Получаем значение ключа
				value, err := client.Get(ctx, key).Result()
				if err != nil {
					log.Printf("Ошибка получения значения для ключа %s: %v\n", key, err)
				} else {
					fmt.Println("scanIterator", key, value)
				}
			}
		}

		keysCount += len(keys)

		// Если новый курсор равен нулю, значит, мы достигли конца
		if newCursor == 0 {
			break
		}

		cursor = newCursor
	}

	fmt.Println("Общее количество ключей:", keysCount)

	// Итерация по полям и значениям хеша 'myhash' с использованием HScan
	var hCursor uint64 = 0

	for {
		// Получаем следующую порцию полей и значений хеша 'myhash'
		results, newCursor, err := client.HScan(ctx, "myhash", hCursor, "*", 100).Result()
		if err != nil {
			log.Printf("Ошибка HSCAN: %v\n", err)
			break
		}

		// Результаты возвращаются в форме map[string]string
		for field, value := range results {
			fmt.Println("scanIterator", field, value)
		}

		// Если новый курсор равен нулю, значит, мы достигли конца
		if newCursor == 0 {
			break
		}

		hCursor = newCursor
	}

	// Итерация по ключам с использованием SCAN с определенными опциями
	// Итерация по ключам с использованием SCAN с определенными опциями
	cursor = uint64(0)
	match := "o*"
	count := int64(100)

	for {
		// Получаем следующую порцию ключей
		keys, newCursor, err := client.Scan(ctx, cursor, match, count).Result()
		if err != nil {
			log.Printf("Ошибка SCAN: %v\n", err)
			break
		}

		for _, key := range keys {
			fmt.Println("scanIterator", key)
		}

		// Если новый курсор равен нулю, значит, мы достигли конца
		if newCursor == 0 {
			break
		}

		cursor = newCursor
	}

	fmt.Println("===================Auto-Pipelining===========================")

	// Пример использования Auto-Pipelining
	pipe = client.Pipeline()

	// Выполнение нескольких команд в одной пайплайн транзакции
	pipe.Set(ctx, "Tm9kZSBSZWRpcw==", "users:1", 0)
	pipe.SAdd(ctx, "users:1:tokens", "Tm9kZSBSZWRpcw==")

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Ошибка выполнения пайплайн транзакции: %v\n", err)
		return
	}

	fmt.Println("===================Exit===========================")
}

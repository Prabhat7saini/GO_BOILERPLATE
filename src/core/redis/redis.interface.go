package redis



type RedisService interface {
	Set(key string, value interface{}) error
	SetWithExpiration(key string, value interface{},expTimeInMinutes int) error
	Get(key string) (string, error)
	Delete(key string) error
}

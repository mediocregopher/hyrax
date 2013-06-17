package config

//param Type
const (
    STRING = iota
    INT
)

type param struct {
    Description string
    Type int
    Default interface{}
}

var params = map[string]param{
    "port":
        param{ Description: "The tcp port to listen for new connections on",
               Type: INT,
               Default: 3400 },

    "redis-addr":
        param{ Description: "The hostname:port (or unix sock location) to connect to redis on",
               Type: STRING,
               Default: "localhost:6379" },
}

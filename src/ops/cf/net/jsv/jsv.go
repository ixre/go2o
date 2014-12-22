package jsv

import (
	"bytes"
	"errors"
	"net"
	"ops/cf/app"
	"reflect"
)

var (
	Context          app.Context
	debugMode        bool //调试模式
	defaultBytesSize int  = 64
	CmdOperateBytes       = []byte(">>")
	cmdDot                = []byte(".")
	invalidBytes          = []byte(`{"error":"Invalid request"}`)
	AutoResetConn         = true
)

func Configure(c app.Context) {
	Context = c
	debugMode = c.Debug()
}

type Server struct {
	services map[string]interface{}
	methods  map[string]map[string]int
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]interface{}),
		methods:  make(map[string]map[string]int),
	}
}

func (this *Server) RegisterName(n string, v interface{}) error {
	if _, exist := this.services[n]; exist {
		return errors.New("exist name :" + n)
	}
	if v == nil {
		return errors.New("service is null")
	}
	this.services[n] = v
	this.registerMethods(n, v)
	return nil
}

func (this *Server) registerMethods(n string, v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		panic("must be point struct")
	}
	methods := make(map[string]int)
	for i, n := 0, t.NumMethod(); i < n; i++ {
		methods[t.Method(i).Name] = i
		Println("[SCAN] - ", t, t.Method(i).Name)
	}
	this.methods[n] = methods
}

// handle socket request, s is struct name
// m is method name, arg as the arguments of method.
// return the result bytes and error
// example socket command like :
// "{"usr":"","pwd":"123"}>>Member.Login"
func (this *Server) handle(s, m, arg []byte) ([]byte, error) {
	Println("[Server][Handle]:%s.%s ", string(s), string(m))
	structName := string(s)
	if v1, e := this.services[structName]; e {
		v := reflect.ValueOf(v1)
		if i, e := this.methods[structName][string(m)]; e {
			method := v.Method(i)
			mt := reflect.Indirect(method).Type()
			var err error

			if mt.NumIn() != 2 {
				return nil, errors.New("method must contain args and result two arguments")
			}
			var aIn, aOut reflect.Value
			var aInT, aOutT reflect.Type

			aInT = mt.In(0)  //arg
			aOutT = mt.In(1) //result

			// reflect.New创建指针
			if aInT.Kind() == reflect.Ptr {
				aIn = reflect.New(aInT.Elem())
			} else {
				aIn = reflect.Indirect(reflect.New(aInT))
			}

			if aOutT.Kind() == reflect.Ptr {
				aOut = reflect.New(aOutT.Elem())
			} else {
				aOut = reflect.Indirect(reflect.New(aOutT))
			}

			err = JsonCodec.Unmarshal(arg, aIn.Interface())
			if err != nil {
				Println("[Unmarshal][Error]:%s ; source %s ;type %#v", err, string(arg), aIn.Interface())
				return nil, err
			}

			vs := v.Method(i).Call([]reflect.Value{aIn, aOut})

			if len(vs) > 0 {
				v := vs[0].Interface()
				if v != nil {
					var cok bool
					err, cok = v.(error)
					if cok {
						return nil, err
					}
				}
			}

			//如果输出类型为String,则直接输出
			v, ok := aOut.Interface().(*string)
			if ok {
				return []byte(*v), nil
			}

			return JsonCodec.Marshal(aOut.Interface())
		}
	}
	//Printf("[Server][ERROR]:%s.%s not registed.", string(s), string(m))
	return invalidBytes, nil
}

// handle the socket request
func (this *Server) HandleRequest(conn net.Conn, d []byte) {
	// example: {"usr":"","pwd":"123"}>>Member.Login
	//	defer func(){
	//		if e := recover();e!= nil{
	//			Println(e)
	//			debug.PrintStack()
	//		}
	//	}()

	if len(d) < len(CmdOperateBytes)+len(cmdDot) {
		conn.Write(invalidBytes)
		return
	}

	i := bytes.LastIndex(d, CmdOperateBytes)
	di := bytes.LastIndex(d, cmdDot)

	if i != -1 && di != -1 {
		rd, err := this.handle(d[i+len(CmdOperateBytes):di],
			d[di+len(cmdDot):], d[:i])

		if err != nil {
			Println("[Server][Output]:%s ", err)
			conn.Write([]byte(err.Error()))
		} else {
			Println("[Server][Output]:%s ", string(rd))
			conn.Write(rd)
		}
	} else {
		conn.Write(invalidBytes)
	}
}

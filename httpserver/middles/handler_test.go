package middles

type Req struct {
	A int
}

type Rsp struct {
	A int
}

type TErr struct {
}

func (e *TErr) Error() string {
	return "terr"
}

// func TestHandler(t *testing.T) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			t.Error(err)
// 		}
// 	}()
// 	handle := func(ctx context.Context, req *Req, rsp *Rsp) *sdkError.CommonError {
// 		rsp.A = req.A
// 		return nil
// 	}
// 	CreateHandlerFunc(handle)

// 	{
// 		handle := func(ctx context.Context, req *Req, rsp *Rsp) error {
// 			return nil
// 		}
// 		CreateHandlerFunc(handle)
// 	}
// }

// func TestHandlerFailed(t *testing.T) {
// 	defer func() {
// 		if err := recover(); err == nil {
// 			t.Error("need panic while handle's ret is not error")
// 		} else {
// 			t.Log(err)
// 		}
// 	}()
// 	{
// 		handle := func(ctx context.Context, req *Req, rsp *Rsp) int {
// 			return 0
// 		}
// 		CreateHandlerFunc(handle)
// 	}
// 	{
// 		handle := func(ctx context.Context, req *Req, rsp *Rsp) *TErr {
// 			return nil
// 		}
// 		CreateHandlerFunc(handle)
// 	}

// }

// func TestCreateParam(t *testing.T) {
// 	testTrueCreateParam(t, WithDataKey(true))
// 	testTrueCreateParam(t)

// 	testFalseCreateParam(t, WithDataKey(false))
// }

// func testTrueCreateParam(t *testing.T, opts ...Option) {
// 	option := getOption(opts...)
// 	rsp := getRetData(reflect.ValueOf(Rsp{}), option)
// 	need := &RespStruct{Data: Rsp{}}
// 	if !reflect.DeepEqual(need, rsp) {
// 		t.Errorf("need: %v  get: %v", need, rsp)
// 	}
// }

// func testFalseCreateParam(t *testing.T, opts ...Option) {
// 	option := getOption(opts...)
// 	rsp := getRetData(reflect.ValueOf(Rsp{}), option)
// 	need := Rsp{}
// 	if !reflect.DeepEqual(need, rsp) {
// 		t.Errorf("need: %v  get: %v", need, rsp)
// 	}
// }

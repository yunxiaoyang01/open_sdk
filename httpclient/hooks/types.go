package hooks

import resty "gopkg.in/resty.v1"

type BeforeRequestHook func(*resty.Client, *resty.Request) error

type PreRequestHook func(*resty.Client, *resty.Request) error

type AfterResponseHook func(*resty.Client, *resty.Response) error

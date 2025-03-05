package commonerrors

import "github.com/valyala/fasthttp"


func FormError(ctx *fasthttp.RequestCtx, errMessage string, errStatusCode int){
	ctx.SetBodyString("{\"error\": \"" + errMessage + "\"}")
	ctx.SetStatusCode(errStatusCode)
}

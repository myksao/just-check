package http

func (justChk *justCheckHandler) MapJustCheckRoute() {
	justChk.router.POST("/check",justChk.Check())
}

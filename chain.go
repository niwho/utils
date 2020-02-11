package utils

import "math"

const (
	abortIndex int8 = math.MaxInt8 / 2
)

type HandlerFunc func(*Chain)
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

type errorMsgs []error

type Chain struct {
	handlers HandlersChain
	index    int8
	Errors   errorMsgs

	Val interface{}
}

func (c *Chain) Handler() HandlerFunc {
	return c.handlers.Last()
}

func (c *Chain) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Chain) Abort(err error) {
	c.index = abortIndex
	if err != nil {
		c.Errors = append(c.Errors, err)
	}
}

func (c *Chain) Uses(h ...HandlerFunc) *Chain {
	c.handlers = append(c.handlers, h...)
	return c
}

func (c *Chain) Use(h HandlerFunc) *Chain {
	c.handlers = append(c.handlers, h)
	return c
}

func (c *Chain) CheckErr() error {
	if len(c.Errors) > 0 {
		return c.Errors[0]
	}

	return nil
}

func (c *Chain) Do(val interface{}) error {
	c.Val = val
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
	if len(c.Errors) > 0 {
		return c.Errors[0]
	}

	return nil
}

/*
// sample
func GetVideoInfo(vid int64) (pv ProcessorVideo, err error) {
	err = (&utils.Chain{}).Use(func(c *utils.Chain) {
		obj := c.Val.(*ProcessorVideo)
		obj.VideoId = vid
		*obj, err = redisGetProcessorVideo(vid)
		if err == nil {
			c.Abort(nil)
			return
		} else if err != xredis.Nil {
			c.Abort(err)
			return
		}

		// redis没有命中
		c.Next() // 重新填充了obj
		err = c.CheckErr()
		if err == nil {
			err = redisSetProcessorVideo(*obj, time.Second*30)
		} else {
			// 30s 空数据
			_ = redisSetProcessorVideo(*obj, time.Second*30)
		}

	}).Use(func(c *utils.Chain) {
		// db处理
		obj := c.Val.(*ProcessorVideo)
		*obj, err = dbGetProcessorVideo(vid)
		if err != nil {
			c.Abort(err)
		}

	}).Do(&pv)

	return
}
*/

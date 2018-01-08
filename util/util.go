package util

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"../vars"
	"github.com/valyala/fasthttp"
)

// Pass 判断该请求是否直接pass（不可缓存）
func Pass(ctx *fasthttp.RequestCtx, passList [][]byte) bool {
	method := ctx.Method()
	if bytes.Compare(method, vars.Get) != 0 && bytes.Compare(method, vars.Head) != 0 {
		return true
	}
	isPass := false
	uri := ctx.RequestURI()
	for _, item := range passList {
		if !isPass && bytes.Contains(uri, item) {
			isPass = true
		}
	}
	return isPass
}

// GetCacheAge 获取s-maxage或者max-age的值
func GetCacheAge(header *fasthttp.ResponseHeader) uint32 {
	cacheControl := header.PeekBytes(vars.CacheControl)
	if len(cacheControl) == 0 {
		return 0
	}
	// 如果设置不可缓存，返回0
	reg, _ := regexp.Compile(`no-cache|no-store|private`)
	match := reg.Match(cacheControl)
	if match {
		return 0
	}

	// 优先从s-maxage中获取
	reg, _ = regexp.Compile(`s-maxage=(\d+)`)
	result := reg.FindSubmatch(cacheControl)
	if len(result) == 2 {
		maxAge, _ := strconv.Atoi(string(result[1]))
		return uint32(maxAge)
	}

	// 从max-age中获取缓存时间
	reg, _ = regexp.Compile(`max-age=(\d+)`)
	result = reg.FindSubmatch(cacheControl)
	if len(result) != 2 {
		return 0
	}
	maxAge, _ := strconv.Atoi(string(result[1]))
	return uint32(maxAge)
}

// SupportGzip 判断是否支持gzip
func SupportGzip(ctx *fasthttp.RequestCtx) bool {
	return ctx.Request.Header.HasAcceptEncodingBytes(vars.Gzip)
}

// SupportBr 判断是否支持brotli压缩
func SupportBr(ctx *fasthttp.RequestCtx) bool {
	return ctx.Request.Header.HasAcceptEncodingBytes(vars.Br)
}

// ConvertUint16ToBytes 将uint16转换为[]byte
func ConvertUint16ToBytes(v uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, v)
	return buf
}

// ConvertBytesToUint16 将[]byte转换为uint16
func ConvertBytesToUint16(buf []byte) uint16 {
	return binary.LittleEndian.Uint16(buf)
}

// ConvertUint32ToBytes 将uint32转换为[]byte
func ConvertUint32ToBytes(v uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	return buf
}

// ConvertBytesToUint32 将[]byte转换为uint32
func ConvertBytesToUint32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

// GetSeconds 获取当前时间的时间戳（秒）
func GetSeconds() uint32 {
	return uint32(time.Now().Unix())
}

// GetNowSecondsBytes 获取当时时间的字节表示(4个字节)
func GetNowSecondsBytes() []byte {
	return ConvertUint32ToBytes(GetSeconds())
}

// ConvertToSeconds 将字节保存的秒转换为整数
func ConvertToSeconds(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

// GenRequestKey 生成请求的key: Method + host + request uri
func GenRequestKey(ctx *fasthttp.RequestCtx) []byte {
	uri := ctx.URI()
	return bytes.Join([][]byte{
		ctx.Method(),
		uri.Host(),
		uri.RequestURI(),
	}, []byte(""))
}

// Gzip 对数据压缩
func Gzip(buf []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}

// Gunzip 解压数据
func Gunzip(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

// GetClientIP 获取客户端IP
func GetClientIP(ctx *fasthttp.RequestCtx) string {
	xFor := ctx.Request.Header.PeekBytes(vars.XForwardedFor)
	if len(xFor) == 0 {
		return ctx.RemoteIP().String()
	}
	arr := bytes.Split(xFor, []byte(","))
	return string(arr[0])
}

// ShouldCompress 判断该响应数据是否应该压缩(针对文本类压缩)
func ShouldCompress(contentType []byte) bool {
	// 检测是否为文本
	reg, _ := regexp.Compile(`text|application/javascript|application/x-javascript|application/json`)
	return reg.Match(contentType)
}

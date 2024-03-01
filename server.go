package mule

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

type Handler interface {
	ServeHTTP(*Response, *Request)
}

type Server struct {
	Addr    string
	Handler Handler
}

func NewServer(addr string, handler Handler) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
	}
}

func (s *Server) ListenAndServe() error {
	lsn, lsnErr := net.Listen("tcp", s.Addr)
	if lsnErr != nil {
		return lsnErr
	}
	defer lsn.Close()
	for {
		conn, err := lsn.Accept()
		if err != nil {
			return err
		}
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		// deal with new connection
		go func() {
			// read request
			req := NewRequest(conn)

			// new response
			rw := NewResponse(conn)

			s.Handler.ServeHTTP(rw, req)
		}()
	}
}

type Request struct {
	Method  string
	Path    string
	Query   map[string][]string
	Headers map[string]string
	Body    string
}

func NewRequest(conn net.Conn) *Request {
	readBF := make([]byte, 1024)
	reqContentBS := make([]byte, 0)
	for {
		readCount, err := conn.Read(readBF)
		if readCount != 0 {
			reqContentBS = append(reqContentBS, readBF[0:readCount]...)
		}
		if err != nil {
			break
		}
	}
	reqContent := string(reqContentBS)
	// 定义正则表达式来匹配HTTP请求的首行、头部和主体
	requestLineRegex := regexp.MustCompile(`^(\w+)\s(.*?)(\?.*)?\sHTTP/\d\.\d`)
	headerRegex := regexp.MustCompile(`^([\w-]+):\s(.*?)$`)
	bodyRegex := regexp.MustCompile(`\r\n\r\n((.|\n)*)$`)

	// 使用正则表达式解析HTTP请求内容
	lines := strings.Split(reqContent, "\r\n")
	var requestMethod, requestPath, body string
	headers := make(map[string]string)
	queryParams := make(map[string][]string)

	// 解析首行
	if matches := requestLineRegex.FindStringSubmatch(lines[0]); matches != nil {
		requestMethod = matches[1]
		requestPath = matches[2]

		// 解析查询参数
		if len(matches) > 3 && matches[3] != "" {
			queryParamsStr := matches[3][1:] // 去除 '?' 符号
			queryParamsArr := strings.Split(queryParamsStr, "&")
			for _, param := range queryParamsArr {
				pair := strings.Split(param, "=")
				if len(pair) == 2 {
					queryParams[pair[0]] = append(queryParams[pair[0]], pair[1])
				}
			}
		}
	}

	// 解析头部
	for _, line := range lines[1:] {
		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			headers[matches[1]] = matches[2]
		} else {
			break
		}
	}

	// 解析主体
	if matches := bodyRegex.FindStringSubmatch(reqContent); matches != nil {
		body = matches[1]
	}
	return &Request{
		Method:  requestMethod,
		Path:    requestPath,
		Query:   queryParams,
		Headers: headers,
		Body:    body,
	}
}

type Response struct {
	conn       net.Conn
	statusCode int
	headers    map[string]string
	body       string
}

func NewResponse(conn net.Conn) *Response {
	return &Response{
		conn:    conn,
		headers: make(map[string]string),
	}
}

func (r *Response) SetStatusCode(statusCode int) {
	r.statusCode = statusCode
}

func (r *Response) SetHeader(key, value string) {
	r.headers[key] = value
}

func (r *Response) SetBody(body string) {
	r.body = body
}

// FIXME(wangli) 结果 body 中多了两个前置空行
func (r *Response) Flush() {
	// build
	resp := make([]string, 0)
	resp = append(resp, fmt.Sprintf("HTTP/1.1 %d", r.statusCode))
	for key, value := range r.headers {
		resp = append(resp, fmt.Sprintf("%s: %s", key, value))
	}
	resp = append(resp, "\r\n"+r.body)

	// write
	r.conn.Write([]byte(strings.Join(resp, "\r\n")))

	// close
	r.conn.Close()
}

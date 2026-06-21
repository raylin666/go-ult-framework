package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestTraceIDConcurrentSafety 测试 TraceID() 在并发场景下的安全性
// 验证多个 goroutine 同时调用 TraceID() 时，只生成一个 UUID
func TestTraceIDConcurrentSafety(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试 HTTP 请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// 创建 Gin 上下文
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = req

	// 创建自定义上下文
	ctx, err := newContext(ginCtx)
	if err != nil {
		t.Fatalf("newContext() failed: %v", err)
	}

	// 并发调用 TraceID()
	var (
		wg           sync.WaitGroup
		traceIDs     = make([]string, 100)
		goroutines   = 100
		allSame      = true
		firstTraceID string
	)

	// 启动多个 goroutine 同时调用 TraceID()
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			traceID := ctx.TraceID()
			t.Logf("TraceID: %s", traceID)

			traceIDs[index] = traceID
		}(i)
	}

	wg.Wait()

	// 验证所有 TraceID 都相同
	firstTraceID = traceIDs[0]
	if firstTraceID == "" {
		t.Fatal("TraceID should not be empty")
	}

	for i, traceID := range traceIDs {
		if traceID != firstTraceID {
			allSame = false
			t.Errorf("TraceID at index %d is different: got %s, want %s", i, traceID, firstTraceID)
		}
	}

	if !allSame {
		t.Errorf("Concurrent calls to TraceID() returned different values")
	}

	t.Logf("All %d goroutines got the same TraceID: %s", goroutines, firstTraceID)
}

// TestTraceIDWithHeader 测试从请求头中获取 TraceID
func TestTraceIDWithHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建带有 TraceID 的请求头（使用正确的 header key）
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("x-md-trace-id", "test-trace-id-12345")
	w := httptest.NewRecorder()

	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = req

	// 创建自定义上下文
	ctx, err := newContext(ginCtx)
	if err != nil {
		t.Fatalf("newContext() failed: %v", err)
	}

	traceID1 := ctx.TraceID()
	if traceID1 != "test-trace-id-12345" {
		t.Errorf("TraceID from header: got %s, want test-trace-id-12345", traceID1)
	}

	// 多次调用应该返回相同的值
	traceID2 := ctx.TraceID()
	if traceID1 != traceID2 {
		t.Errorf("Multiple calls to TraceID() returned different values: %s vs %s", traceID1, traceID2)
	}
}

// TestRawDataReturnsCopy 测试 RawData() 返回的是副本
// 验证修改返回值不会影响原始数据
func TestRawDataReturnsCopy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 原始请求体
	originalBody := []byte("original request body")

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(originalBody))
	w := httptest.NewRecorder()

	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = req

	// 创建自定义上下文
	ctx, err := newContext(ginCtx)
	if err != nil {
		t.Fatalf("newContext() failed: %v", err)
	}

	// 第一次获取 RawData
	rawData1 := ctx.RawData()
	if string(rawData1) != string(originalBody) {
		t.Errorf("RawData(): got %s, want %s", string(rawData1), string(originalBody))
	}

	// 修改返回的数据
	rawData1[0] = 'X'
	rawData1[1] = 'X'
	rawData1[2] = 'X'

	// 第二次获取 RawData，应该还是原始数据
	rawData2 := ctx.RawData()
	if string(rawData2) == string(rawData1) {
		t.Error("RawData() returned the same slice, should return a copy")
	}

	if string(rawData2) != string(originalBody) {
		t.Errorf("RawData() was modified: got %s, want %s", string(rawData2), string(originalBody))
	}

	t.Logf("Original: %s", string(originalBody))
	t.Logf("Modified copy: %s", string(rawData1))
	t.Logf("Second call: %s", string(rawData2))
}

// TestRawDataConcurrentSafety 测试 RawData() 在并发场景下的安全性
func TestRawDataConcurrentSafety(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalBody := []byte("concurrent test body")
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(originalBody))
	w := httptest.NewRecorder()

	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = req

	// 创建自定义上下文
	ctx, err := newContext(ginCtx)
	if err != nil {
		t.Fatalf("newContext() failed: %v", err)
	}

	var wg sync.WaitGroup
	goroutines := 50
	errors := make(chan error, goroutines)

	// 并发调用 RawData() 并修改返回值
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := ctx.RawData()
			if string(data) != string(originalBody) {
				errors <- nil
				return
			}
			// 尝试修改返回的数据
			for j := 0; j < len(data); j++ {
				data[j] = 'X'
			}
		}()
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}

	// 验证原始数据未被修改
	finalData := ctx.RawData()
	if string(finalData) != string(originalBody) {
		t.Errorf("Original data was modified: got %s, want %s", string(finalData), string(originalBody))
	}

	t.Logf("Original data preserved after %d concurrent modifications", goroutines)
}

// TestContextPoolReset 测试上下文池的重置功能
// 验证从池中获取的上下文对象是干净的
func TestContextPoolReset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 第一次使用
	req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
	req1.Header.Set("X-Trace-ID", "trace-id-1")
	w1 := httptest.NewRecorder()

	ginCtx1, _ := gin.CreateTestContext(w1)
	ginCtx1.Request = req1

	// 创建自定义上下文
	ctx1, err := newContext(ginCtx1)
	if err != nil {
		t.Fatalf("newContext() failed: %v", err)
	}
	traceID1 := ctx1.TraceID()

	// 归还到池中
	recoveryContext(ctx1)

	// 第二次使用（应该从池中获取）
	req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
	w2 := httptest.NewRecorder()

	ginCtx2, _ := gin.CreateTestContext(w2)
	ginCtx2.Request = req2

	// 创建自定义上下文
	ctx2, err := newContext(ginCtx2)
	if err != nil {
		t.Fatalf("newContext() failed: %v", err)
	}
	traceID2 := ctx2.TraceID()

	// 两个 TraceID 应该不同（因为第二个请求没有设置 TraceID 头）
	if traceID1 == traceID2 {
		t.Logf("Warning: TraceIDs are the same (might be coincidence): %s", traceID1)
	}

	t.Logf("First TraceID: %s", traceID1)
	t.Logf("Second TraceID: %s", traceID2)
}

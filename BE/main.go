package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mw "kitchen/MW"
	"kitchen/handler"
	"kitchen/store"
)

func main() {
	h := &handler.Handler{Store: store.NewMemoryStore()}

	var app http.Handler = http.TimeoutHandler(h.Routes(), 10*time.Second,
		"ครัวใช้เวลานานเกินไป")
	app = mw.Logging(mw.Recovery(mw.CORS(app)))

	srv := &http.Server{Addr: ":8080", Handler: app}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	log.Println("ครัวเปิดที่ :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // ยืนรอสัญญาณปิดร้าน
	log.Println("ได้รับสัญญาณปิดร้าน ไม่รับใบสั่งใหม่แล้ว")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("ปิดไม่ทัน", err)
	}
	log.Println("ปิดร้านเรียบร้อย ลูกค้ากลับหมดแล้ว")
}

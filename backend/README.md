# Backend Go Fiber

ต้องใช้ Go 1.25+  
คำสั่งตรวจพื้นฐาน:

```sh
go test ./...
go vet ./...
```

`cmd/server` เป็น entry point ขนาดเล็ก ทำหน้าที่โหลด config, รับ OS signals และเรียก `bootstrap.Run`

## โครงสร้าง package

- `internal/config`: โหลดค่า environment, default และ validation
- `internal/bootstrap`: init database, wiring dependencies, start HTTP และ graceful shutdown
- `internal/collector`: UDP/TCP Syslog listeners และ connection lifecycle
- `internal/retention`: cleanup worker ที่รันทันทีตอน startup และหลังจากนั้นทุก 1 นาทีใน demo/test config ปัจจุบัน
- `internal/router`: ตั้งค่า Fiber server และ route table ใน `routes.go`
- `internal/handler`: HTTP handlers แยกตาม feature เช่น auth, ingest, logs, alerts, health
- `internal/middleware`: authentication, Admin check, ingest authorization, Origin check และ error handling
- `internal/model`: shared data structures และ JSON contracts ไม่มี Fiber/SQL dependency
- `internal/normalize`: provider mapping, Syslog normalization และ JSON batch decoding
- `internal/repository`: Store interface, PostgreSQL implementation และ embedded schema

## ลำดับ dependency

```text
bootstrap → router → handler/middleware → repository → model
```

handlers และ collector ใช้:

```text
normalize → model
```

routes คง URL และ middleware ordering ไว้ชัดเจน Tests ใต้ `router` ตรวจ public HTTP behavior ส่วน normalization tests อยู่กับ package `normalize`

## การรัน

ใช้ Docker Compose จาก root project เป็นหลัก หรือถ้าจะรัน backend ตรง ๆ ให้ตั้ง environment ตาม `../docs/setup_appliance.md` แล้วรัน:

```sh
go run ./cmd/server
```

## Schema ฐานข้อมูล

schema เป็น idempotent bootstrap ไม่ใช่ migration framework แบบ versioned ถ้ามี schema change ในอนาคตควรทำ migration แยก

ตอน startup ระบบจะ insert demo accounts/tenants ที่ยังไม่มีอยู่เท่านั้น และจะไม่ rotate credentials เดิมแบบเงียบ ๆ

# Backend

Backend ของ Logdesk เขียนด้วย Go + Fiber และทำหน้าที่เป็นศูนย์กลางของระบบ ได้แก่ authentication, RBAC, log ingestion, normalization, search, dashboard stats, alert และ retention worker

ต้องใช้ Go 1.25+

## คำสั่งที่ใช้บ่อย

```sh
go test ./...
go vet ./...
go run ./cmd/server
```

โดยปกติแนะนำให้รันผ่าน Docker Compose จาก root project

```sh
docker compose up --build -d
```

## โครงสร้าง package

| Package | หน้าที่ |
|---|---|
| `cmd/server` | entry point ของ backend |
| `internal/config` | โหลดและ validate environment variables |
| `internal/bootstrap` | init database, wire dependencies, start HTTP/syslog/worker |
| `internal/collector` | Syslog UDP/TCP listeners |
| `internal/retention` | cleanup worker สำหรับ log เก่าและ expired sessions |
| `internal/router` | Fiber app และ route table |
| `internal/handler` | HTTP handlers แยกตาม feature |
| `internal/middleware` | auth, admin check, ingest authorization, origin check และ error handling |
| `internal/model` | shared data structures และ JSON contracts |
| `internal/normalize` | provider mapping และ normalization pipeline |
| `internal/repository` | Store interface, PostgreSQL implementation และ embedded schema |

Dependency หลัก

```text
bootstrap -> router -> handler/middleware -> repository -> model
collector -> normalize -> model
handler -> normalize -> model
```

## Environment สำคัญ

| Variable | ใช้สำหรับ |
|---|---|
| `DATABASE_URL` | connection string ไป PostgreSQL |
| `ADMIN_PASSWORD` | password ของ admin demo accounts |
| `VIEWER_PASSWORD` | password ของ viewer demo accounts |
| `API_KEY_A` | ingest API key สำหรับ tenant `demo-a` |
| `API_KEY_B` | ingest API key สำหรับ tenant `demo-b` |
| `APP_ORIGIN` | origin ที่อนุญาตสำหรับ cookie-authenticated writes |
| `COOKIE_SECURE` | เปิด Secure cookie เมื่อใช้ HTTPS |
| `RETENTION_DAYS` | จำนวนวันที่เก็บ log ขั้นต่ำ ค่า default คือ 7 |
| `SYSLOG_TENANT` | tenant ที่ใช้กับ Syslog collector |

## Tenant และ RBAC

Backend enforce tenant ทุก request

- UI ใช้ tenant จาก session หลัง login
- HTTP ingest ใช้ tenant จาก `X-API-Key`
- Syslog ใช้ tenant จาก `SYSLOG_TENANT`

Admin เป็น admin เฉพาะ tenant ของตัวเอง ส่วน Viewer อ่านได้อย่างเดียวใน tenant ของตัวเอง

## Alert

Alert ปัจจุบันมีเงื่อนไขหลักเดียว คือพบ log field ที่ rule กำหนดครบ **5 ครั้งภายใน 5 นาที** ภายใน tenant เดียวกัน ตัวอย่าง demo คือ failed login จาก IP เดียวกัน 5 ครั้ง

Alert ถูกแสดงใน UI และ Admin สามารถ acknowledge ได้ ยังไม่มี email/webhook delivery

## Retention

Retention worker รันตอน backend start และรันเป็นรอบตาม interval ปัจจุบันตั้งไว้ 1 นาทีเพื่อ demo/test

- ลบ logs ที่ `ingested_at` เก่ากว่า `RETENTION_DAYS`
- ค่า retention ขั้นต่ำคือ 7 วัน
- ลบ expired sessions ด้วย

## Database schema

Schema เป็น idempotent bootstrap ที่ฝังไว้ใน backend ไม่ใช่ migration framework แบบ versioned ถ้ามี schema change ระยะยาวควรเพิ่ม migration แยก

ตอน startup ระบบ seed demo tenants/users/API keys เฉพาะรายการที่ยังไม่มีอยู่ การเปลี่ยนค่า `.env` ภายหลังไม่ rotate credentials เดิมใน database อัตโนมัติ

# คู่มือทดสอบ API ด้วย Postman / Insomnia

ไฟล์ collection อยู่ที่ [`docs/postman_collection.json`](postman_collection.json) ใช้สำหรับกดทดสอบ API โดยไม่ต้องพิมพ์ request เอง

## 1. Import collection เข้าโปรแกรม

ใน Postman:

1. กด **Import**
2. เลือกไฟล์ `docs/postman_collection.json`
3. เปิด collection ชื่อ **Logdesk Log Management Demo**

ใน Insomnia:

1. กด **Create → Import from File**
2. เลือกไฟล์ `docs/postman_collection.json`

## 2. ตั้งค่า variables

เปิดไฟล์ `.env` แล้วคัดลอกค่าเหล่านี้ไปใส่ใน collection variables:

| Variable | ใส่ค่าอะไร |
|---|---|
| `base_url` | `http://localhost:8080` |
| `admin_email` | `admin.a@demo.local` |
| `admin_password` | ค่า `ADMIN_PASSWORD` จาก `.env` |
| `api_key_a` | ค่า `API_KEY_A` จาก `.env` |
| `api_key_b` | ค่า `API_KEY_B` จาก `.env` |
| `alert_id` | ใส่ภายหลังจาก `GET /api/alerts` |

อย่า commit ค่า password/API key จริงลง Git

## 3. ลำดับทดสอบที่แนะนำ

### ตรวจระบบ

1. เปิดระบบด้วย `.\run.ps1` หรือ `sh run.sh`
2. กด request **Health → GET /api/health**
3. ควรได้ response ว่าระบบพร้อมใช้งาน

### ทดสอบ login/session

1. กด **Session auth → Login as Admin A**
2. Postman จะเก็บ session cookie ให้อัตโนมัติ
3. กด **Session auth → Current user**
4. ควรเห็น user `admin.a@demo.local` และ tenant `demo-a`

### ทดสอบ HTTP ingest ด้วย API key

1. กด **Ingest with API key → POST /ingest single API log to tenant A**
2. กด **Session-protected search and dashboard → GET /api/logs search postman**
3. ควรเจอ event ที่เพิ่งส่งเข้าไป

### ทดสอบ alert

1. กด **POST /ingest batch to trigger alert**
2. กด **GET /api/alerts**
3. ควรเห็น alert จาก IP `203.0.113.77`
4. นำ `id` ของ alert ไปใส่ variable `alert_id`
5. กด **POST /api/alerts/:id/acknowledge**

### ทดสอบ tenant isolation

1. Login เป็น Admin A
2. กด **POST /ingest tenant B example** เพื่อส่ง log เข้า tenant B ด้วย `API_KEY_B`
3. กด **GET /api/logs denied cross-tenant query**
4. ควรได้ `403` เพราะ session ของ Admin A อยู่ `demo-a`

## 4. Payload ตัวอย่าง

### Login

ใช้กับ request **Session auth → Login as Admin A**

```json
{
  "email": "{{admin_email}}",
  "password": "{{admin_password}}"
}
```

หลังส่ง request นี้ Postman จะเก็บ session cookie ไว้อัตโนมัติ

### ส่ง log เดี่ยวเข้า tenant A

ใช้กับ request **POST /ingest single API log to tenant A**

Header:

```http
X-API-Key: {{api_key_a}}
Content-Type: application/json
```

Body:

```json
{
  "source": "api",
  "event_type": "postman_app_login_failed",
  "user": "alice",
  "ip": "203.0.113.7",
  "reason": "wrong_password"
}
```

ผลที่คาดหวัง:

```json
{
  "accepted": 1,
  "tenant": "demo-a"
}
```

### ส่ง batch เพื่อ trigger alert

ใช้กับ request **POST /ingest batch to trigger alert**

Header:

```http
X-API-Key: {{api_key_a}}
Content-Type: application/json
```

Body:

```json
[
  {
    "source": "api",
    "event_type": "login_failed",
    "ip": "203.0.113.77",
    "user": "alice"
  },
  {
    "source": "api",
    "event_type": "login_failed",
    "ip": "203.0.113.77",
    "user": "bob"
  },
  {
    "source": "api",
    "event_type": "login_failed",
    "ip": "203.0.113.77",
    "user": "carol"
  },
  {
    "source": "api",
    "event_type": "login_failed",
    "ip": "203.0.113.77",
    "user": "dave"
  },
  {
    "source": "api",
    "event_type": "login_failed",
    "ip": "203.0.113.77",
    "user": "erin"
  }
]
```

ผลที่คาดหวัง:

```json
{
  "accepted": 5,
  "tenant": "demo-a"
}
```

จากนั้นกด `GET /api/alerts` ควรเห็น alert ของ IP `203.0.113.77`

### ส่ง log เข้า tenant B

ใช้กับ request **POST /ingest tenant B example**

Header:

```http
X-API-Key: {{api_key_b}}
Content-Type: application/json
```

Body:

```json
{
  "source": "api",
  "event_type": "tenant_b_postman_event",
  "user": "beacon-user",
  "ip": "198.51.100.20"
}
```

ผลที่คาดหวัง:

```json
{
  "accepted": 1,
  "tenant": "demo-b"
}
```

### ตัวอย่างที่ควร fail เพราะ tenant ไม่ตรง

ใช้กับ request **POST /ingest mismatched tenant should fail**

Header ใช้ `API_KEY_A` แต่ body ระบุ `tenant` เป็น `demo-b`

```http
X-API-Key: {{api_key_a}}
Content-Type: application/json
```

```json
{
  "tenant": "demo-b",
  "source": "api",
  "event_type": "should_be_rejected"
}
```

ผลที่คาดหวัง:

```json
{
  "error": "Log 1: tenant mismatch"
}
```

หรือได้ error 400 ที่สื่อว่า tenant ไม่ตรงกัน ขึ้นอยู่กับข้อความ validation ปัจจุบันของ backend

## 5. Auth แต่ละแบบใช้ตอนไหน

| การทดสอบ | ใช้อะไร |
|---|---|
| Login, logs, stats, alerts, rule | session cookie หลัง login |
| `POST /ingest` จากระบบภายนอก | `X-API-Key` |
| Syslog UDP/TCP | ไม่ใช้ collection นี้ ใช้ `samples/send_syslog.py` |

## 6. หมายเหตุ

- ถ้า `POST /api/rule` หรือ acknowledge ได้ `403` ให้ login เป็น admin ก่อน
- ถ้า search ไม่เจอ log ให้ดู time range/filter ใน UI หรือกด `GET /api/logs latest`
- API key A เขียนเข้า `demo-a`; API key B เขียนเข้า `demo-b`
- ถ้า payload ใส่ tenant ไม่ตรงกับ API key ระบบจะ reject เพื่อกัน cross-tenant write

# Samples

ไฟล์ในโฟลเดอร์นี้เป็นข้อมูลสังเคราะห์สำหรับ demo สอง tenant หรือสองบริษัท

| Tenant | บริษัทสมมติ | Admin | Viewer | API key |
|---|---|---|---|---|
| `demo-a` | Atlas Technology | `admin.a@demo.local` | `viewer.a@demo.local` | `API_KEY_A` |
| `demo-b` | Beacon Retail | `admin.b@demo.local` | `viewer.b@demo.local` | `API_KEY_B` |

Admin ทั้งสองใช้ `ADMIN_PASSWORD` จาก `.env` ที่ root project ส่วน Viewer ทั้งสองใช้ `VIEWER_PASSWORD`

Admin และ Viewer ถูกผูกกับ tenant ของตัวเองเท่านั้น field `company` ใน sample เป็น label ให้อ่านง่าย ไม่ใช่ตัวกำหนดสิทธิ์จริง

## โครงสร้างไฟล์

แต่ละ tenant folder มีไฟล์หลัก

- `events.json`: events ปกติสำหรับหลาย source เช่น API, Firewall, Network, AWS, Microsoft 365 และ AD
- `abnormal.json`: failed login 5 รายการจาก IP เดียวกัน ใช้ trigger alert เงื่อนไข 5 ครั้งภายใน 5 นาที

ตัวอย่าง tenant A

```text
samples/tenant-a/events.json
samples/tenant-a/abnormal.json
```

ตัวอย่าง tenant B

```text
samples/tenant-b/events.json
samples/tenant-b/abnormal.json
```

ข้อมูลทั้งหมดเป็น synthetic data และตั้งใจ omit timestamp เพื่อให้ระบบใช้ ingestion time ปัจจุบัน

## วิธี seed ข้อมูลสองบริษัท

ต้องเปิด stack อยู่ก่อน จากนั้นรันจาก root project

```powershell
py samples/seed_companies.py
```

script จะทำงานดังนี้

1. อ่าน `.env` ที่ root project
2. ส่ง events ของแต่ละ tenant ด้วย API key ของ tenant นั้น
3. login เป็น Viewer ของแต่ละ tenant เพื่อตรวจ data isolation และ write restrictions

script เป็นแบบ append data ถ้ารันซ้ำจะเพิ่ม batch ใหม่ ไม่ลบ logs หรือ credentials เดิม

## ส่ง sample ผ่าน HTTP API

```powershell
py samples/post_logs.py --file samples/tenant-a/events.json
py samples/post_logs.py --file samples/tenant-a/abnormal.json
```

`post_logs.py` อ่าน API key จาก `.env` ได้ จึงไม่ต้องกรอก key เองถ้าใช้ค่า default

## ส่ง Syslog sample

UDP

```powershell
py samples/send_syslog.py --host localhost --port 5514
```

TCP

```powershell
py samples/send_syslog.py --tcp --host localhost --port 5514
```

ถ้าส่งไป SaaS demo ให้ใช้ host

```text
logdisk.malaysiawest.cloudapp.azure.com
```

## Upload ผ่าน UI

ถ้าจะ upload เองผ่าน UI ให้ใช้ Admin ของ tenant ที่ตรงกับไฟล์

- `demo-a`: ใช้ `admin.a@demo.local`
- `demo-b`: ใช้ `admin.b@demo.local`

Viewer upload ไม่ได้

ถ้าแก้ payload ให้มี `tenant` ไม่ตรงกับ credential ระบบจะ reject

## วิธีตรวจใน UI

เปิด Log Explorer ด้วย user ของ tenant ที่ต้องการ แล้วเลือกช่วงเวลา Last 24 hours จากนั้นค้นหาคำเหล่านี้ได้

```text
company-demo
Atlas Technology
Beacon Retail
```

ถ้า login เป็น Viewer A ต้องเห็นเฉพาะข้อมูลของ `demo-a` และ Viewer B ต้องเห็นเฉพาะข้อมูลของ `demo-b`

## Alert sample

`abnormal.json` ใช้ทดสอบ alert ปัจจุบัน ซึ่งแจ้งเตือนเมื่อเจอ log field ที่ rule กำหนดครบ 5 ครั้งภายใน 5 นาที

ถ้า alert rule อยู่ใน cooldown การส่ง failed-login ซ้ำอาจใช้ alert เดิม หรือถ้า rule ถูกเปลี่ยน/ปิดไว้ อาจไม่ trigger

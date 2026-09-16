# Sample สำหรับสองบริษัท

ไฟล์ในโฟลเดอร์นี้เป็นข้อมูลสังเคราะห์สำหรับ demo สอง tenant หรือสองบริษัท

| Tenant | บริษัทสมมติ | Admin | Viewer | Credential สำหรับ ingest |
|---|---|---|---|---|
| demo-a | Atlas Technology | admin.a@demo.local | viewer.a@demo.local | API_KEY_A |
| demo-b | Beacon Retail | admin.b@demo.local | viewer.b@demo.local | API_KEY_B |

Admin ทั้งสองใช้ `ADMIN_PASSWORD` จาก `.env` ที่ root project ส่วน Viewer ทั้งสองใช้ `VIEWER_PASSWORD`

Admin แต่ละคนอยู่เฉพาะ tenant ของตัวเอง field company เป็น label เพื่อให้ sample อ่านง่ายเท่านั้น การแยกสิทธิ์จริง enforce จาก tenant ของ credential ไม่ได้เชื่อ label ใน payload

## โครงสร้างไฟล์

แต่ละ tenant folder มีไฟล์:

- `events.json`: 7 events ครอบคลุม API, Firewall, Network, CrowdStrike, AWS, Microsoft 365 และ Windows/AD
- `abnormal.json`: 5 failed logins จาก IP เดียวกัน ใช้ trigger default alert rule แบบ 5 ครั้งใน 5 นาที

Atlas ใช้:

- IP ช่วง `10.10.x.x`
- host prefix `atlas-*`
- user domain `atlas.example`
- cloud account `111111111111`

Beacon ใช้:

- IP ช่วง `10.20.x.x`
- host prefix `beacon-*`
- user domain `beacon.example`
- cloud account `222222222222`

ข้อมูลทั้งหมดเป็น synthetic data และตั้งใจ omit timestamp เพื่อให้ระบบใช้ ingestion time ปัจจุบัน

## วิธี seed ข้อมูลสองบริษัท

ต้องเปิด stack อยู่ก่อน จากนั้นรันจาก root project:

```powershell
python samples/seed_companies.py
```

script จะ:

1. อ่าน `.env` ที่ root project
2. ส่ง events ของแต่ละ tenant ด้วย API key ของ tenant นั้น
3. login เป็น Viewer ของแต่ละ tenant เพื่อตรวจ data isolation และ write restrictions

script เป็นแบบ append data ถ้ารันซ้ำจะเพิ่ม batch ใหม่ ไม่ลบ logs หรือ credentials เดิม

ถ้า alert rule อยู่ใน cooldown การส่ง failed-login ซ้ำอาจใช้ alert เดิม ถ้า rule ถูกเปลี่ยนหรือปิดไว้อาจไม่ trigger

## Upload ผ่าน UI แบบ manual

ถ้าจะ upload เองผ่าน UI ให้ใช้ Admin ของ tenant ที่ตรงกับไฟล์:

- `demo-a`: ใช้ `admin.a@demo.local`
- `demo-b`: ใช้ `admin.b@demo.local`

Viewer upload ไม่ได้

ถ้าเปลี่ยนแค่ field `tenant` ใน payload แต่ใช้ credential ของอีก tenant ระบบจะ reject

## วิธีตรวจใน UI

เปิด Log explorer ด้วย Viewer ที่ตรงกับ tenant แล้วเลือก:

- All sources
- Last 24 hours

จากนั้น search คำเหล่านี้ได้:

```text
company-demo
Atlas Technology
Beacon Retail
```

sidebar ยังแสดง tenant ID เดิม เช่น `demo-a` / `demo-b` sample นี้ไม่ได้ rename tenant records

generic samples และ Nginx live feed อาจอยู่ใน `demo-a` ด้วยถ้าไม่ได้ filter แยก

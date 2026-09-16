# CI/CD

โปรเจกต์นี้ใช้ GitHub Actions สำหรับตรวจ build/test อัตโนมัติ และ deploy ไป Azure VM ผ่าน SSH

## Workflow ที่มี

| Workflow | ไฟล์ | ทำงานเมื่อไหร่ | หน้าที่ |
|---|---|---|---|
| CI | `.github/workflows/ci.yml` | push หรือ pull request เข้า `main` | ตรวจ backend, frontend และ Docker Compose config |
| Deploy to Azure VM | `.github/workflows/deploy-azure.yml` | push เข้า branch `deploy` หรือกด manual | SSH เข้า VM แล้ว rebuild/restart container |

## CI

CI ตรวจรายการหลักเหล่านี้

- Backend: `go test ./...`
- Backend: `go vet ./...`
- Frontend: `npm test`
- Frontend: `npm run build`
- Docker Compose: `docker compose config`

ใช้ workflow นี้เพื่อให้มั่นใจว่า commit บน `main` build ผ่านก่อนนำไป deploy

## CD ไป Azure VM

Deploy workflow ทำงานอัตโนมัติเมื่อ push เข้า branch `deploy`

สิ่งที่ workflow ทำ

1. ใช้ SSH key จาก GitHub Secrets
2. SSH เข้า Azure VM
3. เข้า directory `$HOME/LogManagement`
4. ดึง code ล่าสุดจาก `origin/deploy`
5. รัน `docker compose up -d --build`
6. แสดง `docker compose ps`

Workflow ไม่สร้าง `.env` และไม่ reset database volume เพื่อป้องกัน secret รั่วและข้อมูลหาย

## GitHub Secrets ที่ต้องตั้ง

ไปที่ GitHub repository

```text
Settings -> Secrets and variables -> Actions -> New repository secret
```

เพิ่ม secrets เหล่านี้

| Secret | ค่า |
|---|---|
| `AZURE_HOST` | domain หรือ public IP ของ VM เช่น `logdisk.malaysiawest.cloudapp.azure.com` |
| `AZURE_USER` | user SSH เช่น `azureuser` |
| `AZURE_SSH_KEY` | private key ที่ใช้ SSH เข้า VM |

ห้ามใส่ `.env`, password, API key หรือ private key ลงใน Git

## เตรียม Azure VM ให้ deploy ได้

บน VM ต้องมี repository อยู่ที่

```text
$HOME/LogManagement
```

และต้องมี `.env` จริงอยู่บน VM แล้ว ตัวอย่างค่าหลัก

```env
APP_ORIGIN=https://logdisk.malaysiawest.cloudapp.azure.com
COOKIE_SECURE=true
HTTP_BIND=127.0.0.1
SYSLOG_BIND=0.0.0.0
SYSLOG_TENANT=demo-a
```

ตรวจ Docker

```sh
docker compose version
docker ps
```

ถ้า `docker ps` ติด permission ให้เพิ่ม user เข้า group docker แล้ว logout/login ใหม่

```sh
sudo usermod -aG docker $USER
```

## วิธี deploy ผ่าน branch

หลังจากงานบน `main` พร้อมแล้ว ให้ fast-forward หรือ merge เข้า branch `deploy`

```sh
git checkout deploy
git merge main
git push origin deploy
```

เมื่อ push เข้า `deploy` แล้ว GitHub Actions จะ deploy ไป Azure VM อัตโนมัติ

ถ้ายังไม่มี branch `deploy`

```sh
git checkout -b deploy main
git push -u origin deploy
```

## วิธี deploy แบบ manual

ใน GitHub ไปที่

```text
Actions -> Deploy to Azure VM -> Run workflow
```

หลัง deploy เสร็จให้เปิด

```text
https://logdisk.malaysiawest.cloudapp.azure.com
```

## วิธีทดสอบ CI/CD

1. แก้ไฟล์เล็ก ๆ เช่น README หรือข้อความหน้าเว็บ
2. commit และ push ไป `main`
3. ดู GitHub Actions ว่า CI ผ่าน
4. merge หรือ fast-forward ไป branch `deploy`
5. push `deploy`
6. รอ workflow deploy สำเร็จ
7. เปิด SaaS URL เพื่อตรวจผล

## หมายเหตุ

- Workflow ไม่ลบ volume และไม่ reset database
- ถ้าเปลี่ยนค่า `.env` ให้แก้บน VM โดยตรง แล้วรัน deploy หรือ `docker compose up -d --build` ใหม่
- ถ้า deploy fail ให้ดู log ใน GitHub Actions และดู container log บน VM ด้วย `docker compose logs --tail=120 backend`

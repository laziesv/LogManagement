# CI/CD

โปรเจกต์นี้ใช้ GitHub Actions สำหรับตรวจ build/test อัตโนมัติ และ deploy ไป Azure VM แบบกดรันเอง

## CI

ไฟล์ workflow:

```text
.github/workflows/ci.yml
```

CI ทำงานเมื่อ push หรือเปิด pull request เข้า branch `main`

งานที่ตรวจ:

- backend: `go test ./...`
- backend: `go vet ./...`
- frontend: `npm test`
- frontend: `npm run build`
- Docker Compose: `docker compose config`

## CD ไป Azure VM

ไฟล์ workflow:

```text
.github/workflows/deploy-azure.yml
```

workflow นี้เป็นแบบ manual trigger ผ่าน GitHub Actions tab เพื่อกัน deploy ทุกครั้งที่ push

สิ่งที่ workflow ทำ:

1. SSH เข้า Azure VM
2. เข้า directory `$HOME/LogManagement`
3. pull code ล่าสุดจาก `origin/main`
4. รัน `docker compose up -d --build`
5. แสดง `docker compose ps`

## GitHub Secrets ที่ต้องตั้ง

ไปที่ GitHub repository:

```text
Settings → Secrets and variables → Actions → New repository secret
```

เพิ่ม secrets:

| Secret | ค่า |
|---|---|
| `AZURE_HOST` | domain หรือ public IP ของ VM เช่น `logdisk.malaysiawest.cloudapp.azure.com` |
| `AZURE_USER` | user SSH เช่น `azureuser` |
| `AZURE_SSH_KEY` | private key ที่ใช้ SSH เข้า VM |

อย่าใส่ `.env`, password, API key หรือ private key ลง Git

## เตรียม Azure VM ให้ deploy ได้

บน VM ต้องมี repository อยู่ที่:

```text
$HOME/LogManagement
```

และต้องตั้ง `.env` บน VM ไว้แล้ว เช่น:

```env
APP_ORIGIN=https://logdisk.malaysiawest.cloudapp.azure.com
COOKIE_SECURE=true
HTTP_BIND=127.0.0.1
SYSLOG_BIND=0.0.0.0
```

VM ต้องติดตั้ง Docker Compose plugin และ user ที่ SSH เข้าไปต้องรัน Docker ได้:

```sh
docker compose version
docker ps
```

ถ้า `docker ps` ติด permission ให้เพิ่ม user เข้า group docker แล้ว logout/login ใหม่:

```sh
sudo usermod -aG docker $USER
```

## วิธี deploy

ใน GitHub:

```text
Actions → Deploy to Azure VM → Run workflow
```

หลัง deploy เสร็จ เปิด:

```text
https://logdisk.malaysiawest.cloudapp.azure.com
```

## หมายเหตุ

- workflow ไม่สร้าง `.env` บน VM เพื่อไม่ให้ secret ไหลผ่าน GitHub Actions โดยไม่จำเป็น
- workflow ไม่รัน database reset และไม่ลบ volume
- ถ้ามีการเปลี่ยนค่า `.env` ให้แก้บน VM โดยตรง แล้วรัน deploy หรือ `docker compose up -d --build` ใหม่

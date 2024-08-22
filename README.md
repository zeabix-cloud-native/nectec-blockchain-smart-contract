## nectec-blockchain-chaincode

### Prerequisite

- [Nodejs](https://nodejs.org/en)
- [Hyperledger Fabric and Fabric sample](https://hyperledger-fabric.readthedocs.io/en/latest/install.html)

### Installation

##### Hyperledger Test Network

1. ในการรัน Hyperledger Fabric ภายในเครื่องจะต้องทำการสร้าง test-network ขึ้นมาโดยให้ cd ไปยัง Folder `test-network` จากการติดตั้งใน Prerequisite

2. จากนั้นให้รันคำสั่ง `export PATH=${PWD}/../bin:$PATH` และ `export PATH=${PWD}/../bin:$PATH` เพื่อให้เครื่องมีคำสั่งของ Hyperledger

3. เพิ่มตัวแปรคำสั่งสำหรับการสร้าง test-network ให้ครบ

```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

4. จากนั้นใช้คำสั่งในการสร้างและ Deploy chaincode

```
./network.sh up createChannel -ca -s couchdb && ./network.sh deployCC -ccn chaincode -ccp <chaincode-path> -ccl go -ccv 1.0.0
```

5. จากนั้นเมื่อรันคำสั่ง docker ps จะเห็น test-network เพิ่มขึ้นมา
   ![Home screen](public/doc-1.png)

##### Nestjs API

###### ตัวแปรสภาพแวดล้อมที่ใช้ (.env)

```
DATABASE_URL=
JWT_ACCESS_SECRET=
JWT_REFRESH_SECRET=
CHANNEL_NAME=
CHAINCODE_NAME=
MSP_ID=
CRYPTO_PATH=
TLS_CERT_PATH=
PEER_ENDPOINT=
CA_HOST_NAME=
PEER_HOST_ALIAS=
KEY_DIRECTORY_PATH=
CERT_PATH=
WALLET_PATH=
CCP_PATH=
AFFILIATION=
ADMIN_IDENTITY_PATH=
ADMIN_CERT_ID=
PORT=3000
```

```
DATABASE_URL

ตัวแปรนี้ใช้เพื่อกำหนด URL สำหรับการเชื่อมต่อฐานข้อมูล MySQL
ตัวอย่าง: mysql://root:1234@localhost:3306/nstda-dev-db?schema=public
root คือชื่อผู้ใช้, 1234 คือรหัสผ่าน, localhost คือโฮสต์, 3306 คือพอร์ตที่ใช้เชื่อมต่อ, และ nstda-dev-db คือชื่อฐานข้อมูล

JWT_ACCESS_SECRET
ค่าลับ (secret key) ที่ใช้ในการเข้ารหัสและถอดรหัส JWT สำหรับการยืนยันตัวตน (Access Token)

JWT_REFRESH_SECRET
ค่าลับ (secret key) ที่ใช้ในการเข้ารหัสและถอดรหัส JWT สำหรับการรีเฟรชโทเค็น (Refresh Token)

CHANNEL_NAME
ชื่อของช่อง (channel) ที่ใช้ในเครือข่าย Hyperledger Fabric ซึ่งเป็นช่องทางสำหรับการสื่อสารระหว่างเพียร์ (peer)

CHAINCODE_NAME
ชื่อของสมาร์ทคอนแทรค (chaincode) ที่ถูกติดตั้งและทำงานบน Hyperledger Fabric
MSP_ID

ID ของการให้บริการสมาชิก (Membership Service Provider - MSP) ที่ใช้ในการระบุตัวตนขององค์กร (organization) บน Hyperledger Fabric
CRYPTO_PATH

TLS_CERT_PATH
เส้นทางไปยังโฟลเดอร์ที่เก็บข้อมูลการเข้ารหัสลับขององค์กร (cryptographic materials) สำหรับ Hyperledger Fabric

PEER_ENDPOINT
เส้นทางไปยังใบรับรอง TLS ที่ใช้เพื่อการสื่อสารแบบปลอดภัย (secure communication) ระหว่างเพียร์ (peer) บน Hyperledger Fabric

CA_HOST_NAME
ที่อยู่และพอร์ตของเพียร์ (peer) ที่ใช้ในการเชื่อมต่อกับ Hyperledger Fabric

PEER_HOST_ALIAS
ชื่อโฮสต์ของผู้ให้บริการใบรับรอง (Certificate Authority - CA) ที่ใช้ใน Hyperledger Fabric

KEY_DIRECTORY_PATH
ชื่อเล่น (alias) ของเพียร์ (peer) ที่กำหนดให้กับโฮสต์

CERT_PATH
เส้นทางไปยังโฟลเดอร์ที่เก็บคีย์ส่วนตัว (private key) ของผู้ใช้สำหรับการลงชื่อดิจิทัล (digital signing) บน Hyperledger Fabric

WALLET_PATH
เส้นทางไปยังใบรับรองดิจิทัล (certificate) ของผู้ใช้ที่ใช้ในการลงชื่อดิจิทัลบน Hyperledger Fabric

CCP_PATH
เส้นทางไปยังโฟลเดอร์ที่เก็บกระเป๋า (wallet) สำหรับจัดเก็บข้อมูลประจำตัว (identity) ของผู้ใช้ใน Hyperledger Fabric

AFFILIATION
เส้นทางไปยังไฟล์ JSON ที่เก็บข้อมูลเชื่อมต่อ (connection profile) ขององค์กรบน Hyperledger Fabric

ADMIN_IDENTITY_PATH
ชื่อของหน่วยงานหรือแผนก (affiliation) ที่ผู้ใช้อยู่ใน Hyperledger Fabric

ADMIN_CERT_ID
เส้นทางไปยังโฟลเดอร์ที่เก็บข้อมูลประจำตัวของผู้ดูแลระบบ (Admin) บน Hyperledger Fabric

PORT
พอร์ตที่ใช้ในการรันแอปพลิเคชัน NestJS
```

1. สร้าง Local mysql database โดยคำสั่ง `docker compose -f "docker-compose-db.yml" up -d --build `

2. รันคำสั่ง `yarn install` เพื่อทำการลง dependencies ของโปรเจ็ค จากนั้นรันคำสั่ง `npx prisma migrate dev` เพื่อทำการสร้างตารางเริ่มต้น

3. จากนั้นรันคำสั่ง `yarn seed` เพื่อทำการเพิ่มข้อมูลเริ่มต้นเช่นข้อมูลของ Admin ภายในระบบ

4. จากนั้นรันคำสั่ง `yarn start:dev` เพื่อเริ่มต้น Application

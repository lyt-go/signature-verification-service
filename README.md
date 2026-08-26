# 数字签名服务（Digital Signature Service）

纯 Go 标准库实现，零第三方依赖。

## 运行说明

```bash
cd origin
go build ./...
go build -o /tmp/signature-server ./cmd/server
/tmp/signature-server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 端点

| 实体 | 方法 | 路径 | 说明 |
|------|------|------|------|
| KeyPair | POST | `/api/keypairs` | 创建密钥对 |
| KeyPair | GET | `/api/keypairs` | 列表（支持 name/algorithm/status/keyword 筛选 + 分页） |
| KeyPair | GET | `/api/keypairs/{id}` | 获取单个 |
| KeyPair | PUT | `/api/keypairs/{id}` | 更新 |
| KeyPair | DELETE | `/api/keypairs/{id}` | 删除 |
| KeyPair | PATCH | `/api/keypairs/{id}/status` | 更新状态（active→revoked/expired） |
| SignRequest | POST | `/api/signrequests` | 创建签名请求 |
| SignRequest | GET | `/api/signrequests` | 列表（支持 key_pair_id/algorithm/status/request_id 筛选 + 分页） |
| SignRequest | GET | `/api/signrequests/{id}` | 获取单个 |
| SignRequest | PUT | `/api/signrequests/{id}` | 更新 |
| SignRequest | DELETE | `/api/signrequests/{id}` | 删除 |
| SignRequest | PATCH | `/api/signrequests/{id}/status` | 更新状态（pending→signed/failed） |
| SignRequest | POST | `/api/signrequests/{id}/process` | 执行签名流程 |
| Signature | POST | `/api/signatures` | 创建签名结果 |
| Signature | GET | `/api/signatures` | 列表（支持 key_pair_id/algorithm/sign_request_id 筛选 + 分页） |
| Signature | GET | `/api/signatures/{id}` | 获取单个 |
| Signature | PUT | `/api/signatures/{id}` | 更新 |
| Signature | DELETE | `/api/signatures/{id}` | 删除 |
| Signature | POST | `/api/signatures/{id}/verify` | 验签流程 |
| VerifyRecord | POST | `/api/verifyrecords` | 创建验签记录 |
| VerifyRecord | GET | `/api/verifyrecords` | 列表（支持 signature_id/valid/verifier 筛选 + 分页） |
| VerifyRecord | GET | `/api/verifyrecords/{id}` | 获取单个 |
| VerifyRecord | PUT | `/api/verifyrecords/{id}` | 更新 |
| VerifyRecord | DELETE | `/api/verifyrecords/{id}` | 删除 |
| VerifyRecord | POST | `/api/verifyrecords/batch` | 批量创建验签记录 |
| Algorithm | POST | `/api/algorithms` | 创建算法定义 |
| Algorithm | GET | `/api/algorithms` | 列表（支持 name/type/enabled 筛选 + 分页） |
| Algorithm | GET | `/api/algorithms/{id}` | 获取单个 |
| Algorithm | PUT | `/api/algorithms/{id}` | 更新 |
| Algorithm | DELETE | `/api/algorithms/{id}` | 删除 |
| Certificate | POST | `/api/certificates` | 创建证书 |
| Certificate | GET | `/api/certificates` | 列表（支持 key_pair_id/status/subject/issuer 筛选 + 分页） |
| Certificate | GET | `/api/certificates/{id}` | 获取单个 |
| Certificate | PUT | `/api/certificates/{id}` | 更新 |
| Certificate | DELETE | `/api/certificates/{id}` | 删除 |
| Certificate | PATCH | `/api/certificates/{id}/status` | 更新状态（valid→revoked/expired） |
| Stats | GET | `/api/stats/overview` | 统计概览 |

共 **36** 个 API 端点。

## 统一响应格式

```json
{"code":0,"message":"ok","data":...}
```

错误码映射：ValidationError → 400、ErrNotFound → 404、ErrConflict → 409、其他 → 500。

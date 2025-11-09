---
id: DC-0001
title: User ドメインモデル
related:
  features: [F-0001, F-0002]
  apis: [POST /auth/signup, GET /users/me]
---

# 🧩 モデル概要

- システム内でユーザーを一意に識別し、認証・権限・プロフィール情報を管理する。

---

## 🧱 スキーマ定義

| フィールド    | 型           | 必須 | 説明                          |
| ------------- | ------------ | ---- | ----------------------------- |
| id            | string(UUID) | ○    | 一意のユーザー ID             |
| email         | string       | ○    | RFC5322 準拠のメール          |
| password_hash | string       | ○    | Bcrypt などのハッシュ済文字列 |
| name          | string       | △    | 表示名（任意）                |
| role          | enum         | ○    | `user` / `admin`              |
| created_at    | datetime     | ○    | 作成日時                      |
| updated_at    | datetime     | ○    | 更新日時                      |

---

## 🔁 関連するエンティティ

| 関連モデル | 関係 | 説明                 |
| ---------- | ---- | -------------------- |
| Session    | 1:N  | 認証セッション       |
| Profile    | 1:1  | ユーザープロフィール |
| AuditLog   | 1:N  | 監査ログ             |

---

## 📘 備考

- パスワードは絶対に平文保存しない。
- メールはユニーク制約必須。
- 将来的に OAuth ログイン統合予定。

# Open Key By Felica
## 概要
Rapsberry Pi + Goを使用したスマートロックリポジトリです。

## 使用ライブラリ
- [go-rpio](https://github.com/stianeikeland/go-rpio)
- [pasori](https://github.com/bamchoh/pasori)

## 使用部品など
- Raspberry Pi 4B 8GB
- Sony RC-S380
- リードスイッチ
- サーボモーター SG92R
- Nch MOSFET TK7R4A10PL

## 機能
- icカードを用いてドアの開閉を行います。

## 接続方法
### 使用pin
- GPIO
  - PwmPin: `13`
  - MosPin: `17`
  - SwPin : `18`
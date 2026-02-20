# インストールガイド

バイナリはリポジトリには含まれていません（`.gitignore` 対象）．
[GitHub Releases](https://github.com/torifo/astral-drive/releases/latest) から各プラットフォーム向けのビルド済みバイナリを取得してください．

## ダウンロード一覧

| ファイル名 | OS |
|---|---|
| `astral-drive-linux-amd64` | Linux (x86_64) |
| `astral-drive-windows-amd64.exe` | Windows (x86_64) |
| `astral-drive-darwin-amd64` | macOS (Intel) |
| `astral-drive-darwin-arm64` | macOS (Apple Silicon) |

---

## Linux / macOS

macOS は Linux と同じ POSIX パス構造のため，手順は共通です．

### ワンライナーインストール

```bash
curl -fsSL https://raw.githubusercontent.com/torifo/astral-drive/main/scripts/install.sh | bash
```

### 手動インストール

```bash
# Linux の場合
curl -L https://github.com/torifo/astral-drive/releases/latest/download/astral-drive-linux-amd64 \
  -o /tmp/astral-drive
chmod +x /tmp/astral-drive
sudo mv /tmp/astral-drive /usr/local/bin/astral-drive

# macOS (Apple Silicon) の場合
curl -L https://github.com/torifo/astral-drive/releases/latest/download/astral-drive-darwin-arm64 \
  -o /tmp/astral-drive
chmod +x /tmp/astral-drive
sudo mv /tmp/astral-drive /usr/local/bin/astral-drive
```

### 動作確認

```bash
astral-drive -n 10 /home
astral-drive -n 10 /          # ルート全体（時間がかかる場合あり）
```

---

## Windows

### PowerShell でインストール（管理者権限不要）

```powershell
$url = "https://github.com/torifo/astral-drive/releases/latest/download/astral-drive-windows-amd64.exe"
$bin = "$env:USERPROFILE\.local\bin"
New-Item -ItemType Directory -Path $bin -Force | Out-Null
Invoke-WebRequest $url -OutFile "$bin\astral-drive.exe"

# PATH に追加（初回のみ・ターミナル再起動後に有効）
[System.Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$bin", "User")
```

### 動作確認

```powershell
# ターミナルを再起動後
astral-drive -n 20 C:\
astral-drive -n 10 D:\
```

---

## アンインストール

**Linux / macOS:**
```bash
sudo rm /usr/local/bin/astral-drive
```

**Windows:**
```powershell
Remove-Item "$env:USERPROFILE\.local\bin\astral-drive.exe"
```

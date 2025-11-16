# Banishing Null with Non-Nullable References

## 概述

「Banishing Null with Non-Nullable References」（消除 Null 與非空引用）是一個程式設計概念，最初由 C# 8.0 引入，旨在透過型別系統減少 null 引用錯誤（NullReferenceException）。這個概念的核心思想是：**明確區分可為空（nullable）和不可為空（non-nullable）的引用型別**。

## 核心概念

### 1. Null Reference Problem（空引用問題）

Tony Hoare（發明 null 引用的人）稱 null 為「十億美元的錯誤」，因為它導致了無數的程式崩潰和安全漏洞。主要問題包括：

- **運行時錯誤**：在運行時才發現空引用錯誤，而非編譯時
- **不明確性**：無法從型別簽名中得知某個值是否可能為 null
- **防禦性編程**：需要大量的 null 檢查程式碼

### 2. 非空引用（Non-Nullable References）的解決方案

非空引用的核心思想是將引用型別分為兩類：

1. **非空引用（Non-nullable）**：預設情況下，引用型別不能為 null
2. **可空引用（Nullable）**：明確標記可以為 null 的引用型別

## 在不同語言中的實現

### C# 8.0 的實現

```csharp
// 在 C# 8.0+ 中啟用 nullable reference types
#nullable enable

// 非空引用：不能為 null
string name = "John";
name = null; // 編譯器警告

// 可空引用：明確可以為 null
string? optionalName = null; // OK
optionalName = "Jane"; // OK

// 使用時需要檢查
if (optionalName != null)
{
    Console.WriteLine(optionalName.Length);
}
```

### Go 語言中的實現

Go 語言採用了不同但相似的方法來處理 null（在 Go 中稱為 `nil`）：

#### 1. 值型別 vs 指標型別

```go
// 值型別：不能為 nil，總是有值
type User struct {
    Name  string  // 永遠不是 nil，可能是空字串 ""
    Age   int     // 永遠不是 nil，可能是 0
}

// 指標型別：可以為 nil
var user *User = nil  // OK，明確表示可能沒有值

// 使用前檢查
if user != nil {
    fmt.Println(user.Name)
}
```

#### 2. 在本專案中的應用

本專案（Trello Backend）展示了 Go 語言中的 null 安全模式：

**範例 1：資料模型中的非空欄位**

```go
// internal/models/user.go
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key"`
    Email        string    `gorm:"unique;not null"`  // 非空欄位
    Name         string    `gorm:"not null"`         // 非空欄位
    PasswordHash string    `gorm:"not null"`         // 非空欄位
    CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
```

在這個模型中：
- `Email`、`Name`、`PasswordHash` 都標記為 `not null`
- 這些欄位使用值型別（`string`），而非指標型別（`*string`）
- 在資料庫層級強制執行非空約束

**範例 2：處理可能不存在的資料**

```go
// internal/repositories/user_repository.go
func (r *userRepository) FindByEmail(email string) (models.User, error) {
    var user models.User
    result := r.db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return models.User{}, result.Error  // 返回錯誤而非 nil
    }
    return user, nil
}
```

這個模式使用 Go 的多返回值來處理「可能不存在」的情況：
- 成功時返回值和 `nil` 錯誤
- 失敗時返回零值和錯誤
- 呼叫者必須檢查錯誤

**範例 3：服務層的錯誤處理**

```go
// internal/services/auth.go
func (s *authService) Login(req models.LoginRequest) (models.AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return models.AuthResponse{}, errors.New("帳號或密碼錯誤")
    }
    
    // user 在這裡保證是有效的，不需要額外的 nil 檢查
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return models.AuthResponse{}, errors.New("帳號或密碼錯誤")
    }
    
    // ... 繼續處理
}
```

## 最佳實踐

### 1. 優先使用值型別

```go
// ✅ 好：使用值型別表示必須存在的資料
type Card struct {
    ID      uint   `gorm:"primaryKey"`
    Title   string `gorm:"not null"`
    Content string // 可以是空字串，但不會是 nil
}

// ❌ 避免：不必要地使用指標
type Card struct {
    ID      *uint
    Title   *string
    Content *string
}
```

### 2. 使用錯誤處理代替 nil 檢查

```go
// ✅ 好：明確的錯誤處理
func GetUser(id uuid.UUID) (User, error) {
    user, err := repository.FindByID(id)
    if err != nil {
        return User{}, fmt.Errorf("user not found: %w", err)
    }
    return user, nil
}

// ❌ 避免：返回指標和隱含的 nil 檢查
func GetUser(id uuid.UUID) *User {
    user := repository.FindByID(id)
    if user == nil {
        return nil
    }
    return user
}
```

### 3. 在資料庫層級強制非空約束

```go
type Board struct {
    ID     uint   `gorm:"primaryKey"`
    Name   string `gorm:"not null"`          // 資料庫層級的約束
    UserID string `gorm:"type:uuid;not null"` // 確保外鍵不為空
}
```

### 4. API 請求驗證

```go
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"` // 必填
    Name     string `json:"name" binding:"required"`        // 必填
    Password string `json:"password" binding:"required,min=6"` // 必填且有最小長度
}
```

## 何時使用指標（允許 nil）

在某些情況下，使用指標和允許 nil 是合理的：

### 1. 可選欄位

```go
type UpdateCardRequest struct {
    Title   *string `json:"title,omitempty"`   // 可選更新
    Content *string `json:"content,omitempty"` // 可選更新
}

// nil 表示「不更新此欄位」
// 非 nil 表示「更新為此值」（包括空字串）
```

### 2. 延遲初始化

```go
type Service struct {
    cache *Cache // 可能稍後才初始化
}

func (s *Service) GetCache() *Cache {
    if s.cache == nil {
        s.cache = NewCache()
    }
    return s.cache
}
```

### 3. 遞迴資料結構

```go
type TreeNode struct {
    Value int
    Left  *TreeNode // 可能沒有左子節點
    Right *TreeNode // 可能沒有右子節點
}
```

## 效益

實施非空引用和明確的 null 安全模式帶來以下好處：

1. **更少的運行時錯誤**：在編譯時或資料庫層級捕獲錯誤
2. **更清晰的程式碼**：型別簽名明確表達意圖
3. **更好的可維護性**：減少防禦性 null 檢查
4. **更安全的重構**：型別系統協助識別問題

## 總結

「Banishing Null with Non-Nullable References」是一個重要的程式設計原則，透過型別系統明確區分可空和非空值。雖然 Go 語言沒有像 C# 那樣的語法糖，但透過以下方式可以達到相同的目標：

- 優先使用值型別而非指標
- 使用錯誤處理機制
- 在資料庫層級強制約束
- 使用驗證標籤確保 API 輸入的有效性

這些實踐可以顯著提高程式碼的安全性和可維護性，減少因 nil/null 引用導致的錯誤。

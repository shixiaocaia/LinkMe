# 短信功能开发文档

## 1. 目标
为LinkMe项目添加短信功能，以实现用户注册、登录、修改密码、危险操作二次验证等功能

## 2. 功能需求

### 2.1 验证码流程
- **开始**：用户请求发送验证码
- **发送验证码**：系统生成验证码，并通过短信发送给用户
- **再次发送验证码**：
    - 判断一分钟内是否已发送过验证码
    - 如果已发送，则提醒用户稍后重试
    - 如果未发送，则发送新的验证码

### 2.2 登录流程
- **提交验证码**：用户输入收到的验证码并提交
- **验证码校验**：
    - 如果验证码不正确，则提示用户重试
    - 如果验证码正确，则进行下一步
- **用户判断**：
    - 查询数据库中是否存在该手机号对应的用户
    - 如果不存在，则使用该手机号注册新用户
    - 如果存在，则登录成功

## 3. 系统设计

### 3.1 短信模块
- 短信模块应作为一个独立的功能模块，以便于在其他业务场景中复用
- 支持更换短信供应商，不依赖于特定供应商的API

### 3.2 限制条件
- 每个手机号码一分钟内只能发送一次验证码
- 验证码有效期为十分钟
- 验证码一旦被使用，则不能再次使用
- 用户连续输入三次错误的验证码后，该验证码失效

### 3.3 存储
- 使用Redis存储验证码，并设置过期时间为十分钟
- 使用Lua脚本保证原子性操作
- 存储的验证码需要使用合适的加密算法进行加密

## 4. 性能需求
- 整个短信功能的实现需要保证流程的流畅性，不影响用户体验

## 5. 安全需求
- 存储的验证码需要进行加密，确保安全性

## 6. 其他需求
- 系统应具备限流功能，防止短信轰炸等恶意行为

## 7. 依赖关系
- 本功能依赖于短信供应商的API(腾讯SMS)
- 本功能依赖于Redis数据库

## 8. 开发计划
- 功能开发：1周
    - 实现短信发送和验证码生成功能
    - 实现验证码存储和校验功能
    - 实现用户注册和登录功能
- 测试和调试：3天
    - 编写单元测试和集成测试


## 9. 测试计划
- 单元测试：针对每个功能模块编写单元测试
- 集成测试：测试模块之间的集成

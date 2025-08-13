# GitHub repository settings và templates
# Đây là file README cho thư mục .github

## 🚀 GitHub Actions Workflows

Repository này sử dụng GitHub Actions để thực hiện CI/CD pipeline với các workflows sau:

### 📋 **CI Workflow (`ci.yml`)**
- **Trigger**: Pull requests và push vào `master`, `main`, `staging`, `develop`
- **Jobs**:
  - 🔍 **Lint**: Code quality check với golangci-lint, go vet, gofmt
  - 🧪 **Test**: Unit tests với coverage, integration tests với MongoDB/MinIO
  - 🔒 **Security**: Security scan với gosec
  - 🔨 **Build**: Multi-platform builds (Linux, Windows, macOS)
  - 🐳 **Docker**: Docker build validation

### 🚀 **CD Workflow (`cd.yml`)**
- **Trigger**: Push vào `master`/`main`, `staging` hoặc manual dispatch
- **Jobs**:
  - 📦 **Build & Push**: Docker image build và push tới GitHub Container Registry
  - 🎯 **Deploy Staging**: Tự động deploy lên staging environment
  - 🚀 **Deploy Production**: Deploy lên production với manual approval
  - ⏪ **Rollback**: Tự động rollback nếu deployment thất bại

### ✅ **PR Check Workflow (`pr-check.yml`)**
- **Trigger**: Pull request events
- **Jobs**:
  - 📝 **PR Validation**: Kiểm tra title format, description, size
  - 🔍 **Code Analysis**: Complexity check, TODO/FIXME scan
  - 📦 **Dependency Check**: Vulnerability scan, unused dependencies
  - ⚡ **Performance Check**: Benchmark tests, binary size analysis

### 🛡️ **Security Workflow (`security.yml`)**
- **Trigger**: Hàng tuần, push code, hoặc PR
- **Jobs**:
  - 🔒 **Gosec Scan**: Static security analysis
  - 📦 **Dependency Scan**: Vulnerability check với govulncheck
  - 🔐 **Secret Scan**: Secret detection với TruffleHog
  - 🐳 **Docker Security**: Container image security với Trivy
  - 🏗️ **Infrastructure Security**: Config security check

### 🎉 **Release Workflow (`release.yml`)**
- **Trigger**: Push tag với format `v*.*.*`
- **Jobs**:
  - ✅ **Validate Release**: Tag format validation
  - 🔨 **Build Release**: Multi-platform binary builds
  - 🐳 **Build Docker**: Multi-arch Docker images
  - 📝 **Generate Changelog**: Automatic release notes
  - 🚀 **Create Release**: GitHub release với assets
  - 📚 **Update Docs**: Version bumps trong documentation

### 🧹 **Cleanup Workflow (`cleanup.yml`)**
- **Trigger**: Hàng ngày lúc 3:00 AM UTC
- **Jobs**:
  - 🔄 **Cleanup Workflows**: Xóa old workflow runs
  - 📦 **Cleanup Artifacts**: Xóa artifacts cũ hơn 1 tuần
  - 🏷️ **Cleanup Releases**: Giữ lại 10 releases gần nhất
  - 📱 **Cleanup Packages**: Giữ lại 20 container image versions

## 🤖 **Dependabot Configuration**

File `dependabot.yml` cấu hình auto-updates cho:
- **Go modules**: Hàng tuần vào thứ 2
- **Docker dependencies**: Hàng tuần vào thứ 3  
- **GitHub Actions**: Hàng tuần vào thứ 4

## 🛠️ **Setup Instructions**

### 1. Repository Secrets
Cần thiết lập các secrets sau trong repository settings:

```bash
# AWS Deployment (nếu dùng AWS)
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
AWS_REGION=your-aws-region

# Notification webhooks (optional)
SLACK_WEBHOOK_URL=your-slack-webhook-url
DISCORD_WEBHOOK_URL=your-discord-webhook-url
```

### 2. Environment Configuration
Tạo các environments trong repository settings:
- **staging**: Staging environment
- **production**: Production environment với protection rules

### 3. Branch Protection
Thiết lập branch protection cho `master`/`main`:
- Require PR reviews
- Require status checks (CI workflow)
- Restrict pushes to branch

## 📊 **Workflow Status**

| Workflow | Status | Description |
|----------|--------|-------------|
| CI | ![CI](https://github.com/your-org/volcanion-hermes/workflows/CI/badge.svg) | Continuous Integration |
| CD | ![CD](https://github.com/your-org/volcanion-hermes/workflows/CD/badge.svg) | Continuous Deployment |
| Security | ![Security](https://github.com/your-org/volcanion-hermes/workflows/Security/badge.svg) | Security Scanning |

## 🎯 **Best Practices**

### Pull Requests
1. **Title Format**: Sử dụng conventional commits
   - `feat: add new feature`
   - `fix: resolve bug`
   - `docs: update documentation`

2. **Size**: Giữ PR nhỏ gọn (<50 files changed)

3. **Description**: Luôn có mô tả chi tiết

### Releases
1. **Semantic Versioning**: Sử dụng format `v1.0.0`
2. **Pre-releases**: Sử dụng `v1.0.0-beta` cho testing
3. **Changelog**: Tự động generate từ commit messages

### Security
1. **No Secrets**: Không commit secrets vào code
2. **Dependencies**: Thường xuyên update dependencies
3. **Scans**: Monitor security scan results

## 🚨 **Troubleshooting**

### Common Issues

1. **CI Failed - Lint Errors**
   ```bash
   # Fix formatting
   gofmt -w .
   
   # Run linter locally
   golangci-lint run
   ```

2. **Security Scan Failed**
   ```bash
   # Run security scan locally
   gosec ./...
   
   # Check vulnerabilities
   govulncheck ./...
   ```

3. **Docker Build Failed**
   ```bash
   # Test Docker build locally
   docker build -t volcanion-hermes:test .
   ```

4. **Release Failed**
   - Kiểm tra tag format (phải là `v1.0.0`)
   - Đảm bảo có permission write cho packages

## 📞 **Support**

Nếu có vấn đề với workflows:
1. Check workflow logs trong Actions tab
2. Review error messages và troubleshooting guide
3. Tạo issue với label `ci/cd` nếu cần hỗ trợ

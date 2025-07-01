# CI/CD Setup Guide

## 🚨 Current Issue: Permission Errors

The release pipeline is failing with:

```
403 Resource not accessible by integration
```

This happens because `GITHUB_TOKEN` only has access to the current repository, but GoReleaser needs to push to `homebrew-tap` and `scoop-bucket` repositories.

## 🔧 Solution: Personal Access Token (PAT)

### Step 1: Create Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Set these permissions:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `write:packages` (Upload packages to GitHub Package Registry)
   - ✅ `read:org` (Read org and team membership, read org projects)

### Step 2: Add Token to Repository Secrets

1. Go to `https://github.com/VapiAI/cli/settings/secrets/actions`
2. Click "New repository secret"
3. Name: `GORELEASER_GITHUB_TOKEN`
4. Value: [paste your PAT]

### Step 3: Enable Homebrew & Scoop Publishing

Once the PAT is configured, update `.goreleaser.yaml`:

```yaml
# Change these from true to false:
brews:
  - skip_upload: false # Enable Homebrew publishing

scoops:
  - skip_upload: false # Enable Scoop publishing
```

## 🧪 Testing the Fix

1. Create a test release:

   ```bash
   git tag -a v0.0.5 -m "Test release with fixed permissions"
   git push origin v0.0.5
   ```

2. Check the Actions tab for successful completion

3. Verify updates in:
   - `VapiAI/homebrew-tap/Formula/vapi-cli.rb`
   - `VapiAI/scoop-bucket/bucket/vapi-cli.json`

## 📊 Current Status

| Component       | Status         | Notes                           |
| --------------- | -------------- | ------------------------------- |
| GitHub Releases | ✅ Working     | Uses default GITHUB_TOKEN       |
| npm Publishing  | ✅ Working     | Manual publish successful       |
| Homebrew        | ⏳ Pending PAT | Disabled until token configured |
| Scoop           | ⏳ Pending PAT | Disabled until token configured |

## 🔍 Verification Commands

After setting up the PAT, test each installation method:

```bash
# Direct download (should work now)
curl -sSL https://github.com/VapiAI/cli/releases/download/v0.0.4/vapi_$(uname -s)_$(uname -m).tar.gz | tar xz

# npm (working)
npm install -g @vapi-ai/cli

# Homebrew (after PAT setup)
brew tap vapi/tap && brew install vapi-cli

# Scoop (after PAT setup)
scoop bucket add vapi https://github.com/VapiAI/scoop-bucket
scoop install vapi-cli
```

## 🔐 Security Best Practices

1. **Limit PAT scope**: Only grant minimum required permissions
2. **Regular rotation**: Rotate the PAT every 90 days
3. **Monitor usage**: Check GitHub audit logs regularly
4. **Team access**: Consider using GitHub Apps for organization-wide access

## 🚀 Alternative: GitHub Apps (Advanced)

For organizations, consider creating a GitHub App instead of using PATs:

1. More secure and granular permissions
2. Better audit trails
3. No personal account dependency

## 📞 Support

If issues persist:

1. Check Actions logs for detailed error messages
2. Verify PAT permissions include all required scopes
3. Ensure repositories are accessible with the PAT
4. Test PAT manually with GitHub API calls

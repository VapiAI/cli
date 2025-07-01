# Slack Release Notifications Setup

This guide shows how to securely set up Slack notifications for CLI releases without exposing webhook URLs.

## 🔒 **Security First**

- ✅ Webhook URLs stored in GitHub Secrets (never in code)
- ✅ Private repository settings
- ✅ Environment variables for sensitive data

## 📋 **Setup Steps**

### 1. Create Slack Webhook

1. **Go to your Slack workspace**
2. **Navigate to:** https://api.slack.com/apps
3. **Create a new app** → "From scratch"
4. **App Name:** `Vapi Release Bot`
5. **Choose workspace:** Your internal Vapi workspace

### 2. Configure Incoming Webhooks

1. **In your app settings** → "Incoming Webhooks"
2. **Toggle "Activate Incoming Webhooks"** → ON
3. **Click "Add New Webhook to Workspace"**
4. **Choose channel:** `#releases` (or create it)
5. **Copy the webhook URL** (looks like: `https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX`)

### 3. Add Webhook to GitHub Secrets

**🚨 CRITICAL: Keep webhook URL secret!**

1. **Go to:** https://github.com/VapiAI/cli/settings/secrets/actions
2. **Click "New repository secret"**
3. **Name:** `SLACK_WEBHOOK`
4. **Value:** [paste your webhook URL]
5. **Click "Add secret"**

### 4. Customize Notification Settings

Edit `.goreleaser.yaml` to customize:

```yaml
announce:
  slack:
    enabled: true
    channel: "#releases" # Change to your preferred channel
    username: "Vapi Release Bot" # Change bot name
    icon_emoji: ":rocket:" # Change emoji
    message_template: | # Customize message format
      🚀 **Vapi CLI {{.Tag}} Released!**

      📦 **Installation:**
      • Universal: `curl -sSL https://vapi.ai/install.sh | bash`
      • npm: `npm install -g @vapi-ai/cli`

      🔗 **Release:** {{ .ReleaseURL }}
```

## 📱 **Notification Preview**

When a release is published, your team will see:

```
🚀 Vapi CLI v0.0.5 Released!

📦 Distribution:
• Universal: curl -sSL https://vapi.ai/install.sh | bash
• npm: npm install -g @vapi-ai/cli
• Homebrew: brew tap vapi/tap && brew install vapi-cli
• Docker: docker run -it ghcr.io/vapiai/cli:v0.0.5 --help

🔗 Links:
• Release: https://github.com/VapiAI/cli/releases/tag/v0.0.5
• Changelog: https://github.com/VapiAI/cli/releases/tag/v0.0.5#changelog

Built with ❤️ by the Vapi team
```

## 🧪 **Testing**

### Test the Slack App

1. **In Slack app settings** → "Incoming Webhooks"
2. **Scroll down to "Sample curl request"**
3. **Run the curl command** to test posting

### Test with GoReleaser

```bash
# Test locally (without publishing)
goreleaser release --snapshot --clean --skip-publish
```

## 🔧 **Advanced Configuration**

### Multiple Channels

```yaml
announce:
  slack:
    enabled: true
    channel: "#releases,#engineering" # Multiple channels
```

### Conditional Notifications

```yaml
announce:
  slack:
    enabled: true
    message_template: |
      {{if not .Prerelease}}
      🚀 **Vapi CLI {{.Tag}} Released!**
      {{else}}
      🧪 **Vapi CLI {{.Tag}} Pre-release**
      {{end}}
```

### Rich Formatting

````yaml
announce:
  slack:
    enabled: true
    message_template: |
      🚀 *Vapi CLI {{.Tag}} Released!*

      *📦 Quick Install:*
      ```
      curl -sSL https://vapi.ai/install.sh | bash
      ```

      <{{ .ReleaseURL }}|View Release Notes>
````

## 🔐 **Security Best Practices**

### ✅ **Do This:**

- Store webhook URL in GitHub Secrets
- Use repository secrets (not environment secrets)
- Limit webhook permissions to specific channels
- Rotate webhook URLs periodically

### ❌ **Don't Do This:**

- Put webhook URLs in code
- Share webhook URLs in public channels
- Use personal Slack apps for team notifications
- Commit secrets or tokens

## 🛠️ **Troubleshooting**

### **Notifications Not Appearing**

1. **Check GitHub Actions logs** for Slack errors
2. **Verify secret name** matches exactly: `SLACK_WEBHOOK`
3. **Test webhook manually** with curl
4. **Check Slack app permissions**

### **Wrong Channel**

1. **Update `channel` setting** in `.goreleaser.yaml`
2. **Ensure bot has access** to target channel
3. **Public channels:** Use `#channel-name`
4. **Private channels:** Invite bot first

### **Formatting Issues**

1. **Use Slack markdown** format
2. **Test templates** with GoReleaser snapshots
3. **Check GoReleaser docs** for available variables

## 📊 **Monitoring**

Track notification delivery:

- GitHub Actions logs
- Slack app dashboard
- Team feedback

## 🔄 **Maintenance**

- **Monthly:** Review webhook usage in Slack app dashboard
- **Quarterly:** Rotate webhook URLs for security
- **Yearly:** Review and update notification content

---

**Need help?** Check the [GoReleaser Slack documentation](https://goreleaser.com/customization/announce/slack/) or ask in #engineering.

# SDK Auto-Update System

This document explains the automated system for keeping the CLI in sync with the [Vapi Go SDK](https://github.com/VapiAI/server-sdk-go).

## Overview

The CLI depends on the Vapi Go SDK for API data model definitions. When the SDK is updated with new API changes (like new voice options), the CLI needs to be updated to remain compatible. This system automates that process using **webhooks** for instant updates.

## How It Works

### 1. Webhook-Based Updates (`sdk-webhook-update.yml`)

- **Trigger**: Instantly when SDK repository publishes a release
- **Method**: GitHub `repository_dispatch` webhooks
- **Process**:
  1. SDK releases new version → webhook fired
  2. CLI repository receives webhook payload
  3. Workflow extracts version and release info
  4. Creates automated PR with full context

### 2. Automated Updates

When a new SDK version is detected:

1. **Update Dependencies**:

   - Updates `go.mod` with latest SDK version
   - Runs `go mod tidy` to resolve dependencies

2. **Version Bump**:

   - Increments CLI patch version (e.g., `0.0.8` → `0.0.9`)
   - Updates `VERSION` file
   - Updates version references in code

3. **Quality Checks**:

   - Runs `make test` to ensure compatibility
   - Runs `make lint` to maintain code quality

4. **Create PR**:
   - Creates a pull request with detailed changelog
   - Labels it as `sdk-update` for automation tracking
   - Requires manual review and approval

### 3. Automatic Release (`auto-release-after-sdk-update.yml`)

When an SDK update PR is merged:

1. **Tag Creation**:

   - Reads version from `VERSION` file
   - Creates a Git tag (e.g., `v0.0.9`)
   - Pushes tag to trigger release workflow

2. **Release Process**:
   - Existing `release.yml` workflow is triggered by the new tag
   - GoReleaser builds binaries for all platforms
   - Publishes to GitHub Releases, npm, Homebrew, Docker, etc.

## Files and Components

### Workflows

- `.github/workflows/sdk-webhook-update.yml` - Webhook-triggered SDK updates and PR creation
- `.github/workflows/sdk-fallback-check.yml` - Weekly safety net in case webhooks fail
- `.github/workflows/auto-release-after-sdk-update.yml` - Automatic release after PR merge
- `.github/workflows/release.yml` - Existing release workflow (unchanged)

### Scripts

- `scripts/check-sdk-update.sh` - Manual SDK update checking and testing

### Configuration

- `go.mod` - Contains SDK version dependency
- `VERSION` - CLI version number
- `.goreleaser.yaml` - Release configuration (unchanged)

## Manual Testing

You can test the SDK update process locally:

```bash
# Check for updates (read-only)
./scripts/check-sdk-update.sh

# Perform update locally
./scripts/check-sdk-update.sh --update
```

## Manual Workflow Triggers

### Force SDK Update

You can manually trigger an SDK update:

1. Go to **Actions** → **SDK Webhook Update**
2. Click **Run workflow**
3. Optionally specify SDK version or check **Force update**
4. Click **Run workflow**

### Test Webhook

Test the webhook system:

```bash
gh workflow run sdk-webhook-update.yml --field force_update=true
```

### Manual Release

If needed, you can manually create a release:

1. Update `VERSION` file with new version
2. Commit and push changes
3. Create and push a tag: `git tag v0.0.X && git push origin v0.0.X`

## Monitoring and Troubleshooting

### Notification Channels

- **GitHub**: PR comments and workflow status
- **Slack**: Release announcements (configured in `.goreleaser.yaml`)

### Common Issues

1. **Test Failures**: SDK update breaks existing functionality

   - PR will fail CI checks
   - Manual intervention required to fix compatibility

2. **Lint Failures**: Code quality issues after update

   - Fix linting errors before merging PR

3. **Version Conflicts**: Tag already exists
   - Manual version bump may be needed

### Monitoring Commands

```bash
# Check current SDK version
grep "github.com/VapiAI/server-sdk-go" go.mod

# Check latest SDK release
curl -s https://api.github.com/repos/VapiAI/server-sdk-go/releases/latest | jq -r '.tag_name'

# View recent webhook workflows
gh run list --workflow="sdk-webhook-update.yml"

# View fallback check status
gh run list --workflow="sdk-fallback-check.yml"
```

## Benefits

1. **Instant Updates**: No waiting - updates trigger immediately on SDK release
2. **Resource Efficient**: No wasteful daily checks when nothing has changed
3. **Quality Assurance**: Automated testing ensures stability
4. **Rich Context**: PRs include full SDK release notes and changes
5. **Reliable**: Webhook + weekly fallback ensures nothing is missed
6. **Audit Trail**: Clear history of what changed and when

## Configuration

### Required Secrets

- `GITHUB_TOKEN` (automatically provided)
- `GORELEASER_GITHUB_TOKEN` (optional, for enhanced release features)
- `SLACK_WEBHOOK` (optional, for release notifications)

### Webhook Setup

**Important**: To enable instant updates, you need to set up webhooks in the SDK repository.

📖 **See [WEBHOOK_SETUP.md](./WEBHOOK_SETUP.md) for complete setup instructions.**

### Fallback Schedule

The safety net fallback check runs weekly. To change frequency, modify `sdk-fallback-check.yml`:

```yaml
schedule:
  # Currently: Weekly on Sundays at 3 AM UTC
  - cron: "0 3 * * 0"

  # Examples:
  # - cron: '0 3 * * 1,3,5'  # Mon/Wed/Fri at 3 AM
  # - cron: '0 */12 * * *'   # Every 12 hours
```

## Disaster Recovery

If the automation fails:

1. **Manual Update**:

   ```bash
   go get github.com/VapiAI/server-sdk-go@latest
   go mod tidy
   make test
   make lint
   ```

2. **Manual Release**:

   ```bash
   # Bump version
   echo "0.0.X" > VERSION
   git add VERSION
   git commit -m "chore: bump version"
   git tag v0.0.X
   git push origin main v0.0.X
   ```

3. **Reset Automation**:
   - Check GitHub Actions permissions
   - Verify workflow files are correct
   - Re-run failed workflows manually

This system ensures your CLI never falls behind the API changes, maintaining a smooth experience for users.
